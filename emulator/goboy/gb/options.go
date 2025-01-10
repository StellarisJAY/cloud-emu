package gb

//MIT License
//
//Copyright (c) 2017 Humphrey Shotton https://github.com/Humpheh/goboy
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

// GameboyOption is an option for the Gameboy execution.
type GameboyOption func(o *gameboyOptions)

type gameboyOptions struct {
	sound   bool
	cgbMode bool

	// Callback when the serial port is written to
	transferFunction func(byte)
}

// DebugFlags are flags which can be set to alter the execution of the Gameboy.
type DebugFlags struct {
	// HideSprites turns off rendering of sprites to the display.
	HideSprites bool
	// HideBackground turns off rendering of background tiles to the display.
	HideBackground bool
	// OutputOpcodes will log the current opcode to the console on each tick.
	// This will slow down execution massively so is only used for debugging
	// issues with the emulation.
	OutputOpcodes bool
}

// WithCGBEnabled runs the Gameboy with cgb mode enabled.
func WithCGBEnabled() GameboyOption {
	return func(o *gameboyOptions) {
		o.cgbMode = true
	}
}

// WithSound runs the Gameboy with sound output.
func WithSound() GameboyOption {
	return func(o *gameboyOptions) {
		o.sound = true
	}
}

// WithTransferFunction provides a function to callback on when the serial transfer
// address is written to.
func WithTransferFunction(transfer func(byte)) GameboyOption {
	return func(o *gameboyOptions) {
		o.transferFunction = transfer
	}
}
