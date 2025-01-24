package game

import (
	"context"
	"fmt"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"time"
)

func (g *Instance) ListenAudioChan(ctx context.Context, left bool) {
	// 每5毫秒发送一次，根据采样率计算buffer大小
	buffer := make([]float32, 0, int64(float64(g.audioSampleRate)*5.0/1000.0))
	audioChan := g.rightAudioChan
	audioEncoder := g.audioEncoder1
	if left {
		audioChan = g.leftAudioChan
		audioEncoder = g.audioEncoder0
	}
	for {
		select {
		case <-ctx.Done():
			return
		case s := <-audioChan:
			buffer = append(buffer, s)
			if len(buffer) == cap(buffer) {
				data, err := audioEncoder.Encode(buffer)
				if err != nil {
					fmt.Println("encode error:", err)
				}
				g.sendAudioSamples(data, left)
				buffer = buffer[:0]
			}
		}
	}
}

func (g *Instance) sendAudioSamples(data []byte, left bool) {
	sample := media.Sample{Data: data, Timestamp: time.Now(), Duration: time.Millisecond * 5}
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	for _, conn := range g.connections {
		if conn.pc.ConnectionState() != webrtc.PeerConnectionStateConnected {
			continue
		}
		if left {
			err := conn.leftAudioTrack.WriteSample(sample)
			if err != nil {
				fmt.Println("write sample error:", err)
			}
		} else {
			err := conn.rightAudioTrack.WriteSample(sample)
			if err != nil {
				fmt.Println("write sample error:", err)
			}
		}
	}
}
