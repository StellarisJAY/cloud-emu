package emulator

import (
	"bytes"
	"context"
	"encoding/binary"
	_ "encoding/binary"
	"github.com/StellrisJAY/cloud-emu/emulator/dawngb/gb"
	"strings"
)

type DawnGbAdapter struct {
	BaseEmulatorAdapter
	g              *gb.GB
	audioBuffer    *bytes.Buffer
	leftAudioChan  chan float32
	rightAudioChan chan float32
}

func init() {
	supportedEmulators[CodeDawnGB] = Info{
		EmulatorCode:           CodeDawnGB,
		EmulatorType:           "GBC",
		Provider:               "https://github.com/akatsuki105/dawngb",
		Description:            "Go语言实现的GB模拟器",
		Name:                   "DawnGB",
		SupportSave:            true,
		SupportGraphicSettings: true,
	}
}

func newDawnGbAdapter(options IEmulatorOptions) (*DawnGbAdapter, error) {
	buffer := bytes.NewBuffer([]byte{})
	g := gb.New(gb.MODEL_CGB, buffer)
	if err := g.Load(gb.LOAD_ROM, options.GameData()); err != nil {
		return nil, err
	}
	a := &DawnGbAdapter{
		g:                   g,
		audioBuffer:         buffer,
		BaseEmulatorAdapter: newBaseEmulatorAdapter(160, 144, options),
	}
	a.stepFunc = a.step
	return a, nil
}

func (d *DawnGbAdapter) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	d.cancel = cancel
	d.g.Reset(false)
	go d.emulatorLoop(ctx)
	return nil
}

func (d *DawnGbAdapter) step() {
	d.g.RunFrame()
	screen := d.g.Screen()
	d.frame.FromNRGBAColors(screen, d.scale)
	d.frameConsumer(d.frame)
	if len(d.audioBuffer.Bytes()) == 0 {
		return
	}
	reader := bytes.NewReader(d.audioBuffer.Bytes())
	samples := make([]float32, 0, len(d.audioBuffer.Bytes())/2)
	for {
		var sample16 int16
		if err := binary.Read(reader, binary.LittleEndian, &sample16); err != nil {
			break
		}
		samples = append(samples, float32(sample16)/32768.0)
	}
	d.audioConsumer(samples)
	d.audioBuffer.Reset()
}

func (d *DawnGbAdapter) Save() (IEmulatorSave, error) {
	saveData, err := d.g.Dump(gb.DUMP_SAVE)
	if err != nil {
		return nil, err
	}
	return &BaseEmulatorSave{
		Data: saveData,
	}, nil
}

func (d *DawnGbAdapter) LoadSave(save IEmulatorSave) error {
	return d.g.Load(gb.LOAD_SAVE, save.SaveData())
}

func (d *DawnGbAdapter) Restart(options IEmulatorOptions) error {
	_ = d.Stop()
	buffer := bytes.NewBuffer([]byte{})
	g := gb.New(gb.MODEL_CGB, d.audioBuffer)
	if err := g.Load(gb.LOAD_ROM, options.GameData()); err != nil {
		return err
	}
	d.g = g
	d.audioBuffer = buffer
	d.audioConsumer = options.AudioConsumer()
	d.frameConsumer = options.FrameConsumer()
	return d.Start()
}

func (d *DawnGbAdapter) SubmitInput(_ int, keyCode string, pressed bool) {
	d.g.SetKeyInput(strings.ToUpper(keyCode), pressed)
}

func (d *DawnGbAdapter) ControllerInfos() []ControllerInfo {
	return []ControllerInfo{
		{
			ControllerId: 1,
			Label:        "玩家1",
		},
	}
}
