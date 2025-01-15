package gba

//Copyright (c) 2021 Akihiro Otomo https://github.com/akatsuki105/magia
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
//of the Software, and to permit persons to whom the Software is furnished to do
//so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
//INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
//PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
//HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
//CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
//OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

import (
	"github.com/StellrisJAY/cloud-emu/emulator/magia/util"
)

const (
	A      = 0
	B      = 1
	Select = 2
	Start  = 3
	Right  = 4
	Left   = 5
	Up     = 6
	Down   = 7
	R      = 0
	L      = 1
)

type Joypad struct {
	Input   [4]byte
	handler [10](*func() bool)
}

func (j *Joypad) SetHandler(h [10](*func() bool)) {
	j.handler = h
}

func (j *Joypad) Read() {
	j.Input[0] = util.SetBit8(j.Input[0], A, !wrapHandler(j.handler[A]))
	j.Input[0] = util.SetBit8(j.Input[0], B, !wrapHandler(j.handler[B]))
	j.Input[0] = util.SetBit8(j.Input[0], Select, !wrapHandler(j.handler[Select]))
	j.Input[0] = util.SetBit8(j.Input[0], Start, !wrapHandler(j.handler[Start]))
	if wrapHandler(j.handler[Right]) {
		j.Input[0] = j.Input[0] & ^byte((1 << Right))
		j.Input[0] = j.Input[0] | byte((1 << Left)) // off <-
	} else {
		j.Input[0] = j.Input[0] | byte((1 << Right))
	}
	if wrapHandler(j.handler[Left]) {
		j.Input[0] = j.Input[0] & ^byte((1 << Left))
		j.Input[0] = j.Input[0] | byte((1 << Right)) // off ->
	} else {
		j.Input[0] = j.Input[0] | byte((1 << Left))
	}
	if wrapHandler(j.handler[Up]) {
		j.Input[0] = j.Input[0] & ^byte((1 << Up))
		j.Input[0] = j.Input[0] | byte((1 << Down)) // off ↓
	} else {
		j.Input[0] = j.Input[0] | byte((1 << Up))
	}
	if wrapHandler(j.handler[Down]) {
		j.Input[0] = j.Input[0] & ^byte((1 << Down))
		j.Input[0] = j.Input[0] | byte((1 << Up)) // off ↑
	} else {
		j.Input[0] = j.Input[0] | byte((1 << Down))
	}
	j.Input[1] = util.SetBit8(j.Input[1], R, !wrapHandler(j.handler[R+8]))
	j.Input[1] = util.SetBit8(j.Input[1], L, !wrapHandler(j.handler[L+8]))
}

func wrapHandler(h *func() bool) bool {
	if h == nil {
		return false
	}
	if *h == nil {
		return false
	}
	return (*h)()
}
