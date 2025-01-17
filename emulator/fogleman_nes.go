package emulator

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/StellrisJAY/cloud-emu/emulator/fogleman_nes/nes"
	"image"
	"time"
)

type FoglemanNesAdapter struct {
	console       *nes.Console
	cancel        context.CancelFunc
	ticker        *time.Ticker
	frame         *BaseFrame
	frameConsumer func(frame IFrame)
	buttons       [2][8]bool
	game          string
	scale         int
}

func init() {
	supportedEmulators[CodeFoglemanNES] = Info{
		EmulatorType:           "NES",
		Provider:               "https://github.com/fogleman/nes",
		Description:            "Go语言实现的功能完善的NES模拟器",
		Name:                   "fogleman/nes",
		EmulatorCode:           CodeFoglemanNES,
		SupportSave:            true,
		SupportGraphicSettings: true,
	}
}

func newFoglemanNesAdapter(options IEmulatorOptions) (*FoglemanNesAdapter, error) {
	adapter := &FoglemanNesAdapter{}
	adapter.game = options.Game()
	adapter.scale = 1
	adapter.frame = MakeEmptyBaseFrame(image.Rect(0, 0, 256, 240))
	console, err := nes.NewConsole(options.GameData())
	if err != nil {
		return nil, err
	}
	adapter.console = console
	adapter.frameConsumer = options.FrameConsumer()
	return adapter, nil
}

func (f *FoglemanNesAdapter) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	f.cancel = cancel
	go f.emulatorLoop(ctx)
	return nil
}

func (f *FoglemanNesAdapter) emulatorLoop(ctx context.Context) {
	f.ticker = time.NewTicker(FrameInterval)
	defer func() {
		if r := recover(); r != nil {
		}
		f.ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-f.ticker.C:
			start := time.Now()
			f.console.StepFrame()
			frame := f.console.Buffer()
			f.frame.FromRGBA(frame, f.scale)
			f.frameConsumer(f.frame)
			interval := max(FrameInterval-time.Since(start), time.Millisecond*5)
			f.ticker.Reset(interval)
		}
	}
}

func (f *FoglemanNesAdapter) Pause() error {
	f.ticker.Stop()
	return nil
}

func (f *FoglemanNesAdapter) Resume() error {
	f.ticker.Reset(FrameInterval)
	return nil
}

func (f *FoglemanNesAdapter) Save() (IEmulatorSave, error) {
	buffer := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buffer)
	if err := f.console.Save(encoder); err != nil {
		return nil, err
	}
	return &BaseEmulatorSave{
		Game: f.game,
		Data: buffer.Bytes(),
	}, nil
}

func (f *FoglemanNesAdapter) LoadSave(save IEmulatorSave) error {
	buffer := bytes.NewBuffer(save.SaveData())
	decoder := gob.NewDecoder(buffer)
	return f.console.Load(decoder)
}

func (f *FoglemanNesAdapter) Restart(options IEmulatorOptions) error {
	f.cancel()
	console, err := nes.NewConsole(options.GameData())
	if err != nil {
		return err
	}
	f.console = console
	f.frameConsumer = options.FrameConsumer()
	return f.Start()
}

func (f *FoglemanNesAdapter) Stop() error {
	f.cancel()
	return nil
}

func (f *FoglemanNesAdapter) SubmitInput(controllerId int, keyCode string, pressed bool) {
	switch keyCode {
	case "A":
		f.buttons[controllerId-1][nes.ButtonA] = pressed
	case "B":
		f.buttons[controllerId-1][nes.ButtonB] = pressed
	case "Left":
		f.buttons[controllerId-1][nes.ButtonLeft] = pressed
	case "Right":
		f.buttons[controllerId-1][nes.ButtonRight] = pressed
	case "Up":
		f.buttons[controllerId-1][nes.ButtonUp] = pressed
	case "Down":
		f.buttons[controllerId-1][nes.ButtonDown] = pressed
	case "Select":
		f.buttons[controllerId-1][nes.ButtonSelect] = pressed
	case "Start":
		f.buttons[controllerId-1][nes.ButtonStart] = pressed
	}
	f.console.SetButtons1(f.buttons[0])
	f.console.SetButtons2(f.buttons[1])
}

func (f *FoglemanNesAdapter) SetGraphicOptions(options *GraphicOptions) {
	_ = f.Pause()
	if options.HighResolution {
		f.scale = 2
	} else {
		f.scale = 1
	}
	f.frame = MakeEmptyBaseFrame(image.Rect(0, 0, 256*f.scale, 240*f.scale))
	_ = f.Resume()
}

func (f *FoglemanNesAdapter) GetGraphicOptions() *GraphicOptions {
	return &GraphicOptions{
		HighResolution: f.scale > 1,
	}
}

func (f *FoglemanNesAdapter) GetCPUBoostRate() float64 {
	//TODO implement me
	panic("implement me")
}

func (f *FoglemanNesAdapter) SetCPUBoostRate(f2 float64) float64 {
	//TODO implement me
	panic("implement me")
}

func (f *FoglemanNesAdapter) OutputResolution() (width, height int) {
	return 256, 240
}

func (f *FoglemanNesAdapter) MultiController() bool {
	return true
}

func (f *FoglemanNesAdapter) ControllerInfos() []ControllerInfo {
	return []ControllerInfo{
		{
			ControllerId: 1,
			Label:        "玩家1",
		},
		{
			ControllerId: 2,
			Label:        "玩家2",
		},
	}
}
