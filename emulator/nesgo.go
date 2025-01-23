package emulator

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/emulator/nes"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/bus"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/config"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/ppu"
	"time"
)

type NesEmulatorAdapter struct {
	e             *nes.Emulator
	cancelFunc    context.CancelFunc
	ticker        *time.Ticker
	frameConsumer func(IFrame)
	scale         int
	boost         float64
	pauseChan     chan struct{}
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
	return &NesEmulatorAdapter{
		e:             e,
		frameConsumer: options.FrameConsumer(),
		scale:         1,
		boost:         1,
		pauseChan:     make(chan struct{}),
	}, nil
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
	n.cancelFunc = cancelFunc
	n.e.Reset()
	go n.emulatorLoop(ctx)
	return nil
}

func (n *NesEmulatorAdapter) emulatorLoop(ctx context.Context) {
	n.ticker = time.NewTicker(getFrameInterval(n.boost))
	defer func() {
		if r := recover(); r != nil {
		}
		n.ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-n.pauseChan:
			n.ticker.Stop()
		case <-n.ticker.C:
			start := time.Now()
			n.e.StepOneFrame()
			n.frameConsumer(n.e.Frame())
			interval := max(getFrameInterval(n.boost)-time.Since(start), time.Millisecond*5)
			n.ticker.Reset(interval)
		}
	}
}

// Pause 暂停NES模拟器
func (n *NesEmulatorAdapter) Pause() error {
	n.pauseChan <- struct{}{}
	return nil
}

// Resume 恢复NES模拟器
func (n *NesEmulatorAdapter) Resume() error {
	if n.ticker != nil {
		n.ticker.Reset(getFrameInterval(n.boost))
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
	n.frameConsumer = options.FrameConsumer()
	ctx, cancelFunc := context.WithCancel(context.Background())
	go n.emulatorLoop(ctx)
	n.cancelFunc = cancelFunc
	return nil
}

// Stop 关闭NES模拟器
func (n *NesEmulatorAdapter) Stop() error {
	n.cancelFunc()
	close(n.pauseChan)
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
	_ = n.Pause()
	if opts.HighResolution {
		n.scale = 2
	} else {
		n.scale = 1
	}
	n.e.Frame().SetScale(n.scale)
	_ = n.Resume()
}

func (n *NesEmulatorAdapter) GetGraphicOptions() *GraphicOptions {
	return &GraphicOptions{
		HighResolution: n.scale > 1,
	}
}

func (n *NesEmulatorAdapter) GetCPUBoostRate() float64 {
	return n.boost
}

func (n *NesEmulatorAdapter) SetCPUBoostRate(rate float64) float64 {
	n.boost = max(0.5, min(rate, 2.0))
	return n.boost
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
