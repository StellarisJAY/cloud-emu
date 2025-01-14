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

type Mapper7 struct {
	*Cartridge
	prgBank int
}

func NewMapper7(cartridge *Cartridge) Mapper {
	return &Mapper7{cartridge, 0}
}

func (m *Mapper7) Save(encoder *gob.Encoder) error {
	encoder.Encode(m.prgBank)
	return nil
}

func (m *Mapper7) Load(decoder *gob.Decoder) error {
	decoder.Decode(&m.prgBank)
	return nil
}

func (m *Mapper7) Step() {
}

func (m *Mapper7) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		return m.CHR[address]
	case address >= 0x8000:
		index := m.prgBank*0x8000 + int(address-0x8000)
		return m.PRG[index]
	case address >= 0x6000:
		index := int(address) - 0x6000
		return m.SRAM[index]
	default:
		log.Fatalf("unhandled mapper7 read at address: 0x%04X", address)
	}
	return 0
}

func (m *Mapper7) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		m.CHR[address] = value
	case address >= 0x8000:
		m.prgBank = int(value & 7)
		switch value & 0x10 {
		case 0x00:
			m.Cartridge.Mirror = MirrorSingle0
		case 0x10:
			m.Cartridge.Mirror = MirrorSingle1
		}
	case address >= 0x6000:
		index := int(address) - 0x6000
		m.SRAM[index] = value
	default:
		log.Fatalf("unhandled mapper7 write at address: 0x%04X", address)
	}
}
