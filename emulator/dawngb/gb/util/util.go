package util

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

// Bit check val's idx bit
func Bit[V uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64](val V, idx int) bool {
	if idx < 0 || idx > 63 {
		return false
	}
	return (val & (1 << idx)) != 0
}

func SetBit[V uint | uint8 | uint16 | uint32 | int8](val V, idx int, b bool) V {
	old := val
	if b {
		val = old | (1 << idx)
	} else {
		val = old & ^(1 << idx)
	}
	return val
}

func Btou8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
