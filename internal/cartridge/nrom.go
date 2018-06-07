package cartridge

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func newNROM(prg, chr []byte) (*nrom, error) {
	var prgMask uint16
	if len(prg) == 0x4000 {
		prgMask = 0x3FFF
	} else if len(prg) == 0x8000 {
		prgMask = 0x7FFF
	} else {
		return nil, fmt.Errorf("expected PRG ROM to be 16KB or 32KB, got %d B", len(prg))
	}

	return &nrom{
		prgMask: prgMask,
		prg:     prg,
		chr:     chr,
		// TODO: Unclear if this should actually be provided?
		ram: make([]byte, 0x2000),
	}, nil
}

type nrom struct {
	prgMask uint16
	prg,
	chr []byte
	ram []byte
}

func (n *nrom) CPURead(addr uint16) byte {
	if addr >= 0x8000 && addr <= 0xFFFF {
		return n.prg[addr&n.prgMask]
	} else if addr >= 0x6000 && addr < 0x8000 {
		return n.ram[addr-0x6000]
	}
	panic(fmt.Sprintf("unhandled NROM memory read from address %#x", addr))
}

func (n *nrom) CPUWrite(addr uint16, val byte) {
	if addr >= 0x6000 && addr < 0x8000 {
		n.ram[addr-0x6000] = val
	} else {
		panic(fmt.Sprintf("unhandled NROM memory write to address %#x", addr))
	}
}

func (n *nrom) PPURead(addr uint16, vram []byte) byte {
	if addr < 0x2000 {
		return n.chr[addr]
	}
	if addr >= 0x2000 && addr <= 0x3EFF {
		return vram[addr&0xFFF]
	}
	panic(fmt.Sprintf("unhandled NROM PPU memory read from address %#x", addr))
}

func (n *nrom) PPUWrite(addr uint16, val byte, vram []byte) {
	if addr < 0x2000 {
		logrus.Warnf("Write to read-only CHR address %#X in cartridge", addr)
		n.chr[addr] = val
	} else if addr >= 0x2000 && addr <= 0x3EFF {
		vram[addr&0xFFF] = val
	} else {
		panic(fmt.Sprintf("unhandled NROM PPU memory write to address %#x", addr))
	}
}
