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
type mbc2 struct {
	c          *Cartridge
	ramEnabled bool
	romBank    uint8 // 0..15
}

func newMBC2(c *Cartridge) *mbc2 {
	return &mbc2{
		c:       c,
		romBank: 1,
	}
}

func (m *mbc2) reset() {
	m.ramEnabled = false
	m.romBank = 1
}

func (m *mbc2) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3: // ROMバンク0
		return m.c.ROM[addr&0x3FFF]
	case 0x4, 0x5, 0x6, 0x7: // ROMバンク1..15
		return m.c.ROM[(uint(m.romBank)<<14)|uint(addr&0x3FFF)]
	case 0xA, 0xB: // RAM
		if m.ramEnabled {
			return 0xF0 | (m.c.ram[addr&0x1FF] & 0x0F)
		}
	}
	return 0xFF
}

func (m *mbc2) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3: // RAM有効化 or ROMバンク切り替え
		mode := (addr >> 8) & 0x1
		if mode == 0 {
			m.ramEnabled = val == 0x0A
		} else {
			m.romBank = val & 0x0F
			if m.romBank == 0 {
				m.romBank = 1
			}
		}
	case 0xA, 0xB: // RAM書き込み
		if m.ramEnabled {
			m.c.ram[addr&0x1FF] = val & 0x0F
		}
	}
}
