package chip8

import (
	"image"
	"image/color"
)

type Frame struct {
	img *image.YCbCr
}

func newFrame() *Frame {
	f := &Frame{
		img: image.NewYCbCr(image.Rect(0, 0, 64, 32), image.YCbCrSubsampleRatio420),
	}
	for i := range f.img.Y {
		f.img.Y[i] = blackY
	}
	for i := range f.img.Cb {
		f.img.Cb[i] = blackCB
		f.img.Cr[i] = blackCR
	}
	return f
}

var (
	whiteY, whiteCB, whiteCR = color.RGBToYCbCr(255, 255, 255)
	blackY, blackCB, blackCR = color.RGBToYCbCr(0, 0, 0)
)

func (f *Frame) Width() int {
	return 64
}

func (f *Frame) Height() int {
	return 32
}

func (f *Frame) YCbCr() *image.YCbCr {
	return f.img
}

func (f *Frame) Read() (image.Image, func(), error) {
	return f.img, func() {}, nil
}

func (f *Frame) setPixel(x, y int, white bool) {
	yOff, cOff := f.img.YOffset(x, y), f.img.COffset(x, y)
	if white {
		f.img.Y[yOff] = whiteY
		f.img.Cb[cOff] = whiteCB
		f.img.Cr[cOff] = whiteCR
	} else {
		f.img.Y[yOff] = blackY
		f.img.Cb[cOff] = blackCB
		f.img.Cr[cOff] = blackCR
	}
}

func (f *Frame) getPixel(x, y int) bool {
	yOff := f.img.YOffset(x, y)
	return f.img.Y[yOff] != 0
}

func (f *Frame) clear() {
	for i := range f.img.Y {
		f.img.Y[i] = blackY
	}
	for i := range f.img.Cb {
		f.img.Cb[i] = blackCB
		f.img.Cr[i] = blackCR
	}
}
