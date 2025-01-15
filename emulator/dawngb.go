package emulator

import (
	"bytes"
	"context"
	"github.com/StellrisJAY/cloud-emu/emulator/dawngb/gb"
	"image"
	"strings"
	"time"
)

type DawnGbAdapter struct {
	g             *gb.GB
	frame         *BaseFrame
	frameConsumer func(IFrame)
	ticker        *time.Ticker
	cancel        context.CancelFunc
	audioBuffer   *bytes.Buffer
}

func init() {
	supportedEmulators[CodeDawnGB] = Info{
		EmulatorCode:           CodeDawnGB,
		EmulatorType:           "GBC",
		Provider:               "https://github.com/akatsuki105/dawngb",
		Description:            "Go语言实现的GB模拟器",
		Name:                   "DawnGB",
		SupportSave:            false,
		SupportGraphicSettings: false,
	}
}

func newDawnGbAdapter(options IEmulatorOptions) (*DawnGbAdapter, error) {
	buffer := bytes.NewBuffer([]byte{})
	g := gb.New(gb.MODEL_CGB, buffer)
	if err := g.Load(gb.LOAD_ROM, options.GameData()); err != nil {
		return nil, err
	}
	return &DawnGbAdapter{
		g:             g,
		frame:         MakeEmptyBaseFrame(image.Rect(0, 0, 160, 144)),
		frameConsumer: options.FrameConsumer(),
		audioBuffer:   buffer,
	}, nil
}

func (d *DawnGbAdapter) emulatorLoop(ctx context.Context) {
	d.ticker = time.NewTicker(FrameInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-d.ticker.C:
			start := time.Now()
			d.g.RunFrame()
			screen := d.g.Screen()
			d.frame.FromNRGBAColors(screen)
			d.frameConsumer(d.frame)
			interval := max(FrameInterval-time.Since(start), time.Millisecond*5)
			d.ticker.Reset(interval)
		}
	}
}

func (d *DawnGbAdapter) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	d.cancel = cancel
	d.g.Reset(false)
	go d.emulatorLoop(ctx)
	return nil
}

func (d *DawnGbAdapter) Pause() error {
	d.ticker.Stop()
	return nil
}

func (d *DawnGbAdapter) Resume() error {
	d.ticker.Reset(FrameInterval)
	return nil
}

func (d *DawnGbAdapter) Save() (IEmulatorSave, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DawnGbAdapter) LoadSave(save IEmulatorSave) error {
	//TODO implement me
	panic("implement me")
}

func (d *DawnGbAdapter) Restart(options IEmulatorOptions) error {
	buffer := bytes.NewBuffer([]byte{})
	g := gb.New(gb.MODEL_CGB, buffer)
	if err := g.Load(gb.LOAD_ROM, options.GameData()); err != nil {
		return err
	}
	d.cancel()
	d.g = g
	d.frameConsumer = options.FrameConsumer()
	d.g.Reset(false)
	return d.Start()
}

func (d *DawnGbAdapter) Stop() error {
	d.cancel()
	return nil
}

func (d *DawnGbAdapter) SubmitInput(controllerId int, keyCode string, pressed bool) {
	d.g.SetKeyInput(strings.ToUpper(keyCode), pressed)
}

func (d *DawnGbAdapter) SetGraphicOptions(options *GraphicOptions) {
	//TODO implement me
	panic("implement me")
}

func (d *DawnGbAdapter) GetCPUBoostRate() float64 {
	//TODO implement me
	panic("implement me")
}

func (d *DawnGbAdapter) SetCPUBoostRate(f float64) float64 {
	//TODO implement me
	panic("implement me")
}

func (d *DawnGbAdapter) OutputResolution() (width, height int) {
	return 160, 144
}

func (d *DawnGbAdapter) MultiController() bool {
	return true
}

func (d *DawnGbAdapter) ControllerInfos() []ControllerInfo {
	return []ControllerInfo{
		{
			ControllerId: 1,
			Label:        "玩家1",
		},
	}
}
