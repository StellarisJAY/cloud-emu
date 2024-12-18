package emulator

import (
	"context"
	"time"
)

type DummyEmulator struct {
	cancelFunc     context.CancelFunc
	renderCallback func(frame IFrame)
}

func MakeDummyEmulator(options IEmulatorOptions) IEmulator {
	return &DummyEmulator{
		renderCallback: options.FrameConsumer(),
	}
}

func (d *DummyEmulator) Start() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	d.cancelFunc = cancelFunc
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				d.renderCallback(nil)
			}
		}
	}()
	return nil
}

func (d *DummyEmulator) Pause() error {
	panic("dummy模拟器未实现该方法")
}

func (d *DummyEmulator) Resume() error {
	panic("dummy模拟器未实现该方法")
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
	//TODO implement me
	panic("implement me")
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
