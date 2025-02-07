package video

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

import (
	"github.com/StellrisJAY/cloud-emu/emulator/magia/util"
)

type Scalerot [4]float64

type OAM struct {
	*MemoryAligned16
	oam      []uint16
	objs     [128]*Obj
	scalerot [32]Scalerot
	video    *SoftwareRenderer
}

func NewOAM(size uint32) *OAM {
	o := &OAM{}
	mem16 := NewMemoryAligned16(size)
	o.MemoryAligned16 = mem16
	o.oam = mem16.buffer

	o.objs = [128]*Obj{}
	for i := 0; i < 128; i++ {
		o.objs[i] = newObj(o, i)
	}

	o.scalerot = [32]Scalerot{}
	for i := 0; i < 32; i++ {
		o.scalerot[i] = Scalerot{1, 0, 0, 1}
	}

	return o
}

func (o *OAM) Store16(offset uint32, value uint16) {
	index := (offset & 0x3f8) >> 3
	obj := o.objs[index]
	scalerot := &o.scalerot[index>>2]

	switch offset & 0x6 {
	case 0: // OAM atr0
		obj.y = value & 0xff
		wasScalerot := obj.scalerot
		obj.scalerot = util.Bit(value, 8)
		if obj.scalerot {
			obj.scalerotOam = o.scalerot[obj.scalerotParam]
			obj.doublesize = util.Bit(value, 9)
			obj.disable = false
			obj.hflip = false
			obj.vflip = false
		} else {
			obj.doublesize = false
			obj.disable = util.Bit(value, 9)
			if wasScalerot {
				obj.hflip = util.Bit(obj.scalerotParam, 3)
				obj.vflip = util.Bit(obj.scalerotParam, 4)
			}
		}

		obj.mode = (value & 0x0c00) >> 6
		obj.mosaic = util.Bit(value, 12)
		obj.color256 = util.Bit(value, 13)
		obj.shape = ObjShape((value & 0xc000) >> 14)
		obj.recalcSize()

	case 2: // OAM atr1
		obj.x = value & 0x01ff
		if obj.scalerot {
			obj.scalerotParam = value & 0x3e00 >> 9
			obj.scalerotOam = o.scalerot[obj.scalerotParam]
			obj.hflip = false
			obj.vflip = false
		} else {
			obj.hflip = util.Bit(value, 12)
			obj.vflip = util.Bit(value, 13)
		}
		obj.isAffine = obj.scalerot
		obj.size = int(value&0xc000) >> 14
		obj.recalcSize()

	case 4: // OAM atr2
		obj.tileBase = uint32(value & 0x03ff)
		obj.priority = int((value & 0x0c00) >> 10)
		obj.palette = (value & 0xf000) >> 8 // This is shifted up 4 to make pushPixel faster

	case 6: // Scaling/rotation parameter
		// FIXME: using float64
		scalerot[index&3] = float64(uint32(value)<<16) / 0x0100_0000
	}

	o.MemoryAligned16.Store16(offset, value)
}

func (o *OAM) Store32(offset uint32, value uint32) {
	lower := uint16(value)
	o.Store16(offset, lower)
	upper := uint16(value >> 16)
	o.Store16(offset+2, upper)
}
