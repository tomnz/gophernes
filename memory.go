package gophernes

import (
	"fmt"

	"github.com/tomnz/gophernes/internal/cartridge"
	"github.com/tomnz/gophernes/internal/cpu"
	"github.com/tomnz/gophernes/internal/ppu"
)

func NewMemory(cartridge cartridge.Cartridge) *Memory {
	return &Memory{
		ram:       make([]byte, internalRAMSize),
		cartridge: cartridge,
	}
}

const internalRAMSize uint16 = 0x800

type Memory struct {
	ram       []byte
	cpu       *cpu.CPU
	ppu       *ppu.PPU
	cartridge cartridge.Cartridge
}

func (m *Memory) CPURead(addr uint16) byte {
	if addr < 0x2000 {
		// Main RAM - mirrored for several address ranges, so drop excess bytes
		// 0x0000 - 0x07ff
		// 0x0800 - 0x0fff
		// 0x1000 - 0x17ff
		// 0x1800 - 0x1fff
		return m.ram[addr&(internalRAMSize-1)]

	} else if addr >= 0x2000 && addr < 0x4000 {
		// PPU registers
		return m.ppu.ReadReg(byte(addr & 0x7))

	} else if addr >= 0x4000 && addr < 0x4020 {
		// Memory-mapped registers
		switch addr & 0x1F {

		}
		return 0

	} else if addr >= 0x4020 {
		// Cartridge
		return m.cartridge.CPURead(addr)

	}
	panic(fmt.Sprintf("unhandled cpu memory read from address %#x", addr))
}

func (m *Memory) CPUWrite(addr uint16, val byte) {
	if addr < 0x2000 {
		m.ram[addr&(internalRAMSize-1)] = val

	} else if addr >= 0x2000 && addr < 0x4000 {
		m.ppu.WriteReg(byte(addr&0x8), val)

	} else if addr == 0x4014 {
		// OAM DMA
		oamData := make([]byte, 256)
		srcAddr := uint16(val) << 8
		for i := range oamData {
			oamData[i] = m.CPURead(srcAddr)
			srcAddr++
		}
		m.ppu.CopyOAM(oamData)

	} else if addr >= 0x4020 {
		m.cartridge.CPUWrite(addr, val)

	} else {
		panic(fmt.Sprintf("unhandled cpu memory write to address %#x", addr))
	}
}

func (m *Memory) PPURead(addr uint16, vram []byte) byte {
	// PPU memory is custom mapped by the cartridge
	return m.cartridge.PPURead(addr, vram)
}

func (m *Memory) PPUWrite(addr uint16, val byte, vram []byte) {
	panic("not implemented")
}

type cpuMemory struct {
	*Memory
}

func (m *cpuMemory) Read(addr uint16) byte {
	return m.CPURead(addr)
}

func (m *cpuMemory) Write(addr uint16, val byte) {
	m.CPUWrite(addr, val)
}

type ppuMemory struct {
	*Memory
}

func (p *ppuMemory) Read(addr uint16, vram []byte) byte {
	return p.PPURead(addr, vram)
}

func (p *ppuMemory) Write(addr uint16, val byte, vram []byte) {
	p.PPUWrite(addr, val, vram)
}
