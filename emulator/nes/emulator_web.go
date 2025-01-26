package nes

import (
	"github.com/StellrisJAY/cloud-emu/emulator/nes/apu"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/bus"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/cartridge"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/config"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/cpu"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/ppu"
	"io"
)

// Emulator browser render nes
type Emulator struct {
	RawEmulator
}

func NewEmulatorWithGameData(game []byte, conf config.Config, apuSampleRate int, audioBuffer io.Writer) (*Emulator, error) {
	c, err := cartridge.MakeCartridge(game)
	if err != nil {
		return nil, err
	}
	e := &Emulator{
		RawEmulator{
			cartridge: c,
			config:    conf,
		},
	}
	e.joyPad1 = bus.NewJoyPad()
	e.joyPad2 = bus.NewJoyPad()
	e.ppu = ppu.NewPPU(e.cartridge.GetChrBank, e.cartridge.GetMirroring, e.cartridge.WriteCHR)
	e.apu = apu.NewBasicAPU(audioBuffer)
	e.bus = bus.NewBus(e.cartridge, e.ppu, e.joyPad1, e.joyPad2, e.apu)
	e.apu.SetRates(bus.CPUFrequency, float64(apuSampleRate))
	e.apu.SetMemReader(e.bus.ReadMemUint8)
	e.processor = cpu.NewProcessor(e.bus)
	return e, nil
}

func (e *Emulator) Frame() *ppu.Frame {
	return e.ppu.Frame()
}
