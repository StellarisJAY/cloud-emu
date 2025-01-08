package emulator

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/emulator/dummy"
)

type DummyAdapter struct {
	e          *dummy.Emulator
	cancelFunc context.CancelFunc
}

func init() {
	supportedEmulators[TypeDummy] = Info{
		EmulatorType:           TypeDummy,
		Provider:               "https://github.com/StellrisJAY/cloud-emu/",
		Description:            "播放gif文件的dummy模拟器",
		Name:                   "Dummy",
		SupportSave:            false,
		SupportGraphicSettings: false,
	}
}

func MakeDummyAdapter(options IEmulatorOptions) (IEmulator, error) {
	e, err := dummy.MakeDummyEmulator(options.GameData(), func(frame *dummy.Frame) {
		options.FrameConsumer()(frame)
	})
	if err != nil {
		return nil, err
	}
	return &DummyAdapter{
		e: e,
	}, nil
}

func (d *DummyAdapter) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	d.cancelFunc = cancel
	go d.e.Start(ctx)
	return nil
}

func (d *DummyAdapter) Pause() error {
	d.e.Pause()
	return nil
}

func (d *DummyAdapter) Resume() error {
	d.e.Resume()
	return nil
}

func (d *DummyAdapter) Save() (IEmulatorSave, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DummyAdapter) LoadSave(save IEmulatorSave) error {
	//TODO implement me
	panic("implement me")
}

func (d *DummyAdapter) Restart(options IEmulatorOptions) error {
	//TODO implement me
	panic("implement me")
}

func (d *DummyAdapter) Stop() error {
	d.cancelFunc()
	return nil
}

func (d *DummyAdapter) SubmitInput(controllerId int, keyCode string, pressed bool) {
	//TODO implement me
	panic("implement me")
}

func (d *DummyAdapter) SetGraphicOptions(options *GraphicOptions) {
	//TODO implement me
	panic("implement me")
}

func (d *DummyAdapter) GetCPUBoostRate() float64 {
	//TODO implement me
	panic("implement me")
}

func (d *DummyAdapter) SetCPUBoostRate(f float64) float64 {
	//TODO implement me
	panic("implement me")
}

func (d *DummyAdapter) OutputResolution() (width, height int) {
	return d.e.Width, d.e.Height
}

func (d *DummyAdapter) MultiController() bool {
	return false
}

func (d *DummyAdapter) ControllerInfos() []ControllerInfo {
	return []ControllerInfo{
		{
			ControllerId: 1,
			Label:        "控制器",
		},
	}
}
