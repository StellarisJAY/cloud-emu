package nes

import (
	"github.com/StellrisJAY/cloud-emu/emulator/nes/apu"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/bus"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/cartridge"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/config"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/cpu"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/ppu"
	"github.com/StellrisJAY/cloud-emu/emulator/nes/trace"
	"time"
)

type RawEmulator struct {
	processor *cpu.Processor
	cartridge cartridge.Cartridge
	bus       *bus.Bus
	ppu       *ppu.PPU
	joyPad1   *bus.JoyPad
	joyPad2   *bus.JoyPad
	apu       *apu.BasicAPU
	config    config.Config

	lastSnapshotTime time.Time
	snapshots        []Snapshot
}

func (e *RawEmulator) Disassemble() {
	e.processor.Disassemble(trace.PrintDisassemble)
}

func (e *RawEmulator) SetJoyPadButtonPressed(id int, button bus.JoyPadButton, pressed bool) {
	if id == 1 {
		e.joyPad1.SetButtonPressed(button, pressed)
	} else {
		e.joyPad2.SetButtonPressed(button, pressed)
	}
}

func (e *RawEmulator) SetCPUBoostRate(rate float64) float64 {
	return e.bus.SetCPUBoostRate(rate)
}

func (e *RawEmulator) BoostCPU(delta float64) float64 {
	return e.bus.BoostCPU(delta)
}

func (e *RawEmulator) CPUBoostRate() float64 {
	return e.bus.CPUBoostRate()
}

func (e *RawEmulator) StepOneFrame() {
	fc := e.ppu.FrameCount
	for fc == e.ppu.FrameCount {
		e.processor.Step()
	}
	e.ppu.FrameCount = 0
}

func (e *RawEmulator) Reset() {
	e.processor.Reset()
}
