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

// NewMBC1 returns a new MBC1 memory controller.
func NewMBC1(data []byte) BankingController {
	return &MBC1{
		rom:     data,
		romBank: 1,
		ram:     make([]byte, 0x8000),
	}
}

// MBC1 is a GameBoy cartridge that supports rom and ram banking.
type MBC1 struct {
	rom     []byte
	romBank uint32

	ram        []byte
	ramBank    uint32
	ramEnabled bool

	romBanking bool
}

// Read returns a value at a memory address in the ROM or RAM.
func (r *MBC1) Read(address uint16) byte {
	switch {
	case address < 0x4000:
		return r.rom[address] // Bank 0 is fixed
	case address < 0x8000:
		return r.rom[uint32(address-0x4000)+(r.romBank*0x4000)] // Use selected rom bank
	default:
		return r.ram[(0x2000*r.ramBank)+uint32(address-0xA000)] // Use selected ram bank
	}
}

// WriteROM attempts to switch the ROM or RAM bank.
func (r *MBC1) WriteROM(address uint16, value byte) {
	switch {
	case address < 0x2000:
		// RAM enable
		if value&0xF == 0xA {
			r.ramEnabled = true
		} else if value&0xF == 0x0 {
			r.ramEnabled = false
		}
	case address < 0x4000:
		// ROM bank number (lower 5)
		r.romBank = (r.romBank & 0xe0) | uint32(value&0x1f)
		r.updateRomBankIfZero()
	case address < 0x6000:
		// ROM/RAM banking
		if r.romBanking {
			r.romBank = (r.romBank & 0x1F) | uint32(value&0xe0)
			r.updateRomBankIfZero()
		} else {
			r.ramBank = uint32(value & 0x3)
		}
	case address < 0x8000:
		// ROM/RAM select mode
		r.romBanking = value&0x1 == 0x00
		if r.romBanking {
			r.ramBank = 0
		} else {
			r.romBank = r.romBank & 0x1F
		}
	}
}

// Update the romBank if it is on a value which cannot be used.
func (r *MBC1) updateRomBankIfZero() {
	if r.romBank == 0x00 || r.romBank == 0x20 || r.romBank == 0x40 || r.romBank == 0x60 {
		r.romBank++
	}
}

// WriteRAM writes data to the ram if it is enabled.
func (r *MBC1) WriteRAM(address uint16, value byte) {
	if r.ramEnabled {
		r.ram[(0x2000*r.ramBank)+uint32(address-0xA000)] = value
	}
}

// GetSaveData returns the save data for this banking controller.
func (r *MBC1) GetSaveData() []byte {
	data := make([]byte, len(r.ram))
	copy(data, r.ram)
	return data
}

// LoadSaveData loads the save data into the cartridge.
func (r *MBC1) LoadSaveData(data []byte) {
	r.ram = data
}
