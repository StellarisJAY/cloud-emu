package cartridge

// Copyright (c) 2024 Akihiro Otomo https://github.com/akatsuki105/dawngb
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
type mbc0 struct {
	c      *Cartridge
	hasRam bool
}

func newMBC0(c *Cartridge) mbc {
	hasRam := c.ROM[0x147] != 0
	return &mbc0{
		c:      c,
		hasRam: hasRam,
	}
}

func (m *mbc0) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7:
		return m.c.ROM[addr]
	case 0xA, 0xB:
		if m.hasRam {
			return m.c.ram[addr&0x1FFF]
		}
	}
	return 0xFF
}

func (m *mbc0) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0xA, 0xB:
		if m.hasRam {
			m.c.ram[addr&0x1FFF] = val
		}
	}
}
