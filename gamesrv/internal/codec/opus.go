package codec

import (
	"github.com/pion/mediadevices/pkg/codec"
	"github.com/pion/mediadevices/pkg/codec/opus"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/mediadevices/pkg/wave"
)

type AudioEncoder struct {
	params     opus.Params
	enc        codec.ReadCloser
	reader     *AudioReader
	sampleRate int
	channels   int
}
type AudioReader struct {
	chunk *wave.Float32Interleaved
}

func (a *AudioReader) Read() (audio wave.Audio, release func(), err error) {
	return a.chunk, func() {}, nil
}

func NewAudioEncoder(sampleRate int) (IAudioEncoder, error) {
	params, err := opus.NewParams()
	if err != nil {
		return nil, err
	}
	params.BitRate = 512_000_000
	params.Latency = opus.Latency20ms
	float32Chunk := wave.NewFloat32Interleaved(wave.ChunkInfo{
		Len:          sampleRate / 60,
		Channels:     1,
		SamplingRate: sampleRate,
	})
	reader := &AudioReader{float32Chunk}
	enc, err := params.BuildAudioEncoder(reader, prop.Media{
		Audio: prop.Audio{
			SampleRate:    sampleRate,
			ChannelCount:  1,
			IsInterleaved: true,
			IsBigEndian:   false,
			IsFloat:       true,
		},
	})
	if err != nil {
		return nil, err
	}
	return &AudioEncoder{params, enc, reader, sampleRate, 1}, nil
}

func (a *AudioEncoder) Encode(samples []float32) ([]byte, error) {
	a.reader.chunk.Data = samples
	a.reader.chunk.Size = wave.ChunkInfo{
		Len:          len(samples),
		Channels:     1,
		SamplingRate: a.sampleRate,
	}
	data, _, err := a.enc.Read()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (a *AudioEncoder) Close() {
	_ = a.enc.Close()
}
