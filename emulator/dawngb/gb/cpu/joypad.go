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

// P1/JOYP (0xFF00)
type joypad struct {
	irq      func(n int)
	p14, p15 bool  // P1.4, P1.5 (P14, P15 は CPUのpinの名前)
	joyp     uint8 // P1.0-3

	inputs uint8 // ゲームの実際のキー入力を反映したもの(pollで使用), (Dpad << 4) | Buttons, 0 is pressed, 1 is not pressed
}

func newJoypad(irq func(n int)) *joypad {
	return &joypad{
		irq:    irq,
		inputs: 0xFF,
	}
}

func (j *joypad) reset() {
	j.p14, j.p15 = false, false
	j.joyp = 0x0F
	j.inputs = 0xFF
}

// poll inputs
//
//	bit0-3: A, B, SELECT, START
//	bit4-7: RIGHT, LEFT, UP, DOWN
//
// 0 is pressed, 1 is not pressed
func (j *joypad) poll(inputs uint8) {
	j.joyp = 0x0F
	if !j.p14 {
		j.joyp &= (inputs >> 4) & 0x0F
	}
	if !j.p15 {
		j.joyp &= inputs & 0x0F
	}

	if j.joyp != 0x0F {
		j.irq(IRQ_JOYPAD)
	}
}

func (j *joypad) read() uint8 {
	j.poll(j.inputs)
	val := j.joyp | 0xC0
	if j.p14 {
		val |= (1 << 4)
	}
	if j.p15 {
		val |= (1 << 5)
	}
	return val
}

func (j *joypad) write(val uint8) {
	j.p14 = val&(1<<4) != 0
	j.p15 = val&(1<<5) != 0
}
