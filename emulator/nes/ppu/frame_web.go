package ppu

import (
	"image"
	"image/color"
)

type Frame struct {
	ycbcr *image.YCbCr
	scale int
}

const (
	WIDTH  = 256
	HEIGHT = 240
)

func NewFrame() *Frame {
	return &Frame{
		image.NewYCbCr(image.Rect(0, 0, WIDTH, HEIGHT), image.YCbCrSubsampleRatio420),
		1,
	}
}

func (f *Frame) setPixel(x0, y0 uint32, c Color) {
	Y, cb, cr := color.RGBToYCbCr(c.R, c.G, c.B)
	dx, dy := int(x0)*f.scale, int(y0)*f.scale
	for x := 0; x < f.scale; x++ {
		for y := 0; y < f.scale; y++ {
			yOff := f.ycbcr.YOffset(x+dx, y+dy)
			cOff := f.ycbcr.COffset(x+dx, y+dy)
			if yOff < len(f.ycbcr.Y) && cOff < len(f.ycbcr.Cb) && cOff < len(f.ycbcr.Cr) {
				f.ycbcr.Y[yOff] = Y
				f.ycbcr.Cb[cOff] = cb
				f.ycbcr.Cr[cOff] = cr
			}
		}
	}
}

func (f *Frame) SetScale(scale int) {
	f.scale = scale
	f.ycbcr = image.NewYCbCr(image.Rect(0, 0, WIDTH*scale, HEIGHT*scale), image.YCbCrSubsampleRatio420)
}

func (f *Frame) Width() int {
	return WIDTH * f.scale
}

func (f *Frame) Height() int {
	return HEIGHT * f.scale
}

func (f *Frame) YCbCr() *image.YCbCr {
	return f.ycbcr
}

func (f *Frame) Read() (image.Image, func(), error) {
	return f.ycbcr, func() {}, nil
}
