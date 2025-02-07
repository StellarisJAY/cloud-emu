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
	"math"
)

type SysCall byte

const (
	SoftReset       SysCall = 0x00
	Div             SysCall = 0x06
	DivArm          SysCall = 0x07
	Sqrt            SysCall = 0x08
	ArcTan          SysCall = 0x09
	ArcTan2         SysCall = 0x0a
	CpuSet          SysCall = 0x0b
	FastCpuSet      SysCall = 0x0c
	GetBiosChecksum SysCall = 0x0d
	BgAffineSet     SysCall = 0x0e
	ObjAffineSet    SysCall = 0x0f
	MidiKey2Freq    SysCall = 0x1f
)

func (g *GBA) resetSP() {
	g.setPrivMode(SWI)
	g.R[13] = 0x3007fe0
	g.setPrivMode(IRQ)
	g.R[13] = 0x3007fa0
	g.setPrivMode(SYS)
	g.R[13] = 0x3007f00
}

func (g *GBA) swi(nn SysCall) {
	switch nn {
	case SoftReset:
		flag := g._getRAM(0x0300_7FFA)
		for i := uint32(0x0300_7e00); i < 0x0300_8000; i += 4 {
			g._setRAM(i, 0, 4)
		}
		g.resetSP()
		g.R[14] = 0x02000000
		if flag > 0 {
			g.R[14] = 0x08000000
		}
		g.SetCPSRFlag(flagT, false)
		g.R[15] = g.R[14]
		g.pipelining()
	case Div:
		r0, r1, r3 := util.Div(int32(g.R[0]), int32(g.R[1]))
		g.R[0], g.R[1], g.R[3] = r0, r1, r3
	case DivArm:
		r0, r1, r3 := util.Div(int32(g.R[1]), int32(g.R[0]))
		g.R[0], g.R[1], g.R[3] = r0, r1, r3
	case Sqrt:
		g.R[0] = util.Sqrt(g.R[0])
	case ArcTan:
		r0, r1, r3 := util.ArcTan(int32(g.R[0]))
		g.R[0], g.R[1], g.R[3] = r0, r1, r3
	case ArcTan2:
		r0, r1 := util.ArcTan2(int32(g.R[0]), int32(g.R[1]))
		g.R[0], g.R[1], g.R[3] = r0, r1, 0x170
	case CpuSet:
		source, dest := g.R[0], g.R[1]
		mode := g.R[2]
		count := mode & 0x000fffff
		fill := util.Bit(mode, 24)
		wordsize := 2
		if util.Bit(mode, 26) {
			wordsize = 4
		}
		if fill {
			if wordsize == 4 {
				source &= 0xfffffffc
				dest &= 0xfffffffc
				word := g._getRAM(source)
				for i := uint32(0); i < count; i++ {
					g._setRAM(dest+(i<<2), word, 4)
				}
			} else {
				source &= 0xfffffffe
				dest &= 0xfffffffe
				word := uint16(g._getRAM(source))
				for i := uint32(0); i < count; i++ {
					g._setRAM(dest+(i<<1), uint32(word), 2)
				}
			}
		} else {
			if wordsize == 4 {
				source &= 0xfffffffc
				dest &= 0xfffffffc
				for i := uint32(0); i < count; i++ {
					word := g._getRAM(source + (i << 2))
					g._setRAM(dest+(i<<2), word, 4)
				}
			} else {
				source &= 0xfffffffe
				dest &= 0xfffffffe
				for i := uint32(0); i < count; i++ {
					word := uint16(g._getRAM(source + (i << 1)))
					g._setRAM(dest+(i<<1), uint32(word), 2)
				}
			}
		}

	case FastCpuSet:
		source, dest := g.R[0]&0xffff_fffc, g.R[1]&0xffff_fffc
		mode := g.R[2]
		count := ((mode&0x000f_ffff + 7) >> 3) << 3
		fill := util.Bit(mode, 24)
		if fill {
			word := g._getRAM(source)
			for i := uint32(0); i < count; i++ {
				g._setRAM(dest+(i<<2), word, 4)
			}
		} else {
			for i := uint32(0); i < count; i++ {
				word := g._getRAM(source + (i << 2))
				g._setRAM(dest+(i<<2), word, 4)
			}
		}
	case GetBiosChecksum:
		g.R[0], g.R[1], g.R[3] = 0xbaae187f, 1, 0x00004000
	case BgAffineSet:
		i := g.R[2]
		var ox, oy float64
		var cx, cy float64
		var sx, sy float64
		var theta float64
		offset, destination := g.R[0], g.R[1]
		var a, b, c, d float64
		var rx, ry float64
		for ; i > 0; i-- {
			// [ sx   0  0 ]   [ cos(theta)  -sin(theta)  0 ]   [ 1  0  cx - ox ]   [ A B rx ]
			// [  0  sy  0 ] * [ sin(theta)   cos(theta)  0 ] * [ 0  1  cy - oy ] = [ C D ry ]
			// [  0   0  1 ]   [     0            0       1 ]   [ 0  0     1    ]   [ 0 0  1 ]
			ox = float64(g._getRAM(offset)) / 256
			oy = float64(g._getRAM(offset+4)) / 256
			cx = float64(uint16(g._getRAM(offset + 8)))
			cy = float64(uint16(g._getRAM(offset + 10)))
			sx = float64(uint16(g._getRAM(offset+12))) / 256
			sy = float64(uint16(g._getRAM(offset+14))) / 256
			theta = (float64(g._getRAM(offset+16)>>8) / 128) * math.Pi
			offset += 20

			// Rotation
			a = math.Cos(theta)
			d = a
			b = math.Sin(theta)
			c = b

			// Scale
			a *= sx
			b *= -sx
			c *= sy
			d *= sy

			// Translate
			rx = ox - (a*cx + b*cy)
			ry = oy - (c*cx + d*cy)

			g._setRAM(destination, uint32(a*256), 2)
			g._setRAM(destination+2, uint32(b*256), 2)
			g._setRAM(destination+4, uint32(c*256), 2)
			g._setRAM(destination+6, uint32(d*256), 2)
			g._setRAM(destination+8, uint32(rx*256), 4)
			g._setRAM(destination+12, uint32(ry*256), 4)
			destination += 16
		}
	case ObjAffineSet:
		i := g.R[2]
		var sx, sy float64
		var theta float64
		offset := g.R[0]
		destination := g.R[1]
		diff := g.R[3]
		var a, b, c, d float64
		for ; i > 0; i-- {
			// [ sx   0 ]   [ cos(theta)  -sin(theta) ]   [ A B ]
			// [  0  sy ] * [ sin(theta)   cos(theta) ] = [ C D ]
			sx = float64(uint16(g._getRAM(offset))) / 256
			sy = float64(uint16(g._getRAM(offset+2))) / 256
			theta = (float64(uint16(g._getRAM(offset+4))>>8) / 128) * math.Pi
			offset += 6

			// Rotation
			a = math.Cos(theta)
			d = a
			b = math.Sin(theta)
			c = b

			// Scale
			a *= sx
			b *= -sx
			c *= sy
			d *= sy

			g._setRAM(destination, uint32(a*256), 2)
			g._setRAM(destination+diff, uint32(b*256), 2)
			g._setRAM(destination+diff*2, uint32(c*256), 2)
			g._setRAM(destination+diff*3, uint32(d*256), 2)
			destination += diff * 4
		}
	case MidiKey2Freq:
		key := float64(g._getRAM(g.R[0] + 4))
		g.R[0] = uint32(key / math.Pow(2, (float64(180-g.R[1]-g.R[2])/256)/12))
	default:
		g.exception(swiVec, SWI)
	}
}
