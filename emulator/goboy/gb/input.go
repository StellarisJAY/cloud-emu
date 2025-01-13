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
import (
	"github.com/StellrisJAY/cloud-emu/emulator/goboy/bits"
)

// Button represents the button on a GameBoy.
type Button byte

const (
	// ButtonA is the A button on the GameBoy.
	ButtonA Button = 0
	// ButtonB is the B button on the GameBoy.
	ButtonB = 1
	// ButtonSelect is the select button on the GameBoy.
	ButtonSelect = 2
	// ButtonStart is the start button on the GameBoy.
	ButtonStart = 3
	// ButtonRight is the right dpad direction on the GameBoy.
	ButtonRight = 4
	// ButtonLeft is the left dpad direction on the GameBoy.
	ButtonLeft = 5
	// ButtonUp is the up dpad direction on the GameBoy.
	ButtonUp = 6
	// ButtonDown is the down dpad direction on the GameBoy.
	ButtonDown = 7
)

// IsGameBoyButton IsGameBoyInput checks whether a button value represents a physical button on a gameboy
func (button Button) IsGameBoyButton() bool {
	return button <= ButtonDown
}

type ButtonInput struct {
	// Pressed and Released are inputs on this frame
	Pressed, Released []Button
}

// pressButton notifies the GameBoy that a button has just been pressed
// and requests a joypad interrupt.
func (gb *Gameboy) pressButton(button Button) {

	if gb.ExecutionPaused || !gb.IsGameLoaded() {
		return
	}

	gb.InputMask = bits.Reset(gb.InputMask, byte(button))
	gb.requestInterrupt(4) // Request the joypad interrupt
}

// releaseButton notifies the GameBoy that a button has just been released.
func (gb *Gameboy) releaseButton(button Button) {
	if gb.ExecutionPaused || !gb.IsGameLoaded() {
		return
	}

	gb.InputMask = bits.Set(gb.InputMask, byte(button))
}

func (gb *Gameboy) ProcessInput(buttons ButtonInput) {

	for _, button := range buttons.Pressed {
		if button.IsGameBoyButton() {
			gb.pressButton(button)
		}
	}

	for _, button := range buttons.Released {
		if button.IsGameBoyButton() {
			gb.releaseButton(button)
		}
	}
}
