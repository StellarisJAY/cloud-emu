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

type FlashMode uint32

const (
	Idle        FlashMode = 0x00
	EraseEntire FlashMode = 0x10
	Erase       FlashMode = 0x80
	EnterID     FlashMode = 0x90
	Write       FlashMode = 0xa0
	BankSwitch  FlashMode = 0xb0
	TerminateID FlashMode = 0xf0
)

type GamePak struct {
	GamePak0, GamePak1, GamePak2 [32 * mb]byte
	SRAM                         [64 * kb]byte
	Flash                        [128 * kb]byte // multiple banks
	flashMode                    FlashMode
	flashBank                    uint32
	flashIDMode                  bool
	HasFlash                     bool
}

func (r *RAM) FlashRead(addr uint32) byte {
	if r.flashIDMode { // This is the Flash ROM ID, we return Sanyo ID code
		switch addr {
		case 0x0e000000:
			return 0x62
		case 0x0e000001:
			return 0x13
		}
	} else if r.HasFlash {
		return r.Flash[r.flashBank|(addr&0xffff)]
	} else {
		return r.SRAM[addr&0xffff]
	}
	return 0
}

func (r *RAM) FlashWrite(addr uint32, b byte) {
	switch {
	case r.flashMode == Write:
		r.Flash[r.flashBank|(addr&0xffff)] = b
		r.flashMode = Idle
	case r.flashMode == BankSwitch && addr == 0x0e000000:
		r.flashBank = uint32(b&1) << 16
		r.flashMode = Idle
	case r.SRAM[0x5555] == 0xaa && r.SRAM[0x2aaa] == 0x55:
		if addr == 0xe005555 { // Command for Flash ROM
			switch FlashMode(b) {
			case EraseEntire:
				if r.flashMode == Erase {
					for idx := 0; idx < 0x20000; idx++ {
						r.Flash[idx] = 0xff
					}
					r.flashMode = Idle
				}
			case Erase:
				r.flashMode = Erase
			case EnterID:
				r.flashIDMode = true
			case Write:
				r.flashMode = Write
			case BankSwitch:
				r.flashMode = BankSwitch
			case TerminateID:
				r.flashIDMode = false
			}

			if r.flashMode != Idle || r.flashIDMode {
				r.HasFlash = true
			}
		} else if r.flashMode == Erase && b == 0x30 {
			bankStart := addr & 0xf000
			bankEnd := bankStart + 0x1000
			for idx := bankStart; idx < bankEnd; idx++ {
				r.Flash[r.flashBank|idx] = 0xff
			}
			r.flashMode = Idle
		}
	}
	r.SRAM[addr&0xffff] = b
}
