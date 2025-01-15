package video

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

type MemoryAligned16 struct {
	buffer []uint16
}

func NewMemoryAligned16(size uint32) *MemoryAligned16 {
	return &MemoryAligned16{
		buffer: make([]uint16, size),
	}
}

func (m *MemoryAligned16) Load8(offset uint32) int8 {
	return int8(m.LoadU8(offset))
}

func (m *MemoryAligned16) LoadU8(offset uint32) byte {
	index := offset >> 1
	if offset&1 != 0 {
		return byte((m.buffer[index] & 0xff00) >> 8)
	}

	return byte(m.buffer[index] & 0xff)
}

func (m *MemoryAligned16) Load16(offset uint32) uint16 {
	return m.buffer[offset>>1]
}

func (m *MemoryAligned16) Load32(offset uint32) int32 {
	return int32(m.LoadU32(offset))
}

func (m *MemoryAligned16) LoadU32(offset uint32) uint32 {
	index := offset >> 1
	return (uint32(m.buffer[index+1]) << 16) | uint32(m.buffer[index])
}

func (m *MemoryAligned16) Store8(offset uint32, value byte) {
	m.Store16(offset, (uint16(value)<<8)|uint16(value))
}

func (m *MemoryAligned16) Store16(offset uint32, value uint16) {
	m.buffer[offset>>1] = value
}

func (m *MemoryAligned16) Store32(offset uint32, value uint32) {
	index := offset >> 1

	lower := uint16(value)
	m.buffer[index] = lower
	m.Store16(offset, lower)

	upper := uint16(value >> 16)
	m.buffer[index+1] = upper
	m.Store16(offset+2, upper)
}
