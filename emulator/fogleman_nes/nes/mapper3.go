package nes

//Copyright (C) 2015 Michael Fogleman https://github.com/fogleman/nes
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
import (
	"encoding/gob"
	"log"
)

type Mapper3 struct {
	*Cartridge
	chrBank  int
	prgBank1 int
	prgBank2 int
}

func NewMapper3(cartridge *Cartridge) Mapper {
	prgBanks := len(cartridge.PRG) / 0x4000
	return &Mapper3{cartridge, 0, 0, prgBanks - 1}
}

func (m *Mapper3) Save(encoder *gob.Encoder) error {
	encoder.Encode(m.chrBank)
	encoder.Encode(m.prgBank1)
	encoder.Encode(m.prgBank2)
	return nil
}

func (m *Mapper3) Load(decoder *gob.Decoder) error {
	decoder.Decode(&m.chrBank)
	decoder.Decode(&m.prgBank1)
	decoder.Decode(&m.prgBank2)
	return nil
}

func (m *Mapper3) Step() {
}

func (m *Mapper3) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		index := m.chrBank*0x2000 + int(address)
		return m.CHR[index]
	case address >= 0xC000:
		index := m.prgBank2*0x4000 + int(address-0xC000)
		return m.PRG[index]
	case address >= 0x8000:
		index := m.prgBank1*0x4000 + int(address-0x8000)
		return m.PRG[index]
	case address >= 0x6000:
		index := int(address) - 0x6000
		return m.SRAM[index]
	default:
		log.Fatalf("unhandled mapper3 read at address: 0x%04X", address)
	}
	return 0
}

func (m *Mapper3) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		index := m.chrBank*0x2000 + int(address)
		m.CHR[index] = value
	case address >= 0x8000:
		m.chrBank = int(value & 3)
	case address >= 0x6000:
		index := int(address) - 0x6000
		m.SRAM[index] = value
	default:
		log.Fatalf("unhandled mapper3 write at address: 0x%04X", address)
	}
}
