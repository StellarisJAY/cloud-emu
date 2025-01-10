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

// NewMBC3 returns a new MBC3 memory controller.
func NewMBC3(data []byte) BankingController {
	return &MBC3{
		rom:        data,
		romBank:    1,
		ram:        make([]byte, 0x8000),
		rtc:        make([]byte, 0x10),
		latchedRtc: make([]byte, 0x10),
	}
}

// MBC3 is a GameBoy cartridge that supports rom and ram banking and possibly
// a real time clock (RTC).
type MBC3 struct {
	rom     []byte
	romBank uint32

	ram        []byte
	ramBank    uint32
	ramEnabled bool

	rtc        []byte
	latchedRtc []byte
	latched    bool
}

// Read returns a value at a memory address in the ROM.
func (r *MBC3) Read(address uint16) byte {
	switch {
	case address < 0x4000:
		return r.rom[address] // Bank 0 is fixed
	case address < 0x8000:
		return r.rom[uint32(address-0x4000)+(r.romBank*0x4000)] // Use selected rom bank
	default:
		if r.ramBank >= 0x4 {
			if r.latched {
				return r.latchedRtc[r.ramBank]
			}
			return r.rtc[r.ramBank]
		}
		return r.ram[(0x2000*r.ramBank)+uint32(address-0xA000)] // Use selected ram bank
	}
}

// WriteROM attempts to switch the ROM or RAM bank.
func (r *MBC3) WriteROM(address uint16, value byte) {
	switch {
	case address < 0x2000:
		// RAM enable
		r.ramEnabled = (value & 0xA) != 0
	case address < 0x4000:
		// ROM bank number (lower 5)
		r.romBank = uint32(value & 0x7F)
		if r.romBank == 0x00 {
			r.romBank++
		}
	case address < 0x6000:
		r.ramBank = uint32(value)
	case address < 0x8000:
		if value == 0x1 {
			r.latched = false
		} else if value == 0x0 {
			r.latched = true
			copy(r.rtc, r.latchedRtc)
		}
	}
}

// WriteRAM writes data to the ram or RTC if it is enabled.
func (r *MBC3) WriteRAM(address uint16, value byte) {
	if r.ramEnabled {
		if r.ramBank >= 0x4 {
			r.rtc[r.ramBank] = value
		} else {
			r.ram[(0x2000*r.ramBank)+uint32(address-0xA000)] = value
		}
	}
}

// GetSaveData returns the save data for this banking controller.
func (r *MBC3) GetSaveData() []byte {
	data := make([]byte, len(r.ram))
	copy(data, r.ram)
	return data
}

// LoadSaveData loads the save data into the cartridge.
func (r *MBC3) LoadSaveData(data []byte) {
	r.ram = data
}
