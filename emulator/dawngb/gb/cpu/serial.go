package cpu

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
	"math"
)

type serial struct {
	irq    func(int)
	until  int64
	sb, sc uint8
}

func newSerial(irq func(int)) *serial {
	return &serial{irq: irq}
}

func (s *serial) reset() {
	s.until = math.MaxInt64
	s.sb, s.sc = 0, 0
}

func (s *serial) run(cycles8MHz int64) {
	for i := int64(0); i < cycles8MHz; i++ {
		s.until--
		if s.until <= 0 {
			s.until = math.MaxInt64
			s.dummyTransfer()
		}
	}
}

func (s *serial) setSC(val uint8) {
	s.sc = val
	// 一部のソフト(ポケモンクリスタル など)は起動時にシリアル通信を実装してないと動かないのでダミーで実装
	if (val & (1 << 7)) != 0 {
		s.until = math.MaxInt64
		if (s.sc & (1 << 0)) != 0 {
			s.until = 512 * 8
		}
	}
}

// ポケモンクリスタルの起動にシリアル通信機能が必要なので暫定措置
func (s *serial) dummyTransfer() {
	s.sc &= 0x7F
	s.sb = 0xFF
	s.irq(IRQ_SERIAL)
}
