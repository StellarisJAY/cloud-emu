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
	"encoding/binary"
	"github.com/StellrisJAY/cloud-emu/emulator/dawngb/gb/util"
	"io"
)

var timaClock = [4]int64{64, 1, 4, 16}

type timer struct {
	irq            func(n int)
	clock          *int64
	cycles         int64 // CPUから見て遅れているマスターサイクル数
	tima, tma, tac uint8
	counter        int64 // 524288Hz(一番細かいのが524288Hzなのであとはそれの倍数で数えれば良い)
}

func newTimer(irq func(n int), clock *int64) *timer {
	return &timer{
		irq:   irq,
		clock: clock,
	}
}

func (t *timer) reset() {
	t.cycles = 0
	t.counter, t.tima, t.tma, t.tac = 0, 0, 0, 0
}

func (t *timer) run(cycles8MHz int64) {
	t.cycles += cycles8MHz
	for t.cycles >= 16 {
		t.update()
		t.cycles -= 16
	}
}

// 524288Hz
func (t *timer) update() {
	x := (*t.clock) / 4

	t.counter++
	if util.Bit(t.tac, 2) {
		if (t.counter % (timaClock[t.tac&0b11] * x)) == 0 {
			t.tima++
			if t.tima == 0 {
				t.tima = t.tma
				t.irq(IRQ_TIMER)
			}
		}
	}
}

func (t *timer) Read(addr uint16) uint8 {
	x := (*t.clock) / 4
	switch addr {
	case 0xFF04:
		div := t.counter / (16 * x)
		return uint8(div)
	case 0xFF05:
		return t.tima
	case 0xFF06:
		return t.tma
	case 0xFF07:
		return t.tac
	}
	return 0
}

func (t *timer) Write(addr uint16, val uint8) {
	switch addr {
	case 0xFF04:
		t.counter = 0
	case 0xFF05:
		t.tima = val
	case 0xFF06:
		t.tma = val
	case 0xFF07:
		t.tac = val & 0b111
	}
}

func (t *timer) Serialize(s io.Writer) {
	data := []uint8{}
	binary.LittleEndian.PutUint64(data, uint64(t.cycles))  // 8
	data = append(data, t.tima, t.tma, t.tac)              // 3
	binary.LittleEndian.PutUint64(data, uint64(t.counter)) // 8
	s.Write(data)
}

func (t *timer) Deserialize(s io.Reader) {
	data := make([]uint8, 19)
	s.Read(data)
	t.cycles = int64(binary.LittleEndian.Uint64(data[0:8]))
	t.tima, t.tma, t.tac = data[8], data[9], data[10]
	t.counter = int64(binary.LittleEndian.Uint64(data[11:19]))
}
