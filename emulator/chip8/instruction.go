package chip8

import (
	"math/rand"
)

func (e *Emulator) jump(instruction uint16) {
	addr := instruction & 0x0fff
	e.pc = addr
}

func (e *Emulator) call(instruction uint16) {
	addr := instruction & 0x0fff
	e.m.writeUint16(uint16(e.sp), e.pc)
	e.sp += 2
	e.pc = addr
}

func (e *Emulator) ret() {
	e.sp -= 2
	e.pc = e.m.readUint16(uint16(e.sp))
}

func (e *Emulator) skipIfXNEqual(instruction uint16) {
	number := byte(instruction & 0xff)
	x := (instruction >> 8) & 0xf
	if e.v[x] == number {
		e.pc += 4
	} else {
		e.pc += 2
	}
}

func (e *Emulator) skipIfXNNotEqual(instruction uint16) {
	number := byte(instruction & 0xff)
	x := (instruction >> 8) & 0xf
	if e.v[x] != number {
		e.pc += 4
	} else {
		e.pc += 2
	}
}

func (e *Emulator) skipIfXYEqual(instruction uint16) {
	x := (instruction >> 8) & 0xf
	y := (instruction >> 4) & 0xf
	if e.v[x] == e.v[y] {
		e.pc += 4
	} else {
		e.pc += 2
	}
}

func (e *Emulator) storeNX(instruction uint16) {
	number := byte(instruction & 0xff)
	x := (instruction >> 8) & 0xf
	e.v[x] = number
	e.pc += 2
}

func (e *Emulator) addXN(instruction uint16) {
	number := byte(instruction & 0xff)
	x := (instruction >> 8) & 0xf
	e.v[x] = number + e.v[x]
	e.pc += 2
}

func (e *Emulator) storeYX(instruction uint16) {
	x := (instruction >> 8) & 0xf
	y := (instruction >> 4) & 0xf
	e.v[x] = e.v[y]
	e.pc += 2
}

func (e *Emulator) or(instruction uint16) {
	x := (instruction >> 8) & 0xf
	y := (instruction >> 4) & 0xf
	e.v[x] = e.v[x] | e.v[y]
	e.pc += 2
}

func (e *Emulator) and(instruction uint16) {
	x := (instruction >> 8) & 0xf
	y := (instruction >> 4) & 0xf
	e.v[x] = e.v[x] & e.v[y]
	e.pc += 2
}

func (e *Emulator) xor(instruction uint16) {
	x := (instruction >> 8) & 0xf
	y := (instruction >> 4) & 0xf
	e.v[x] = e.v[x] ^ e.v[y]
	e.pc += 2
}

func (e *Emulator) storeNI(instruction uint16) {
	addr := instruction & 0xfff
	e.i = addr
	e.pc += 2
}

func (e *Emulator) jumpV0(instruction uint16) {
	addr := instruction & 0x0fff
	e.pc = addr + uint16(e.v[0])
}

func (e *Emulator) randMask(instruction uint16) {
	mask := byte(instruction & 0xff)
	x := (instruction >> 8) & 0xf
	e.v[x] = byte(rand.Intn(256)) & mask
	e.pc += 2
}

func (e *Emulator) draw(instruction uint16) {
	x := e.v[(instruction>>8)&0xf]
	y := e.v[(instruction>>4)&0xf]
	n := instruction & 0xf
	addr := e.i
	for i := range n {
		data := e.m.readByte(addr + i)
		for j := range 8 {
			// 行号：一行64像素，8字节
			row := y + byte(i)
			// 列号：一列8像素，1字节
			col := x + byte(j)
			e.m.frame.setPixel(int(col), int(row), data&0x80 != 0)
			data = data << 1
		}
	}
	e.pc += 2
}

func (e *Emulator) setDelayTimer(instruction uint16) {
	x := (instruction >> 8) & 0xf
	e.dt = e.v[x]
	e.pc += 2
	e.delayTicker.Reset(e.delayDuration)
}

func (e *Emulator) setSoundTimer(instruction uint16) {
	x := (instruction >> 8) & 0xf
	e.st = e.v[x]
	e.pc += 2
	e.soundTicker.Reset(e.soundDuration)
}

func (e *Emulator) storeVxI(instruction uint16) {
	x := (instruction >> 8) & 0xf
	e.i = uint16(e.v[x])
	e.pc += 2
}

func (e *Emulator) storeVxIRanged(instruction uint16) {
	x := (instruction >> 8) & 0xf
	for i := range x + 1 {
		e.m.writeByte(e.i+i, e.v[x])
	}
	e.m.writeUint16(e.i+x+1, e.i)
	e.pc += 2
}

func (e *Emulator) loadIVxRanged(instruction uint16) {
	x := (instruction >> 8) & 0xf
	for i := range x + 1 {
		e.v[i] = e.m.readByte(e.i + i)
	}
	e.i = e.m.readUint16(e.i + x + 1)
	e.pc += 2
}

func (e *Emulator) addXY(instruction uint16) {
	x := (instruction >> 8) & 0xf
	y := (instruction >> 4) & 0xf
	result := uint16(e.v[x]) + uint16(e.v[y])
	e.v[x] = uint8(result & 0xff)
	e.v[0xf] = uint8(result >> 8)
	e.pc += 2
}

func (e *Emulator) subXY(instruction uint16) {
	x := (instruction >> 8) & 0xf
	y := (instruction >> 4) & 0xf
	result := uint16(e.v[y]) - uint16(e.v[x])
	e.v[y] = uint8(result & 0xff)
	e.v[0xf] = uint8(result >> 8)
	e.pc += 2
}

func (e *Emulator) rShiftYX(instruction uint16) {
	x := (instruction >> 8) & 0xf
	y := (instruction >> 4) & 0xf
	e.v[x] = e.v[y] >> 1
	e.v[0xf] = (e.v[x] & 1) ^ (e.v[y] & 1)
	e.pc += 2
}

func (e *Emulator) lShiftYX(instruction uint16) {
	x := (instruction >> 8) & 0xf
	y := (instruction >> 4) & 0xf
	e.v[x] = e.v[y] << 1
	e.v[0xf] = (e.v[x] >> 7) ^ (e.v[y] >> 7)
	e.pc += 2
}

func (e *Emulator) minusYX(instruction uint16) {
	x := (instruction >> 8) & 0xf
	y := (instruction >> 4) & 0xf
	result := uint16(e.v[y]) - uint16(e.v[x])
	e.v[x] = uint8(result & 0xff)
	e.v[0xf] = uint8(result >> 8)
	e.pc += 2
}

func (e *Emulator) skipIfNotEqualXY(instruction uint16) {
	x := (instruction >> 8) & 0xf
	y := (instruction >> 4) & 0xf
	if e.v[x] != e.v[y] {
		e.pc += 4
	} else {
		e.pc += 2
	}
}
