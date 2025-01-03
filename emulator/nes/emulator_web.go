//go:build !sdl

package nes

import (
	"context"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/apu"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/bus"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/cartridge"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/config"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/cpu"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/ppu"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/trace"
	"sync"
)

// Emulator browser render nes
type Emulator struct {
	RawEmulator
}

func NewEmulator(game string, conf config.Config, callback bus.RenderCallback, audioSampleChan chan float32, apuSampleRate int) (*Emulator, error) {
	nesData, err := ReadGameFile(game)
	if err != nil {
		return nil, err
	}
	return NewEmulatorWithGameData(nesData, conf, callback, audioSampleChan, apuSampleRate)
}

func NewEmulatorWithGameData(game []byte, conf config.Config, callback bus.RenderCallback, audioSampleChan chan float32, apuSampleRate int) (*Emulator, error) {
	c, err := cartridge.MakeCartridge(game)
	if err != nil {
		return nil, err
	}
	e := &Emulator{
		RawEmulator{
			cartridge: c,
			config:    conf,
			m:         &sync.Mutex{},
		},
	}
	e.joyPad1 = bus.NewJoyPad()
	e.joyPad2 = bus.NewJoyPad()
	e.ppu = ppu.NewPPU(e.cartridge.GetChrBank, e.cartridge.GetMirroring, e.cartridge.WriteCHR)
	e.apu = apu.NewBasicAPU()
	e.bus = bus.NewBus(e.cartridge, e.ppu, callback, e.joyPad1, e.joyPad2, e.apu)
	e.apu.SetRates(bus.CPUFrequency, float64(apuSampleRate))
	e.apu.SetOutputChan(audioSampleChan)
	e.apu.SetMemReader(e.bus.ReadMemUint8)
	e.processor = cpu.NewProcessor(e.bus)
	return e, nil
}

func (e *Emulator) LoadAndRun(ctx context.Context, enableTrace bool) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	if enableTrace {
		e.processor.LoadAndRunWithCallback(ctx, trace.Trace, nil)
	} else {
		e.processor.LoadAndRunWithCallback(ctx, nil, nil)
	}
}

func (e *Emulator) Frame() *ppu.Frame {
	return e.ppu.Frame()
}
