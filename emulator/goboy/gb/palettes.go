package gb

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
import (
	"fmt"
	"github.com/StellrisJAY/cloud-emu/emulator/goboy/bits"
)

const (
	// PaletteGreyscale is the default greyscale gameboy colour palette.
	PaletteGreyscale = byte(iota)
	// PaletteOriginal is more authentic looking green tinted gameboy
	// colour palette  as it would have been on the GameBoy
	PaletteOriginal
	// PaletteBGB used by default in the BGB emulator.
	PaletteBGB
)

// CurrentPalette is the global current DMG palette.
var CurrentPalette = PaletteBGB

// Palettes is an mapping from colour palettes to their colour values
// to be used by the emulator.
var Palettes = [][][]byte{
	// PaletteGreyscale
	{
		{0xFF, 0xFF, 0xFF},
		{0xCC, 0xCC, 0xCC},
		{0x77, 0x77, 0x77},
		{0x00, 0x00, 0x00},
	},
	// PaletteOriginal
	{
		{0x9B, 0xBC, 0x0F},
		{0x8B, 0xAC, 0x0F},
		{0x30, 0x62, 0x30},
		{0x0F, 0x38, 0x0F},
	},
	// PaletteBGB
	{
		{0xE0, 0xF8, 0xD0},
		{0x88, 0xC0, 0x70},
		{0x34, 0x68, 0x56},
		{0x08, 0x18, 0x20},
	},
}

// GetPaletteColour returns the colour based on the colour index and the currently
// selected palette.
func GetPaletteColour(index byte) (uint8, uint8, uint8) {
	col := Palettes[CurrentPalette][index]
	return col[0], col[1], col[2]
}

// NewPalette makes a new CGB colour palette.
func NewPalette() *cgbPalette {
	pal := make([]byte, 0x40)
	for i := range pal {
		pal[i] = 0xFF
	}

	return &cgbPalette{
		Palette: pal,
	}
}

// Palette for cgb containing information tracking the palette colour info.
type cgbPalette struct {
	// Palette colour information.
	Palette []byte
	// Current index the palette is referencing.
	Index byte
	// If to auto increment on write.
	Inc bool
}

// Update the index the palette is indexing and set
// auto increment if bit 7 is set.
func (pal *cgbPalette) updateIndex(value byte) {
	pal.Index = value & 0x3F
	pal.Inc = bits.Test(value, 7)
}

// Read the palette information stored at the current
// index.
func (pal *cgbPalette) read() byte {
	return pal.Palette[pal.Index]
}

// Read the current index.
func (pal *cgbPalette) readIndex() byte {
	return pal.Index
}

// Write a value to the palette at the current index.
func (pal *cgbPalette) write(value byte) {
	pal.Palette[pal.Index] = value
	if pal.Inc {
		pal.Index = (pal.Index + 1) & 0x3F
	}
}

// Get the rgb colour for a palette at a colour number.
func (pal *cgbPalette) get(palette byte, num byte) (uint8, uint8, uint8) {
	idx := (palette * 8) + (num * 2)
	colour := uint16(pal.Palette[idx]) | uint16(pal.Palette[idx+1])<<8
	r := uint8(colour & 0x1F)
	g := uint8((colour >> 5) & 0x1F)
	b := uint8((colour >> 10) & 0x1F)
	return colArr[r], colArr[g], colArr[b]
}

// String output of the palette values.
func (pal *cgbPalette) String() string {
	out := ""
	for i := 0; i < len(pal.Palette); i += 2 {
		out += fmt.Sprintf("%02x%02x ", pal.Palette[i+1], pal.Palette[i])
		if (i+2)%8 == 0 {
			out += "\n"
		}
	}
	return out
}

// Mapping of the 5 bit colour value to a 8 bit value.
var colArr = []uint8{
	0x0,
	0x8,
	0x10,
	0x18,
	0x20,
	0x29,
	0x31,
	0x39,
	0x41,
	0x4a,
	0x52,
	0x5a,
	0x62,
	0x6a,
	0x73,
	0x7b,
	0x83,
	0x8b,
	0x94,
	0x9c,
	0xa4,
	0xac,
	0xb4,
	0xbd,
	0xc5,
	0xcd,
	0xd5,
	0xde,
	0xe6,
	0xee,
	0xf6,
	0xff,
}
