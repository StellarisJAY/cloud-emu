package apu

// Copyright (c) 2021 Akihiro Otomo https://github.com/akatsuki105/magia
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is furnished to do
// so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
// CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
const base = -0x60

// Sound IO
const (
	SOUND1CNT_L, SOUND1CNT_H, SOUND1CNT_X = base + 0x60, base + 0x62, base + 0x64
	SOUND2CNT_L, SOUND2CNT_H              = base + 0x68, base + 0x6c
	SOUND3CNT_L, SOUND3CNT_H, SOUND3CNT_X = base + 0x70, base + 0x72, base + 0x74
	SOUND4CNT_L, SOUND4CNT_H              = base + 0x78, base + 0x7c
	SOUNDCNT_L, SOUNDCNT_H, SOUNDCNT_X    = base + 0x80, base + 0x82, base + 0x84
	SOUNDBIAS                             = base + 0x88
	WAVE_RAM                              = base + 0x90
	FIFO_A, FIFO_B                        = base + 0xa0, base + 0xa4
)
