package emulator

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/emulator/goboy/gb"
	"image"
	"image/color"
	"time"
)

type GoboyAdapter struct {
	gb            *gb.Gameboy // origin:  https://github.com/Humpheh/goboy
	cancel        context.CancelFunc
	ticker        *time.Ticker
	frame         *GoboyFrame
	frameConsumer func(frame IFrame)
	game          string
	boost         float64
	pauseChan     chan struct{}
	resumeChan    chan struct{}
}

func init() {
	supportedEmulators[CodeGoboy] = Info{
		EmulatorType:           "GB",
		EmulatorCode:           CodeGoboy,
		Provider:               "https://github.com/Humpheh/goboy",
		Description:            "Go语言实现的GameBoy模拟器",
		Name:                   "GoBoy",
		SupportSave:            false,
		SupportGraphicSettings: false,
	}
}

func newGoboyAdapter(options IEmulatorOptions) (IEmulator, error) {
	e, err := gb.NewGameboy(options.Game(), options.GameData())
	if err != nil {
		return nil, err
	}
	g := &GoboyAdapter{
		gb:            e,
		frameConsumer: options.FrameConsumer(),
		frame:         newGoboyFrame(),
		game:          options.Game(),
		boost:         1.0,
		pauseChan:     make(chan struct{}),
		resumeChan:    make(chan struct{}),
	}
	return g, nil
}

func (g *GoboyAdapter) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	g.cancel = cancel
	go g.emulatorLoop(ctx)
	return nil
}

func (g *GoboyAdapter) emulatorLoop(ctx context.Context) {
	g.ticker = time.NewTicker(getFrameInterval(g.boost))
	defer func() {
		if r := recover(); r != nil {
		}
		g.ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			close(g.resumeChan)
			close(g.pauseChan)
			return
		case <-g.pauseChan:
			g.ticker.Stop()
			g.gb.ExecutionPaused = true
		case <-g.resumeChan:
			g.ticker.Reset(getFrameInterval(g.boost))
			g.gb.ExecutionPaused = false
		case <-g.ticker.C:
			start := time.Now()
			g.gb.Update()
			g.frame.reset(g.gb.PreparedData)
			g.frameConsumer(g.frame)
			interval := max(getFrameInterval(g.boost)-time.Since(start), time.Millisecond*5)
			g.ticker.Reset(interval)
		}
	}
}

func (g *GoboyAdapter) Pause() error {
	g.pauseChan <- struct{}{}
	return nil
}

func (g *GoboyAdapter) Resume() error {
	g.resumeChan <- struct{}{}
	return nil
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
	e, err := gb.NewGameboy(options.Game(), options.GameData())
	if err != nil {
		return err
	}
	g.cancel()
	g.gb = e
	g.frameConsumer = options.FrameConsumer()
	ctx, cancel := context.WithCancel(context.Background())
	g.cancel = cancel
	go g.emulatorLoop(ctx)
	return nil
}

func (g *GoboyAdapter) Stop() error {
	g.cancel()
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

func (g *GoboyAdapter) setButtonBit(bit byte, pressed bool) {
	if pressed {
		g.gb.InputMask &= ^(1 << bit)
	} else {
		g.gb.InputMask |= 1 << bit
	}
}

func (g *GoboyAdapter) SetGraphicOptions(options *GraphicOptions) {
	//TODO implement me
	panic("implement me")
}

func (g *GoboyAdapter) GetGraphicOptions() *GraphicOptions {
	panic("not implemented")
}

func (g *GoboyAdapter) GetCPUBoostRate() float64 {
	return g.boost
}

func (g *GoboyAdapter) SetCPUBoostRate(f float64) float64 {
	g.boost = max(0.5, min(f, 2.0))
	return g.boost
}

func (g *GoboyAdapter) OutputResolution() (width, height int) {
	return 160, 144
}

func (g *GoboyAdapter) MultiController() bool {
	return true
}

func (g *GoboyAdapter) ControllerInfos() []ControllerInfo {
	return []ControllerInfo{
		{
			ControllerId: 1,
			Label:        "玩家1",
		},
	}
}

type GoboyFrame struct {
	data [160][144][3]uint8
	yuv  *image.YCbCr
}

func newGoboyFrame() *GoboyFrame {
	return &GoboyFrame{
		yuv: image.NewYCbCr(image.Rect(0, 0, 160, 144), image.YCbCrSubsampleRatio420),
	}
}

func (g *GoboyFrame) reset(data [160][144][3]uint8) {
	g.data = data
	for x := 0; x < 160; x++ {
		for y := 0; y < 144; y++ {
			yOff := g.yuv.YOffset(x, y)
			cOff := g.yuv.COffset(x, y)
			y, cb, cr := color.RGBToYCbCr(data[x][y][0], data[x][y][1], data[x][y][2])
			g.yuv.Y[yOff] = y
			g.yuv.Cb[cOff] = cb
			g.yuv.Cr[cOff] = cr
		}
	}
}

func (g *GoboyFrame) Width() int {
	return 160
}

func (g *GoboyFrame) Height() int {
	return 144
}

func (g *GoboyFrame) YCbCr() *image.YCbCr {
	return g.yuv
}

func (g *GoboyFrame) Read() (image.Image, func(), error) {
	return g.yuv, func() {}, nil
}
