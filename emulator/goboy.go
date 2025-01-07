package emulator

import (
	"context"
	"github.com/Humpheh/goboy/pkg/gb"
	"image"
	"image/color"
	"os"
	"time"
)

type GoboyAdapter struct {
	gb            *gb.Gameboy
	cancel        context.CancelFunc
	ticker        *time.Ticker
	frame         *GoboyFrame
	frameConsumer func(frame IFrame)
	game          string
}

func newGoboyAdapter(options IEmulatorOptions) (IEmulator, error) {
	// goboy 只能从文件系统读取rom，所以需要先把rom文件缓存到本地
	cacheGameFile(options.Game(), options.GameData())
	e, err := gb.NewGameboy(options.Game())
	if err != nil {
		return nil, err
	}
	g := &GoboyAdapter{
		gb:            e,
		frameConsumer: options.FrameConsumer(),
		frame:         newGoboyFrame(),
		game:          options.Game(),
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
	ticker := time.NewTicker(time.Millisecond * 16)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			g.gb.Update()
			g.frame.reset(g.gb.PreparedData)
			g.frameConsumer(g.frame)
		}
	}
}

func cacheGameFile(name string, data []byte) {
	_ = os.WriteFile(name, data, 0666)
}

func (g *GoboyAdapter) Pause() error {
	g.gb.ExecutionPaused = true
	return nil
}

func (g *GoboyAdapter) Resume() error {
	g.gb.ExecutionPaused = false
	return nil
}

func (g *GoboyAdapter) Save() (IEmulatorSave, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GoboyAdapter) LoadSave(save IEmulatorSave) error {
	//TODO implement me
	panic("implement me")
}

func (g *GoboyAdapter) Restart(options IEmulatorOptions) error {
	cacheGameFile(options.Game(), options.GameData())
	e, err := gb.NewGameboy(options.Game())
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

func (g *GoboyAdapter) SubmitInput(controllerId int, keyCode string, pressed bool) {
	directionMask, buttonMask := byte(0b11101111), byte(0b11011111)
	switch keyCode {
	case "Right":
		g.gb.InputMask &= directionMask
		g.setButtonBit(0, pressed)
	case "Left":
		g.gb.InputMask &= directionMask
		g.setButtonBit(1, pressed)
	case "Down":
		g.gb.InputMask &= directionMask
		g.setButtonBit(2, pressed)
	case "Up":
		g.gb.InputMask &= directionMask
		g.setButtonBit(3, pressed)
	case "A":
		g.gb.InputMask &= buttonMask
		g.setButtonBit(0, pressed)
	case "B":
		g.gb.InputMask &= buttonMask
		g.setButtonBit(1, pressed)
	case "Select":
		g.gb.InputMask &= buttonMask
		g.setButtonBit(2, pressed)
	case "Start":
		g.gb.InputMask &= buttonMask
		g.setButtonBit(3, pressed)
	}
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

func (g *GoboyAdapter) GetCPUBoostRate() float64 {
	//TODO implement me
	panic("implement me")
}

func (g *GoboyAdapter) SetCPUBoostRate(f float64) float64 {
	//TODO implement me
	panic("implement me")
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
		{
			ControllerId: 2,
			Label:        "玩家2",
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
