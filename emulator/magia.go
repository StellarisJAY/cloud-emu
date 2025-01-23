package emulator

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/emulator/magia/gba"
	"image"
	"time"
)

type MagiaAdapter struct {
	e             *gba.GBA
	cancel        context.CancelFunc
	game          string
	frame         *BaseFrame
	frameConsumer func(frame IFrame)
	ticker        *time.Ticker
	buttons       [10]bool
	scale         int
	boost         float64
	pauseChan     chan struct{}
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
		e:             g,
		game:          options.Game(),
		frame:         MakeEmptyBaseFrame(image.Rect(0, 0, 240, 160)),
		frameConsumer: options.FrameConsumer(),
		scale:         1,
		boost:         1.0,
		pauseChan:     make(chan struct{}),
	}
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

func (m *MagiaAdapter) emulatorLoop(ctx context.Context) {
	m.ticker = time.NewTicker(getFrameInterval(m.boost))
	defer func() {
		if r := recover(); r != nil {
		}
		m.ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.pauseChan:
			m.ticker.Stop()
		case <-m.ticker.C:
			start := time.Now()
			m.e.Update()
			m.frame.FromRGBARaw(m.e.Draw(), m.scale)
			m.frameConsumer(m.frame)
			processTime := time.Since(start)
			interval := max(getFrameInterval(m.boost)-processTime, time.Millisecond*5)
			m.ticker.Reset(interval)
		}
	}
}

func (m *MagiaAdapter) Pause() error {
	m.pauseChan <- struct{}{}
	return nil
}

func (m *MagiaAdapter) Resume() error {
	m.ticker.Reset(getFrameInterval(m.boost))
	return nil
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
	g := gba.New(options.GameData(), false)
	m.cancel()
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

func (m *MagiaAdapter) Stop() error {
	m.cancel()
	close(m.pauseChan)
	return nil
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

func (m *MagiaAdapter) SetGraphicOptions(options *GraphicOptions) {
	_ = m.Pause()
	if options.HighResolution {
		m.scale = 2
	} else {
		m.scale = 1
	}
	m.frame = MakeEmptyBaseFrame(image.Rect(0, 0, 240*m.scale, 160*m.scale))
	_ = m.Resume()
}

func (m *MagiaAdapter) GetGraphicOptions() *GraphicOptions {
	return &GraphicOptions{
		HighResolution: m.scale > 1,
	}
}

func (m *MagiaAdapter) GetCPUBoostRate() float64 {
	return m.boost
}

func (m *MagiaAdapter) SetCPUBoostRate(f float64) float64 {
	m.boost = max(0.5, min(f, 2.0))
	return m.boost
}

func (m *MagiaAdapter) OutputResolution() (width, height int) {
	return 240 * m.scale, 160 * m.scale
}

func (m *MagiaAdapter) MultiController() bool {
	return true
}

func (m *MagiaAdapter) ControllerInfos() []ControllerInfo {
	return []ControllerInfo{
		{
			ControllerId: 1,
			Label:        "玩家1",
		},
	}
}
