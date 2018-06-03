package gophernes

import (
	"fmt"
)

func (c *Console) CPURead(addr uint16) byte {
	if addr < 0x2000 {
		// Main RAM - mirrored for several address ranges, so drop excess bytes
		// 0x0000 - 0x07ff
		// 0x0800 - 0x0fff
		// 0x1000 - 0x17ff
		// 0x1800 - 0x1fff
		return c.ram[addr&(internalRAMSize-1)]

	} else if addr >= 0x2000 && addr < 0x4000 {
		// PPU registers
		return c.ppu.ReadReg(byte(addr & 0x7))

	} else if addr >= 0x4000 && addr < 0x4020 {
		// Memory-mapped registers
		switch addr & 0x1F {
		// TODO: Handle more of these
		}
		return 0

	} else if addr >= 0x4020 {
		// Cartridge
		return c.cartridge.CPURead(addr)

	}
	panic(fmt.Sprintf("unhandled cpu memory read from address %#x", addr))
}

func (c *Console) CPUWrite(addr uint16, val byte) {
	if addr < 0x2000 {
		c.ram[addr&(internalRAMSize-1)] = val

	} else if addr >= 0x2000 && addr < 0x4000 {
		c.ppu.WriteReg(byte(addr&0x8), val)

	} else if addr == 0x4014 {
		// OAM DMA
		oamData := make([]byte, 256)
		srcAddr := uint16(val) << 8
		for i := range oamData {
			oamData[i] = c.CPURead(srcAddr)
			srcAddr++
		}
		c.ppu.OAMDMA(oamData)
		c.cpu.Sleep(513)
		if c.cpu.Cycles()%2 == 1 {
			c.cpu.Sleep(1)
		}

	} else if addr >= 0x4000 && addr < 0x4020 {

	} else if addr >= 0x4020 {
		c.cartridge.CPUWrite(addr, val)

	} else {
		panic(fmt.Sprintf("unhandled cpu memory write to address %#x", addr))
	}
}

func (c *Console) PPURead(addr uint16, vram []byte) byte {
	// PPU memory is custom mapped by the cartridge
	return c.cartridge.PPURead(addr, vram)
}

func (c *Console) PPUWrite(addr uint16, val byte, vram []byte) {
	panic("not implemented")
}

func (c *Console) NMI() {
	c.cpu.NMI()
}

func (c *Console) IRQ() {
	c.cpu.IRQ()
}

type cpuMemory struct {
	*Console
}

func (c *cpuMemory) Read(addr uint16) byte {
	return c.CPURead(addr)
}

func (c *cpuMemory) Write(addr uint16, val byte) {
	c.CPUWrite(addr, val)
}

type ppuMemory struct {
	*Console
}

func (p *ppuMemory) Read(addr uint16, vram []byte) byte {
	return p.PPURead(addr, vram)
}

func (p *ppuMemory) Write(addr uint16, val byte, vram []byte) {
	p.PPUWrite(addr, val, vram)
}
