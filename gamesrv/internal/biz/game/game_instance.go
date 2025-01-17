package game

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/emulator"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/codec"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pion/webrtc/v3"
	"image"
	"sync"
	"time"
)

// 消息类型定义
const (
	MsgPlayerControlButtonPressed byte = iota
	MsgPlayerControlButtonReleased
	MsgChat
	MsgNewConn
	MsgCloseConn
	MsgPeerConnected
	MsgPauseEmulator
	MsgResumeEmulator
	MsgSaveGame
	MsgLoadSave
	MsgRestartEmulator
	MsgPing              byte = 14
	MsgSetGraphicOptions byte = 15
	MsgSetEmulatorSpeed  byte = 16
	MsgEmulatorRestart   byte = 17
	MsgSetController     byte = 18
	MsgResetController   byte = 19
	MsgGetController     byte = 20
	MsgGetGraphicOptions byte = 21
)

const (
	EnhanceFrameScale      = 2
	DefaultAudioSampleRate = 48000
)

type ConsumerResult struct {
	Success bool
	Error   error
	Message string
	Data    any
}

type Message struct {
	Type      byte  `json:"type"`
	From      int64 `json:"from"`
	To        int64 `json:"to"`
	Timestamp int64 `json:"timestamp"`
	Data      any   `json:"data"`

	resultChan chan ConsumerResult
}

type GraphicOptions struct {
	HighResOpen  bool `json:"highResOpen"`
	ReverseColor bool `json:"reverseColor"`
	Grayscale    bool `json:"grayscale"`
}

type Instance struct {
	RoomId               int64
	EmulatorCode         string             // 模拟器唯一标识码,用于区分不同模拟器
	e                    emulator.IEmulator // 模拟器接口
	game                 string
	videoEncoder         codec.IVideoEncoder
	audioSampleRate      int
	audioEncoder         codec.IAudioEncoder
	audioSampleChan      chan float32             // audioSampleChan 音频输出channel
	messageChan          chan *Message            // 消息接收通道，单线程处理多个客户端发送的消息
	Cancel               context.CancelFunc       // 消息接收和音频接收取消函数
	connections          map[int64]*Connection    // 连接列表
	mutex                *sync.RWMutex            // 连接列表mutex
	createTime           time.Time                // 实例创建时间
	allConnCloseCallback func(instance *Instance) // 所有连接关闭后回调，用于异步释放房间会话
	enhancedFrame        *image.YCbCr             // 高分辨率画面缓存
	enhanceFrameOpen     bool
	frameEnhancer        func(frame emulator.IFrame) emulator.IFrame // 高分辨率画面生成器
	reverseColorOpen     bool
	grayscaleOpen        bool
	lastFrameTime        time.Time

	roomInstanceId int64
	sessionKey     string
	DoneChan       chan struct{}
	emulatorId     int64
	gameId         int64
	controllerMap  map[int]int64
}

// MakeGameInstance 创建初始的游戏实例，其中运行dummy模拟器，该模拟器只输出一个提示玩家选择游戏的单帧画面（后期考虑动画）
func MakeGameInstance(roomId, emulatorId, gameId int64, emulatorCode string, gameData []byte) (*Instance, error) {
	instance := &Instance{
		RoomId:      roomId,
		messageChan: make(chan *Message),
		connections: make(map[int64]*Connection),
		mutex:       &sync.RWMutex{},
		createTime:  time.Now(),
		frameEnhancer: func(frame emulator.IFrame) emulator.IFrame {
			return frame
		},
		allConnCloseCallback: func(instance *Instance) {},
		lastFrameTime:        time.Now(),
		emulatorId:           emulatorId,
		gameId:               gameId,
		EmulatorCode:         emulatorCode,
	}

	// 创建dummy模拟器，输出静止介绍画面
	options, err := emulator.MakeBaseEmulatorOptions("gopher.gif", gameData, 0, instance.audioSampleChan, func(frame emulator.IFrame) {
		instance.RenderCallback(frame, log.NewHelper(log.DefaultLogger))
	})
	e, err := emulator.MakeEmulator(emulatorCode, options)
	if err != nil {
		return nil, err
	}
	width, height := e.OutputResolution()
	// 创建视频和音频编码器
	videoEncoder, err := codec.NewVideoEncoder("vp8", width, height)
	if err != nil {
		return nil, err
	}
	audioEncoder, err := codec.NewAudioEncoder(DefaultAudioSampleRate)
	if err != nil {
		return nil, err
	}
	instance.videoEncoder = videoEncoder
	instance.audioEncoder = audioEncoder
	instance.e = e
	if err := e.Start(); err != nil {
		return nil, err
	}
	// 暂停，等待新连接继续
	_ = e.Pause()
	return instance, nil
}

