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
type mbc5 struct {
	c          *Cartridge
	hasRam     bool
	ramEnabled bool
	romBank    uint16 // 0..511
	ramBank    uint8  // 0..15
}

func newMBC5(c *Cartridge) mbc {
	hasRam := c.ROM[0x147] != 25 && c.ROM[0x147] != 28
	return &mbc5{
		c:       c,
		hasRam:  hasRam,
		romBank: 1,
	}
}

func (m *mbc5) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3:
		return m.c.ROM[addr]
	case 0x4, 0x5, 0x6, 0x7:
		return m.c.ROM[(uint32(m.romBank)<<14)|(uint32(addr&0x3FFF))]
	case 0xA, 0xB:
		if m.hasRam && m.ramEnabled {
			n := int((uint(m.ramBank) << 13) | uint(addr&0x1FFF))
			if n >= len(m.c.ram) {
				n &= len(m.c.ram) - 1
			}
			return m.c.ram[n]
		}
	}
	return 0xFF
}

func (m *mbc5) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0x0, 0x1:
		m.ramEnabled = (val&0x0F == 0x0A)
	case 0x2:
		m.romBank &= 0x100
		m.romBank |= uint16(val)
	case 0x3:
		m.romBank &= 0xFF
		m.romBank |= uint16(val&0b1) << 8
	case 0x4, 0x5:
		m.ramBank = (val & 0b1111)
	case 0xA, 0xB:
		if m.hasRam && m.ramEnabled {
			n := int((uint(m.ramBank) << 13) | uint(addr&0x1FFF))
			if n >= len(m.c.ram) {
				n &= len(m.c.ram) - 1
			}
			m.c.ram[n] = val
		}
	}
}
