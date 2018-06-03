package ppu

import "fmt"

const (
	internalVRAMSize = 0x800
	oamSize          = 0x100
)

func NewPPU(mem Memory, opts ...Option) *PPU {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	ppu := &PPU{
		config: config,
		mem:    mem,
		vram:   make([]byte, internalVRAMSize),
		oam:    make([]byte, oamSize),
	}
	return ppu
}

type Memory interface {
	Read(addr uint16, vram []byte) byte
	Write(addr uint16, val byte, vram []byte)
}

type PPU struct {
	config  *config
	cycles  uint64
	mem     Memory
	vram    []byte
	oam     []byte
	regs    Registers
	flags   flags
	portBus byte
	// addrLatch is used to funnel writes to the Scroll and Address registers
	addrLatch bool
	rendering bool
}

type flags struct {
	spriteOverflow bool
	sprite0Hit     bool
}

func (p *PPU) ReadReg(reg byte) byte {
	val := p.portBus
	switch reg {
	case regStatus:
		// Take bits 0-4 from last written value
		val = val & 0x1F
		if p.flags.spriteOverflow {
			val |= 1 << 5
		}
		if p.flags.sprite0Hit {
			val |= 1 << 6
		}
		if p.regs.NMIOccurred {
			val |= 1 << 7
		}
		// Resets NMI flag as a side effect
		p.regs.NMIOccurred = false
		// Resets address latch as a side effect
		p.addrLatch = false

	case regOAMData:
		val = p.oam[p.regs.OAMAddr]

	case regData:
		// TODO: Handle reads during rendering correctly?
		val = p.vram[p.regs.VRAMAddr]
		p.regs.VRAMAddr += p.regs.VRAMAddressIncrement

	// Write-only registers just return the current bus value
	case regController, regMask, regOAMAddress, regScroll, regAddress:
		val = p.portBus

	default:
		panic(fmt.Sprintf("read from unknown PPU register %#x", reg))
	}

	p.portBus = val
	return val
}

func (p *PPU) WriteReg(reg byte, val byte) {
	p.portBus = val
	switch reg {
	case regController:
		switch val & 0x3 {
		case 0:
			p.regs.BaseNametableAddress = addrNametable0
		case 1:
			p.regs.BaseNametableAddress = addrNametable1
		case 2:
			p.regs.BaseNametableAddress = addrNametable2
		case 3:
			p.regs.BaseNametableAddress = addrNametable3
		}

		if (val>>2)&1 == 1 {
			p.regs.VRAMAddressIncrement = 32
		} else {
			p.regs.VRAMAddressIncrement = 1
		}

		if (val>>3)&1 == 1 {
			p.regs.SpritePatternTableAddress = addrPatternTable1
		} else {
			p.regs.SpritePatternTableAddress = addrPatternTable0
		}

		if (val>>4)&1 == 1 {
			p.regs.BackgroundPatternTableAddress = addrPatternTable1
		} else {
			p.regs.BackgroundPatternTableAddress = addrPatternTable0
		}

		p.regs.TallSprites = (val>>5)&1 == 1

		p.regs.NMIGenerate = (val>>7)&1 == 1

	case regMask:
		p.regs.Grayscale = val&1 == 1
		p.regs.ShowLeftBackground = (val>>1)&1 == 1
		p.regs.ShowLeftSprites = (val>>2)&1 == 1
		p.regs.ShowBackground = (val>>3)&1 == 1
		p.regs.ShowSprites = (val>>4)&1 == 1
		p.regs.EmphasizeRed = (val>>5)&1 == 1
		p.regs.EmphasizeGreen = (val>>6)&1 == 1
		p.regs.EmphasizeBlue = (val>>7)&1 == 1

	case regOAMAddress:
		p.regs.OAMAddr = val

	case regOAMData:
		// TODO: Handle "glitchy" writes during rendering?
		// http://wiki.nesdev.com/w/index.php/PPU_registers
		if !p.rendering {
			p.oam[p.regs.OAMAddr] = val
		}

	case regScroll:
		if !p.addrLatch {
			p.regs.ScrollX = val
		} else {
			// TODO: Handle values higher than 239 correctly
			p.regs.ScrollY = val
		}
		p.addrLatch = !p.addrLatch

	case regAddress:
		if !p.addrLatch {
			p.regs.VRAMAddr = (uint16(val) << 8) | (p.regs.VRAMAddr & 0xFF)
		} else {
			p.regs.VRAMAddr = uint16(val) | (p.regs.VRAMAddr & 0xFF00)
		}
		p.addrLatch = !p.addrLatch

	case regData:
		// TODO: Handle writes during rendering correctly?
		p.vram[p.regs.VRAMAddr] = val
		p.regs.VRAMAddr += p.regs.VRAMAddressIncrement

	default:
		panic(fmt.Sprintf("write to unknown PPU register %#x", reg))
	}
}

func (p *PPU) CopyOAM(src []byte) {
	if len(src) != oamSize {
		panic("attempted to write OAM data with incorrect number of bytes")
	}
	copy(p.oam, src)
}

func (p *PPU) read8(addr uint16) byte {
	return p.mem.Read(addr, p.vram)
}

func (p *PPU) read16(addr uint16) uint16 {
	result := uint16(p.read8(addr))
	result |= uint16(p.read8(addr+1)) << 8
	return result
}

func (p *PPU) write8(addr uint16, val byte) {
	p.mem.Write(addr, val, p.vram)
}

func (p *PPU) write16(addr, val uint16) {
	p.write8(addr, byte(val&0xFF))
	p.write8(addr+1, byte(val>>8))
}