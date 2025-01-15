package sm83

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
import "github.com/StellrisJAY/cloud-emu/emulator/dawngb/gb/util"

type pair struct {
	Lo, Hi uint8
}

func (p *pair) pack() uint16 {
	return uint16(p.Hi)<<8 | uint16(p.Lo)
}

func (p *pair) Unpack(val uint16) {
	p.Lo = uint8(val)
	p.Hi = uint8(val >> 8)
}

type Registers struct {
	A          uint8
	F          psr
	BC, DE, HL pair
	SP, PC     uint16
}

func (r *Registers) reset() {
	r.A = 0x00
	r.F.Unpack(0x00)
	r.BC.Unpack(0x0000)
	r.DE.Unpack(0x0000)
	r.HL.Unpack(0x0000)
	r.SP, r.PC = 0x0000, 0x0000
}

// ZNHC----
type psr struct {
	z, n, h, c bool
}

func (p *psr) pack() uint8 {
	packed := uint8(0)
	packed = util.SetBit(packed, 7, p.z)
	packed = util.SetBit(packed, 6, p.n)
	packed = util.SetBit(packed, 5, p.h)
	packed = util.SetBit(packed, 4, p.c)
	return packed
}

func (p *psr) Unpack(val uint8) {
	p.z = util.Bit(val, 7)
	p.n = util.Bit(val, 6)
	p.h = util.Bit(val, 5)
	p.c = util.Bit(val, 4)
}