// MessageHandler 消息处理循环，游戏实例启动后在单独的goroutine中消费并处理消息
// 使模拟器的各种事件单线程处理，避免复杂的多线程控制
func (g *Instance) MessageHandler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-g.messageChan:
			switch msg.Type {
			case MsgPlayerControlButtonPressed:
				fallthrough
			case MsgPlayerControlButtonReleased:
				keyCode := msg.Data.(string)
				g.handlePlayerControl(keyCode, msg.Type, msg.From)
			case MsgNewConn:
				msg.resultChan <- g.handleMsgNewConn(msg.Data.(*Connection))
			case MsgPeerConnected:
				g.handlePeerConnected(msg.Data.(*Connection))
			case MsgCloseConn:
				g.handleMsgCloseConn(msg.Data.(*Connection))
			case MsgSaveGame:
				msg.resultChan <- g.handleSaveGame()
			case MsgLoadSave:
				msg.resultChan <- g.handleLoadSave(msg.Data.(*emulatorLoadSaveRequest))
			case MsgRestartEmulator:
				msg.resultChan <- g.handleRestartEmulator(msg.Data.(*emulatorRestartRequest))
			case MsgChat:
				g.handleChat(msg)
			case MsgSetGraphicOptions:
				msg.resultChan <- g.setGraphicOptions(msg.Data.(*GraphicOptions))
			case MsgGetGraphicOptions:
				msg.resultChan <- ConsumerResult{Success: true, Data: g.getGraphicOptions()}
			case MsgSetEmulatorSpeed:
				msg.resultChan <- g.setEmulatorSpeed(msg.Data.(float64))
			case MsgSetController:
				g.handleSetController(msg.Data.([]ControllerPlayer))
				msg.resultChan <- ConsumerResult{
					Success: true,
					Data:    g.getControllerPlayer(),
				}
			case MsgResetController:
			case MsgGetController:
				players := g.getControllerPlayer()
				msg.resultChan <- ConsumerResult{
					Success: true,
					Data:    players,
				}
			default:
			}
		}
	}
}

func (g *Instance) GetConnection(userId int64) (*Connection, bool) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	conn, ok := g.connections[userId]
	return conn, ok
}

// onConnected webrtc连接状态变为connected时调用，用于控制模拟器的自启动和暂停
func (g *Instance) onConnected(conn *Connection) {
	g.messageChan <- &Message{
		Type: MsgPeerConnected,
		Data: conn,
	}
}

// onDataChannelMessage 用户连接数据通道接收消息回调，消息交给消息处理循环处理
func (g *Instance) onDataChannelMessage(userId int64, msg *Message) {
	msg.From = userId
	g.messageChan <- msg
}

// closeConnection webrtc连接状态变为disconnect时调用，用于关闭连接和模拟器自动暂停
func (g *Instance) closeConnection(conn *Connection) {
	g.messageChan <- &Message{
		Type: MsgCloseConn,
		Data: conn,
	}
}

// handleMsgNewConn 添加新连接，踢下线同一个用户的旧连接
func (g *Instance) handleMsgNewConn(conn *Connection) ConsumerResult {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	// 踢下线旧连接，保证同一游戏实例每个用户连接唯一
	if old, ok := g.connections[conn.userId]; ok {
		old.Close()
		delete(g.connections, conn.userId)
	}
	g.connections[conn.userId] = conn
	var controllers []int
	for c, v := range g.controllerMap {
		if v == conn.userId {
			controllers = append(controllers, c)
		}
	}
	for _, c := range controllers {
		delete(g.controllerMap, c)
	}
	return ConsumerResult{Success: true}
}

// handlePeerConnected 新连接建立成功，判断是否需要启动模拟器
func (g *Instance) handlePeerConnected(_ *Connection) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	active := g.filterConnection(func(c *Connection) bool {
		return c.pc.ConnectionState() == webrtc.PeerConnectionStateConnected
	})
	// 首个活跃连接，开启模拟器
	if len(active) == 1 {
		_ = g.e.Resume()
	}
}

// handleMsgCloseConn 连接关闭，判断是否需要暂停模拟器
// TODO 如果长时间没有活跃连接，销毁游戏实例，避免一台服务器运行过多实例
func (g *Instance) handleMsgCloseConn(conn *Connection) {
	g.mutex.Lock()
	// 被动关闭连接，可能是因为新连接挤掉旧连接，需要避免删除新连接
	if cur, ok := g.connections[conn.userId]; ok {
		if cur.pc.ConnectionState() == webrtc.PeerConnectionStateClosed ||
			cur.pc.ConnectionState() == webrtc.PeerConnectionStateFailed ||
			cur.pc.ConnectionState() == webrtc.PeerConnectionStateDisconnected {
			delete(g.connections, conn.userId)
		}
	}
	active := g.filterConnection(func(conn *Connection) bool {
		return conn.pc.ConnectionState() == webrtc.PeerConnectionStateConnected
	})
	// 在Pause之前必须释放连接列表锁，避免 模拟器goroutine和messageConsumer死锁
	// 死锁循环等待：模拟器RenderCallback等待获取g.mutex, 之后消费processor.channel(无缓冲通道)
	//             closeConn获取到了g.mutex, 之后向processor.channel发送消息。
	// 模拟器等待g.mutex, closeConn等待processor.channel，循环等待导致死锁
	g.mutex.Unlock()
	// 没有活跃连接，暂停模拟器
	if len(active) == 0 {
		_ = g.e.Pause()
		g.allConnCloseCallback(g)
	}
}

func (g *Instance) filterConnection(filter func(*Connection) bool) []*Connection {
	result := make([]*Connection, 0, len(g.connections))
	for _, conn := range g.connections {
		if filter(conn) {
			result = append(result, conn)
		}
	}
	return result
}

func (g *Instance) GetOnlineUsers() []int64 {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	result := make([]int64, 0, len(g.connections))
	for userId := range g.connections {
		result = append(result, userId)
	}
	return result
}
