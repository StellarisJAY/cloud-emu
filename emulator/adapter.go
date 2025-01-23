package emulator

import (
	"context"
	"image"
	"sync"
	"time"
)

// BaseEmulatorAdapter 基础模拟器适配器，实现了IEmulator的部分常用的方法
type BaseEmulatorAdapter struct {
	frame         *BaseFrame         // 最新的一帧画面
	frameConsumer func(IFrame)       // 画面输出到视频函数
	ticker        *time.Ticker       // 模拟器固定帧率运行ticker
	cancel        context.CancelFunc // 停止模拟器的context
	scale         int                // 输出分辨率缩放
	boost         float64            // 模拟加速倍率
	pauseChan     chan struct{}      // 暂停channel
	resumeChan    chan struct{}      // 恢复channel
	wg            *sync.WaitGroup    // 等待emulatorLoop结束

	width  int // 模拟器原始分辨率W
	height int // 模拟器原始分辨率H

	stepFunc func() // 推进一帧画面的函数
}

func newBaseEmulatorAdapter(width, height int, options IEmulatorOptions) BaseEmulatorAdapter {
	return BaseEmulatorAdapter{
		frame:         MakeEmptyBaseFrame(image.Rect(0, 0, width, height)),
		frameConsumer: options.FrameConsumer(),
		scale:         1,
		boost:         1.0,
		pauseChan:     make(chan struct{}),
		resumeChan:    make(chan struct{}),
		wg:            &sync.WaitGroup{},
		width:         width,
		height:        height,
		stepFunc:      func() {},
	}
}

func (b *BaseEmulatorAdapter) Pause() error {
	// 直接stop可能会被emulatorLoop中每帧后的reset覆盖，这里使用pauseChan将stop操作交给emulatorLoop
	b.pauseChan <- struct{}{}
	return nil
}

func (b *BaseEmulatorAdapter) Resume() error {
	b.resumeChan <- struct{}{}
	return nil
}

func (b *BaseEmulatorAdapter) Stop() error {
	b.cancel()
	// 必须等待之前的emulatorLoop goroutine结束
	b.wg.Wait()
	return nil
}

func (b *BaseEmulatorAdapter) SetGraphicOptions(options *GraphicOptions) {
	if options.HighResolution {
		b.scale = 2
	} else {
		b.scale = 1
	}
	b.frame = MakeEmptyBaseFrame(image.Rect(0, 0, b.width*b.scale, b.height*b.scale))
}

func (b *BaseEmulatorAdapter) GetGraphicOptions() *GraphicOptions {
	return &GraphicOptions{
		HighResolution: b.scale > 1,
	}
}

func (b *BaseEmulatorAdapter) GetCPUBoostRate() float64 {
	return b.boost
}

func (b *BaseEmulatorAdapter) SetCPUBoostRate(f float64) float64 {
	b.boost = max(0.5, min(f, 2.0))
	return b.boost
}

func (b *BaseEmulatorAdapter) OutputResolution() (width, height int) {
	return b.width * b.scale, b.height * b.scale
}

func (b *BaseEmulatorAdapter) emulatorLoop(ctx context.Context) {
	b.ticker = time.NewTicker(getFrameInterval(b.boost))
	b.wg.Add(1)
	defer func() {
		if r := recover(); r != nil {
			// TODO logging
		}
		b.ticker.Stop()
		b.wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-b.pauseChan:
			b.ticker.Stop()
		case <-b.resumeChan:
			b.ticker.Reset(getFrameInterval(b.boost))
		case <-b.ticker.C:
			start := time.Now()
			// 调用模拟器的step()方法推进一帧画面
			b.stepFunc()
			// 计算一帧画面的耗时，根据耗时调整下一次的ticker时间
			interval := max(getFrameInterval(b.boost)-time.Since(start), time.Millisecond*5)
			b.ticker.Reset(interval)
		}
	}
}
