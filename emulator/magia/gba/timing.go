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
	"github.com/StellrisJAY/cloud-emu/emulator/magia/gba/ram"
)

var (
	wsN  = [4]int{4, 3, 2, 8}
	wsS0 = [2]int{2, 1}
	wsS1 = [2]int{4, 1}
	wsS2 = [2]int{8, 1}
)

func (g *GBA) cycleN(addr uint32) int {
	switch {
	case ram.EWRAM(addr):
		return 3
	case ram.GamePak0(addr):
		offset := ram.IOOffset(ram.WAITCNT)
		idx := g.RAM.IO[offset] >> 2 & 0b11
		return wsN[idx] + 1
	case ram.GamePak1(addr):
		idx := g._getRAM(ram.WAITCNT) >> 5 & 0b11
		return wsN[idx] + 1
	case ram.GamePak2(addr):
		idx := g._getRAM(ram.WAITCNT) >> 8 & 0b11
		return wsN[idx] + 1
	case ram.SRAM(addr):
		idx := g._getRAM(ram.WAITCNT) & 0b11
		return wsN[idx] + 1
	}
	return 1
}

func (g *GBA) cycleS(addr uint32) int {
	switch {
	case ram.EWRAM(addr):
		return 3
	case ram.GamePak0(addr):
		offset := ram.IOOffset(ram.WAITCNT)
		idx := g.RAM.IO[offset] >> 4 & 0b1
		return wsS0[idx] + 1
	case ram.GamePak1(addr):
		idx := g._getRAM(ram.WAITCNT) >> 7 & 0b1
		return wsS1[idx] + 1
	case ram.GamePak2(addr):
		idx := g._getRAM(ram.WAITCNT) >> 10 & 0b1
		return wsS2[idx] + 1
	case ram.SRAM(addr):
		idx := g._getRAM(ram.WAITCNT) & 0b11
		return wsN[idx] + 1
	}
	return 1
}

func (g *GBA) waitBus(addr uint32, size int, s bool) int {
	busWidth := ram.BusWidth(addr)
	if busWidth == 8 {
		return 5 * (size / 8)
	}

	if size > busWidth {
		if s {
			return 2 * g.cycleS(addr)
		}
		return g.cycleN(addr) + g.cycleS(addr+2)
	}

	if s {
		return g.cycleS(addr)
	}
	return g.cycleN(addr)
}
