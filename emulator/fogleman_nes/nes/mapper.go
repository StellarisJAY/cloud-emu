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
	"encoding/gob"
	"fmt"
)

type Mapper interface {
	Read(address uint16) byte
	Write(address uint16, value byte)
	Step()
	Save(encoder *gob.Encoder) error
	Load(decoder *gob.Decoder) error
}

func NewMapper(console *Console) (Mapper, error) {
	cartridge := console.Cartridge
	switch cartridge.Mapper {
	case 0:
		return NewMapper2(cartridge), nil
	case 1:
		return NewMapper1(cartridge), nil
	case 2:
		return NewMapper2(cartridge), nil
	case 3:
		return NewMapper3(cartridge), nil
	case 4:
		return NewMapper4(console, cartridge), nil
	case 7:
		return NewMapper7(cartridge), nil
	case 40:
		return NewMapper40(console, cartridge), nil
	case 225:
		return NewMapper225(cartridge), nil
	}
	err := fmt.Errorf("unsupported mapper: %d", cartridge.Mapper)
	return nil, err
}
