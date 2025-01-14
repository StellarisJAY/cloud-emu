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
import "encoding/gob"

type Cartridge struct {
	PRG     []byte // PRG-ROM banks
	CHR     []byte // CHR-ROM banks
	SRAM    []byte // Save RAM
	Mapper  byte   // mapper type
	Mirror  byte   // mirroring mode
	Battery byte   // battery present
}

func NewCartridge(prg, chr []byte, mapper, mirror, battery byte) *Cartridge {
	sram := make([]byte, 0x2000)
	return &Cartridge{prg, chr, sram, mapper, mirror, battery}
}

func (cartridge *Cartridge) Save(encoder *gob.Encoder) error {
	encoder.Encode(cartridge.PRG)
	encoder.Encode(cartridge.CHR)
	encoder.Encode(cartridge.SRAM)
	encoder.Encode(cartridge.Mirror)
	return nil
}

func (cartridge *Cartridge) Load(decoder *gob.Decoder) error {
	decoder.Decode(&cartridge.PRG)
	decoder.Decode(&cartridge.CHR)
	decoder.Decode(&cartridge.SRAM)
	decoder.Decode(&cartridge.Mirror)
	return nil
}
