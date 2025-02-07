package gba

//Copyright (c) 2021 Akihiro Otomo https://github.com/akatsuki105/magia
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
//of the Software, and to permit persons to whom the Software is furnished to do
//so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
//INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
//PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
//HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
//CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
//OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

import (
	"github.com/StellrisJAY/cloud-emu/emulator/magia/gba/apu"
	"github.com/StellrisJAY/cloud-emu/emulator/magia/util"
)

type dmaTiming uint32

const (
	dmaImmediate dmaTiming = iota
	dmaVBlank
	dmaHBlank
	dmaSpecial
)

type DMA struct {
	io                  [12]byte
	src, dst            uint32
	count, defaultCount int
}

func NewDMA() [4]*DMA {
	return [4]*DMA{{defaultCount: 0x4000}, {defaultCount: 0x4000}, {defaultCount: 0x4000}, {defaultCount: 0x10000}}
}
func (ch *DMA) cnt() uint32 { return util.LE32(ch.io[8:]) }
func (ch *DMA) setCnt(v uint32) {
	ch.io[8], ch.io[9], ch.io[10], ch.io[11] = byte(v), byte(v>>8), byte(v>>16), byte(v>>24)
}
func isDMA0IO(addr uint32) bool { return 0x0400_00B0 <= addr && addr <= 0x0400_00BB }
func isDMA1IO(addr uint32) bool { return 0x0400_00BC <= addr && addr <= 0x0400_00C7 }
func isDMA2IO(addr uint32) bool { return 0x0400_00C8 <= addr && addr <= 0x0400_00D3 }
func isDMA3IO(addr uint32) bool { return 0x0400_00D4 <= addr && addr <= 0x0400_00DF }

func (ch *DMA) get(ofs uint32) uint32 { return util.LE32(ch.io[ofs:]) }
func (ch *DMA) set(ofs uint32, b byte) bool {
	old := byte(ch.cnt() >> 24)
	ch.io[ofs] = b
	if ofs == 11 {
		turnon := !util.Bit(old, 7) && util.Bit(b, 7)
		if turnon {
			ch.src, ch.dst = util.LE32(ch.io[0:]), util.LE32(ch.io[4:])
			ch.count = ch.wordCount()
			switch ch.size() {
			case 32:
				ch.src &= ^uint32(3)
				ch.dst &= ^uint32(3)
			case 16:
				ch.src &= ^uint32(1)
				ch.dst &= ^uint32(1)
			}
		}

		return turnon && ch.timing() == 0
	}
	return false
}

func (ch *DMA) dstCnt() int64 {
	switch (ch.cnt() >> (16 + 5)) & 0b11 {
	case 0, 3:
		return int64(ch.size()) / 8
	case 1:
		return -int64(ch.size()) / 8
	default:
		return 0
	}
}
func (ch *DMA) dstReload() bool { return (ch.cnt()>>(16+5))&0b11 == 3 }
func (ch *DMA) srcCnt() int64 {
	switch (ch.cnt() >> (16 + 7)) & 0b11 {
	case 0:
		return int64(ch.size()) / 8
	case 1:
		return -int64(ch.size()) / 8
	default:
		return 0
	}
}
func (ch *DMA) repeat() bool      { return util.Bit(ch.cnt(), 16+9) }
func (ch *DMA) size() int         { return map[bool]int{true: 32, false: 16}[util.Bit(ch.cnt(), 16+10)] }
func (ch *DMA) timing() dmaTiming { return dmaTiming((ch.cnt() >> (16 + 12)) & 0b11) }
func (ch *DMA) irq() bool         { return util.Bit(ch.cnt(), 16+14) }
func (ch *DMA) enabled() bool     { return util.Bit(ch.cnt(), 16+15) }
func (ch *DMA) disable()          { ch.setCnt(ch.cnt() & 0x7fff_ffff) }
func (ch *DMA) wordCount() int {
	wordCount := ch.cnt() & 0xffff
	if wordCount == 0 {
		return ch.defaultCount
	}
	return int(wordCount)
}

func (g *GBA) dmaTransfer(t dmaTiming) {
	for i, ch := range g.dma {
		if !ch.enabled() || ch.timing() != t {
			continue
		}

		size := ch.size()
		srcInc, dstInc := ch.srcCnt(), ch.dstCnt()
		for ch.count > 0 {
			g._setRAM(ch.dst, g._getRAM(ch.src), size/8)

			ch.dst, ch.src = uint32(int64(ch.dst)+dstInc), uint32(int64(ch.src)+srcInc)
			ch.count--
		}

		if ch.irq() {
			g.triggerIRQ(irqDMA0 + IRQID(i))
		}

		if ch.repeat() {
			ch.count = ch.wordCount()
			if ch.dstReload() {
				ch.dst = util.LE32(ch.io[4:])
			}
		} else {
			ch.disable()
		}
	}
}

// Receive 4 x 32bit (16 bytes) per DMA
func (g *GBA) dmaTransferFifo(ch int) {
	if !g.apu.IsSoundMasterEnable() || !g.dma[ch].enabled() || g.dma[ch].timing() != dmaSpecial {
		return
	}

	cnt := g.dma[ch].cnt()
	for i := 0; i < 4; i++ { // 32bit × 4 = 4 words
		val := g._getRAM(g.dma[ch].src)
		g._setRAM(g.dma[ch].dst, val, 4)

		if ch == 1 {
			apu.FifoACopy(val)
		} else {
			apu.FifoBCopy(val)
		}

		switch (cnt >> (16 + 7)) & 0b11 {
		case 0:
			g.dma[ch].src += 4
		case 1:
			g.dma[ch].src -= 4
		}
	}

	if g.dma[ch].irq() {
		g.triggerIRQ(irqDMA0 + IRQID(ch))
	}
}
