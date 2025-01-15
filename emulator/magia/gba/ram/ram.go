package ram

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
	_ "embed"
	"github.com/StellrisJAY/cloud-emu/emulator/magia/util"
)

const (
	_       = iota
	kb uint = 1 << (10 * iota)
	mb
	gb
)

//go:embed bios.gba
var sBIOS []byte

// RAM struct
type RAM struct {
	BIOS  [16 * kb]byte
	EWRAM [256 * kb]byte
	IWRAM [32 * kb]byte
	IO    [2 * kb]byte
	GamePak
	ROMSize int
}

func New(src []byte) *RAM {
	bios := [16 * kb]byte{}
	for i, b := range sBIOS {
		bios[i] = b
	}

	gamePak0 := [32 * mb]byte{}
	for i, b := range src {
		gamePak0[i] = b
	}

	return &RAM{
		BIOS: bios,
		GamePak: GamePak{
			GamePak0: gamePak0,
		},
		ROMSize: len(src),
	}
}

func (r *RAM) Get(addr uint32) uint32 {
	switch {
	case BIOS(addr):
		offset := BIOSOffset(addr)
		if offset > 0x3FFF {
			return 0
		}
		return util.LE32(r.BIOS[offset:])
	case EWRAM(addr):
		offset := EWRAMOffset(addr)
		return util.LE32(r.EWRAM[offset:])
	case IWRAM(addr):
		offset := IWRAMOffset(addr)
		return util.LE32(r.IWRAM[offset:])
	case IO(addr):
		offset := IOOffset(addr)
		if offset > 0x3fe {
			return 0
		}
		return util.LE32(r.IO[offset:])
	case GamePak0(addr):
		offset := GamePak0Offset(addr)
		return util.LE32(r.GamePak0[offset:])
	case GamePak1(addr):
		offset := GamePak1Offset(addr)
		return util.LE32(r.GamePak0[offset:])
	case GamePak2(addr):
		offset := GamePak2Offset(addr)
		return util.LE32(r.GamePak0[offset:])
	case SRAM(addr):
		return uint32(r.FlashRead(addr))
	}
	return 0
}

// Set8 sets byte into addr
func (r *RAM) Set8(addr uint32, b byte) {
	switch {
	case BIOS(addr): // write only
		return
	case EWRAM(addr):
		r.EWRAM[EWRAMOffset(addr)] = b
	case IWRAM(addr):
		r.IWRAM[IWRAMOffset(addr)] = b
	case IO(addr):
		offset := IOOffset(addr)
		if offset > 0x3fe {
			return
		}
		r.IO[offset] = b
	case GamePak0(addr):
		return
	case GamePak1(addr):
		return
	case GamePak2(addr):
		return
	case SRAM(addr):
		r.FlashWrite(addr, b)
	}
}

var busWidthMap = map[uint32]int{0x0: 32, 0x3: 32, 0x4: 32, 0x7: 32, 0x2: 16, 0x5: 16, 0x6: 16, 0x8: 16, 0x9: 16, 0xa: 16, 0xb: 16, 0xc: 16, 0xd: 16, 0xe: 8, 0xf: 8}

func BusWidth(addr uint32) int { return busWidthMap[addr>>24] }
