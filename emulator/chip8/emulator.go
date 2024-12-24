package chip8

import (
	"context"
	"time"
)

type Emulator struct {
	v  [16]byte // 通用寄存器v0-vf
	pc uint16   // 程序计数器
	sp uint16   // 栈指针
	i  uint16   // 地址寄存器
	dt byte     // 延迟定时器
	st byte     // 声音定时器

	clock              uint64 // 时钟，记录从开始到现在的时间ns
	instructionCounter uint64 // 指令计数器，记录从开始到现在执行的指令数量

	delayTicker *time.Ticker
	soundTicker *time.Ticker
	m           *memory

	delayDuration time.Duration
	soundDuration time.Duration

	renderCallback func(*Frame)

	frameTicker *time.Ticker

	pauseChan chan struct{}
}

var (
	defaultDelayDuration = time.Second / 60
	defaultSoundDuration = time.Second / 60
)

func NewEmulator(game []byte, renderCallback func(*Frame)) *Emulator {
	e := &Emulator{
		m:              newMemory(),
		delayTicker:    time.NewTicker(defaultDelayDuration),
		soundTicker:    time.NewTicker(defaultSoundDuration),
		delayDuration:  defaultDelayDuration,
		soundDuration:  defaultSoundDuration,
		renderCallback: renderCallback,
		frameTicker:    time.NewTicker(time.Second / 60),
	}
	copy(e.m.programData[:], game)
	return e
}

func (e *Emulator) reset() {
	e.v = [16]byte{}
	e.pc = 0x200
	e.sp = dataStart
	e.i = 0
	e.dt = 0
	e.st = 0
	e.clock = 0
	e.instructionCounter = 0
}

func (e *Emulator) Run(ctx context.Context) {
	e.reset()
	for {
		select {
		case <-e.pauseChan:
			e.waitResume()
		case <-ctx.Done():
			return
		case <-e.delayTicker.C:
			if e.dt > 0 {
				e.dt--
			} else {
				e.delayTicker.Stop()
			}
		case <-e.soundTicker.C:
			if e.st > 0 {
				e.st--
			} else {
				e.soundTicker.Stop() //
			}
		case <-e.frameTicker.C:
			e.renderCallback(e.m.frame)
		default:
			e.step()
		}
	}
}

func (e *Emulator) waitResume() {
	for {
		select {
		case <-e.pauseChan:
			return
		}
	}
}

func (e *Emulator) Pause() {
	e.pauseChan <- struct{}{}
}

func (e *Emulator) Resume() {
	e.pauseChan <- struct{}{}
}

func (e *Emulator) step() {
	instruction := e.m.readUint16(e.pc)
	code := instruction & 0xF000
	switch code {
	case 0x0000:
		e.sys(instruction)
	case 0x1000:
		e.jump(instruction)
	case 0x2000:
		e.call(instruction)
	case 0x3000:
		e.skipIfXNEqual(instruction)
	case 0x4000:
		e.skipIfXNNotEqual(instruction)
	case 0x5000:
		e.skipIfXYEqual(instruction)
	case 0x6000:
		e.storeNX(instruction)
	case 0x7000:
		e.addXN(instruction)
	case 0x8000:
		e.arithmetic(instruction)
	case 0x9000:
		e.skipIfNotEqualXY(instruction)
	case 0xA000:
		e.storeNI(instruction)
	case 0xB000:
		e.jumpV0(instruction)
	case 0xC000:
		e.randMask(instruction)
	case 0xD000:
		e.draw(instruction)
	case 0xE000:
		// TODO keyboard input
	case 0xF000:
		e.misc(instruction)
	}
}

func (e *Emulator) sys(instruction uint16) {
	switch instruction {
	case 0x00E0:
		e.m.frame.clear()
		e.pc += 2
	case 0x00EE:
		e.ret()
	default:
		e.call(instruction)
	}
}

func (e *Emulator) arithmetic(instruction uint16) {
	switch instruction & 0xf {
	case 0x0:
		e.storeYX(instruction)
	case 0x1:
		e.or(instruction)
	case 0x2:
		e.and(instruction)
	case 0x3:
		e.xor(instruction)
	case 0x4:
		e.addXY(instruction)
	case 0x5:
		e.subXY(instruction)
	case 0x6:
		e.rShiftYX(instruction)
	case 0x7:
		e.minusYX(instruction)
	case 0xE:
		e.lShiftYX(instruction)
	}
}

func (e *Emulator) misc(instruction uint16) {
	switch instruction & 0xff {
	case 0x07:
	case 0x0A:
	case 0x15:
		e.setDelayTimer(instruction)
	case 0x18:
		e.setSoundTimer(instruction)
	case 0x1E:
	case 0x29:
	case 0x33:
	case 0x55:
		e.storeVxIRanged(instruction)
	case 0x65:
		e.loadIVxRanged(instruction)
	}
}
