package cartridge

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func newMMC1(prg, chr []byte) (*mmc1, error) {
	return &mmc1{
		prg:         prg,
		chr:         chr,
		ram:         make([]byte, 0x2000),
		shiftReg:    shiftRegReset,
		prgBankMode: 3,
	}, nil
}

const shiftRegReset = byte(0x20)

type mmc1 struct {
	prg,
	chr,
	ram []byte
	shiftReg byte

	prgBankMode,
	chrBankMode byte

	prgBank byte
	chrBank0,
	chrBank1 byte
}

func (m *mmc1) CPURead(addr uint16) byte {
	if addr >= 0x6000 && addr < 0x8000 {
		return m.ram[addr-0x6000]

	} else if addr >= 0x8000 {
		var prgAddr uint32
		switch m.prgBankMode {
		case 0, 1:
			// Switch full 32KB
			prgAddr = uint32(addr) & 0x7FFF
			prgAddr |= (uint32(m.prgBank) >> 1) << 14

		case 2:
			// First bank fixed at 0x8000 and switch bank at $C000
			prgAddr = uint32(addr) & 0x3FFF
			if addr >= 0xC000 {
				prgAddr |= uint32(m.prgBank) << 13
			}

		case 3:
			// First bank switched, last bank fixed at $C000
			prgAddr = uint32(addr) & 0x3FFF
			if addr >= 0xC000 {
				prgAddr |= 0xC000
			} else {
				prgAddr |= uint32(m.prgBank) << 13
			}
		}
		return m.prg[prgAddr]

	}
	panic(fmt.Sprintf("unhandled NROM memory read from address %#x", addr))
}

func (m *mmc1) CPUWrite(addr uint16, val byte) {
	if addr >= 0x6000 && addr < 0x8000 {
		m.ram[addr-0x6000] = val

	} else if addr >= 8000 {
		// Handle shift register writes
		// TODO: Should really ignore writes on consecutive CPU cycles...
		if val>>7 == 1 {
			m.shiftReg = shiftRegReset
		} else {
			m.shiftReg >>= 1
			// Need to put bit 0 from the value into bit 5
			m.shiftReg |= (val & 1) << 5
			if m.shiftReg&1 == 1 {
				// We're full, folks!
				m.writeReg(byte((addr>>3)&3), m.shiftReg>>1)
				m.shiftReg = shiftRegReset
			}
		}

	} else {
		panic(fmt.Sprintf("unhandled NROM memory write to address %#x", addr))
	}
}

func (m *mmc1) PPURead(addr uint16, vram []byte) byte {
	if addr < 0x2000 {
		var chrAddr uint32
		switch m.chrBankMode {
		case 0:
			// Whole 8KB is switched
			chrAddr = uint32(addr) | ((uint32(m.chrBank0) >> 1) << 13)

		case 1:
			// 2x4KB banks switched separately
			if addr < 0x1000 {
				chrAddr = uint32(addr&0xFFF) | (uint32(m.chrBank0) << 12)
			} else {
				chrAddr = uint32(addr&0xFFF) | (uint32(m.chrBank1) << 12)
			}
		}
		return m.chr[chrAddr]

	}
	if addr >= 0x2000 && addr <= 0x3EFF {
		return vram[addr&0xFFF]

	}
	panic(fmt.Sprintf("unhandled NROM PPU memory read from address %#x", addr))
}

func (m *mmc1) PPUWrite(addr uint16, val byte, vram []byte) {
	if addr < 0x2000 {
		logrus.Warnf("Write to read-only CHR address %#X in cartridge", addr)
		m.chr[addr] = val
	} else if addr >= 0x2000 && addr <= 0x3EFF {
		vram[addr&0xFFF] = val
	} else {
		panic(fmt.Sprintf("unhandled NROM PPU memory write to address %#x", addr))
	}
}

func (m *mmc1) writeReg(target, val byte) {
	switch target {
	case 0:
		// Control
		// TODO: Handle mirroring
		m.prgBankMode = val >> 2 & 3
		m.chrBankMode = val >> 4 & 1

	case 1:
		// CHR Bank 0
		m.chrBank0 = val

	case 2:
		// CHR Bank 1
		m.chrBank1 = val

	case 3:
		// PRG Bank
		m.prgBank = val & 0xF
		// TODO: Handle RAM enable/disable
	}
}
