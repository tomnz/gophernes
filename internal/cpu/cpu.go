package cpu

import (
	"log"
)

func NewCPU(mem Memory, opts ...Option) *CPU {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	cpu := &CPU{
		config: config,
		mem:    mem,
	}
	cpu.initInstructions()
	return cpu
}

type Memory interface {
	Read(addr uint16) byte
	Write(addr uint16, val byte)
}

// CPU is the main CPU of the NES system.
type CPU struct {
	config *config
	pc     uint16
	cycles uint64
	insts  [256]*inst
	mem    Memory
	regs   registers
	flags  flags
}

type registers struct {
	StackPtr,
	Accumulator,
	IndexX,
	IndexY byte
}

type flags struct {
	Negative,
	Overflow,
	BreakCmd,
	InterruptDisable,
	Zero,
	Carry bool
}

const (
	resetVector = uint16(0xFFFE)
)

func (c *CPU) Reset() {
	c.flags = flags{
		BreakCmd:         true,
		InterruptDisable: true,
	}

	c.pc = c.read16(resetVector)

	if c.config.trace {
		log.Printf("Reset CPU to PC %#x", c.pc)
	}
}

func (c *CPU) Run(instructions int) {
	for i := 0; i < instructions; i++ {
		c.doInstruction()
	}
}

func (c *CPU) doInstruction() {
	opCode := c.prgRead8()
	inst := c.insts[opCode]

	addr, cross := c.resolve(inst.addressMode)
	inst.op(addr)

	cycles := inst.cycles
	if inst.pageCrossCycle && cross {
		cycles++
	}
	if c.config.trace {
		log.Printf("Op: %s | (%d)", inst.fullName(), cycles)
	}
	c.cycles += cycles
}

func (c *CPU) prgRead8() byte {
	result := c.read8(c.pc)
	c.pc++
	return result
}

func (c *CPU) prgRead16() uint16 {
	result := c.read16(c.pc)
	c.pc += 2
	return result
}

func (c *CPU) read8(addr uint16) byte {
	return c.mem.Read(addr)
}

func (c *CPU) read16(addr uint16) uint16 {
	result := uint16(c.mem.Read(addr))
	result |= uint16(c.mem.Read(addr+1)) << 8
	return result
}

func (c *CPU) write8(addr uint16, val byte) {
	c.mem.Write(addr, val)
}

func (c *CPU) write16(addr uint16, val uint16) {
	c.mem.Write(addr, byte(val&0xFF))
	c.mem.Write(addr+1, byte(val>>8))
}
