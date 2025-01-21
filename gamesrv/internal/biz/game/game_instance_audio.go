package game

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"time"
)

func (g *Instance) AudioSampleListener(ctx context.Context, logger *log.Helper) {
	// 每5毫秒发送一次，根据采样率计算buffer大小
	buffer := make([]float32, 0, g.audioSampleRate*5/1000)
	for {
		select {
		case <-ctx.Done():
			return
		case s := <-g.audioSampleChan:
			buffer = append(buffer, s)
			if len(buffer) == cap(buffer) {
				g.sendAudioSamples(buffer, logger)
				buffer = buffer[:0]
			}
		}
	}
}

func (g *Instance) sendAudioSamples(buffer []float32, logger *log.Helper) {
	data, err := g.audioEncoder.Encode(buffer)
	if err != nil {
		logger.Error("encode audio samples error: ", err)
	}
	sample := media.Sample{Data: data, Timestamp: time.Now(), Duration: time.Millisecond * 5}
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	for _, conn := range g.connections {
		if conn.pc.ConnectionState() != webrtc.PeerConnectionStateConnected {
			continue
		}
		err := conn.audioTrack.WriteSample(sample)
		if err != nil {
			logger.Errorf("write sample error: %v", err)
		}
	}
}
