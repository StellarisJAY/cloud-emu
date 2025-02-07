package psg

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
	"encoding/binary"
	"io"
)

type noise struct {
	enabled bool // NR52.3

	length uint8 // 音の残り再生時間
	stop   bool  // NR44.6

	envelope *envelope

	lfsr uint16 // Noiseの疑似乱数(lfsr: Linear Feedback Shift Register = 疑似乱数生成アルゴリズム)

	// この2つでノイズの周波数(疑似乱数の生成頻度)を決める
	divisor uint8 // NR43.0-2; ノイズ周波数1(カウント指定)
	octave  uint8 // NR43.4-7; ノイズ周波数2(オクターブ指定)
	period  uint32

	narrow bool // NR43.3; 0: 15bit, 1: 7bit

	output uint8 // 0..15
}

func newNoiseChannel() *noise {
	return &noise{
		envelope: newEnvelope(),
		lfsr:     0,
	}
}

func (ch *noise) reset() {
	ch.enabled = false
	ch.length, ch.stop = 0, false
	ch.envelope.reset()
	ch.lfsr = 0
	ch.divisor, ch.octave = 0, 0
	ch.period = 0
	ch.narrow = false
	ch.output = 0
}

func (ch *noise) reload() {
	ch.enabled = ch.dacEnable()
	ch.envelope.reload()
	ch.lfsr = 0x7FFF
	if ch.length == 0 {
		ch.length = 64
	}
}

func (ch *noise) clock64Hz() {
	if ch.enabled {
		ch.envelope.update()
	}
}

func (ch *noise) clock256Hz() {
	if ch.stop && ch.length > 0 {
		ch.length--
		if ch.length == 0 {
			ch.enabled = false
		}
	}
}

func (ch *noise) clockTimer() {
	// ch.enabledに関わらず、乱数は生成される
	ch.period--
	if ch.period == 0 {
		ch.period = ch.calcFreqency()
		ch.update()
	}

	result := uint8(0)
	if (ch.lfsr & 1) == 0 {
		result = ch.envelope.volume
	}

	if !ch.enabled {
		result = 0
	}

	ch.output = result
}

func (ch *noise) update() {
	if ch.octave < 14 {
		bit := ((ch.lfsr ^ (ch.lfsr >> 1)) & 1)
		if ch.narrow {
			ch.lfsr = (ch.lfsr >> 1) ^ (bit << 6)
		} else {
			ch.lfsr = (ch.lfsr >> 1) ^ (bit << 14)
		}
	}
}

var noisePeriodTable = []uint8{4, 8, 16, 24, 32, 40, 48, 56}

func (ch *noise) calcFreqency() uint32 {
	return uint32(noisePeriodTable[ch.divisor]) << ch.octave
}

func (ch *noise) getOutput() uint8 {
	if ch.enabled {
		return ch.output
	}
	return 0
}

func (ch *noise) dacEnable() bool {
	return ((ch.envelope.initialVolume != 0) || ch.envelope.direction)
}

func (ch *noise) serialize(s io.Writer) {
	binary.Write(s, binary.LittleEndian, ch.enabled)
	binary.Write(s, binary.LittleEndian, ch.length)
	binary.Write(s, binary.LittleEndian, ch.stop)
	ch.envelope.serialize(s)
	binary.Write(s, binary.LittleEndian, ch.lfsr)
	binary.Write(s, binary.LittleEndian, ch.octave)
	binary.Write(s, binary.LittleEndian, ch.divisor)
	binary.Write(s, binary.LittleEndian, ch.period)
	binary.Write(s, binary.LittleEndian, ch.narrow)
}

func (ch *noise) deserialize(s io.Reader) {
	binary.Read(s, binary.LittleEndian, &ch.enabled)
	binary.Read(s, binary.LittleEndian, &ch.length)
	binary.Read(s, binary.LittleEndian, &ch.stop)
	ch.envelope.deserialize(s)
	binary.Read(s, binary.LittleEndian, &ch.lfsr)
	binary.Read(s, binary.LittleEndian, &ch.octave)
	binary.Read(s, binary.LittleEndian, &ch.divisor)
	binary.Read(s, binary.LittleEndian, &ch.period)
	binary.Read(s, binary.LittleEndian, &ch.narrow)
}
