package util

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

// LSL logical left shift
func LSL(val uint32, shiftAmount uint) uint32 { return val << shiftAmount }

// LSR logical right shift
func LSR(val uint32, shiftAmount uint) uint32 { return val >> shiftAmount }

// ASR arithmetic right shift
func ASR(val uint32, shiftAmount uint) uint32 {
	msb := val & 0x8000_0000
	for i := uint(0); i < shiftAmount; i++ {
		val = (val >> 1) | msb
	}
	return val
}

// ROR rotate val's bit by shift
func ROR(val uint32, shiftAmount uint) uint32 {
	shiftAmount %= 32
	tmp0 := (val) >> (shiftAmount)        // XX00YY -> 00XX00
	tmp1 := (val) << (32 - (shiftAmount)) // XX00YY -> YY0000
	return tmp0 | tmp1                    // YYXX00
}
