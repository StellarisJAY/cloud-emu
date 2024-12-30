package emulator

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/emulator/nes"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/bus"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/config"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/ppu"
	"image"
)

type NesEmulatorAdapter struct {
	e                *nes.Emulator
	options          IEmulatorOptions
	cancelFunc       context.CancelFunc
	reverseColorOpen bool
	grayscaleOpen    bool
}

type NesEmulatorOptions struct {
	NesGame               string
	NesGameData           []byte
	OutputAudioSampleRate int
	OutputAudioSampleChan chan float32
	NesRenderCallback     func(frame *ppu.Frame)
}

type NesFrameAdapter struct {
	frame *ppu.Frame
}

func MakeNesEmulatorOptions(game string, gameData []byte, audioSampleRate int, audioChan chan float32, renderCallback func(frame IFrame)) IEmulatorOptions {
	return &NesEmulatorOptions{
		NesGame:               game,
		NesGameData:           gameData,
		OutputAudioSampleRate: audioSampleRate,
		OutputAudioSampleChan: audioChan,
		NesRenderCallback: func(frame *ppu.Frame) {
			f := MakeNESFrameAdapter(frame)
			renderCallback(f)
		},
	}
}

func (n *NesEmulatorOptions) GameData() []byte {
	return n.NesGameData
}

func (n *NesEmulatorOptions) AudioSampleRate() int {
	return n.OutputAudioSampleRate
}

func (n *NesEmulatorOptions) AudioSampleChan() chan float32 {
	return n.OutputAudioSampleChan
}

func (n *NesEmulatorOptions) Game() string {
	return n.NesGame
}

func (n *NesEmulatorOptions) FrameConsumer() func(frame IFrame) {
	return nil
}

// Start 启动NES模拟器，创建单独的goroutine运行CPU循环，使用context打断
func (n *NesEmulatorAdapter) Start() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go n.e.LoadAndRun(ctx, false)
	n.cancelFunc = cancelFunc
	return nil
}

// Pause 暂停NES模拟器
func (n *NesEmulatorAdapter) Pause() error {
	n.e.Pause()
	return nil
}

// Resume 恢复NES模拟器
func (n *NesEmulatorAdapter) Resume() error {
	n.e.Resume()
	return nil
}

// Restart 重启NES模拟器，结束旧模拟器goroutine，创建并运行新模拟器
func (n *NesEmulatorAdapter) Restart(options IEmulatorOptions) error {
	n.cancelFunc()
	e, err := makeNESEmulator(options)
	if err != nil {
		return err
	}
	n.e = e
	ctx, cancelFunc := context.WithCancel(context.Background())
	go n.e.LoadAndRun(ctx, false)
	n.cancelFunc = cancelFunc
	n.options = options
	return nil
}

// Stop 关闭NES模拟器
func (n *NesEmulatorAdapter) Stop() error {
	n.cancelFunc()
	return nil
}

// Save 获取模拟器存档数据
func (n *NesEmulatorAdapter) Save() (IEmulatorSave, error) {
	n.e.Pause()
	defer n.e.Resume()
	s, err := n.e.GetSaveData()
	if err != nil {
		return nil, err
	}
	return &BaseEmulatorSave{
		Data: s,
	}, nil
}

func (n *NesEmulatorAdapter) LoadSave(save IEmulatorSave) error {
	n.e.Pause()
	defer n.e.Resume()
LOOP:
	for {
		select {
		case <-n.options.AudioSampleChan():
		default:
			break LOOP
		}
	}
	return n.e.Load(save.SaveData())
}

func (n *NesEmulatorAdapter) SubmitInput(controlId int, keyCode string, pressed bool) {
	switch keyCode {
	case "Up":
		n.e.SetJoyPadButtonPressed(controlId, bus.Up, pressed)
	case "Down":
		n.e.SetJoyPadButtonPressed(controlId, bus.Down, pressed)
	case "Left":
		n.e.SetJoyPadButtonPressed(controlId, bus.Left, pressed)
	case "Right":
		n.e.SetJoyPadButtonPressed(controlId, bus.Right, pressed)
	case "A":
		n.e.SetJoyPadButtonPressed(controlId, bus.ButtonA, pressed)
	case "B":
		n.e.SetJoyPadButtonPressed(controlId, bus.ButtonB, pressed)
	case "Select":
		n.e.SetJoyPadButtonPressed(controlId, bus.Select, pressed)
	case "Start":
		n.e.SetJoyPadButtonPressed(controlId, bus.Start, pressed)
	}
}

func (n *NesEmulatorAdapter) SetGraphicOptions(opts *GraphicOptions) {
	if opts.ReverseColor != n.reverseColorOpen {
		n.reverseColorOpen = opts.ReverseColor
		if n.reverseColorOpen {
			n.e.Frame().UseReverseColorPreprocessor()
		} else {
			n.e.Frame().RemoveReverseColorPreprocessor()
		}

	}

	if opts.Grayscale != n.grayscaleOpen {
		n.grayscaleOpen = opts.Grayscale
		if n.grayscaleOpen {
			n.e.Frame().UseGrayscalePreprocessor()
		} else {
			n.e.Frame().RemoveGrayscalePreprocessor()
		}
	}

}

func (n *NesEmulatorAdapter) GetCPUBoostRate() float64 {
	return n.e.CPUBoostRate()
}

func (n *NesEmulatorAdapter) SetCPUBoostRate(rate float64) float64 {
	return n.e.SetCPUBoostRate(rate)
}

func (n *NesEmulatorAdapter) OutputResolution() (width, height int) {
	return ppu.WIDTH, ppu.HEIGHT
}

func makeNESEmulator(options IEmulatorOptions) (*nes.Emulator, error) {
	configs := config.Config{
		Game:               options.Game(),
		EnableTrace:        false,
		Disassemble:        false,
		SnapshotSerializer: "json",
		MuteApu:            false,
		Debug:              false,
	}
	renderCallback := func(p *ppu.PPU) {
		options.(*NesEmulatorOptions).NesRenderCallback(p.Frame())
	}
	e, err := nes.NewEmulatorWithGameData(options.GameData(), configs, renderCallback, options.AudioSampleChan(), options.AudioSampleRate())
	if err != nil {
		return nil, err
	}
	return e, nil
}

func MakeNESFrameAdapter(frame *ppu.Frame) IFrame {
	return &NesFrameAdapter{
		frame: frame,
	}
}

func (f *NesFrameAdapter) Width() int {
	return ppu.WIDTH
}

func (f *NesFrameAdapter) Height() int {
	return ppu.HEIGHT
}

func (f *NesFrameAdapter) YCbCr() *image.YCbCr {
	return f.frame.YCbCr()
}

func (f *NesFrameAdapter) Read() (image.Image, func(), error) {
	return f.frame.Read()
}
