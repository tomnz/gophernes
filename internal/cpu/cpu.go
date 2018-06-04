package cpu

import (
	"errors"

	"fmt"
	"github.com/sirupsen/logrus"
)

var ErrHalted = errors.New("cpu halted")

func NewCPU(mem Memory, opts ...Option) *CPU {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	cpu := &CPU{
		config:  config,
		opQueue: &opQueue{},
		mem:     mem,
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
	// Operations to perform for the next cycles - the next instruction is executed when
	// this is exhausted
	opQueue *opQueue
	insts   [256]*inst
	mem     Memory
	regs    Registers
	flags   Flags
	shouldNMI,
	shouldIRQ bool
	halted bool
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
		flags |= 1 << 7
	}
	if f.Overflow {
		flags |= 1 << 6
	}
	flags |= 1 << 5
	if f.BreakCmd {
		flags |= 1 << 4
	}
	if f.InterruptDisable {
		flags |= 1 << 2
	}
	if f.Zero {
		flags |= 1 << 1
	}
	if f.Carry {
		flags |= 1 << 0
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
	c.flags.Negative = (flags>>7)&1 == 1
	c.flags.Overflow = (flags>>6)&1 == 1
	c.flags.BreakCmd = (flags>>4)&1 == 1
	c.flags.InterruptDisable = (flags>>2)&1 == 1
	c.flags.Zero = (flags>>1)&1 == 1
	c.flags.Carry = (flags>>0)&1 == 1
}

const (
	nmiVector   = uint16(0xFFFA)
	resetVector = uint16(0xFFFC)
	irqVector   = uint16(0xFFFE)
)

func (c *CPU) Reset() {
	c.flags = Flags{
		InterruptDisable: true,
	}

	c.regs.StackPtr = 0xFD
	c.pc = c.read16(resetVector)
	c.halted = false

	if c.config.trace {
		logrus.Debugf("CPU: Reset to PC %#x", c.pc)
	}
}

func (c *CPU) RunTilHalt() uint64 {
	var cycles uint64
	for {
		if c.halted {
			return cycles
		}
		c.Step()
		cycles++
	}
}

func (c *CPU) Sleep(cycles uint64) {
	for i := uint64(0); i < cycles; i++ {
		c.opQueue.push(nil)
	}
}

func (c *CPU) NMI() {
	c.shouldNMI = true
}

func (c *CPU) IRQ() {
	if !c.flags.InterruptDisable {
		c.shouldIRQ = true
	}
}

func (c *CPU) Cycles() uint64 {
	return c.cycles
}

func (c *CPU) Step() {
	if c.halted {
		panic("cpu halted")
	}

	if c.opQueue.empty() {
		if c.shouldNMI {
			// TODO: Concurrent interrupt behavior
			c.nmi()
			c.Sleep(7)
			c.shouldNMI = false
		} else if c.shouldIRQ {
			c.irq()
			c.Sleep(7)
			c.shouldIRQ = false
		} else {
			opCode := c.prgRead8()
			inst := c.insts[opCode]

			if c.config.trace {
				// TODO: Better tracing! Let's store this as objects instead of logging
				c.trace(inst)
			}
			addr, cross := c.resolve(inst.addressMode)
			op := inst.op(c, addr, inst.addressMode)

			cycles := inst.cycles
			if inst.pageCrossCycle && cross {
				cycles++
			}

			if cycles > 0 {
				c.Sleep(cycles - 1)
			}
			// TODO: More advanced/correct behavior than just running the op at the end
			c.opQueue.push(op)
		}
	}

	nextOp := c.opQueue.pop()
	if nextOp != nil {
		nextOp()
	}
	c.cycles++
}

func (c *CPU) trace(inst *inst) {
	spaces := ""
	for i := 0xFF; i > int(c.regs.StackPtr); i-- {
		spaces += " "
	}

	fmt.Printf(
		// "CPU: c%d  A:%02X X:%02X Y:%02X S:%02X %s$%04X  %s\n",
		"A:%02X X:%02X Y:%02X S:%02X  %s$%04X %s\n",
		c.regs.Accumulator,
		c.regs.IndexX,
		c.regs.IndexY,
		c.regs.StackPtr,
		spaces,
		// Op read advanced the PC - decrement one for trace
		c.pc-1,
		inst.name,
	)
}

func (c *CPU) nmi() {
	c.stackPush16(c.pc)
	c.stackPush8(c.flags.asByte())
	c.flags.InterruptDisable = true
	c.pc = c.read16(nmiVector)
}

func (c *CPU) irq() {
	c.stackPush16(c.pc)
	c.stackPush8(c.flags.asByte())
	c.flags.InterruptDisable = true
	c.pc = c.read16(irqVector)
}

func (c *CPU) branch(offset uint16) {
	page := c.pc >> 8
	c.pc += offset
	if offset >= 0x80 {
		c.pc -= 0x100
	}
	// Special case - need to manually add cycles
	c.Sleep(1)
	if page != c.pc>>8 {
		// Page cross
		c.Sleep(2)
	}
}

func (c *CPU) compare(a, b byte) {
	diff := a - b
	c.flags.Carry = a >= b
	c.flags.Zero = diff == 0
	c.flags.Negative = (diff>>7)&1 == 1
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
	result := uint16(c.read8(addr))
	result |= uint16(c.read8(addr+1)) << 8
	return result
}

func (c *CPU) read16Wrap(addr uint16) (uint16, bool) {
	// Handle incorrect case where original 6502 wraps using high byte from the same page
	// http://obelisk.me.uk/6502/reference.html#JMP
	page := addr >> 8
	highAddr := addr + 1
	cross := page != highAddr>>8
	highAddr &= 0xFF
	highAddr |= page << 8

	lowVal := uint16(c.read8(addr))
	highVal := uint16(c.read8(highAddr))
	return highVal<<8 | lowVal, cross
}

func (c *CPU) write8(addr uint16, val byte) {
	c.mem.Write(addr, val)
}

func (c *CPU) write16(addr, val uint16) {
	c.write8(addr, byte(val&0xFF))
	c.write8(addr+1, byte(val>>8))
}

func (c *CPU) stackPush8(val byte) {
	stackAddr := uint16(c.regs.StackPtr)
	stackAddr |= 0x100
	c.write8(stackAddr, val)
	c.regs.StackPtr--
}

func (c *CPU) stackPull8() byte {
	c.regs.StackPtr++
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
