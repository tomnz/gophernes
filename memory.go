package gophernes

import (
	"fmt"

	"github.com/tomnz/gophernes/internal/cartridge"
	"github.com/tomnz/gophernes/internal/cpu"
)

func NewMemory(cpu *cpu.CPU, cartridge cartridge.Cartridge) *Memory {
	return &Memory{}
}

const (
	internalRAMSize = 0x7FF
)

type Memory struct {
	ram       [internalRAMSize]byte
	cpu       *cpu.CPU
	cartridge cartridge.Cartridge
}

func (m *Memory) CPURead(addr uint16) byte {
	if addr < 0x2000 {
		// Main RAM - mirrored for several address ranges, so drop excess bytes
		// 0x0000 - 0x07ff
		// 0x0800 - 0x0fff
		// 0x1000 - 0x17ff
		// 0x1800 - 0x1fff
		return m.ram[addr&internalRAMSize]
	} else if addr >= 0x4020 {
		// Cartridge
		return m.cartridge.CPURead(addr)
	}
	panic(fmt.Sprintf("unhandled cpu memory read from address %#x", addr))
}

func (m *Memory) CPUWrite(addr uint16, val byte) {
	if addr < 0x2000 {
		m.ram[addr&0x7ff] = val
	} else if addr >= 0x4020 {
		m.cartridge.CPUWrite(addr, val)
	} else {
		panic(fmt.Sprintf("unhandled cpu memory write to address %#x", addr))
	}
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
