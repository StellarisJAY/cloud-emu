package emulator

type chip8EmulatorAdapter struct {
}

func (c *chip8EmulatorAdapter) Start() error {
	panic("not implemented")
}

func (c *chip8EmulatorAdapter) Pause() error {
	panic("not implemented")
}

func (c *chip8EmulatorAdapter) Resume() error {
	panic("not implemented")
}

func (c *chip8EmulatorAdapter) Save() (IEmulatorSave, error) {
	panic("not implemented")
}

func (c *chip8EmulatorAdapter) LoadSave(save IEmulatorSave, gameFileRepo IGameFileRepo) error {
	panic("not implemented")
}

func (c *chip8EmulatorAdapter) Restart(options IEmulatorOptions) error {
	panic("not implemented")
}

func (c *chip8EmulatorAdapter) Stop() error {
	panic("not implemented")
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
	panic("not implemented")
}
