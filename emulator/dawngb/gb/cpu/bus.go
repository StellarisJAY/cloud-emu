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

func (c *CPU) Read(addr uint16) uint8 {
	if c.bios.ff50 {
		if addr < 0x100 {
			return c.bios.data[addr]
		}
		if len(c.bios.data) == 2048 && (addr >= 0x200 && addr < 0x900) {
			return c.bios.data[addr-0x100]
		}
	}
	if addr >= 0xFF80 && addr <= 0xFFFE { // HRAM
		return c.HRAM[addr&0x7F]
	}
	return c.bus.Read(addr)
}

func (c *CPU) Write(addr uint16, val uint8) {
	if addr >= 0xFF80 && addr <= 0xFFFE { // HRAM
		c.HRAM[addr&0x7F] = val
		return
	}
	c.bus.Write(addr, val)
}
