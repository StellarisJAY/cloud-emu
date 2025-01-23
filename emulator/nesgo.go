package emulator

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/emulator/nes"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/bus"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/config"
)

type NesEmulatorAdapter struct {
	BaseEmulatorAdapter
	e *nes.Emulator
}

func init() {
	supportedEmulators[CodeNesgo] = Info{
		EmulatorCode:           CodeNesgo,
		EmulatorType:           "NES",
		Provider:               "https://github.com/StellrisJAY/cloud-emu",
		Description:            "Go语言实现的NES模拟器，部分游戏运行存在Bug，但该模拟器是CloudEmu的起源，具有纪念意义所以没有被删除和替换。",
		Name:                   "NESGO",
		SupportSave:            true,
		SupportGraphicSettings: true,
	}
}

func makeNESEmulatorAdapter(options IEmulatorOptions) (IEmulator, error) {
	e, err := makeNESEmulator(options)
	if err != nil {
		return nil, err
	}
	n := &NesEmulatorAdapter{
		e:                   e,
		BaseEmulatorAdapter: newBaseEmulatorAdapter(256, 240, options),
	}
	n.stepFunc = n.step
	return n, nil
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
	e, err := nes.NewEmulatorWithGameData(options.GameData(), configs, options.AudioSampleChan(), options.AudioSampleRate())
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Start 启动NES模拟器，创建单独的goroutine运行CPU循环，使用context打断
func (n *NesEmulatorAdapter) Start() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	n.cancel = cancelFunc
	n.e.Reset()
	go n.emulatorLoop(ctx)
	return nil
}

func (n *NesEmulatorAdapter) step() {
	n.e.StepOneFrame()
	n.frameConsumer(n.e.Frame())
}

// Restart 重启NES模拟器，结束旧模拟器goroutine，创建并运行新模拟器
func (n *NesEmulatorAdapter) Restart(options IEmulatorOptions) error {
	_ = n.Stop()
	e, err := makeNESEmulator(options)
	if err != nil {
		return err
	}
	n.e = e
	n.e.Reset()
	n.frameConsumer = options.FrameConsumer()
	ctx, cancelFunc := context.WithCancel(context.Background())
	go n.emulatorLoop(ctx)
	n.cancel = cancelFunc
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
	if opts.HighResolution {
		n.scale = 2
	} else {
		n.scale = 1
	}
	n.e.Frame().SetScale(n.scale)
}

func (n *NesEmulatorAdapter) GetGraphicOptions() *GraphicOptions {
	return &GraphicOptions{
		HighResolution: n.scale > 1,
	}
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
