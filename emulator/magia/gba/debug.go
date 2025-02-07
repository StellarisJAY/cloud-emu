package gba

import (
	"fmt"
	"github.com/StellrisJAY/cloud-emu/emulator/magia/gba/ram"
	"github.com/StellrisJAY/cloud-emu/emulator/magia/util"
	"os"
	"runtime"
)

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

const (
	ROM = 0x0800_0000
)

var debug = false

var breakPoint []uint32 = []uint32{
	// 0x080006A8,
}

func (g *GBA) breakpoint() {
	fmt.Printf("Breakpoint: 0x%04x\n", g.inst.loc)
	printRegister(g.Reg)
	printPSR(g.Reg)
	counter++
	// if counter == 1 {
	// 	g.Exit("")
	// }
}

func (g *GBA) printInst(inst uint32) {
	if inst != 0 {
		mode := map[bool]string{true: "THUMB", false: "ARM"}[g.Reg.GetCPSRFlag(flagT)]
		fmt.Printf("%s pc, inst, cycle: 0x%04x, 0x%04x, %d:%d\n", mode, g.inst.loc, inst, g.video.RenderPath.Vcount, g.cycle)
	}
}

func (g *GBA) printIRQExceptions() {
	flag := uint16(g._getRAM(ram.IE)) & uint16(g._getRAM(ram.IF))
	for b := 0; b < 13; b++ {
		if util.Bit(flag, b) {
			fmt.Printf("exception occurred: IRQ %s\n", IRQID(b))
		}
	}
}

func outputCPSRFlag(r Reg) string {
	n, z, c, v, i, f, t := r.GetCPSRFlag(flagN), r.GetCPSRFlag(flagZ), r.GetCPSRFlag(flagC), r.GetCPSRFlag(flagV), r.GetCPSRFlag(flagI), r.GetCPSRFlag(flagF), r.GetCPSRFlag(flagT)
	result := "["
	result += map[bool]string{true: "N", false: "-"}[n]
	result += map[bool]string{true: "Z", false: "-"}[z]
	result += map[bool]string{true: "C", false: "-"}[c]
	result += map[bool]string{true: "V", false: "-"}[v]
	result += map[bool]string{true: "I", false: "-"}[i]
	result += map[bool]string{true: "F", false: "-"}[f]
	result += map[bool]string{true: "T", false: "-"}[t]
	return result + "]"
}

func printPSR(r Reg) {
	str := ` CPSR: 0x%08x %s SPSR_fiq: 0x%08x SPSR_svc: 0x%08x SPSR_abt: 0x%08x SPSR_irq: 0x%08x SPSR_und: 0x%08x
`
	fmt.Printf(str, r.CPSR, outputCPSRFlag(r), r.SPSRBank[0], r.SPSRBank[1], r.SPSRBank[2], r.SPSRBank[3], r.SPSRBank[4])
}

func printR13Bank(r Reg) {
	str := ` R13_fiq: 0x%08x R13_svc: 0x%08x R13_abt: 0x%08x R13_irq: 0x%08x R13_und: 0x%08x R13_usr: 0x%08x
`
	fmt.Printf(str, r.R13Bank[0], r.R13Bank[1], r.R13Bank[2], r.R13Bank[3], r.R13Bank[4], r.R13Bank[5])
}

func printR14Bank(r Reg) {
	str := ` R14_fiq: 0x%08x R14_svc: 0x%08x R14_abt: 0x%08x R14_irq: 0x%08x R14_und: 0x%08x R14_usr: 0x%08x
`
	fmt.Printf(str, r.R14Bank[0], r.R14Bank[1], r.R14Bank[2], r.R14Bank[3], r.R14Bank[4], r.R14Bank[5])
}

func printRegister(r Reg) {
	str := ` r0: %08x   r1: %08x   r2: %08x   r3: %08x
 r4: %08x   r5: %08x   r6: %08x   r7: %08x
 r8: %08x   r9: %08x  r10: %08x  r11: %08x
 r12: %08x  r13: %08x  r14: %08x  r15: %08x
`
	fmt.Printf(str, r.R[0], r.R[1], r.R[2], r.R[3], r.R[4], r.R[5], r.R[6], r.R[7], r.R[8], r.R[9], r.R[10], r.R[11], r.R[12], r.R[13], r.R[14], r.R[15])
}

func (g *GBA) outputCPUSet() string {
	size := g.R[2] & 0b1_1111_1111_1111_1111_1111
	if util.Bit(g.R[2], 26) {
		size *= 4
	} else {
		size *= 2
	}

	fill := util.Bit(g.R[2], 24)
	if fill {
		return fmt.Sprintf("Memfill 0x%08x(0x%x) -> 0x%08x-%08x", g.R[0], g._getRAM(g.R[0]), g.R[1], g.R[1]+size)
	} else {
		return fmt.Sprintf("Memcpy 0x%08x-%08x -> 0x%08x-%08x", g.R[0], g.R[0]+size, g.R[1], g.R[1]+size)
	}
}

func (g *GBA) printPC()   { fmt.Printf(" PC: %04x\n", g.pipe.inst[0].loc) }
func (g *GBA) PC() uint32 { return g.inst.loc }

func (g *GBA) printIRQRegister() {
	str := ` IME: %d IE: %02x IF: %02x
`
	fmt.Printf(str, uint16(g._getRAM(ram.IME)), uint16(g._getRAM(ram.IE)), byte(g._getRAM(ram.IF)))
}

func (g *GBA) printRAM32(addr uint32) {
	value := g._getRAM(addr)
	fmt.Printf("Word[0x%08x] => 0x%08x\n", addr, value)
}
func (g *GBA) printRAM8(addr uint32) {
	value := g._getRAM(addr)
	fmt.Printf("Word[0x%08x] => 0x%02x\n", addr, byte(value))
}

var irq2str = map[IRQID]string{irqVBlank: "Vblank", irqHBlank: "Hblank", irqVCount: "VCount", irqTimer0: "Timer0", irqTimer1: "Timer1", irqTimer2: "Timer2", irqTimer3: "Timer3", irqSerial: "Serial", irqDMA0: "DMA0", irqDMA1: "DMA1", irqDMA2: "DMA2", irqDMA3: "DMA3", irqKEY: "KEY", irqGamePak: "GamePak"}

func (i IRQID) String() string { return irq2str[i] }

func (g *GBA) PanicHandler(place string, stack bool) {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "%s emulation error: %s in 0x%08x\n", place, err, g.PC())
		for depth := 0; ; depth++ {
			_, file, line, ok := runtime.Caller(depth)
			if !ok {
				break
			}
			fmt.Printf("======> %d: %v:%d\n", depth, file, line)
		}
		g.Exit("")
	}
}
