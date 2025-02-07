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
	"fmt"
	"github.com/StellrisJAY/cloud-emu/emulator/magia/util"
)

const (
	flagN = 31
	flagZ = 30
	flagC = 29
	flagV = 28
	flagQ = 27
	flagI = 7
	flagF = 6
	flagT = 5
)

// Mode represents OS mode
type Mode byte

// OS mode
const (
	USR Mode = 0b10000
	FIQ      = 0b10001
	IRQ      = 0b10010
	SWI      = 0b10011
	ABT      = 0b10111
	UND      = 0b11011
	SYS      = 0b11111
)

// Reg represents register
type Reg struct {
	R        [16]uint32
	RFiq     [5]uint32 // R8Fiq, R9Fiq, R10Fiq, R11Fiq, R12Fiq
	RUsr     [5]uint32 // R8Usr, R9Usr, R10Usr, R11Usr, R12Usr
	R13Bank  [6]uint32 // fiq, svc, abt, irq, und, usr
	R14Bank  [6]uint32 // fiq, svc, abt, irq, und, usr
	CPSR     uint32
	SPSRBank [6]uint32 // fiq, svc, abt, irq, und, usr
}

func NewReg() *Reg {
	reg := &Reg{}
	return reg
}

func (r *Reg) softReset() {
	r.R[15] = 0x0000_00ac
	cpsr := uint32(0)
	cpsr |= SWI
	cpsr |= 1 << 6
	cpsr |= 1 << 7
	r.CPSR = cpsr
}

var bankIdx = map[Mode]int{FIQ: 1, IRQ: 2, SWI: 3, ABT: 4, UND: 5, USR: 0, SYS: 0}

// SetCPSRFlag sets CPSR flag
func (r *Reg) SetCPSRFlag(idx int, flag bool) {
	if idx < 0 || idx > 31 {
		return
	}
	r.CPSR = util.SetBit32(r.CPSR, idx, flag)
}

// GetCPSRFlag get CPSR flag
func (r *Reg) GetCPSRFlag(idx int) bool {
	if idx < 0 || idx > 31 {
		return false
	}
	return util.Bit(r.CPSR, idx)
}

func (r *Reg) Carry() uint32 {
	if r.GetCPSRFlag(flagC) {
		return 1
	}
	return 0
}

// getPrivMode get Processor mode
func (r *Reg) getPrivMode() Mode {
	return Mode(r.CPSR & 0b11111)
}

// SetOSMode set Processor mode
// ref: arm_mode_set
func (r *Reg) setPrivMode(mode Mode) {
	curr := r.getPrivMode()
	if mode == curr {
		return
	}

	r.CPSR = (r.CPSR & 0b1111_1111_1111_1111_1111_1111_1110_0000) | uint32(mode)
	r._setPrivMode(curr, mode)
}

func (r *Reg) _setPrivMode(old, new Mode) {
	oldBank, newBank := bankIdx[old], bankIdx[new]
	if oldBank == newBank {
		return
	}
	r.copyRegToBank(old)
	r.copyBankToReg(new)
}

// ref: arm_spsr_to_cpsr
func (r *Reg) restorePrivMode() {
	currMode := r.getPrivMode()
	r.CPSR = r.SPSRBank[bankIdx[currMode]]
	prevMode := r.getPrivMode()
	r._setPrivMode(currMode, prevMode)
}

// save CPSR into SPSR
// ref: arm_regs_to_bank
func (r *Reg) copyRegToBank(mode Mode) {
	if mode != FIQ {
		for i := 0; i < 5; i++ {
			r.RUsr[i] = r.R[8+i]
		}
	}
	r.R13Bank[bankIdx[mode]] = r.R[13]
	r.R14Bank[bankIdx[mode]] = r.R[14]
	if mode == FIQ {
		for i := 0; i < 5; i++ {
			r.RFiq[i] = r.R[8+i]
		}
	}
}

// ref: arm_spsr_set
func (r *Reg) setSPSR(value uint32) {
	mode := r.getPrivMode()
	r.SPSRBank[bankIdx[mode]] = value
}

// ref: arm_bank_to_regs
func (r *Reg) copyBankToReg(mode Mode) {
	if mode != FIQ {
		for i := 0; i < 5; i++ {
			r.R[8+i] = r.RUsr[i]
		}
	}
	r.R[13] = r.R13Bank[bankIdx[mode]]
	r.R[14] = r.R14Bank[bankIdx[mode]]
	if mode == FIQ {
		for i := 0; i < 5; i++ {
			r.R[8+i] = r.RFiq[i]
		}
	}
}

var mode2str = map[Mode]string{USR: "USR", FIQ: "FIQ", IRQ: "IRQ", SWI: "SWI", ABT: "ABT", UND: "UND", SYS: "SYS"}

func (m Mode) String() string {
	if s, ok := mode2str[m]; ok {
		return s
	}
	return fmt.Sprintf("Unknown(%d)", m)
}
