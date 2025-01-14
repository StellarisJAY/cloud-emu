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
	"fmt"
	"log"
)

type Mapper40 struct {
	*Cartridge
	console *Console
	bank    int
	cycles  int
}

func NewMapper40(console *Console, cartridge *Cartridge) Mapper {
	return &Mapper40{cartridge, console, 0, 0}
}

func (m *Mapper40) Save(encoder *gob.Encoder) error {
	encoder.Encode(m.bank)
	encoder.Encode(m.cycles)
	return nil
}

func (m *Mapper40) Load(decoder *gob.Decoder) error {
	decoder.Decode(&m.bank)
	decoder.Decode(&m.cycles)
	return nil
}

func (m *Mapper40) Step() {
	if m.cycles < 0 {
		return
	}
	m.cycles++
	if m.cycles%(4096*3) == 0 {
		m.cycles = 0
		m.console.CPU.triggerIRQ()
	}
}

func (m *Mapper40) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		return m.CHR[address]
	case address >= 0x6000 && address < 0x8000:
		return m.PRG[address-0x6000+0x2000*6]
	case address >= 0x8000 && address < 0xa000:
		return m.PRG[address-0x8000+0x2000*4]
	case address >= 0xa000 && address < 0xc000:
		return m.PRG[address-0xa000+0x2000*5]
	case address >= 0xc000 && address < 0xe000:
		return m.PRG[address-0xc000+0x2000*uint16(m.bank)]
	case address >= 0xe000:
		return m.PRG[address-0xe000+0x2000*7]
	default:
		log.Fatalf("unhandled mapper40 read at address: 0x%04X", address)
	}
	return 0
}

func (m *Mapper40) Write(address uint16, value byte) {
	switch {
	case address < 0x2000:
		m.CHR[address] = value
	case address >= 0x8000 && address < 0xa000:
		m.cycles = -1
	case address >= 0xa000 && address < 0xc000:
		m.cycles = 0
	case address >= 0xe000:
		m.bank = int(value)
	default:
		// log.Fatalf("unhandled mapper40 write at address: 0x%04X", address)
		fmt.Printf("unhandled mapper40 write at address: 0x%04X\n", address)
	}
}
