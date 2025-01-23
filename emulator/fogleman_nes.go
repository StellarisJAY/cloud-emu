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
	boost         float64
	pauseChan     chan struct{}
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
	adapter.boost = 1.0
	adapter.pauseChan = make(chan struct{})
	adapter.frame = MakeEmptyBaseFrame(image.Rect(0, 0, 256, 240))
	console, err := nes.NewConsole(options.GameData())
	if err != nil {
		return nil, err
	}
	console.SetAudioSampleRate(float64(options.AudioSampleRate()))
	console.SetAudioChannel(options.AudioSampleChan())
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
	f.ticker = time.NewTicker(getFrameInterval(f.boost))
	defer func() {
		if r := recover(); r != nil {
		}
		f.ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-f.pauseChan:
			f.ticker.Stop()
		case <-f.ticker.C:
			start := time.Now()
			f.console.StepFrame()
			frame := f.console.Buffer()
			f.frame.FromRGBA(frame, f.scale)
			f.frameConsumer(f.frame)
			interval := max(getFrameInterval(f.boost)-time.Since(start), time.Millisecond*5)
			f.ticker.Reset(interval)
		}
	}
}

func (f *FoglemanNesAdapter) Pause() error {
	// 直接stop可能会被emulatorLoop中每帧后的reset覆盖，这里使用pauseChan将stop操作交给emulatorLoop
	// f.ticker.Stop()
	f.pauseChan <- struct{}{}
	return nil
}

func (f *FoglemanNesAdapter) Resume() error {
	f.ticker.Reset(getFrameInterval(f.boost))
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
	close(f.pauseChan)
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
	return f.boost
}

func (f *FoglemanNesAdapter) SetCPUBoostRate(r float64) float64 {
	f.boost = max(0.5, min(r, 2.0))
	return f.boost
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
