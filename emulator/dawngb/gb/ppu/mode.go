package ppu

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
	"github.com/StellrisJAY/cloud-emu/emulator/dawngb/gb/util"
)

// Mode 0
func (p *PPU) hblank() {
	oldStat := p.stat
	p.stat = (p.stat & 0xFC)
	if util.Bit(p.lcdc, 7) && !p.enableLatch {
		p.r.DrawScanline(p.ly, p.screen[p.ly*160:(p.ly+1)*160])
	}
	if !statIRQAsserted(oldStat) && statIRQAsserted(p.stat) {
		p.cpu.IRQ(1)
	}
	p.cpu.HBlank()
}

// Mode 1
func (p *PPU) vblank() {
	oldStat := p.stat
	p.stat = (p.stat & 0xFC) | 1
	p.cpu.IRQ(0)

	if !statIRQAsserted(oldStat) && statIRQAsserted(p.stat) {
		p.cpu.IRQ(1)
	}
}

// Mode 2
func (p *PPU) scanOAM() {
	oldStat := p.stat
	p.stat = (p.stat & 0xFC) | 2
	if !statIRQAsserted(oldStat) && statIRQAsserted(p.stat) {
		p.cpu.IRQ(1)
	}
}

// Mode 3
func (p *PPU) drawing() {
	p.stat = (p.stat & 0xFC) | 3

	// Count scanline objects
	h := 8
	if util.Bit(p.lcdc, 2) {
		h = 16
	}
	o := uint8(0)
	for i := 0; i < 40; i++ {
		y := int(p.OAM[i*4]) - 16
		if y <= p.ly && p.ly < y+h {
			o++
		}
	}
	p.objCount = o
}
