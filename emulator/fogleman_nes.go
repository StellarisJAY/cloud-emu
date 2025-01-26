package emulator

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"github.com/StellrisJAY/cloud-emu/emulator/fogleman_nes/nes"
)

type FoglemanNesAdapter struct {
	BaseEmulatorAdapter
	console     *nes.Console
	buttons     [2][8]bool
	game        string
	audioBuffer *bytes.Buffer
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
	adapter := &FoglemanNesAdapter{
		BaseEmulatorAdapter: newBaseEmulatorAdapter(256, 240, options),
		game:                options.Game(),
	}
	audioBuffer := bytes.NewBuffer([]byte{})
	adapter.stepFunc = adapter.step
	adapter.audioBuffer = audioBuffer
	console, err := nes.NewConsole(options.GameData(), audioBuffer)
	if err != nil {
		return nil, err
	}
	console.SetAudioSampleRate(float64(options.AudioSampleRate()))
	adapter.console = console
	return adapter, nil
}

func (f *FoglemanNesAdapter) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	f.cancel = cancel
	go f.emulatorLoop(ctx)
	return nil
}

func (f *FoglemanNesAdapter) step() {
	f.console.StepFrame()
	frame := f.console.Buffer()
	f.frame.FromRGBA(frame, f.scale)
	f.frameConsumer(f.frame)
	buffer := f.audioBuffer.Bytes()
	if len(buffer) == 0 {
		return
	}
	samples := make([]float32, 0, len(buffer)/4)
	reader := bytes.NewReader(buffer)
	for {
		var sample float32
		if err := binary.Read(reader, binary.LittleEndian, &sample); err != nil {
			break
		}
		samples = append(samples, sample)
	}
	f.audioBuffer.Reset()
	f.audioConsumer(samples)
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
	_ = f.Stop()
	f.audioBuffer = bytes.NewBuffer([]byte{})
	console, err := nes.NewConsole(options.GameData(), f.audioBuffer)
	if err != nil {
		return err
	}
	f.console = console
	f.frameConsumer = options.FrameConsumer()
	f.audioConsumer = options.AudioConsumer()
	return f.Start()
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
