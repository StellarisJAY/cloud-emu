package game

import (
	"github.com/StellrisJAY/cloud-emu/emulator"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/codec"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"image"
	"time"
)

func (g *Instance) enhanceFrame(frame emulator.IFrame) emulator.IFrame {
	if g.enhancedFrame == nil {
		g.enhancedFrame = image.NewYCbCr(image.Rect(0, 0, frame.Width()*EnhanceFrameScale, frame.Height()*EnhanceFrameScale), image.YCbCrSubsampleRatio420)
	}
	original := frame.YCbCr()
	enhancedFrame := g.enhancedFrame
	for y := 0; y < frame.Height(); y++ {
		for x := 0; x < frame.Width(); x++ {
			// 分辨率放大到原来的两倍，每个像素变成四个像素
			offset := original.YOffset(x, y)
			cOffset := original.COffset(x, y)
			dx, dy := x*EnhanceFrameScale, y*EnhanceFrameScale
			for i := 0; i < EnhanceFrameScale; i++ {
				for j := 0; j < EnhanceFrameScale; j++ {
					enhancedFrame.Y[enhancedFrame.YOffset(dx+i, dy+j)] = original.Y[offset]
					enhancedFrame.Cb[enhancedFrame.COffset(dx+i, dy+j)] = original.Cb[cOffset]
					enhancedFrame.Cr[enhancedFrame.COffset(dx+i, dy+j)] = original.Cr[cOffset]
				}
			}
		}
	}
	return emulator.MakeBaseFrame(enhancedFrame, frame.Width()*EnhanceFrameScale, frame.Height()*EnhanceFrameScale)
}

func (g *Instance) RenderCallback(frame emulator.IFrame, logger *log.Helper) {
	frame = g.frameEnhancer(frame)
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

func (g *Instance) setGraphicOptions(options *GraphicOptions) ConsumerResult {
	_ = g.e.Pause()
	defer g.e.Resume()
	if g.enhanceFrameOpen != options.HighResOpen {
		g.enhanceFrameOpen = options.HighResOpen
		g.videoEncoder.Close()
		enhanceRate := 1
		if options.HighResOpen {
			enhanceRate = 2
			g.frameEnhancer = g.enhanceFrame
		} else {
			g.frameEnhancer = func(frame emulator.IFrame) emulator.IFrame {
				return frame
			}
		}
		width, height := g.e.OutputResolution()
		enc, err := codec.NewVideoEncoder("vp8", width*enhanceRate, height*enhanceRate)
		if err != nil {
			return ConsumerResult{Error: err}
		}
		g.videoEncoder = enc
	}
	g.e.SetGraphicOptions(&emulator.GraphicOptions{
		Grayscale:    options.Grayscale,
		ReverseColor: options.ReverseColor,
	})
	options.HighResOpen = g.enhanceFrameOpen
	return ConsumerResult{Success: true, Error: nil}
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
