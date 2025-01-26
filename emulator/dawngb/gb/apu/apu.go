package apu

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
	"github.com/StellrisJAY/cloud-emu/emulator/dawngb/gb/apu/psg"
	"io"
)

// SoCに組み込まれているため、`/cpu`にある方が正確ではある
type APU struct {
	cycles int64 // 8MHzのマスターサイクル単位
	*psg.PSG
	sampleWriter io.Writer

	samples     [547 * 2]int16 // [[left, right]...], 547 = 32768 / 60
	sampleCount uint16
}

func New(writer io.Writer) *APU {
	return &APU{
		PSG:          psg.New(psg.MODEL_GB),
		sampleWriter: writer,
	}
}

func (a *APU) Reset() {
	a.PSG.Reset()
	a.cycles = 0
	clear(a.samples[:])
	a.sampleCount = 0
}

func (a *APU) Run(cycles8MHz int64) {
	for i := int64(0); i < cycles8MHz; i++ {
		a.cycles++
		if a.cycles%2 == 0 {
			a.PSG.Step()
		}
		if a.cycles%256 == 0 { // 32768Hzにダウンサンプリングしたい = 32768Hzごとにサンプルを生成したい = 256マスターサイクルごとにサンプルを生成する (8MHz / 32768Hz = 256)
			if int(a.sampleCount) < len(a.samples)/2 {
				left, right := a.PSG.Sample()
				lvolume, rvolume := a.PSG.Volume()
				lsample, rsample := (int(left)*512)-16384, (int(right)*512)-16384
				lsample, rsample = (lsample*int(lvolume+1))/8, (rsample*int(rvolume+1))/8
				a.samples[a.sampleCount*2] = int16(lsample) / 2
				a.samples[a.sampleCount*2+1] = int16(rsample) / 2
				a.sampleCount++
			}
		}
	}
}

func (a *APU) FlushSamples() {
	_ = binary.Write(a.sampleWriter, binary.LittleEndian, a.samples[:a.sampleCount*2])
	a.sampleCount = 0
}
