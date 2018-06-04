package gophernes

import (
	"fmt"
	"github.com/sirupsen/logrus"
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
		case 0x16:
			return 0
		case 0x17:
			return 0x03
		default:
			logrus.Infof("Read from unhandled IO addr: %#X", addr)
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
		c.ppu.WriteReg(byte(addr&0x7), val)

	} else if addr >= 0x4000 && addr < 0x4020 {
		switch addr & 0x1F {
		case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0xA, 0xB, 0xC, 0xE, 0xF, 0x10, 0x11, 0x12, 0x13, 0x15:
			// APU
		case 0x16, 0x17:
			// Controllers
		case 0x14:
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

		default:
			logrus.Infof("Write to unhandled IO addr: %#X", addr)
		}

	} else if addr >= 0x4020 {
		c.cartridge.CPUWrite(addr, val)

	} else {
		panic(fmt.Sprintf("unhandled cpu memory write to address %#x", addr))
	}
}

func (c *Console) PPURead(addr uint16, vram []byte) byte {
	// Palette is a special case unaffected by the cartridge
	if addr >= 0x3F00 && addr <= 0x3FFF {
		return c.ppu.ReadPalette(addr % 0x20)
	}

	// The rest of PPU memory is custom mapped by the cartridge
	return c.cartridge.PPURead(addr, vram)
}

func (c *Console) PPUWrite(addr uint16, val byte, vram []byte) {
	if addr >= 0x3F00 && addr <= 0x3FFF {
		c.ppu.WritePalette(addr%0x20, val)
	} else {
		c.cartridge.PPUWrite(addr, val, vram)
	}
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
