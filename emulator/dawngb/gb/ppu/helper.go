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

type rgb555 = uint16 // Litte Endian (0b0BBBBBGGGGGRRRRR); e.g. 0x6180 is RGB555(0, 12, 24)

var dmgPalette = [4]rgb555{
	0x7FFF, // (31, 31, 31)
	0x56B5, // (15, 15, 15)
	0x294A, // (10, 10, 10)
	0x0000, // (0, 0, 0)
}

// DMGのゲームをCGBで起動した場合のデフォルトのパレット
var cgbPalette = [8]rgb555{
	// BG (Green)
	0x7FFF, // (31, 31, 31)
	0x1BEF, // (15, 31, 6)
	0x6180, // (0, 12, 24)
	0x0000, // (0, 0, 0)

	// OBJ (Red)
	0x7FFF, // (31, 31, 31)
	0x421F, // (31, 15, 15)
	0x1CF2, // (18, 7, 7)
	0x0000, // (0, 0, 0)
}
