package cpu

import (
	"context"
	"errors"

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
		opQueue: make([]func(), 0),
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
	opQueue []func()
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
	nmiVector   = uint16(0xFFFA)
	resetVector = uint16(0xFFFC)
	irqVector   = uint16(0xFFFE)
)

func (c *CPU) Reset() {
	c.flags = Flags{
		BreakCmd:         true,
		InterruptDisable: true,
	}

	c.pc = c.read16(resetVector)
	c.halted = false

	if c.config.trace {
		logrus.Debugf("CPU: Reset to PC %#x", c.pc)
	}
}

const clockDivisor = 12

func (c *CPU) Run(ctx context.Context, clock <-chan struct{}) {
	subCycles := 0
	cycles := uint64(0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-clock:
			subCycles++
			if subCycles >= clockDivisor {
				subCycles = 0
				cycles++
				_, err := c.step()
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func (c *CPU) RunTilHalt() (uint64, error) {
	var cycles uint64
	for {
		stepCycles, err := c.step()
		cycles += stepCycles
		if err != nil {
			if err == ErrHalted {
				return cycles, nil
			}
			return cycles, err
		}
	}
}

func (c *CPU) Sleep(cycles uint64) {
	for i := uint64(0); i < cycles; i++ {
		c.opQueue = append(c.opQueue, nil)
	}
}

func (c *CPU) NMI() {
	c.shouldNMI = true
}

func (c *CPU) IRQ() {
	c.shouldIRQ = true
}

func (c *CPU) Cycles() uint64 {
	return c.cycles
}

func (c *CPU) step() (uint64, error) {
	if c.halted {
		return 0, ErrHalted
	}
	defer func() { c.cycles++ }()

	if len(c.opQueue) == 0 {
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
				logrus.Debugf("CPU: PC: %#x", c.pc)
				logrus.Debugf("CPU: Op: %s", inst.fullName())
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
			c.opQueue = append(c.opQueue, op)
		}
	}

	nextOp := c.opQueue[0]
	if nextOp != nil {
		nextOp()
	}
	c.opQueue = c.opQueue[1:]
	return 1, nil
}

func (c *CPU) nmi() {
	c.stackPush16(c.pc - 1)
	c.stackPush8(c.flags.asByte())
	c.pc = c.read16(nmiVector)
}

func (c *CPU) irq() {
	c.stackPush16(c.pc - 1)
	c.stackPush8(c.flags.asByte())
	c.pc = c.read16(irqVector)
}

func (c *CPU) branch(offset uint16) {
	page := c.pc >> 2
	if offset < 0x80 {
		c.pc += offset
	} else {
		c.pc += offset - 0x100
	}
	// Special case - need to manually add cycles
	c.Sleep(1)
	if page != c.pc>>2 {
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
