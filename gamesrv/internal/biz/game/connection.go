package game

import (
	"encoding/json"
	"fmt"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pion/webrtc/v3"
	"sync"
)

type Connection struct {
	pc          *webrtc.PeerConnection         // webrtc 连接
	videoTrack  *webrtc.TrackLocalStaticSample // 视频轨道，发送模拟器画面输出
	audioTrack  *webrtc.TrackLocalStaticSample // 音频轨道，发送模拟器音频输出
	dataChannel *webrtc.DataChannel            // 数据通道，接收客户端操作、心跳，发送服务端事件（重启、下线）
	userId      int64                          // 连接的用户id

	localCandidates []*webrtc.ICECandidate // 本地ICE候选地址
	mutex           *sync.Mutex            // ICE候选地址 锁
}

type ConnectionFactory struct {
	api        *webrtc.API
	iceServers []webrtc.ICEServer
}

func NewConnectionFactory(c *conf.WebRTC) *ConnectionFactory {
	if c.PortMin <= 0 || c.PortMax > 65535 {
		panic("invalid port range")
	}
	servers := make([]webrtc.ICEServer, len(c.IceServers))
	for i, server := range c.IceServers {
		servers[i] = webrtc.ICEServer{
			URLs:       []string{server.Url},
			Username:   server.Username,
			Credential: server.Credential,
		}
	}
	setting := webrtc.SettingEngine{}
	// 设置允许的UDP端口范围
	if err := setting.SetEphemeralUDPPortRange(uint16(c.PortMin), uint16(c.PortMax)); err != nil {
		panic(err)
	}
	media := webrtc.MediaEngine{}
	if err := media.RegisterDefaultCodecs(); err != nil {
		panic(err)
	}
	api := webrtc.NewAPI(webrtc.WithSettingEngine(setting), webrtc.WithMediaEngine(&media))
	return &ConnectionFactory{
		api:        api,
		iceServers: servers,
	}
}

// NewConnection PeerConn工厂创建新连接
func (cf *ConnectionFactory) NewConnection(userId int64, g *Instance) (*Connection, string, error) {
	pc, err := cf.api.NewPeerConnection(webrtc.Configuration{
		ICEServers: cf.iceServers,
	})
	if err != nil {
		return nil, "", fmt.Errorf("new peer connection error: %v", err)
	}

	defer func() {
		if err := recover(); err != nil {
			log.NewHelper(log.DefaultLogger).Errorf("new peer connection error: %v", err)
			_ = pc.Close()
		}
	}()

	// 创建音视频轨道，TODO 可配置编码器
	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "CloudEmuVideo")
	if err != nil {
		panic(fmt.Errorf("create video track error: %v", err))
	}
	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "CloudEmuAudio")
	if err != nil {
		panic(fmt.Errorf("create audio track error: %v", err))
	}
	if _, err := pc.AddTrack(videoTrack); err != nil {
		panic(fmt.Errorf("add video track error: %v", err))
	}
	if _, err := pc.AddTrack(audioTrack); err != nil {
		panic(fmt.Errorf("add audio track error: %v", err))
	}
	// 创建数据通道，用于发送客户端操作、心跳，接收服务端事件（重启、下线）
	dataChannel, err := pc.CreateDataChannel("control-channel", nil)
	if err != nil {
		panic(fmt.Errorf("create data channel error: %v", err))
	}
	// 创建SDP Offer
	offer, err := pc.CreateOffer(nil)
	if err != nil {
		panic(fmt.Errorf("create sdp offer error: %v", err))
	}
	if err := pc.SetLocalDescription(offer); err != nil {
		panic(fmt.Errorf("set local description error: %v", err))
	}

	conn := &Connection{
		pc:          pc,
		videoTrack:  videoTrack,
		audioTrack:  audioTrack,
		dataChannel: dataChannel,
		userId:      userId,
		mutex:       &sync.Mutex{},
	}
	// 本地获取到ICE候选地址
	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			g.mutex.Lock()
			conn.localCandidates = append(conn.localCandidates, candidate)
			g.mutex.Unlock()
		}
	})

	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		conn.OnPeerConnectionState(state, g)
	})
	pc.OnICEConnectionStateChange(conn.OnICEStateChange)
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		conn.OnDataChannelMessage(g, msg)
	})
	// 向游戏实例添加连接
	rch := make(chan ConsumerResult)
	g.messageChan <- &Message{
		Type:       MsgNewConn,
		Data:       conn,
		resultChan: rch,
	}
	<-rch
	return conn, offer.SDP, nil
}

// OnPeerConnectionState 处理连接状态变化
func (c *Connection) OnPeerConnectionState(state webrtc.PeerConnectionState, instance *Instance) {
	switch state {
	case webrtc.PeerConnectionStateConnected:
		instance.onConnected(c)
	case webrtc.PeerConnectionStateFailed:
		instance.closeConnection(c)
	case webrtc.PeerConnectionStateDisconnected:
		instance.closeConnection(c)
	case webrtc.PeerConnectionStateClosed:
		instance.closeConnection(c)
	default:
	}
}

// OnICEStateChange 处理ICE连接状态变化
func (c *Connection) OnICEStateChange(_ webrtc.ICEConnectionState) {

}

// OnDataChannelMessage 处理数据通道消息
func (c *Connection) OnDataChannelMessage(instance *Instance, msg webrtc.DataChannelMessage) {
	m := &Message{}
	err := json.Unmarshal(msg.Data, m)
	if err != nil {
		return
	}
	if m.Type == MsgPing {
		// 心跳消息，直接返回
		_ = c.dataChannel.Send(msg.Data)
	} else {
		instance.onDataChannelMessage(c.userId, m)
	}
}

func (c *Connection) Close() {
	_ = c.dataChannel.Close()
	_ = c.pc.Close()
}

func (c *Connection) GetLocalICECandidates() []string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	result := make([]string, len(c.localCandidates))
	for i, candidate := range c.localCandidates {
		bytes, _ := json.Marshal(candidate.ToJSON())
		result[i] = string(bytes)
	}
	return result
}

func (c *Connection) SetRemoteDescription(description webrtc.SessionDescription) error {
	return c.pc.SetRemoteDescription(description)
}

func (c *Connection) AddICECandidate(candidate webrtc.ICECandidateInit) error {
	return c.pc.AddICECandidate(candidate)
}
