package cart

//MIT License
//
//Copyright (c) 2017 Humphrey Shotton https://github.com/Humpheh/goboy
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

// NewMBC2 returns a new MBC2 memory controller.
func NewMBC2(data []byte) BankingController {
	return &MBC2{
		rom:     data,
		romBank: 1,
		ram:     make([]byte, 0x2000),
	}
}

// MBC2 is a basic Gameboy cartridge.
type MBC2 struct {
	rom     []byte
	romBank uint32

	ram        []byte
	ramEnabled bool
}

// Read returns a value at a memory address in the ROM or RAM.
func (r *MBC2) Read(address uint16) byte {
	switch {
	case address < 0x4000:
		return r.rom[address] // Bank 0 is fixed
	case address < 0x8000:
		return r.rom[uint32(address-0x4000)+(r.romBank*0x4000)] // Use selected rom bank
	default:
		return r.ram[address-0xA000] // Use ram
	}
}

// WriteROM attempts to switch the ROM or RAM bank.
func (r *MBC2) WriteROM(address uint16, value byte) {
	switch {
	case address < 0x2000:
		// RAM enable
		if address&0x100 == 0 {
			if value&0xF == 0xA {
				r.ramEnabled = true
			} else if value&0xF == 0x0 {
				r.ramEnabled = false
			}
		}
		return
	case address < 0x4000:
		// ROM bank number (lower 4)
		if address&0x100 == 0x100 {
			r.romBank = uint32(value & 0xF)
			if r.romBank == 0x00 || r.romBank == 0x20 || r.romBank == 0x40 || r.romBank == 0x60 {
				r.romBank++
			}
		}
		return
	}
}

// WriteRAM writes data to the ram if it is enabled.
func (r *MBC2) WriteRAM(address uint16, value byte) {
	if r.ramEnabled {
		r.ram[address-0xA000] = value & 0xF
	}
}

// GetSaveData returns the save data for this banking controller.
func (r *MBC2) GetSaveData() []byte {
	data := make([]byte, len(r.ram))
	copy(data, r.ram)
	return data
}

// LoadSaveData loads the save data into the cartridge.
func (r *MBC2) LoadSaveData(data []byte) {
	r.ram = data
}
