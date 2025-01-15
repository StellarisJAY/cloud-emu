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
import "github.com/StellrisJAY/cloud-emu/emulator/dawngb/gb/util"

func (c *CPU) ReadIO(addr uint16) uint8 {
	switch addr {
	case 0xFF00:
		return c.joypad.read()
	case 0xFF01:
		return c.serial.sb
	case 0xFF02:
		return c.serial.sc
	case 0xFF04, 0xFF05, 0xFF06, 0xFF07:
		return c.timer.Read(addr)
	case 0xFF0F:
		val := uint8(0)
		for i := 0; i < 5; i++ {
			val |= (util.Btou8(c.interrupt[i]) << i)
		}
		return val
	case 0xFF4C:
		return c.key0
	case 0xFF4D:
		if c.isCGB {
			key1 := c.key1 | 0x7E
			if c.Clock == 4 { // 2x
				key1 |= 1 << 7
			}
			return key1
		}
	case 0xFF50:
		return 1
	case 0xFF55:
		if c.isCGB {
			return c.DMA.Read(addr)
		}
	case 0xFF72:
		return c.ff72
	case 0xFF73:
		return c.ff73
	case 0xFF74:
		return c.ff74
	}
	return 0
}

func (c *CPU) WriteIO(addr uint16, val uint8) {
	switch addr {
	case 0xFF00:
		c.joypad.write(val)
	case 0xFF01:
		c.serial.sb = val
	case 0xFF02:
		c.serial.setSC(val)
	case 0xFF04, 0xFF05, 0xFF06, 0xFF07:
		c.timer.Write(addr, val)
	case 0xFF0F:
		for i := 0; i < 5; i++ {
			c.interrupt[i] = (val & (1 << i)) != 0
		}
	case 0xFF4C:
		if c.key0 == 0 {
			c.key0 = val
		}
	case 0xFF4D:
		if c.isCGB {
			c.key1 = (c.key1 & 0x80) | (val & 0x01)
		}
	case 0xFF50: // BANK
		c.bios.ff50 = false
	case 0xFF51, 0xFF52, 0xFF53, 0xFF54:
		if c.isCGB {
			c.DMA.Write(addr, val)
		}
	case 0xFF55:
		if c.isCGB {
			cycles := c.DMA.Write(addr, val)
			c.Cycles += cycles
		}
	case 0xFF72:
		c.ff72 = val
	case 0xFF73:
		c.ff73 = val
	case 0xFF74:
		c.ff74 = val
	}
}
