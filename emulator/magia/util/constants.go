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

type Region uint32

const (
	BIOS    Region = 0x0
	EWRAM   Region = 0x2
	IWRAM   Region = 0x3
	IO      Region = 0x4
	PALETTE Region = 0x5
	VRAM    Region = 0x6
	OAM     Region = 0x7
	CART0   Region = 0x8
	CART1   Region = 0xa
	CART2   Region = 0xc
	SRAM    Region = 0xe
)

var RegionSize = map[string]uint32{
	"BIOS":     0x00004000,
	"EWRAM":    0x00040000,
	"IWRAM":    0x00008000,
	"IO":       0x00000400,
	"PALETTE":  0x00000400,
	"VRAM":     0x00018000,
	"OAM":      0x00000400,
	"CART0":    0x02000000,
	"CART1":    0x02000000,
	"CART2":    0x02000000,
	"SRAM":     0x00008000,
	"FLASH512": 0x00010000,
	"FLASH1M":  0x00020000,
	"EEPROM":   0x00002000,
}
