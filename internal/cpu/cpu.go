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
	regs   Registers
	flags  Flags
}

type Registers struct {
	StackPtr,
	Accumulator,
	IndexX,
	IndexY byte
}

type Flags struct {
	Negative,
	Overflow,
	BreakCmd,
	InterruptDisable,
	Zero,
	Carry bool
}

func (c *CPU) Registers() Registers {
	return c.regs
}

func (c *CPU) Flags() Flags {
	return c.flags
}

const resetVector = uint16(0xFFFE)

func (c *CPU) Reset() {
	c.flags = Flags{
		BreakCmd:         true,
		InterruptDisable: true,
	}

	c.pc = c.read16(resetVector)

	if c.config.trace {
		log.Printf("Reset CPU to PC %#x", c.pc)
	}
}

func (c *CPU) Run(steps int) uint64 {
	var cycles uint64
	for i := 0; i < steps; i++ {
		cycles += c.step()
	}
	return cycles
}

func (c *CPU) step() uint64 {
	opCode := c.prgRead8()
	inst := c.insts[opCode]

	addr, cross := c.resolve(inst.addressMode)
	inst.op(c, addr)

	cycles := inst.cycles
	if inst.pageCrossCycle && cross {
		cycles++
	}
	if c.config.trace {
		// TODO: Better tracing! Let's store this as objects instead of logging
		log.Printf("Op: %s | (%d)", inst.fullName(), cycles)
	}
	c.cycles += cycles
	return cycles
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
