package dummy

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"time"
)

type Frame struct {
	img      *image.YCbCr
	duration time.Duration
}

func parseGif(file []byte) ([]*Frame, int, int, error) {
	gifData, err := gif.DecodeAll(bytes.NewReader(file))
	if err != nil {
		return nil, 0, 0, err
	}
	frames := make([]*Frame, len(gifData.Image))
	width, height := 0, 0
	for _, img := range gifData.Image {
		width = max(width, img.Bounds().Dx())
		height = max(height, img.Bounds().Dy())
	}

	for i, img := range gifData.Image {
		frames[i] = &Frame{
			img:      gifFrameToYCbCrImage(img, width, height),
			duration: time.Duration(gifData.Delay[i]) * 10 * time.Millisecond,
		}
	}

	return frames, width, height, nil
}

func gifFrameToYCbCrImage(g *image.Paletted, w, h int) *image.YCbCr {
	yuv := image.NewYCbCr(image.Rect(0, 0, w, h), image.YCbCrSubsampleRatio420)
	bY, bcb, bcr := color.RGBToYCbCr(0, 0, 0)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			yOff := yuv.YOffset(x, y)
			cOff := yuv.COffset(x, y)
			if yOff < len(yuv.Y) && cOff < len(yuv.Cb) && cOff < len(yuv.Cr) {
				yuv.Y[yOff] = bY
				yuv.Cb[cOff] = bcb
				yuv.Cr[cOff] = bcr
			}
		}
	}

	for y := g.Bounds().Min.Y; y < g.Bounds().Max.Y; y++ {
		for x := g.Bounds().Min.X; x < g.Bounds().Max.X; x++ {
			r, g, b, _ := g.At(x, y).RGBA()
			Y, cb, cr := color.RGBToYCbCr(uint8(r), uint8(g), uint8(b))
			yOff := yuv.YOffset(x, y)
			cOff := yuv.COffset(x, y)
			if yOff < len(yuv.Y) && cOff < len(yuv.Cb) && cOff < len(yuv.Cr) {
				yuv.Y[yOff] = Y
				yuv.Cb[cOff] = cb
				yuv.Cr[cOff] = cr
			}
		}
	}
	return yuv
}

func (d *Frame) Width() int {
	return d.img.Bounds().Dx()
}

func (d *Frame) Height() int {
	return d.img.Bounds().Dy()
}

func (d *Frame) YCbCr() *image.YCbCr {
	return d.img
}

func (d *Frame) Read() (image.Image, func(), error) {
	return d.img, func() {}, nil
}
