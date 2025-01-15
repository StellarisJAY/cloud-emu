package dummy

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
	"image/color"
)

type Renderer struct{}

func New() *Renderer {
	return &Renderer{}
}

func (r *Renderer) DrawScanline(y int, scanline []color.NRGBA) {}
func (r *Renderer) SetLCDC(val uint8)                          {}
func (r *Renderer) SetBGP(val uint8)                           {}
func (r *Renderer) SetOBP0(val uint8)                          {}
func (r *Renderer) SetOBP1(val uint8)                          {}
func (r *Renderer) SetSCX(val uint8)                           {}
func (r *Renderer) SetSCY(val uint8)                           {}
func (r *Renderer) SetWX(val uint8)                            {}
func (r *Renderer) SetWY(val uint8)                            {}
