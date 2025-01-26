package game

import (
	"fmt"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"time"
)

func (g *Instance) sendAudioSamples(data []byte, left bool) {
	sample := media.Sample{Data: data, Timestamp: time.Now(), Duration: time.Millisecond * 15}
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

func (g *Instance) AudioConsumer(samples []float32) {
	data, err := g.audioEncoder0.Encode(samples)
	if err != nil {
		fmt.Println("encode error:", err)
		return
	}
	sample := media.Sample{Data: data, Timestamp: time.Now(), Duration: time.Millisecond * 15}
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	for _, conn := range g.connections {
		if conn.pc.ConnectionState() != webrtc.PeerConnectionStateConnected {
			continue
		}
		if err := conn.leftAudioTrack.WriteSample(sample); err != nil {
			fmt.Println("write sample error:", err)
		}
	}
}
