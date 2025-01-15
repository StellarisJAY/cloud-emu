package cart

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
import "fmt"

// Header represents GBA Cartridge header
type Header struct {
	Entry     [4]byte
	Title     string
	GameCode  string
	MakerCode string
}

// New cartridge header
func New(src []byte) *Header {
	return &Header{
		Entry:     [4]byte{src[0], src[1], src[2], src[3]},
		Title:     string(src[0xa0 : 0xa0+12]),
		GameCode:  string(src[0xac : 0xac+4]),
		MakerCode: string(src[0xb0 : 0xb0+2]),
	}
}

func (h *Header) String() string {
	str := `Title: %s
GameCode: %s
MakerCode: %s`
	return fmt.Sprintf(str, h.Title, h.GameCode, h.MakerCode)
}
