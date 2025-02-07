package cpu

//Copyright (c) 2024 Akihiro Otomo https://github.com/akatsuki105/dawngb
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.
import (
	"errors"
	"github.com/StellrisJAY/cloud-emu/emulator/dawngb/gb/cpu/sm83"
	"github.com/StellrisJAY/cloud-emu/emulator/dawngb/gb/util"
)

const (
	IRQ_VBLANK = iota
	IRQ_LCDSTAT
	IRQ_TIMER
	IRQ_SERIAL
	IRQ_JOYPAD
)

// DMG-CPU, CGB-CPU
type CPU struct {
	isCGB  bool  // ハードがCGBかどうか
	Cycles int64 // 8MHzのマスターサイクル単位
	*sm83.SM83
	bus              sm83.Bus
	Clock            int64 // 8(1x) or 4(2x)
	timer            *timer
	DMA              *DMA
	joypad           *joypad
	serial           *serial
	bios             BIOS
	HRAM             [0x7F]uint8
	halted           bool
	IE               uint8
	interrupt        [5]bool // IF
	key0             uint8   // FF4C
	key1             uint8   // FF4D
	ff72, ff73, ff74 uint8
}

// a.k.a. Boot ROM
type BIOS struct {
	ff50 bool
	data []uint8
}

func New(isCGB bool, bus sm83.Bus) *CPU {
	c := &CPU{
		isCGB: isCGB,
		bus:   bus,
	}
	c.SM83 = sm83.New(c, c.halt, c.stop, c.wait)
	c.timer = newTimer(c.IRQ, &c.Clock)
	c.joypad = newJoypad(c.IRQ)
	c.DMA = newDMA(c)
	c.serial = newSerial(c.IRQ)
	return c
}

func (c *CPU) Reset() {
	c.Cycles = 0
	c.SM83.Reset()
	c.Clock = 8
	c.timer.reset()
	c.DMA.reset()
	c.joypad.reset()
	c.serial.reset()
	clear(c.HRAM[:])
	c.halted = false
	c.IE, c.interrupt = 0, [5]bool{}
	c.bios.ff50 = true
	c.key0, c.key1 = 0, 0
	c.ff72, c.ff73, c.ff74 = 0, 0, 0
}

func (c *CPU) SkipBIOS() {
	c.bios.ff50 = false
	c.timer.tac = 0xF8
	c.joypad.write(0x30)
	c.joypad.write(0xCF)
	c.DMA.skipBIOS()

	if c.isCGB {
		c.R.A = 0x11
		c.R.F.Unpack(0x80)
		c.R.BC.Unpack(0x0000)
		c.R.DE.Unpack(0xFF56)
		c.R.HL.Unpack(0x000D)
	} else {
		c.R.A = 0x01
		c.R.F.Unpack(0x80)
		c.R.BC.Unpack(0x0013)
		c.R.DE.Unpack(0x00D8)
		c.R.HL.Unpack(0x014D)
	}
	c.R.SP, c.R.PC = 0xFFFE, 0x0100
}

func (c *CPU) wait(n int64) {
	c.Cycles += n * c.Clock
}

func (c *CPU) LoadBIOS(bios []uint8) error {
	c.bios.ff50 = false

	switch len(bios) {
	case 256: // DMG, MGB, SGB
		c.bios.data = make([]uint8, 256)
		copy(c.bios.data[:], bios)
	case 2048: // CGB, AGB
		c.bios.data = make([]uint8, 2048)
		copy(c.bios.data[:], bios)
	case 2048 + 256: // CGB, AGB (0x100..200 is padded)
		c.bios.data = make([]uint8, 2048)
		copy(c.bios.data[:256], bios[:256]) // 0x000..100
		copy(c.bios.data[256:], bios[512:]) // 0x200..900
	default:
		return errors.New("invalid BIOS size")
	}

	c.bios.ff50 = true
	return nil
}

func (c *CPU) stop() {
	if c.key1&(1<<0) != 0 {
		if c.Clock == 4 {
			c.Clock = 8
		} else {
			c.Clock = 4
		}
		c.key1 &^= 1 << 0
	}
}

func (c *CPU) HBlank() {
	if c.isCGB {
		c.DMA.startHDMA()
	}
}

func (c *CPU) Step() int64 {
	cycles := c.step()
	c.timer.run(cycles)
	c.serial.run(cycles)
	return cycles
}

func (c *CPU) step() int64 {
	prev := c.Cycles
	if c.DMA.doHDMA {
		c.DMA.doHDMA = false
		c.DMA.runHDMA()
		c.Cycles += 64
		return c.Cycles - prev
	}

	irqID := c.checkInterrupt()
	if irqID >= 0 {
		c.halted = false
		if c.IME {
			c.interrupt[irqID] = false
			c.Interrupt(irqID)
		} else {
			c.SM83.Step()
		}
	} else if c.halted {
		c.Cycles++
	} else {
		c.SM83.Step()
	}

	return c.Cycles - prev
}

func (c *CPU) IRQ(id int) { c.interrupt[id] = true }

func (c *CPU) SendInputs(inputs uint8) {
	c.joypad.inputs = inputs
}

func (c *CPU) checkInterrupt() int {
	for i := 0; i < 5; i++ {
		if util.Bit(c.IE, i) && c.interrupt[i] {
			return i
		}
	}
	return -1
}

func (c *CPU) halt() {
	if c.IME {
		c.halted = true
	} else {
		if c.checkInterrupt() < 0 {
			c.halted = true
		}
	}
}

func (c *CPU) IsCGBMode() bool {
	return c.isCGB && c.key0 != 4
}
