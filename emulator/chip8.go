package emulator

import (
	"context"
	"errors"
	"github.com/StellrisJAY/cloud-emu/emulator/chip8"
)

type chip8EmulatorAdapter struct {
	e      *chip8.Emulator
	cancel context.CancelFunc
}

func init() {
	//supportedEmulators[TypeChip8] = Info{
	//	EmulatorType:           TypeChip8,
	//	Provider:               "https://github.com/StellrisJAY/cloud-emu",
	//	Description:            "CHIP-8 模拟器",
	//	Name:                   "CHIP-8",
	//	SupportSave:            false,
	//	SupportGraphicSettings: false,
	//}
}

func makeChip8EmulatorAdapter(options IEmulatorOptions) IEmulator {
	e := chip8.NewEmulator(options.GameData(), func(frame *chip8.Frame) {
		options.FrameConsumer()(frame)
	})
	return &chip8EmulatorAdapter{
		e: e,
	}
}

func (c *chip8EmulatorAdapter) Start() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	c.cancel = cancelFunc
	go c.e.Run(ctx)
	return nil
}

func (c *chip8EmulatorAdapter) Pause() error {
	c.e.Pause()
	return nil
}

func (c *chip8EmulatorAdapter) Resume() error {
	c.e.Resume()
	return nil
}

func (c *chip8EmulatorAdapter) Save() (IEmulatorSave, error) {
	return nil, errors.New("不支持存档")
}

func (c *chip8EmulatorAdapter) LoadSave(save IEmulatorSave) error {
	return errors.New("不支持存档")
}

func (c *chip8EmulatorAdapter) Restart(options IEmulatorOptions) error {
	panic("not implemented")
}

func (c *chip8EmulatorAdapter) Stop() error {
	c.cancel()
	return nil
}

func (c *chip8EmulatorAdapter) SubmitInput(controllerId int, keyCode string, pressed bool) {
	panic("not implemented")
}

func (c *chip8EmulatorAdapter) SetGraphicOptions(options *GraphicOptions) {
	panic("not implemented")
}

func (c *chip8EmulatorAdapter) GetCPUBoostRate() float64 {
	panic("not implemented")
}

func (c *chip8EmulatorAdapter) SetCPUBoostRate(f float64) float64 {
	panic("not implemented")
}

func (c *chip8EmulatorAdapter) OutputResolution() (width, height int) {
	return 64, 32
}

func (c *chip8EmulatorAdapter) MultiController() bool {
	return false
}

func (c *chip8EmulatorAdapter) ControllerInfos() []ControllerInfo {
	return []ControllerInfo{
		{
			ControllerId: 1,
			Label:        "控制器",
		},
	}
}

func (c *chip8EmulatorAdapter) GetGraphicOptions() *GraphicOptions {
	panic("not implemented")
}
