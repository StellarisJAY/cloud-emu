package cartridge

//Copyright (c) 2024 Akihiro Otomo https://github.com/akatsuki105/dawngb
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
	"github.com/StellrisJAY/cloud-emu/emulator/dawngb/gb/util"
)

const KB, MB = 1024, 1024 * 1024

var RAM_SIZES = map[uint8]uint{
	2: 8 * KB,
	3: 32 * KB,
	4: 128 * KB,
	5: 64 * KB,
}

type mbc interface {
	read(addr uint16) uint8
	write(addr uint16, val uint8)
}

type Cartridge struct {
	title string
	ROM   []uint8
	ram   []uint8 // SRAM
	mbc           // mapper
}

func New(rom []uint8) (*Cartridge, error) {
	isCGB := util.Bit(rom[0x143], 7)

	title := string(rom[0x134:0x13F])
	if !isCGB {
		title = string(rom[0x134:0x144])
	}

	c := &Cartridge{
		title: title,
		ram:   make([]uint8, 0),
	}

	c.ROM = make([]uint8, (32*KB)<<rom[0x148])
	copy(c.ROM, rom)

	ramSize, ok := RAM_SIZES[rom[0x149]] // これはSRAMチップのサイズであって、MBC2のようなMBCチップにRAMが内蔵されている場合は0になるっぽい
	if ok {
		c.ram = make([]uint8, ramSize)
	}

	mbc, err := createMBC(c)
	if err != nil {
		return nil, err
	}
	c.mbc = mbc
	// fmt.Println("MapperID:", c.ROM[0x147])
	return c, nil
}

func createMBC(c *Cartridge) (mbc, error) {
	mbcType := c.ROM[0x147]
	switch mbcType {
	case 0:
		return newMBC0(c), nil
	case 1, 2, 3:
		return newMBC1(c), nil
	case 5, 6:
		c.ram = make([]uint8, 512)
		return newMBC2(c), nil
	case 16, 19:
		return newMBC3(c), nil
	case 25, 26, 27:
		return newMBC5(c), nil
	default:
		return nil, fmt.Errorf("unsupported mbc type: 0x%02X", mbcType)
	}
}

func (c *Cartridge) Title() string {
	return c.title
}

func (c *Cartridge) Read(addr uint16) uint8 {
	return c.mbc.read(addr)
}

func (c *Cartridge) Write(addr uint16, val uint8) {
	c.mbc.write(addr, val)
}

func (c *Cartridge) CGBFlag() uint8 {
	return c.ROM[0x143]
}

func (c *Cartridge) LoadSRAM(data []uint8) error {
	copy(c.ram, data)
	return nil
}

func (c *Cartridge) SRAM() []uint8 {
	return c.ram
}
