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
		page := addr >> 8
		addr += uint16(c.regs.IndexX)
		return addr, page != addr>>8

	case AddressAbsoluteY:
		addr := c.prgRead16()
		page := addr >> 8
		addr += uint16(c.regs.IndexY)
		return addr, page != addr>>8

	case AddressIndirect:
		return c.read16Wrap(c.prgRead16())

	case AddressIndirectX:
		addr := uint16(c.prgRead8())
		addr += uint16(c.regs.IndexX)
		// Wrap around if we overflow the first page
		// addr &= 0xFF
		return c.read16Wrap(addr)

	case AddressIndirectY:
		addr, _ := c.read16Wrap(uint16(c.prgRead8()))
		page := addr >> 8
		addr += uint16(c.regs.IndexY)
		return addr, page != addr>>8
	}
	panic(fmt.Sprintf("couldn't resolve for address mode %d", mode))
}
