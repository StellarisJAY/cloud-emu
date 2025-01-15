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
	"github.com/StellrisJAY/cloud-emu/emulator/magia/util"
)

func (g *GBA) armLSL(val uint32, is uint32, carryMut bool, imm bool) uint32 {
	switch {
	case is == 0 && imm:
		return val
	case is > 32:
		if carryMut {
			g.SetCPSRFlag(flagC, false)
		}
		return 0
	default:
		carry := val&(1<<(32-is)) > 0
		if is > 0 && carryMut {
			g.SetCPSRFlag(flagC, carry)
		}
		return util.LSL(val, uint(is))
	}
}

func (g *GBA) armLSR(val uint32, is uint32, carryMut bool, imm bool) uint32 {
	if is == 0 && imm {
		is = 32
	}
	carry := val&(1<<(is-1)) > 0
	if is > 0 && carryMut {
		g.SetCPSRFlag(flagC, carry)
	}
	return util.LSR(val, uint(is))
}

func (g *GBA) armASR(val uint32, is uint32, carryMut bool, imm bool) uint32 {
	if (is == 0 && imm) || is > 32 {
		is = 32
	}
	carry := val&(1<<(is-1)) > 0
	if is > 0 && carryMut {
		g.SetCPSRFlag(flagC, carry)
	}
	return util.ASR(val, uint(is))
}

func (g *GBA) armROR(val uint32, is uint32, carryMut bool, imm bool) uint32 {
	if is == 0 && imm {
		c := g.Carry()
		g.SetCPSRFlag(flagC, util.Bit(val, 0))
		return util.ROR(((val & ^(uint32(1))) | c), 1)
	}
	carry := (val>>(is-1))&0b1 > 0
	if is > 0 && carryMut {
		g.SetCPSRFlag(flagC, carry)
	}
	return util.ROR(val, uint(is))
}
