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
type mbc1 struct {
	c                *Cartridge
	ramEnabled       bool
	romBank, ramBank uint8
	mode             uint8
}

func newMBC1(c *Cartridge) *mbc1 {
	return &mbc1{
		c:       c,
		romBank: 1,
	}
}

func (m *mbc1) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3:
		return m.c.ROM[addr&0x3FFF]
	case 0x4, 0x5, 0x6, 0x7:
		romBank := uint(m.romBank)
		if m.mode == 0 {
			if len(m.c.ROM) >= int(1*MB) {
				romBank |= (uint(m.ramBank) << 5)
			}
		}
		return m.c.ROM[(romBank<<14)|uint(addr&0x3FFF)]
	case 0xA, 0xB:
		if m.ramEnabled {
			ramBank := uint(0)
			if m.mode == 1 {
				if len(m.c.ram) >= int(32*KB) {
					ramBank = uint(m.ramBank)
				}
			}
			return m.c.ram[(ramBank<<13)|uint(addr&0x1FFF)]
		}
	}
	return 0xFF
}

func (m *mbc1) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0x0, 0x1:
		m.ramEnabled = (val&0x0F == 0x0A)
	case 0x2, 0x3:
		m.romBank = val & 0b11111
		if m.romBank == 0 {
			m.romBank = 1
		}
	case 0x4, 0x5:
		m.ramBank = val & 0b11
	case 0x6, 0x7:
		m.mode = val & 0b1
	case 0xA, 0xB:
		if m.ramEnabled {
			ramBank := uint(0)
			if m.mode == 1 {
				ramBank = uint(m.ramBank)
			}
			bank := m.c.ram[(8*KB)*ramBank:]
			addr &= 0x1FFF
			if len(bank) > int(addr) {
				bank[addr] = val
			}
		}
	}
}
