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
	"github.com/StellrisJAY/cloud-emu/emulator/dawngb/gb/ppu/renderer"
	"github.com/StellrisJAY/cloud-emu/emulator/dawngb/gb/ppu/renderer/software"
	"github.com/StellrisJAY/cloud-emu/emulator/dawngb/gb/util"
	"image/color"
)

const KB = 1024

const CYCLE = 2

// 便宜的にPPU構造体に入れているが、正確にはボード上にある
type VRAM struct {
	data [16 * KB]uint8
	bank uint8 // 0 or 1; VBK(0xFF4F)
}

type CPU interface {
	Read(addr uint16) uint8
	IRQ(id int)
	HBlank()
	IsCGBMode() bool // CGBモードかどうか
}

// OAM DMA
type DMA struct {
	active bool
	src    uint16
	until  int64
}

// SoCに組み込まれているため、`/cpu`にある方が正確ではある
type PPU struct {
	cpu             CPU
	cycles          int64 // 遅れているサイクル数(8.38MHzのマスターサイクル単位)
	screen          [160 * 144]color.NRGBA
	FrameCounter    uint64
	lx, ly          int
	r               renderer.Renderer
	ram             VRAM
	DMA             *DMA
	lcdc, stat, lyc uint8
	OAM             [160]uint8
	Palette         [(4 * 8) * 2]uint16 // 4bppの8パレットが BG と OBJ　の1つずつ
	ioreg           [0x30]uint8
	enableLatch     bool // LCDC.7をセットしてPPUを有効にすると、次のフレームから表示が開始される そうじゃないとゴミが表示される
	objCount        uint8
	bgpi, obpi      uint8
}

func New(cpu CPU) *PPU {
	p := &PPU{
		cpu: cpu,
		DMA: &DMA{},
	}
	return p
}

func (p *PPU) Reset() {
	p.r = software.New(p.ram.data[:], p.Palette[:], p.OAM[:], p.cpu.IsCGBMode)
	p.lx, p.ly = 0, 0
	p.stat = 0x80
	p.ram.bank = 0
	p.objCount = 0
	p.DMA.active, p.DMA.src, p.DMA.until = false, 0, 0
	p.bgpi, p.obpi = 0, 0
	clear(p.Palette[:])
}

func (p *PPU) SkipBIOS() {
	p.Write(0xFF40, 0x91) // LCDC
	p.Write(0xFF47, 0xFC) // BGP
	copy(p.Palette[:4], dmgPalette[:])
	copy(p.Palette[32:36], dmgPalette[:])
}

func (p *PPU) Screen() []color.NRGBA {
	return p.screen[:]
}

func (p *PPU) Run(cycles8MHz int64) {
	if p.DMA.active {
		p.runDMA(cycles8MHz)
	}

	p.cycles += cycles8MHz
	for p.cycles >= 2 { // 1dot = 4MHz
		p.step()
		p.cycles -= 2
	}
}

func (p *PPU) step() {
	if util.Bit(p.lcdc, 7) {
		if p.ly < 144 {
			switch p.lx {
			case 0:
				p.scanOAM()
			case 80:
				p.drawing()
			case 252 + (int(p.objCount) * 6):
				p.hblank()
			}
		}
		p.lx++
		if p.lx == 456 {
			p.lx = 0
			p.incrementLY()
		}
	}
}

func (p *PPU) incrementLY() {
	p.objCount = 0
	p.ly++
	switch p.ly {
	case 144:
		p.vblank()
	case 154:
		p.ly = 0
		p.enableLatch = false
		p.FrameCounter++
	}
	p.compareLYC()
}

func (p *PPU) compareLYC() {
	oldStat := p.stat
	p.stat = util.SetBit(p.stat, 2, p.ly == int(p.lyc))
	if !statIRQAsserted(oldStat) && statIRQAsserted(p.stat) {
		p.cpu.IRQ(1)
	}
}

// GBCのBIOSがやる、DMGゲームに対する色付け処理
func (p *PPU) ColorizeDMG() {
	copy(p.Palette[:4], cgbPalette[:])
	copy(p.Palette[32:36], cgbPalette[4:])
}

func (p *PPU) runDMA(cycles8MHz int64) {
	p.DMA.until -= cycles8MHz
	if p.DMA.until <= 0 {
		for i := uint16(0); i < 160; i++ {
			p.Write(0xFE00+i, p.cpu.Read(p.DMA.src+i))
		}
		p.DMA.active = false
	}
}

func (p *PPU) TriggerDMA(src uint16, m int64) {
	if !p.DMA.active {
		p.DMA.active = true
		p.DMA.src = src
		p.DMA.until = 160 * m
	}
}

func statIRQAsserted(stat uint8) bool {
	if util.Bit(stat, 6) && util.Bit(stat, 2) {
		return true
	}
	switch stat & 0b11 {
	case 0:
		return util.Bit(stat, 3)
	case 1:
		return util.Bit(stat, 4)
	case 2:
		return util.Bit(stat, 5)
	}
	return false
}
