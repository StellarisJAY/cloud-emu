package sm83

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
	"fmt"
)

type Bus interface {
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
}

type SM83 struct {
	R    Registers
	bus  Bus
	inst struct {
		opcode uint8
		addr   uint16
		cb     bool
	}
	IME        bool
	halt, stop func()
	tick       func(clockCycles int64)
}

func New(bus Bus, halt, stop func(), tick func(int64)) *SM83 {
	if tick == nil {
		panic("tick function is required")
	}
	return &SM83{
		bus:  bus,
		halt: halt,
		stop: stop,
		tick: tick,
	}
}

func (c *SM83) Reset() {
	c.R.reset()
	c.IME = false
}

func (c *SM83) Step() {
	pc := c.R.PC
	c.inst.addr = pc
	opcode := c.fetch()
	c.inst.opcode = opcode
	c.inst.cb = false

	fn := opTable[opcode]
	if fn != nil {
		// fmt.Printf("0x%02X in 0x%04X\n", opcode, pc)
		fn(c)
	} else {
		panic(fmt.Sprintf("illegal opcode: 0x%02X in 0x%04X", opcode, pc))
	}

	c.tick(opCycles[opcode])
}

func (c *SM83) fetch() uint8 {
	pc := c.R.PC
	c.R.PC++
	return c.bus.Read(pc)
}
