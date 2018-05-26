package cpu

import "fmt"

type AddressMode byte

const (
	AddressUnknown AddressMode = iota
	AddressImplicit
	AddressImmediate
	AddressAccumulator
	AddressZeroPage
	AddressZeroPageX
	AddressZeroPageY
	AddressRelative
	AddressAbsolute
	AddressAbsoluteX
	AddressAbsoluteY
	AddressIndirect
	AddressIndirectX
	AddressIndirectY
)

func (a AddressMode) name() string {
	switch a {
	case AddressUnknown:
		return "???"
	case AddressImplicit:
		return "imp"
	case AddressImmediate:
		return "imm"
	case AddressAccumulator:
		return "acc"
	case AddressZeroPage:
		return "zpg"
	case AddressZeroPageX:
		return "zpx"
	case AddressZeroPageY:
		return "zpy"
	case AddressRelative:
		return "rel"
	case AddressAbsolute:
		return "abs"
	case AddressAbsoluteX:
		return "abx"
	case AddressAbsoluteY:
		return "aby"
	case AddressIndirect:
		return "ind"
	case AddressIndirectX:
		return "idx"
	case AddressIndirectY:
		return "idy"
	}
	panic(fmt.Sprintf("unknown address mode %d", a))
}

// resolve resolves the given address using the specified address mode.
// The second return value is used to indicate true if a page was crossed, which can add CPU
// cycles in some cases.
func (c *CPU) resolve(mode AddressMode) (uint16, bool) {
	switch mode {
	case AddressImplicit:
		// Implicit ops should ignore the address
		return 0, false

	case AddressImmediate:
		addr := c.pc
		c.pc++
		return addr, false

	case AddressAccumulator:
		// Accumulator ops should ignore the address
		return 0, false

	case AddressZeroPage:
		return uint16(c.prgRead8()), false

	case AddressZeroPageX:
		addr := uint16(c.prgRead8())
		addr += uint16(c.regs.IndexX)
		// Wrap around if we overflow the first page
		addr &= 0xFF
		return addr, false

	case AddressZeroPageY:
		addr := uint16(c.prgRead8())
		addr += uint16(c.regs.IndexY)
		addr &= 0xFF
		return addr, false

	case AddressRelative:
		// Return the relative (signed) offset - ops should handle this appropriately
		return uint16(c.prgRead8()), false

	case AddressAbsolute:
		return c.prgRead16(), false

	case AddressAbsoluteX:
		addr := c.prgRead16()
		page := addr >> 2
		addr += uint16(c.regs.IndexX)
		return addr, page != addr>>2

	case AddressAbsoluteY:
		addr := c.prgRead16()
		page := addr >> 2
		addr += uint16(c.regs.IndexY)
		return addr, page != addr>>2

	case AddressIndirect:
		// TODO: Handle incorrect case where original 6502 wraps high byte from the same page?
		// http://obelisk.me.uk/6502/reference.html#JMP
		return c.read16(c.prgRead16()), false

	case AddressIndirectX:
		addr := uint16(c.prgRead8())
		addr += uint16(c.regs.IndexX)
		// Wrap around if we overflow the first page
		addr &= 0xFF
		return c.read16(addr), false

	case AddressIndirectY:
		addr := c.read16(uint16(c.prgRead8()))
		page := addr >> 2
		addr += uint16(c.regs.IndexY)
		return c.read16(addr), page != addr>>2
	}
	panic(fmt.Sprintf("couldn't resolve for address mode %d", mode))
}
