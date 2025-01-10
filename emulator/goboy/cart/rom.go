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

// NewROM returns a new ROM cartridge.
func NewROM(data []byte) BankingController {
	return &ROM{
		rom: data,
	}
}

// ROM is a basic Gameboy cartridge.
type ROM struct {
	rom []byte
}

// Read returns a value at a memory address in the ROM.
func (r *ROM) Read(address uint16) byte {
	return r.rom[address]
}

// WriteROM would switch between cartridge banks, however a ROM cart does
// not support banking.
func (r *ROM) WriteROM(address uint16, value byte) {}

// WriteRAM would write data to the cartridge RAM, however a ROM cart does
// not contain RAM so this is a noop.
func (r *ROM) WriteRAM(address uint16, value byte) {}

// GetSaveData returns the save data for this banking controller. As RAM is
// not supported on this memory controller, this is a noop.
func (r *ROM) GetSaveData() []byte {
	return []byte{}
}

// LoadSaveData loads the save data into the cartridge. As RAM is not supported
// on this memory controller, this is a noop.
func (r *ROM) LoadSaveData([]byte) {}
