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
	supportedEmulators[CodeDummy] = Info{
		EmulatorType:           "DUMMY",
		Provider:               "https://github.com/StellrisJAY/cloud-emu/",
		Description:            "播放gif文件的dummy模拟器",
		Name:                   "Dummy",
		EmulatorCode:           CodeDummy,
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
	d.e.Stop()
	e, err := dummy.MakeDummyEmulator(options.GameData(), func(frame *dummy.Frame) {
		options.FrameConsumer()(frame)
	})
	if err != nil {
		return err
	}
	d.e = e
	return d.Start()
}

func (d *DummyAdapter) Stop() error {
	d.cancelFunc()
	return nil
}

func (d *DummyAdapter) SubmitInput(controllerId int, keyCode string, pressed bool) {
}

func (d *DummyAdapter) SetGraphicOptions(_ *GraphicOptions) {

}

func (d *DummyAdapter) GetGraphicOptions() *GraphicOptions {
	return &GraphicOptions{
		HighResolution: false,
	}
}

func (d *DummyAdapter) GetCPUBoostRate() float64 {
	return 1.0
}

func (d *DummyAdapter) SetCPUBoostRate(f float64) float64 {
	return 1.0
}

func (d *DummyAdapter) OutputResolution() (width, height int) {
	return d.e.Width, d.e.Height
}

func (d *DummyAdapter) ControllerInfos() []ControllerInfo {
	return []ControllerInfo{
		{
			ControllerId: 1,
			Label:        "控制器",
		},
	}
}
