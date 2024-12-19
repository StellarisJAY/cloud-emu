package emulator

import (
	"context"
	"image"
	"image/color"
	"time"
)

type DummyEmulator struct {
	cancelFunc     context.CancelFunc
	renderCallback func(frame IFrame)
	pauseChan      chan struct{}
	resumeChan     chan struct{}
}

func MakeDummyEmulator(options IEmulatorOptions) IEmulator {
	return &DummyEmulator{
		renderCallback: options.FrameConsumer(),
		pauseChan:      make(chan struct{}),
		resumeChan:     make(chan struct{}),
	}
}

const dummyEmulatorFrameTick = 100 * time.Millisecond

var (
	g = image.NewYCbCr(image.Rect(0, 0, 256, 240), image.YCbCrSubsampleRatio420)
)

type dummyFrameAdapter struct{}

func (d *dummyFrameAdapter) Width() int {
	return 256
}

func (d *dummyFrameAdapter) Height() int {
	return 240
}

func (d *dummyFrameAdapter) YCbCr() *image.YCbCr {
	return g
}

func (d *dummyFrameAdapter) Read() (image.Image, func(), error) {
	return g, func() {}, nil
}

func init() {
	for y := 0; y < 256; y++ {
		for x := 0; x < 240; x++ {
			Y, cb, cr := color.RGBToYCbCr(0, 255, 0)
			yOff := g.YOffset(x, y)
			cOff := g.COffset(x, y)
			if yOff < len(g.Y) && cOff < len(g.Cb) && cOff < len(g.Cr) {
				g.Y[yOff] = Y
				g.Cb[cOff] = cb
				g.Cr[cOff] = cr
			}
		}
	}
}

func (d *DummyEmulator) Start() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	d.cancelFunc = cancelFunc
	frame := &dummyFrameAdapter{}
	go func() {
		ticker := time.NewTicker(dummyEmulatorFrameTick)
		for {
			select {
			case <-ctx.Done():
				return
			case <-d.pauseChan:
				ticker.Stop()
			case <-d.resumeChan:
				ticker.Reset(dummyEmulatorFrameTick)
			case <-ticker.C:
				d.renderCallback(frame)
			}
		}
	}()
	return nil
}

func (d *DummyEmulator) Pause() error {
	d.pauseChan <- struct{}{}
	return nil
}

func (d *DummyEmulator) Resume() error {
	d.resumeChan <- struct{}{}
	return nil
}

func (d *DummyEmulator) Save() (IEmulatorSave, error) {
	panic("dummy模拟器未实现该方法")
}

func (d *DummyEmulator) LoadSave(save IEmulatorSave, gameFileRepo IGameFileRepo) error {
	panic("dummy模拟器未实现该方法")
}

func (d *DummyEmulator) Restart(options IEmulatorOptions) error {
	panic("dummy模拟器未实现该方法")
}

func (d *DummyEmulator) Stop() error {
	d.cancelFunc()
	return nil
}

func (d *DummyEmulator) SubmitInput(controllerId int, keyCode string, pressed bool) {
	panic("dummy模拟器未实现该方法")
}

func (d *DummyEmulator) SetGraphicOptions(options *GraphicOptions) {
	panic("dummy模拟器未实现该方法")
}

func (d *DummyEmulator) GetCPUBoostRate() float64 {
	panic("dummy模拟器未实现该方法")
}

func (d *DummyEmulator) SetCPUBoostRate(f float64) float64 {
	panic("dummy模拟器未实现该方法")
}

func (d *DummyEmulator) OutputResolution() (width, height int) {
	return 256, 240
}

type DummyOptions struct {
	frameConsumer func(frame IFrame)
}

func MakeDummyOptions(frameConsumer func(frame IFrame)) IEmulatorOptions {
	return &DummyOptions{
		frameConsumer: frameConsumer,
	}
}

func (b *DummyOptions) Game() string {
	return ""
}

func (b *DummyOptions) GameData() []byte {
	return nil
}

func (b *DummyOptions) AudioSampleRate() int {
	return 0
}

func (b *DummyOptions) AudioSampleChan() chan float32 {
	return nil
}

func (b *DummyOptions) FrameConsumer() func(frame IFrame) {
	return b.frameConsumer
}
