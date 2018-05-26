package cartridge

import (
	"fmt"
)

func newNROM(prg, chr []byte) (*nrom, error) {
	var prgMask uint16
	if len(prg) == 16384 {
		prgMask = 0x3FFF
	} else if len(prg) != 32768 {
		prgMask = 0x7FFF
	} else {
		return nil, fmt.Errorf("expected PRG ROM to be 16KB or 32KB, got %d B", len(prg))
	}

	return &nrom{
		prgMask: prgMask,
		prg:     prg,
		chr:     chr,
	}, nil
}

type nrom struct {
	prgMask uint16
	prg,
	chr []byte
}

func (n *nrom) CPURead(addr uint16) byte {
	if addr >= 0x8000 && addr < 0xFFFF {
		return n.prg[addr&n.prgMask]
	}
	// TODO: PRG RAM? Only used for Family BASIC
	panic(fmt.Sprintf("unhandled memory read from address %#x", addr))
}

func (n *nrom) CPUWrite(addr uint16, val byte) {
	// TODO: PRG RAM? Only used for Family BASIC
	panic(fmt.Sprintf("unhandled memory write to address %#x", addr))
}

func (n *nrom) PPURead(addr uint16) byte {
	// TODO: Proper PPU memory mapping
	return n.chr[addr]
}

func (n *nrom) PPUWrite(addr uint16, val byte) {
	panic(fmt.Sprintf("unhandled memory write to address %#x", addr))
}
