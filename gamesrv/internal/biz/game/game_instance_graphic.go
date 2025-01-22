package game

import (
	"github.com/StellrisJAY/cloud-emu/emulator"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/codec"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"time"
)

func (g *Instance) RenderCallback(frame emulator.IFrame, logger *log.Helper) {
	data, release, err := g.videoEncoder.Encode(frame)
	if err != nil {
		logger.Error("encode frame error", "err", err)
		return
	}
	defer release()
	frameTime := time.Now()
	duration := frameTime.Sub(g.lastFrameTime)
	// 模拟器暂停后可能导致上一帧距离当前帧时间过长，将该帧的时长设置为固定的20ms
	if duration >= 500*time.Millisecond {
		duration = 20 * time.Millisecond
	}
	g.lastFrameTime = frameTime
	sample := media.Sample{Data: data, Duration: duration, Timestamp: time.Now()}
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	for _, conn := range g.connections {
		if conn.pc.ConnectionState() != webrtc.PeerConnectionStateConnected {
			continue
		}
		if err := conn.videoTrack.WriteSample(sample); err != nil {
			logger.Errorf("write sample error: %v", err)
		}
	}
}

func (g *Instance) SetGraphicOptions(options *GraphicOptions) {
	resultCh := make(chan ConsumerResult)
	g.messageChan <- &Message{Type: MsgSetGraphicOptions, Data: options, resultChan: resultCh}
	<-resultCh
	close(resultCh)
}

func (g *Instance) GetGraphicOptions() *GraphicOptions {
	resultCh := make(chan ConsumerResult)
	g.messageChan <- &Message{Type: MsgGetGraphicOptions, resultChan: resultCh}
	res := <-resultCh
	close(resultCh)
	return res.Data.(*GraphicOptions)
}

func (g *Instance) setGraphicOptions(options *GraphicOptions) ConsumerResult {
	_ = g.e.Pause()
	defer g.e.Resume()

	g.e.SetGraphicOptions(&emulator.GraphicOptions{
		//Grayscale:    options.Grayscale,
		//ReverseColor: options.ReverseColor,
		HighResolution: options.HighResOpen,
	})
	w, h := g.e.OutputResolution()
	enc, err := codec.NewVideoEncoder("vp8", w, h)
	if err != nil {
		return ConsumerResult{Success: false, Error: err}
	}
	g.videoEncoder.Close()
	g.videoEncoder = enc
	return ConsumerResult{Success: true, Error: nil}
}

func (g *Instance) getGraphicOptions() *GraphicOptions {
	options := g.e.GetGraphicOptions()
	return &GraphicOptions{
		HighResOpen: options.HighResolution,
	}
}

func (g *Instance) SetEmulatorSpeed(boostRate float64) float64 {
	resultCh := make(chan ConsumerResult)
	g.messageChan <- &Message{Type: MsgSetEmulatorSpeed, Data: boostRate, resultChan: resultCh}
	result := <-resultCh
	close(resultCh)
	return result.Data.(float64)
}

func (g *Instance) setEmulatorSpeed(boostRate float64) ConsumerResult {
	_ = g.e.Pause()
	defer g.e.Resume()
	rate := g.e.SetCPUBoostRate(boostRate)
	return ConsumerResult{Success: true, Error: nil, Data: rate}
}

func (g *Instance) GetEmulatorSpeed() float64 {
	return g.e.GetCPUBoostRate()
}
