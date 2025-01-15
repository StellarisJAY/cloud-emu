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

// Cond represents condition
type Cond byte

// condition
const (
	EQ Cond = iota // ==
	NE             // !=
	CS             // carry
	CC             // not carry
	MI             // minus
	PL             // plus
	VS             // overflow
	VC             // not-overflow
	HI             // higher
	LS             // not carry
	GE             // >=
	LT             // <
	GT             // >
	LE             // <=
	AL
	NV
)

var condLut = [16]uint16{
	0xF0F0, 0x0F0F, 0xCCCC, 0x3333, 0xFF00, 0x00FF, 0xAAAA, 0x5555, 0x0C0C, 0xF3F3, 0xAA55, 0x55AA, 0x0A05, 0xF5FA, 0xFFFF, 0x0000,
}

var cond2str = map[Cond]string{EQ: "eq", NE: "ne", CS: "cs", CC: "cc", MI: "mi", PL: "pl", VS: "vs", VC: "vc", HI: "hi", LS: "ls", GE: "ge", LT: "lt", GT: "gt", LE: "le", AL: "al", NV: "nv"}

func (c Cond) String() string {
	if s, ok := cond2str[c]; ok {
		return s
	}
	return ""
}

// Check returns if instruction condition is ok
func (g *GBA) Check(cond Cond) bool {
	flags := g.Reg.CPSR >> 28
	return (condLut[cond] & (1 << flags)) != 0
}
