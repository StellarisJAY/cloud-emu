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
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const iNESFileMagic = 0x1a53454e

type iNESFileHeader struct {
	Magic    uint32  // iNES magic number
	NumPRG   byte    // number of PRG-ROM banks (16KB each)
	NumCHR   byte    // number of CHR-ROM banks (8KB each)
	Control1 byte    // control bits
	Control2 byte    // control bits
	NumRAM   byte    // PRG-RAM size (x 8KB)
	_        [7]byte // unused padding
}

// LoadNESFile reads an iNES file (.nes) and returns a Cartridge on success.
// http://wiki.nesdev.com/w/index.php/INES
// http://nesdev.com/NESDoc.pdf (page 28)
func LoadNESFile(gameData []byte) (*Cartridge, error) {
	file := bytes.NewReader(gameData)
	// read file header
	header := iNESFileHeader{}
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	// verify header magic number
	if header.Magic != iNESFileMagic {
		return nil, errors.New("invalid .nes file")
	}

	// mapper type
	mapper1 := header.Control1 >> 4
	mapper2 := header.Control2 >> 4
	mapper := mapper1 | mapper2<<4

	// mirroring type
	mirror1 := header.Control1 & 1
	mirror2 := (header.Control1 >> 3) & 1
	mirror := mirror1 | mirror2<<1

	// battery-backed RAM
	battery := (header.Control1 >> 1) & 1

	// read trainer if present (unused)
	if header.Control1&4 == 4 {
		trainer := make([]byte, 512)
		if _, err := io.ReadFull(file, trainer); err != nil {
			return nil, err
		}
	}

	// read prg-rom bank(s)
	prg := make([]byte, int(header.NumPRG)*16384)
	if _, err := io.ReadFull(file, prg); err != nil {
		return nil, err
	}

	// read chr-rom bank(s)
	chr := make([]byte, int(header.NumCHR)*8192)
	if _, err := io.ReadFull(file, chr); err != nil {
		return nil, err
	}

	// provide chr-rom/ram if not in file
	if header.NumCHR == 0 {
		chr = make([]byte, 8192)
	}

	// success
	return NewCartridge(prg, chr, mapper, mirror, battery), nil
}
