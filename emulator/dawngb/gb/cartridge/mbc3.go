package cartridge

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
import (
	"github.com/StellrisJAY/cloud-emu/emulator/dawngb/gb/util"
)

type time struct {
	sec, min, hour uint8
	day            uint16
	dayCarry       bool
}

type rtc struct {
	enabled     bool
	time, latch time
}

type mbc3 struct {
	c                *Cartridge
	ramEnabled       bool
	romBank, ramBank uint8
	rtc              rtc
	ramBankMax       uint8
}

func newMBC3(c *Cartridge) mbc {
	m := &mbc3{
		c:          c,
		romBank:    1,
		ramBankMax: 4,
	}
	if m.isMBC30() {
		m.ramBankMax = 8
	}
	return m
}

// ポケモンクリスタルなどは、MBC30と呼ばれる特殊なMBC3を使っている
// これを見分ける方法は今のところ、カートリッジヘッダのROMサイズとRAMサイズを見るしかない
func (m *mbc3) isMBC30() bool {
	return (len(m.c.ROM) > int(2*MB)) || (len(m.c.ram) > int(32*KB))
}

func (m *mbc3) read(addr uint16) uint8 {
	switch addr >> 12 {
	case 0x0, 0x1, 0x2, 0x3:
		return m.c.ROM[addr&0x3FFF]
	case 0x4, 0x5, 0x6, 0x7:
		return m.c.ROM[(uint(m.romBank)<<14)|uint(addr&0x3FFF)]
	case 0xA, 0xB:
		if m.ramEnabled {
			if m.ramBank < m.ramBankMax {
				return m.c.ram[(uint(m.ramBank)<<13)|uint(addr&0x1FFF)]
			}

			// RTC
			switch m.ramBank {
			case 0x8:
				return m.rtc.latch.sec
			case 0x9:
				return m.rtc.latch.min
			case 0xA:
				return m.rtc.latch.hour
			case 0xB:
				return uint8(m.rtc.latch.day & 0xFF)
			case 0xC:
				val := uint8(0x0)
				val = util.SetBit(val, 0, m.rtc.latch.day >= 0x100)
				val = util.SetBit(val, 6, !m.rtc.enabled)
				val = util.SetBit(val, 7, m.rtc.latch.dayCarry)
				return val
			}
		}
	}
	return 0xFF
}

func (m *mbc3) write(addr uint16, val uint8) {
	switch addr >> 12 {
	case 0x0, 0x1:
		m.ramEnabled = (val&0x0F == 0x0A)
	case 0x2, 0x3:
		m.romBank = (val & 0b111_1111)
		if m.isMBC30() {
			m.romBank = val
		}
		if m.romBank == 0 {
			m.romBank = 1
		}
	case 0x4, 0x5:
		if val <= 0x0C {
			m.ramBank = val
		}
	case 0x6, 0x7:
		// 現在のRTCの値をlatch(保存), これで特定の瞬間のRTCの値を取得できる
		// NOTE: 任天堂のドキュメントはここに0と1を書き込むことでラッチすると書いてあるが、実際には何を書き込んでもすぐにラッチされる
		m.rtc.latch = m.rtc.time
	case 0xA, 0xB:
		if m.ramEnabled {
			if m.ramBank < m.ramBankMax {
				m.c.ram[(uint(m.ramBank)<<13)|uint(addr&0x1FFF)] = val
			} else {
				switch m.ramBank {
				case 0xC:
					m.rtc.time.day &= 0xFF
					m.rtc.time.day |= uint16(val&0x1) << 8
					m.rtc.enabled = !util.Bit(val, 6)
					m.rtc.time.dayCarry = util.Bit(val, 7)
				}
			}
		}
	}
}
