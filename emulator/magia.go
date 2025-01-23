package emulator

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/emulator/magia/gba"
)

type MagiaAdapter struct {
	BaseEmulatorAdapter
	e       *gba.GBA
	game    string
	buttons [10]bool
}

func init() {
	supportedEmulators[CodeMagia] = Info{
		EmulatorCode:           CodeMagia,
		EmulatorType:           "GBA",
		Provider:               "https://github.com/akatsuki105/magia",
		Description:            "Go语言实现的GBA模拟器",
		Name:                   "Magia",
		SupportSave:            false,
		SupportGraphicSettings: true,
	}
}

func newMagiaAdapter(options IEmulatorOptions) (*MagiaAdapter, error) {
	g := gba.New(options.GameData(), false)
	adapter := &MagiaAdapter{
		e:                   g,
		game:                options.Game(),
		BaseEmulatorAdapter: newBaseEmulatorAdapter(240, 160, options),
	}
	adapter.stepFunc = adapter.step
	handlers := [10]func() bool{}
	for i := 0; i < 10; i++ {
		handlers[i] = func() bool {
			return adapter.buttons[i]
		}
	}
	g.SetJoypadHandler(handlers)
	return adapter, nil
}

func (m *MagiaAdapter) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	m.e.SoftReset()
	go m.emulatorLoop(ctx)
	return nil
}

func (m *MagiaAdapter) step() {
	m.e.Update()
	m.frame.FromRGBARaw(m.e.Draw(), m.scale)
	m.frameConsumer(m.frame)
}

func (m *MagiaAdapter) Save() (IEmulatorSave, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MagiaAdapter) LoadSave(save IEmulatorSave) error {
	//TODO implement me
	panic("implement me")
}

func (m *MagiaAdapter) Restart(options IEmulatorOptions) error {
	_ = m.Stop()
	g := gba.New(options.GameData(), false)
	m.e = g
	handlers := [10]func() bool{}
	for i := 0; i < 10; i++ {
		handlers[i] = func() bool {
			return m.buttons[i]
		}
	}
	g.SetJoypadHandler(handlers)
	m.frameConsumer = options.FrameConsumer()
	return m.Start()
}

func (m *MagiaAdapter) SubmitInput(_ int, keyCode string, pressed bool) {
	switch keyCode {
	case "A":
		m.buttons[gba.A] = pressed
	case "B":
		m.buttons[gba.B] = pressed
	case "Left":
		m.buttons[gba.Left] = pressed
	case "Right":
		m.buttons[gba.Right] = pressed
	case "Up":
		m.buttons[gba.Up] = pressed
	case "Down":
		m.buttons[gba.Down] = pressed
	case "Start":
		m.buttons[gba.Start] = pressed
	case "Select":
		m.buttons[gba.Select] = pressed
	case "L":
		m.buttons[gba.L] = pressed
	case "R":
		m.buttons[gba.R] = pressed
	}
}

func (m *MagiaAdapter) ControllerInfos() []ControllerInfo {
	return []ControllerInfo{
		{
			ControllerId: 1,
			Label:        "玩家1",
		},
	}
}
