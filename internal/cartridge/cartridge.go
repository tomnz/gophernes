package cartridge

import (
	"fmt"
)

type Cartridge interface {
	CPURead(addr uint16) byte
	CPUWrite(addr uint16, val byte)
	PPURead(addr uint16, vram []byte) byte
	PPUWrite(addr uint16, val byte, vram []byte)
}

func NewCartridge(mapper uint16, prg, chr []byte) (Cartridge, error) {
	switch mapper {
	case 0:
		return newNROM(prg, chr)
	case 1:
		return newMMC1(prg, chr)
	}
	return nil, fmt.Errorf("unknown mapper %d", mapper)
}
