package cpu

import "fmt"

type addressMode byte

const (
	addressUnknown addressMode = iota
	addressImplicit
	addressImmediate
	addressAccumulator
	addressZeroPage
	addressZeroPageX
	addressZeroPageY
	addressRelative
	addressAbsolute
	addressAbsoluteX
	addressAbsoluteY
	addressIndirect
	addressIndexedIndirect
	addressIndirectIndexed
)

func (a addressMode) name() string {
	switch a {
	case addressUnknown:
		return "u"
	case addressImplicit:
		return "imp"
	case addressImmediate:
		return "imm"
	case addressZeroPage:
		return "zp"
	case addressZeroPageX:
		return "zpx"
	case addressZeroPageY:
		return "zpy"
	case addressRelative:
		return "rel"
	case addressAbsolute:
		return "abs"
	case addressAbsoluteX:
		return "absx"
	case addressAbsoluteY:
		return "absy"
	case addressIndirect:
		return "ind"
	case addressIndexedIndirect:
		return "indx"
	case addressIndirectIndexed:
		return "indy"
	}
	panic(fmt.Sprintf("unknown address mode %d", a))
}

// resolve resolves the given address using the specified address mode.
// The second return value is used to indicate true if a page was crossed, which can add CPU
// cycles in some cases.
func (c *CPU) resolve(mode addressMode) (uint16, bool) {
	switch mode {
	case addressImplicit:
		// Implicit ops should ignore the address
		return 0, false

	case addressImmediate:
		addr := c.pc
		c.pc++
		return addr, false

	case addressZeroPage:
		return uint16(c.prgRead8()), false

	case addressZeroPageX:
		addr := uint16(c.prgRead8())
		addr += uint16(c.regs.IndexX)
		// Wrap around if we overflow the first page
		addr &= 0xff
		return addr, false

	case addressZeroPageY:
		addr := uint16(c.prgRead8())
		addr += uint16(c.regs.IndexY)
		addr &= 0xff
		return addr, false

	case addressAbsolute:
		return c.prgRead16(), false

	case addressAbsoluteX:
		addr := c.prgRead16()
		page := addr >> 2
		addr += uint16(c.regs.IndexX)
		return addr, page != addr>>2

	case addressAbsoluteY:
		addr := c.prgRead16()
		page := addr >> 2
		addr += uint16(c.regs.IndexY)
		return addr, page != addr>>2
	}
	panic(fmt.Sprintf("couldn't resolve for address mode %d", mode))
}
