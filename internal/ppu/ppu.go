package ppu

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

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
		config:   config,
		mem:      mem,
		vram:     make([]byte, internalVRAMSize),
		oam:      make([]byte, oamSize),
		scanLine: 261,
	}
	return ppu
}

type Memory interface {
	Read(addr uint16, vram []byte) byte
	Write(addr uint16, val byte, vram []byte)
	NMI()
}

type PPU struct {
	config *config
	cycles uint64

	mem  Memory
	vram []byte
	oam  []byte

	regs    Registers
	portBus byte

	spriteOverflow,
	sprite0Hit,
	nmiOccurred,
	nmiPrevious bool

	frames uint64
	scanLine,
	lineCycle int

	vramAddr,
	vramTempAddr uint16
	// addrLatch is used to funnel writes to the Scroll and Address registers
	addrLatch bool
}

func (p *PPU) Reset() {
	// TODO: Reset sequence
}

func (p *PPU) Frames() uint64 {
	return p.frames
}

func (p *PPU) Step() {
	// Perform the vertical blank handling as needed
	p.stepNMI()

	renderEnabled := p.regs.ShowBackground || p.regs.ShowSprites

	if renderEnabled {
		if p.lineCycle >= 1 && p.lineCycle <= 256 {
			switch p.lineCycle % 8 {
			case 1:
				p.readNametable()
			}
			if p.lineCycle == 256 {
				p.incrementY()
			} else if p.lineCycle%8 == 0 {
				p.incrementX()
			}

		} else if p.lineCycle == 257 {
			p.resetX()
		}

		if p.scanLine == 261 {
			if p.lineCycle >= 280 && p.lineCycle <= 304 {
				p.resetY()
			}
		}
	}

	// Move to the next cell/scanline as needed
	p.stepScan()

	p.cycles++
}

func (p *PPU) readNametable() {
	p.read8(p.vramAddr)
}

func (p *PPU) incrementX() {
	if p.vramAddr&0x1F == 31 {
		p.vramAddr &= ^uint16(0x1F)
		p.vramAddr ^= 0x400
	} else {
		p.vramAddr++
	}
}

func (p *PPU) incrementY() {
	// TODO: Make this cleaner?
	// http://wiki.nesdev.com/w/index.php/PPU_scrolling#Wrapping_around
	if p.vramAddr&0x7000 != 0x7000 {
		p.vramAddr += 0x1000
	} else {
		p.vramAddr &= ^uint16(0x7000)
		y := (p.vramAddr & 0x3E0) >> 5
		if y == 29 {
			y = 0
			p.vramAddr ^= 0x800
		} else if y == 31 {
			y = 0
		} else {
			y++
		}
		p.vramAddr = p.vramAddr & ^uint16(0x3E0) | (y << 5)
	}
}

func (p *PPU) resetX() {
	p.vramAddr = (p.vramAddr & 0xFBE0) | (p.vramTempAddr & 0x41F)
}

func (p *PPU) resetY() {
	p.vramAddr = (p.vramAddr & 0x841F) | (p.vramTempAddr & 0x7BE0)
}

func (p *PPU) stepNMI() {
	if p.scanLine == 241 && p.lineCycle == 1 {
		p.triggerNMI()
	}
	if p.scanLine == 261 && p.lineCycle == 1 {
		p.resetNMI()
		p.sprite0Hit = false
		p.spriteOverflow = false
	}

	shouldNMI := p.regs.NMIGenerate && p.nmiOccurred
	if shouldNMI && !p.nmiPrevious {
		if p.config.trace {
			logrus.Debug("PPU: Sending NMI to CPU")
		}
		p.mem.NMI()
	}
	p.nmiPrevious = shouldNMI
}

func (p *PPU) stepScan() {
	p.lineCycle++
	if p.lineCycle > 340 {
		p.lineCycle = 0
		p.scanLine++
		if p.scanLine > 261 {
			p.scanLine = 0
			p.frames++
			if p.frames%2 == 1 {
				p.lineCycle = 1
			}
		}
	}
}

func (p *PPU) triggerNMI() {
	p.nmiOccurred = true
}

func (p *PPU) resetNMI() {
	p.nmiOccurred = false
	p.nmiPrevious = false
}

func (p *PPU) ReadReg(reg byte) byte {
	val := p.portBus
	switch reg {
	case regStatus:
		// Take bits 0-4 from last written value
		val = val & 0x1F
		if p.spriteOverflow {
			val |= 1 << 5
		}
		if p.sprite0Hit {
			val |= 1 << 6
		}
		if p.nmiOccurred {
			val |= 1 << 7
		}
		// Resets NMI flag as a side effect
		p.resetNMI()
		// Resets address latch as a side effect
		p.addrLatch = false

	case regOAMData:
		val = p.oam[p.regs.OAMAddr]

	case regData:
		// TODO: Handle reads during renderEnable correctly?
		val = p.vram[p.regs.VRAMAddr]
		p.regs.VRAMAddr += p.regs.VRAMAddressIncrement
		p.regs.VRAMAddr %= internalVRAMSize

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
		p.oam[p.regs.OAMAddr] = val

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
		// TODO: Handle writes during renderEnable correctly?
		p.vram[p.regs.VRAMAddr] = val
		p.regs.VRAMAddr += p.regs.VRAMAddressIncrement

	default:
		panic(fmt.Sprintf("write to unknown PPU register %#x", reg))
	}
}

func (p *PPU) OAMDMA(src []byte) {
	// TODO: Break this up somehow so it happens over multiple clock cycles
	if len(src) != oamSize {
		panic("attempted to write OAM data with incorrect number of bytes")
	}
	for i := 0; i < oamSize; i++ {
		p.oam[(int(p.regs.OAMAddr)+i)%oamSize] = src[i]
	}
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
