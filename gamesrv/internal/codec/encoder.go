package codec

import (
	"errors"
	"github.com/StellrisJAY/cloud-emu/emulator"
	"image"
	"log"

	"github.com/pion/mediadevices/pkg/codec"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
)

type IVideoEncoder interface {
	Encode(frame emulator.IFrame) ([]byte, func(), error)
	Close()
}

type FrameReader struct {
	frame emulator.IFrame
}

func (f *FrameReader) Read() (img image.Image, release func(), err error) {
	return f.frame.Read()
}

func (f *FrameReader) SetFrame(frame emulator.IFrame) {
	f.frame = frame
}

type VideoEncoder struct {
	enc         codec.ReadCloser
	params      any
	frameReader *FrameReader
}

func NewVideoEncoder(codec string, width, height int) (IVideoEncoder, error) {
	media := prop.Media{
		DeviceID: "nesgo-video",
		Video: prop.Video{
			Width:       width,
			Height:      height,
			FrameRate:   60,
			FrameFormat: frame.FormatI420,
		},
	}
	switch codec {
	case "h264":
		return NewX264Encoder(media)
	case "vp8":
		return NewVpxEncoder(media, 8)
	case "vp9":
		return NewVpxEncoder(media, 9)
	default:
		panic(errors.New("codec not available"))
	}
}

func (v *VideoEncoder) Encode(frame emulator.IFrame) ([]byte, func(), error) {
	v.frameReader.SetFrame(frame)
	return v.enc.Read()
}

func (v *VideoEncoder) Close() {
	if err := v.enc.Close(); err != nil {
		log.Println("close video encoder error:", err)
	}
}

type IAudioEncoder interface {
	// Encode PCM to opus packet, Emulator outputs float32 PCM
	Encode(pcm []float32) ([]byte, error)
	Close()
}
