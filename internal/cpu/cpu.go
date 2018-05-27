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

func (f Flags) asByte() byte {
	var flags byte
	if f.Negative {
		flags |= 0x1 << 6
	}
	if f.Overflow {
		flags |= 0x1 << 5
	}
	if f.BreakCmd {
		flags |= 0x1 << 4
	}
	if f.InterruptDisable {
		flags |= 0x1 << 2
	}
	if f.Zero {
		flags |= 0x1 << 1
	}
	if f.Carry {
		flags |= 0x1 << 0
	}
	return flags
}

func (c *CPU) Registers() Registers {
	return c.regs
}

func (c *CPU) Flags() Flags {
	return c.flags
}

func (c *CPU) setFlagsFromByte(flags byte) {
	c.flags.Negative = (flags>>6)&1 == 1
	c.flags.Overflow = (flags>>5)&1 == 1
	c.flags.BreakCmd = (flags>>4)&1 == 1
	c.flags.InterruptDisable = (flags>>1)&1 == 1
	c.flags.Zero = (flags>>1)&1 == 1
	c.flags.Carry = (flags>>0)&1 == 1
}

const (
	resetVector = uint16(0xFFFC)
	irqVector   = uint16(0xFFFE)
)

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
	inst.op(c, addr, inst.addressMode)

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

func (c *CPU) branch(offset uint16) {
	page := c.pc >> 2
	c.pc += offset
	// Special case - need to manually add cycles
	c.cycles++
	if page != c.pc>>2 {
		// Page cross
		c.cycles += 2
	}
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

func (c *CPU) stackPush8(val byte) {
	stackAddr := uint16(c.regs.StackPtr)
	stackAddr |= 0x100
	c.write8(stackAddr, val)
	c.regs.StackPtr++
}

func (c *CPU) stackPull8() byte {
	c.regs.StackPtr--
	stackAddr := uint16(c.regs.StackPtr)
	stackAddr |= 0x100
	return c.read8(stackAddr)
}

func (c *CPU) stackPush16(val uint16) {
	c.stackPush8(byte(val >> 8))
	c.stackPush8(byte(val & 0xFF))
}

func (c *CPU) stackPull16() uint16 {
	val := uint16(c.stackPull8())
	val |= uint16(c.stackPull8()) << 8
	return val
}
