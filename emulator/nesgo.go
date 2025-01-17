package emulator

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/emulator/nes"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/bus"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/config"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/ppu"
	"image"
	"time"
)

type NesEmulatorAdapter struct {
	e                *nes.Emulator
	options          IEmulatorOptions
	cancelFunc       context.CancelFunc
	reverseColorOpen bool
	grayscaleOpen    bool
	ticker           *time.Ticker
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

func init() {
	supportedEmulators[CodeNesgo] = Info{
		EmulatorCode:           CodeNesgo,
		EmulatorType:           "NES",
		Provider:               "https://github.com/StellrisJAY/cloud-emu",
		Description:            "Go语言实现的NES模拟器，部分游戏运行存在Bug，但该模拟器是CloudEmu的起源，具有纪念意义所以没有被删除和替换。",
		Name:                   "NESGO",
		SupportSave:            true,
		SupportGraphicSettings: false,
	}
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
	n.e.Reset()
	go n.emulatorLoop(ctx)
	n.cancelFunc = cancelFunc
	return nil
}

func (n *NesEmulatorAdapter) emulatorLoop(ctx context.Context) {
	n.ticker = time.NewTicker(FrameInterval)
	defer func() {
		if r := recover(); r != nil {
		}
		n.ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-n.ticker.C:
			start := time.Now()
			n.e.StepOneFrame()
			interval := max(FrameInterval-time.Since(start), time.Millisecond*5)
			n.ticker.Reset(interval)
		}
	}
}

// Pause 暂停NES模拟器
func (n *NesEmulatorAdapter) Pause() error {
	if n.ticker != nil {
		n.ticker.Stop()
	}
	return nil
}

// Resume 恢复NES模拟器
func (n *NesEmulatorAdapter) Resume() error {
	if n.ticker != nil {
		n.ticker.Reset(FrameInterval)
	}
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
	n.e.Reset()
	ctx, cancelFunc := context.WithCancel(context.Background())
	go n.emulatorLoop(ctx)
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
	_ = n.Pause()
	defer n.Resume()
	s, err := n.e.GetSaveData()
	if err != nil {
		return nil, err
	}
	return &BaseEmulatorSave{
		Data: s,
	}, nil
}

func (n *NesEmulatorAdapter) LoadSave(save IEmulatorSave) error {
	_ = n.Pause()
	defer n.Resume()
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
	panic("not implemented")
}

func (n *NesEmulatorAdapter) GetGraphicOptions() *GraphicOptions {
	panic("not implemented")
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

func (n *NesEmulatorAdapter) MultiController() bool {
	return true
}

func (n *NesEmulatorAdapter) ControllerInfos() []ControllerInfo {
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
