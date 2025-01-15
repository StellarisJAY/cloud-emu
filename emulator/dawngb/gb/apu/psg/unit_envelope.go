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

// 音の三要素 のうち、 音の大きさ (振幅)
type envelope struct {
	initialVolume uint8 // 初期音量(リスタート時にセット)
	volume        uint8 // 0..15
	direction     bool  // 音量変更の方向(trueで大きくなっていく)

	// 音量変更の速さ(0に近いほど速い、ただし0だと変化なし)
	// speed が n のとき、 音量変更は (n / 64) 秒 ごとに行われる
	speed uint8
	step  uint8 // 音量変更を行うタイミングをカウントするためのカウンタ
}

func newEnvelope() *envelope {
	return &envelope{
		step: 8,
	}
}

func (e *envelope) reset() {
	e.initialVolume, e.volume = 0, 0
	e.direction = false
	e.speed, e.step = 0, 8
}

func (e *envelope) reload() {
	e.volume = e.initialVolume
	e.step = e.speed
	if e.speed == 0 {
		e.step = 8
	}
}

func (e *envelope) update() {
	if e.speed != 0 {
		e.step--
		if e.step == 0 {
			e.updateVolume()
			e.step = e.speed
		}
	}
}

func (e *envelope) updateVolume() {
	if e.direction {
		if e.volume < 15 {
			e.volume++
		}
	} else {
		if e.volume > 0 {
			e.volume--
		}
	}
}

func (e *envelope) serialize(s io.Writer) {
	binary.Write(s, binary.LittleEndian, e.initialVolume)
	binary.Write(s, binary.LittleEndian, e.volume)
	binary.Write(s, binary.LittleEndian, e.direction)
	binary.Write(s, binary.LittleEndian, e.speed)
	binary.Write(s, binary.LittleEndian, e.step)
}

func (e *envelope) deserialize(d io.Reader) {
	binary.Read(d, binary.LittleEndian, &e.initialVolume)
	binary.Read(d, binary.LittleEndian, &e.volume)
	binary.Read(d, binary.LittleEndian, &e.direction)
	binary.Read(d, binary.LittleEndian, &e.speed)
	binary.Read(d, binary.LittleEndian, &e.step)
}
