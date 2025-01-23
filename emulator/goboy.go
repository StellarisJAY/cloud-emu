package emulator

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/emulator/goboy/gb"
)

type GoboyAdapter struct {
	BaseEmulatorAdapter
	gb   *gb.Gameboy // origin:  https://github.com/Humpheh/goboy
	game string
	buf  [][][3]byte
}

func init() {
	supportedEmulators[CodeGoboy] = Info{
		EmulatorType:           "GB",
		EmulatorCode:           CodeGoboy,
		Provider:               "https://github.com/Humpheh/goboy",
		Description:            "Go语言实现的GameBoy模拟器",
		Name:                   "GoBoy",
		SupportSave:            false,
		SupportGraphicSettings: true,
	}
}

func newGoboyAdapter(options IEmulatorOptions) (IEmulator, error) {
	e, err := gb.NewGameboy(options.Game(), options.GameData())
	if err != nil {
		return nil, err
	}
	g := &GoboyAdapter{
		gb:                  e,
		game:                options.Game(),
		BaseEmulatorAdapter: newBaseEmulatorAdapter(160, 144, options),
		buf:                 make([][][3]byte, 160),
	}
	g.stepFunc = g.step
	return g, nil
}

func (g *GoboyAdapter) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	g.cancel = cancel
	go g.emulatorLoop(ctx)
	return nil
}

func (g *GoboyAdapter) step() {
	g.gb.Update()
	for i := range g.buf {
		g.buf[i] = g.gb.PreparedData[i][:]
	}
	g.frame.FromRGBRaw(g.buf, g.scale)
	g.frameConsumer(g.frame)
}

func (g *GoboyAdapter) Save() (IEmulatorSave, error) {
	save := g.gb.Memory.Cart.BankingController.GetSaveData()
	if len(save) == 0 {
		return nil, nil
	}
	return &BaseEmulatorSave{
		Game: g.game,
		Data: save,
	}, nil
}

func (g *GoboyAdapter) LoadSave(save IEmulatorSave) error {
	g.gb.Memory.Cart.BankingController.LoadSaveData(save.SaveData())
	return nil
}

func (g *GoboyAdapter) Restart(options IEmulatorOptions) error {
	_ = g.Stop()
	e, err := gb.NewGameboy(options.Game(), options.GameData())
	if err != nil {
		return err
	}
	g.gb = e
	g.frameConsumer = options.FrameConsumer()
	ctx, cancel := context.WithCancel(context.Background())
	g.cancel = cancel
	go g.emulatorLoop(ctx)
	return nil
}

func (g *GoboyAdapter) SubmitInput(_ int, keyCode string, pressed bool) {
	input := gb.ButtonInput{}
	key := byte(255)
	switch keyCode {
	case "Right":
		key = gb.ButtonRight
	case "Left":
		key = gb.ButtonLeft
	case "Down":
		key = gb.ButtonDown
	case "Up":
		key = gb.ButtonUp
	case "A":
		key = byte(gb.ButtonA)
	case "B":
		key = gb.ButtonB
	case "Select":
		key = gb.ButtonSelect
	case "Start":
		key = gb.ButtonStart
	}
	if key != 255 {
		if pressed {
			input.Pressed = []gb.Button{gb.Button(key)}
		} else {
			input.Released = []gb.Button{gb.Button(key)}
		}
	}
	g.gb.ProcessInput(input)
}

func (g *GoboyAdapter) ControllerInfos() []ControllerInfo {
	return []ControllerInfo{
		{
			ControllerId: 1,
			Label:        "玩家1",
		},
	}
}
