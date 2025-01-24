package opus

import (
	"gopkg.in/hraban/opus.v2"
	"log"
)

type Encoder struct {
	enc *opus.Encoder

	SampleRate int
	Channels   int
}

func NewEncoder(sampleRate int) (*Encoder, error) {
	enc, err := opus.NewEncoder(sampleRate, 1, opus.AppVoIP)
	if err != nil {
		log.Println("new encoder error: ", err)
		return nil, err
	}
	return &Encoder{
		enc:        enc,
		SampleRate: sampleRate,
		Channels:   1,
	}, nil
}

func (e *Encoder) Encode(pcm []float32) ([]byte, error) {
	buffer := make([]byte, 1024)
	n, err := e.enc.EncodeFloat32(pcm, buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}
