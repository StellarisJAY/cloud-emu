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

// 27-26: 00
func IsArmALU(inst uint32) bool {
	return inst&0b0000_1100_0000_0000_0000_0000_0000_0000 == 0
}

// 27-24: 1011
func IsArmBL(inst uint32) bool {
	return inst&0b0000_1111_0000_0000_0000_0000_0000_0000 == 0b0000_1011_0000_0000_0000_0000_0000_0000
}

// 27-24: 1010
func IsArmB(inst uint32) bool {
	return inst&0b0000_1111_0000_0000_0000_0000_0000_0000 == 0b0000_1010_0000_0000_0000_0000_0000_0000
}

// 27-8: 0001_0010_1111_1111_1111 && 7-4: 0001
func IsArmBX(inst uint32) bool {
	return inst&0b0000_1111_1111_1111_1111_1111_1111_0000 == 0b0000_0001_0010_1111_1111_1111_0001_0000
}

// 27-24: 1111
func IsArmSWI(inst uint32) bool {
	return inst&0b0000_1111_0000_0000_0000_0000_0000_0000 == 0b0000_1111_0000_0000_0000_0000_0000_0000
}

// 27-25: 011
func IsArmUND(inst uint32) bool {
	return inst&0b0000_1110_0000_0000_0000_0000_0000_0000 == 0b0000_0110_0000_0000_0000_0000_0000_0000
}

// multiply

// 27-25: 000 & 7-4: 1001
func IsArmMPY(inst uint32) bool {
	return inst&0b0000_1110_0000_0000_0000_0000_1111_0000 == 0b0000_0000_0000_0000_0000_0000_1001_0000
}

// 27-25: 000 & 20: 0 & 7: 1 & 4: 0
func IsArmMPY16(inst uint32) bool {
	return inst&0b0000_1110_0001_0000_0000_0000_1001_0000 == 0b0000_0000_0000_0000_0000_0000_1000_0000
}

// loadstore

// 27-26: 01
func IsArmLDR(inst uint32) bool {
	return inst&0b0000_1100_0001_0000_0000_0000_0000_0000 == 0b0000_0100_0001_0000_0000_0000_0000_0000
}

// 27-26: 01
func IsArmSTR(inst uint32) bool {
	return inst&0b0000_1100_0001_0000_0000_0000_0000_0000 == 0b0000_0100_0000_0000_0000_0000_0000_0000
}

// 27-25: 000 & 20: 1 & 7-4: 1011
func IsArmLDRH(inst uint32) bool {
	return inst&0b0000_1110_0001_0000_0000_0000_1111_0000 == 0b0000_0000_0001_0000_0000_0000_1011_0000
}

// 27-25: 000 & 20: 1 & 7-4: 1101
func IsArmLDRSB(inst uint32) bool {
	return inst&0b0000_1110_0001_0000_0000_0000_1111_0000 == 0b0000_0000_0001_0000_0000_0000_1101_0000
}

// 27-25: 000 & 20: 1 & 7-4: 1111
func IsArmLDRSH(inst uint32) bool {
	return inst&0b0000_1110_0001_0000_0000_0000_1111_0000 == 0b0000_0000_0001_0000_0000_0000_1111_0000
}

// 27-25: 000 & 20: 0 & 7-4: 1011
func IsArmSTRH(inst uint32) bool {
	return inst&0b0000_1110_0001_0000_0000_0000_1111_0000 == 0b0000_0000_0000_0000_0000_0000_1011_0000
}

// 27-25: 100 & 20: 1
func IsArmLDM(inst uint32) bool {
	return inst&0b0000_1110_0001_0000_0000_0000_0000_0000 == 0b0000_1000_0001_0000_0000_0000_0000_0000
}

// 27-25: 100 & 20: 0
func IsArmSTM(inst uint32) bool {
	return inst&0b0000_1110_0001_0000_0000_0000_0000_0000 == 0b0000_1000_0000_0000_0000_0000_0000_0000
}

// 27-23: 0001_0 & 21-20: 00 & 11-4: 0000_1001
func IsArmSWP(inst uint32) bool {
	return inst&0b0000_1111_1011_0000_0000_1111_1111_0000 == 0b0000_0001_0000_0000_0000_0000_1001_0000
}

// 27-23: 0001_0 & 21-16: 00_1111 & 11-0: 0000_0000_0000
func IsArmMRS(inst uint32) bool {
	return inst&0b0000_1111_1011_1111_0000_1111_1111_1111 == 0b0000_0001_0000_1111_0000_0000_0000_0000
}

// 27-26: 00 & 24-23: 10 & 21-20: 10 & 15-12: 1111
func IsArmMSR(inst uint32) bool {
	return inst&0b0000_1101_1011_0000_1111_0000_0000_0000 == 0b0000_0001_0010_0000_1111_0000_0000_0000
}
