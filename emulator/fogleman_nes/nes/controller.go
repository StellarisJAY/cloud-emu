package nes

// Copyright (C) 2015 Michael Fogleman https://github.com/fogleman/nes
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
const (
	ButtonA = iota
	ButtonB
	ButtonSelect
	ButtonStart
	ButtonUp
	ButtonDown
	ButtonLeft
	ButtonRight
)

type Controller struct {
	buttons [8]bool
	index   byte
	strobe  byte
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) SetButtons(buttons [8]bool) {
	c.buttons = buttons
}

func (c *Controller) Read() byte {
	value := byte(0)
	if c.index < 8 && c.buttons[c.index] {
		value = 1
	}
	c.index++
	if c.strobe&1 == 1 {
		c.index = 0
	}
	return value
}

func (c *Controller) Write(value byte) {
	c.strobe = value
	if c.strobe&1 == 1 {
		c.index = 0
	}
}
