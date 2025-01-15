package ram

//Copyright (c) 2021 Akihiro Otomo https://github.com/akatsuki105/magia
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
//of the Software, and to permit persons to whom the Software is furnished to do
//so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
//INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
//PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
//HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
//CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
//OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

func BIOS(addr uint32) bool         { return (addr >> 24) == 0x0 }
func BIOSOffset(addr uint32) uint32 { return addr }

func EWRAM(addr uint32) bool         { return (addr >> 24) == 0x2 }
func EWRAMOffset(addr uint32) uint32 { return (addr - 0x0200_0000) % 0x40000 }

func IWRAM(addr uint32) bool         { return (addr >> 24) == 0x3 }
func IWRAMOffset(addr uint32) uint32 { return (addr - 0x0300_0000) % 0x8000 }

func IO(addr uint32) bool         { return (addr >> 24) == 0x4 }
func IOOffset(addr uint32) uint32 { return addr - 0x0400_0000 }

func Palette(addr uint32) bool         { return (addr >> 24) == 0x5 }
func PaletteOffset(addr uint32) uint32 { return (addr - 0x0500_0000) % 0x400 }

func VRAM(addr uint32) bool { return (addr >> 24) == 0x6 }
func VRAMOffset(addr uint32) uint32 {
	offset := addr - 0x0600_0000
	switch {
	case offset < 0x18000:
		return offset
	case offset < 0x20000:
		return offset - 0x8000
	default:
		offset = offset % 0x20000
		if offset >= 0x18000 && offset < 0x20000 {
			return offset - 0x8000
		}
		return offset
	}
}

func OAM(addr uint32) bool         { return (addr >> 24) == 0x7 }
func OAMOffset(addr uint32) uint32 { return (addr - 0x0700_0000) % 0x400 }

func GamePak0(addr uint32) bool         { return 0x0800_0000 <= addr && addr < 0x0a00_0000 }
func GamePak0Offset(addr uint32) uint32 { return addr - 0x0800_0000 }

func GamePak1(addr uint32) bool         { return 0x0a00_0000 <= addr && addr < 0x0c00_0000 }
func GamePak1Offset(addr uint32) uint32 { return addr - 0x0a00_0000 }

func GamePak2(addr uint32) bool         { return 0x0c00_0000 <= addr && addr < 0x0e00_0000 }
func GamePak2Offset(addr uint32) uint32 { return addr - 0x0c00_0000 }

func SRAM(addr uint32) bool         { return 0x0e00_0000 <= addr && addr < 0xffff_ffff }
func SRAMOffset(addr uint32) uint32 { return (addr - 0x0e00_0000) % 0x10000 }
