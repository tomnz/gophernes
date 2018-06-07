package ppu

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	internalVRAMSize = 0x1000
	oamSize          = 0x100
	DisplayWidth     = 256
	DisplayHeight    = 240
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

	mem         Memory
	vram        []byte
	oam         []byte
	paletteData [32]byte

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
	fineX byte
	// addrLatch is used to funnel writes to the Scroll and Address registers
	addrLatch    bool
	bufferedData byte

	// Background data
	nameTableByte,
	attributeTableByte,
	lowTileByte,
	highTileByte byte
	tileData uint64

	// Sprite data
	spriteCount      int
	spritePatterns   [8]uint32
	spritePositions  [8]byte
	spritePriorities [8]byte
	spriteIndexes    [8]byte

	// Display buffers
	backBuffer  [DisplayHeight][DisplayWidth]byte
	frontBuffer [DisplayHeight][DisplayWidth]byte
}

func (p *PPU) Reset() {
	p.lineCycle = 340
	p.scanLine = 240
	p.frames = 0
	// TODO: Reset registers
}

func (p *PPU) Frames() uint64 {
	return p.frames
}

func (p *PPU) Buffer() [DisplayHeight][DisplayWidth]byte {
	return p.frontBuffer
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
		p.triggerVertBlank()
	}
	if p.scanLine == 261 && p.lineCycle == 1 {
		p.clearVertBlank()
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

func (p *PPU) triggerVertBlank() {
	p.nmiOccurred = true
	p.frontBuffer, p.backBuffer = p.backBuffer, p.frontBuffer
}

func (p *PPU) clearVertBlank() {
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
		p.clearVertBlank()
		// Resets address latch as a side effect
		p.addrLatch = false

	case regOAMData:
		val = p.oam[p.regs.OAMAddr]

	case regData:
		// TODO: Handle reads during renderEnable correctly?
		val = p.read8(p.vramAddr)
		if p.vramAddr%0x4000 < 0x3F00 {
			buffered := p.bufferedData
			p.bufferedData = val
			val = buffered
		} else {
			p.bufferedData = p.read8(p.vramAddr - 0x1000)
		}
		p.vramAddr += p.regs.VRAMAddressIncrement
		// p.vramAddr %= internalVRAMSize

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
		p.regs.OAMAddr++

	case regScroll:
		if !p.addrLatch {
			p.vramTempAddr = (p.vramTempAddr & 0xFFE0) | (uint16(val) >> 3)
			p.fineX = val & 0x7
		} else {
			// TODO: Handle values higher than 239 correctly
			p.vramTempAddr = (p.vramTempAddr & 0x8FFF) | ((uint16(val) & 0x7) << 12)
			p.vramTempAddr = (p.vramTempAddr & 0xFC1F) | ((uint16(val) & 0xF8) << 2)
		}
		p.addrLatch = !p.addrLatch

	case regAddress:
		if !p.addrLatch {
			p.vramTempAddr = (p.vramTempAddr & 0x80FF) | ((uint16(val) & 0x3F) << 8)
		} else {
			p.vramTempAddr = (p.vramTempAddr & 0xFF00) | uint16(val)
			p.vramAddr = p.vramTempAddr
		}
		p.addrLatch = !p.addrLatch

	case regData:
		// TODO: Handle writes during renderEnable correctly?
		p.write8(p.vramAddr, val)
		p.vramAddr += p.regs.VRAMAddressIncrement

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
