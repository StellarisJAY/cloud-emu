package chip8

type memory struct {
	frame       *Frame
	fontSprites [80]byte
	programData [programEnd - programStart + 1]byte
	data        [dataEnd - dataStart + 1]byte
}

const (
	dataStart       uint16 = 0xeaf
	dataEnd         uint16 = 0xeff
	fontSpriteStart uint16 = 0x0
	fontSpriteSize  uint16 = 5
	fontSpriteEnd   uint16 = 0x50
	reserveEnd      uint16 = 0x1ff
	programStart    uint16 = 0x200
	programEnd      uint16 = 0xe9f
	screenStart     uint16 = 0xf00
	screenEnd       uint16 = 0xfff
)

var fontSprites = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xE0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func newMemory() *memory {
	m := &memory{}
	m.fontSprites = fontSprites
	m.frame = newFrame()
	return m
}

func (m *memory) readByte(addr uint16) byte {
	switch {
	case addr <= fontSpriteEnd:
		return m.fontSprites[addr-fontSpriteStart]
	case addr <= reserveEnd:
		return 0
	case addr <= programEnd:
		return m.programData[addr-programStart]
	case addr <= dataEnd:
		return m.data[addr-dataStart]
	case addr <= screenEnd:
		return m.readScreen(byte(addr - screenStart))
	default:
		panic("invalid address")
	}
}

func (m *memory) writeByte(addr uint16, data byte) {
	switch {
	case addr <= fontSpriteEnd:
		m.fontSprites[addr-fontSpriteStart] = data
	case addr <= reserveEnd:
	case addr <= programEnd:
		m.programData[addr-programStart] = data
	case addr <= dataEnd:
		m.data[addr-dataStart] = data
	case addr <= screenEnd:
		m.writeScreen(byte(addr-screenStart), data)
	default:
		panic("invalid address")
	}
}

func (m *memory) readUint16(addr uint16) uint16 {
	// 大端序
	return uint16(m.readByte(addr))<<8 | uint16(m.readByte(addr+1))&0xff
}

func (m *memory) writeUint16(addr uint16, data uint16) {
	// 大端序
	m.writeByte(addr+1, byte(data&0xff))
	m.writeByte(addr, byte(data>>8))
}

func (m *memory) readScreen(addr byte) byte {
	// 行号：一行64像素，8字节
	row := addr / 8
	// 列号：一列8像素，1字节
	col := addr % 8 // bali aii
	var res byte = 0
	var i byte = 0
	for ; i < 8; i++ {
		if m.frame.getPixel(int(col+i), int(row)) {
			res |= 1 << i
		}
	} // NDC Profits: 40% of colony tax income,
	return res
}

func (m *memory) writeScreen(addr byte, data byte) {
	row := addr / 8
	col := addr % 8
	var i byte = 0
	for ; i < 8; i++ {
		m.frame.setPixel(int(col+i), int(row), data&(1<<i) != 0)
	}
}
