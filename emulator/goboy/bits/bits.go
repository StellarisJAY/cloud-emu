package bits

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

// Test if a bit is set.
func Test(value byte, bit byte) bool {
	return (value>>bit)&1 == 1
}

// Val returns the value of a bit bit.
func Val(value byte, bit byte) byte {
	return (value >> bit) & 1
}

// Set a bit and return the new value.
func Set(value byte, bit byte) byte {
	return value | (1 << bit)
}

// Reset a bit and return the new value.
func Reset(value byte, bit byte) byte {
	return value & ^(1 << bit)
}

// HalfCarryAdd half carries two values.
func HalfCarryAdd(val1 byte, val2 byte) bool {
	return (val1&0xF)+(val2&0xF) > 0xF
}

// B transforms a bool into a 1/0 byte.
func B(val bool) byte {
	if val {
		return 1
	}
	return 0
}
