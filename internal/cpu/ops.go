package cpu

import "github.com/sirupsen/logrus"

type op func(cpu *CPU, addr uint16, mode AddressMode) func()

type inst struct {
	name           string
	cycles         uint64
	op             op
	addressMode    AddressMode
	pageCrossCycle bool
}

func (i *inst) fullName() string {
	return i.name + " " + i.addressMode.name()
}

var insts = [256]*inst{
	// 00
	{"BRK", 7, brk, AddressImplicit, false},
	{"ORA", 6, ora, AddressIndirectX, false},
	instHalt(),
	{"TBD", 8, tbd, AddressIndirectX, false},
	{"NOP", 3, nop, AddressZeroPage, false},
	{"ORA", 3, ora, AddressZeroPage, false},
	{"ASL", 5, asl, AddressZeroPage, false},
	{"TBD", 5, tbd, AddressZeroPage, false},
	{"PHP", 3, php, AddressImplicit, false},
	{"ORA", 2, ora, AddressImmediate, false},
	{"ASL", 2, asl, AddressAccumulator, false},
	{"TBD", 2, tbd, AddressImmediate, false},
	{"NOP", 4, nop, AddressAbsolute, false},
	{"ORA", 4, ora, AddressAbsolute, false},
	{"ASL", 6, asl, AddressAbsolute, false},
	{"TBD", 6, tbd, AddressAbsolute, false},

	// 10
	{"BPL", 2, bpl, AddressRelative, false},
	{"ORA", 5, ora, AddressIndirectY, true},
	instHalt(),
	{"TBD", 8, tbd, AddressIndirectY, false},
	{"NOP", 4, nop, AddressZeroPageX, false},
	{"ORA", 4, ora, AddressZeroPageX, false},
	{"ASL", 6, asl, AddressZeroPageX, false},
	{"TBD", 6, tbd, AddressZeroPageX, false},
	{"CLC", 2, clc, AddressImplicit, false},
	{"ORA", 4, ora, AddressAbsoluteY, true},
	{"NOP", 2, nop, AddressImplicit, false},
	{"TBD", 7, tbd, AddressAbsoluteY, false},
	{"NOP", 4, nop, AddressAbsoluteX, true},
	{"ORA", 4, ora, AddressAbsoluteX, true},
	{"ASL", 7, asl, AddressAbsoluteX, false},
	{"TBD", 7, tbd, AddressAbsoluteX, false},

	// 20
	{"JSR", 6, jsr, AddressAbsolute, false},
	{"AND", 6, and, AddressIndirectX, false},
	instHalt(),
	{"TBD", 8, tbd, AddressIndirectX, false},
	{"BIT", 3, bit, AddressZeroPage, false},
	{"AND", 3, and, AddressZeroPage, false},
	{"ROL", 5, rol, AddressZeroPage, false},
	{"TBD", 5, tbd, AddressZeroPage, false},
	{"PLP", 4, plp, AddressImplicit, false},
	{"AND", 2, and, AddressImmediate, false},
	{"ROL", 2, rol, AddressAccumulator, false},
	{"TBD", 2, tbd, AddressImmediate, false},
	{"BIT", 4, bit, AddressAbsolute, false},
	{"AND", 4, and, AddressAbsolute, false},
	{"ROL", 6, rol, AddressAbsolute, false},
	{"TBD", 6, tbd, AddressAbsolute, false},

	// 30
	{"BMI", 2, bmi, AddressRelative, false},
	{"AND", 5, and, AddressIndirectY, true},
	instHalt(),
	{"TBD", 8, tbd, AddressIndirectY, false},
	{"NOP", 4, nop, AddressZeroPageX, false},
	{"AND", 4, and, AddressZeroPageX, false},
	{"ROL", 6, rol, AddressZeroPageX, false},
	{"TBD", 6, tbd, AddressZeroPageX, false},
	{"SEC", 2, sec, AddressImplicit, false},
	{"AND", 4, and, AddressAbsoluteY, true},
	{"NOP", 2, nop, AddressImplicit, false},
	{"TBD", 7, tbd, AddressAbsoluteY, false},
	{"NOP", 4, nop, AddressAbsoluteX, true},
	{"AND", 4, and, AddressAbsoluteX, true},
	{"ROL", 7, rol, AddressAbsoluteX, false},
	{"TBD", 7, tbd, AddressAbsoluteX, false},

	// 40
	{"RTI", 6, rti, AddressImplicit, false},
	{"EOR", 6, eor, AddressIndirectX, false},
	instHalt(),
	{"TBD", 8, tbd, AddressIndirectX, false},
	{"NOP", 3, nop, AddressZeroPage, false},
	{"EOR", 3, eor, AddressZeroPage, false},
	{"LSR", 5, lsr, AddressZeroPage, false},
	{"TBD", 5, tbd, AddressZeroPage, false},
	{"PHA", 3, pha, AddressImplicit, false},
	{"EOR", 2, eor, AddressImmediate, false},
	{"LSR", 2, lsr, AddressAccumulator, false},
	{"ALR", 2, alr, AddressImmediate, false},
	{"JMP", 3, jmp, AddressAbsolute, false},
	{"EOR", 4, eor, AddressAbsolute, false},
	{"LSR", 6, lsr, AddressAbsolute, false},
	{"TBD", 6, tbd, AddressAbsolute, false},

	// 50
	{"BVC", 2, bvc, AddressRelative, false},
	{"EOR", 5, eor, AddressIndirectX, true},
	instHalt(),
	{"TBD", 8, tbd, AddressIndirectY, false},
	{"NOP", 4, nop, AddressZeroPageX, false},
	{"EOR", 4, eor, AddressZeroPageX, false},
	{"LSR", 6, lsr, AddressZeroPageX, false},
	{"TBD", 6, tbd, AddressZeroPageX, false},
	{"CLI", 2, cli, AddressImplicit, false},
	{"EOR", 4, eor, AddressAbsoluteY, true},
	{"NOP", 2, nop, AddressImplicit, false},
	{"TBD", 7, tbd, AddressAbsoluteY, false},
	{"NOP", 4, nop, AddressAbsoluteX, true},
	{"EOR", 4, eor, AddressAbsoluteX, true},
	{"LSR", 7, lsr, AddressAbsoluteX, false},
	{"TBD", 7, tbd, AddressAbsoluteX, false},

	// 60
	{"RTS", 6, rts, AddressImplicit, false},
	{"ADC", 6, adc, AddressIndirectX, false},
	instHalt(),
	{"TBD", 8, tbd, AddressIndirectX, false},
	{"NOP", 3, nop, AddressZeroPage, false},
	{"ADC", 3, adc, AddressZeroPage, false},
	{"ROR", 5, ror, AddressZeroPage, false},
	{"TBD", 5, tbd, AddressZeroPage, false},
	{"PLA", 4, pla, AddressImplicit, false},
	{"ADC", 2, adc, AddressImmediate, false},
	{"ROR", 2, ror, AddressAccumulator, false},
	{"TBD", 2, tbd, AddressImmediate, false},
	{"JMP", 5, jmp, AddressIndirect, false},
	{"ADC", 4, adc, AddressAbsolute, false},
	{"ROR", 6, ror, AddressAbsolute, false},
	{"TBD", 6, tbd, AddressAbsoluteX, false},

	// 70
	{"BVS", 2, bvs, AddressRelative, false},
	{"ADC", 5, adc, AddressIndirectY, true},
	instHalt(),
	{"TBD", 8, tbd, AddressIndirectY, false},
	{"NOP", 4, nop, AddressZeroPageX, false},
	{"ADC", 4, adc, AddressZeroPageX, false},
	{"ROR", 6, ror, AddressZeroPageX, false},
	{"TBD", 6, tbd, AddressZeroPageX, false},
	{"SEI", 2, sei, AddressImplicit, false},
	{"ADC", 4, adc, AddressAbsoluteY, true},
	{"NOP", 2, nop, AddressImplicit, false},
	{"TBD", 7, tbd, AddressAbsoluteY, false},
	{"NOP", 4, nop, AddressAbsoluteX, true},
	{"ADC", 4, adc, AddressAbsoluteX, true},
	{"ROR", 7, ror, AddressAbsoluteX, false},
	{"TBD", 7, tbd, AddressAbsoluteX, false},

	// 80
	{"NOP", 2, nop, AddressImmediate, false},
	{"STA", 6, sta, AddressIndirectX, false},
	{"NOP", 2, nop, AddressImmediate, false},
	{"TBD", 6, tbd, AddressIndirectX, false},
	{"STY", 3, sty, AddressZeroPage, false},
	{"STA", 3, sta, AddressZeroPage, false},
	{"STX", 3, stx, AddressZeroPage, false},
	{"TBD", 3, tbd, AddressZeroPage, false},
	{"DEY", 2, dey, AddressImplicit, false},
	{"NOP", 2, nop, AddressImmediate, false},
	{"TXA", 2, txa, AddressImplicit, false},
	{"TBD", 2, tbd, AddressImmediate, false},
	{"STY", 4, sty, AddressAbsolute, false},
	{"STA", 4, sta, AddressAbsolute, false},
	{"STX", 4, stx, AddressAbsolute, false},
	{"TBD", 4, tbd, AddressAbsoluteX, false},

	// 90
	{"BCC", 2, bcc, AddressRelative, false},
	{"STA", 6, sta, AddressIndirectY, false},
	instHalt(),
	{"TBD", 6, tbd, AddressIndirectY, false},
	{"STY", 4, sty, AddressZeroPageX, false},
	{"STA", 4, sta, AddressZeroPageX, false},
	{"STX", 4, stx, AddressZeroPageY, false},
	{"TBD", 4, tbd, AddressZeroPageY, false},
	{"TYA", 2, tya, AddressImplicit, false},
	{"STA", 5, sta, AddressAbsoluteY, false},
	{"TXS", 2, txs, AddressImplicit, false},
	{"TBD", 5, tbd, AddressAbsoluteY, false},
	{"TBD", 5, tbd, AddressAbsoluteX, false},
	{"STA", 5, sta, AddressAbsoluteX, false},
	{"TBD", 5, tbd, AddressAbsoluteY, false},
	{"TBD", 5, tbd, AddressAbsoluteY, false},

	// A0
	{"LDY", 2, ldy, AddressImmediate, false},
	{"LDA", 6, lda, AddressIndirectX, false},
	{"LDX", 2, ldx, AddressImmediate, false},
	{"TBD", 6, tbd, AddressIndirectX, false},
	{"LDY", 3, ldy, AddressZeroPage, false},
	{"LDA", 3, lda, AddressZeroPage, false},
	{"LDX", 3, ldx, AddressZeroPage, false},
	{"TBD", 3, tbd, AddressZeroPage, false},
	{"TAY", 2, tay, AddressImplicit, false},
	{"LDA", 2, lda, AddressImmediate, false},
	{"TAX", 2, tax, AddressImplicit, false},
	{"TBD", 2, tbd, AddressImmediate, false},
	{"LDY", 4, ldy, AddressAbsolute, false},
	{"LDA", 4, lda, AddressAbsolute, false},
	{"LDX", 4, ldx, AddressAbsolute, false},
	{"TBD", 4, tbd, AddressAbsolute, false},

	// B0
	{"BCS", 2, bcs, AddressRelative, false},
	{"LDA", 5, lda, AddressIndirectY, true},
	instHalt(),
	{"TBD", 5, tbd, AddressIndirectY, true},
	{"LDY", 4, ldy, AddressZeroPageX, false},
	{"LDA", 4, lda, AddressZeroPageX, false},
	{"LDX", 4, ldx, AddressZeroPageY, false},
	{"TBD", 4, tbd, AddressZeroPageY, false},
	{"CLV", 2, clv, AddressImplicit, false},
	{"LDA", 4, lda, AddressAbsoluteY, true},
	{"TSX", 2, tsx, AddressImplicit, false},
	{"TBD", 4, tbd, AddressAbsoluteY, true},
	{"LDY", 4, ldy, AddressAbsoluteX, true},
	{"LDA", 4, lda, AddressAbsoluteX, true},
	{"LDX", 4, ldx, AddressAbsoluteY, true},
	{"TBD", 4, tbd, AddressAbsoluteY, true},

	// C0
	{"CPY", 2, cpy, AddressImmediate, false},
	{"CMP", 6, cmp, AddressIndirectX, false},
	{"NOP", 2, nop, AddressImmediate, false},
	{"TBD", 8, tbd, AddressIndirectX, false},
	{"CPY", 3, cpy, AddressZeroPage, false},
	{"CMP", 3, cmp, AddressZeroPage, false},
	{"DEC", 5, dec, AddressZeroPage, false},
	{"TBD", 5, tbd, AddressZeroPageY, false},
	{"INY", 2, iny, AddressImplicit, false},
	{"CMP", 2, cmp, AddressImmediate, false},
	{"DEX", 2, dex, AddressImplicit, false},
	{"TBD", 2, tbd, AddressImmediate, false},
	{"CPY", 4, cpy, AddressAbsolute, false},
	{"CMP", 4, cmp, AddressAbsolute, false},
	{"DEC", 6, dec, AddressAbsolute, false},
	{"TBD", 6, tbd, AddressAbsolute, false},

	// D0
	{"BNE", 2, bne, AddressRelative, false},
	{"CMP", 5, cmp, AddressIndirectY, true},
	instHalt(),
	{"TBD", 8, tbd, AddressIndirectY, false},
	{"NOP", 4, nop, AddressZeroPageX, false},
	{"CMP", 4, cmp, AddressZeroPageX, false},
	{"DEC", 6, dec, AddressZeroPageX, false},
	{"TBD", 6, tbd, AddressZeroPageX, false},
	{"CLD", 2, cld, AddressImplicit, false},
	{"CMP", 4, cmp, AddressAbsoluteY, true},
	{"NOP", 2, nop, AddressImplicit, false},
	{"TBD", 7, tbd, AddressAbsoluteY, false},
	{"NOP", 4, nop, AddressAbsoluteX, true},
	{"CMP", 4, cmp, AddressAbsoluteX, true},
	{"DEC", 7, dec, AddressAbsoluteX, false},
	{"TBD", 7, tbd, AddressAbsoluteX, false},

	// E0
	{"CPX", 2, cpx, AddressImmediate, false},
	{"SBC", 6, sbc, AddressIndirectX, false},
	{"NOP", 2, nop, AddressImmediate, false},
	{"TBD", 8, tbd, AddressIndirectX, false},
	{"CPX", 3, cpx, AddressZeroPage, false},
	{"SBC", 3, sbc, AddressZeroPage, false},
	{"INC", 5, inc, AddressZeroPage, false},
	{"TBD", 5, tbd, AddressZeroPage, false},
	{"INX", 2, inx, AddressImplicit, false},
	{"SBC", 2, sbc, AddressImmediate, false},
	{"NOP", 2, nop, AddressImplicit, false},
	{"TBD", 2, tbd, AddressImmediate, false},
	{"CPX", 4, cpx, AddressAbsolute, false},
	{"SBC", 4, sbc, AddressAbsolute, false},
	{"INC", 6, inc, AddressAbsolute, false},
	{"TBD", 6, tbd, AddressAbsolute, false},

	// F0
	{"BEQ", 2, beq, AddressRelative, false},
	{"SBC", 5, sbc, AddressIndirectY, true},
	instHalt(),
	{"TBD", 8, tbd, AddressIndirectY, false},
	{"NOP", 4, nop, AddressZeroPageX, false},
	{"SBC", 4, sbc, AddressZeroPageX, false},
	{"INC", 6, inc, AddressZeroPageX, false},
	{"TBD", 6, tbd, AddressZeroPageX, false},
	{"SED", 2, sed, AddressImplicit, false},
	{"SBC", 4, sbc, AddressAbsoluteY, true},
	{"NOP", 2, nop, AddressImplicit, false},
	{"TBD", 7, tbd, AddressAbsoluteY, false},
	{"NOP", 4, nop, AddressAbsoluteX, true},
	{"SBC", 4, sbc, AddressAbsoluteX, true},
	{"INC", 7, inc, AddressAbsoluteX, false},
	{"TBD", 7, tbd, AddressAbsoluteX, false},
}

func OpCodes() map[string]map[AddressMode]byte {
	opCodes := make(map[string]map[AddressMode]byte)
	for code, inst := range insts {
		if _, ok := opCodes[inst.name]; !ok {
			opCodes[inst.name] = make(map[AddressMode]byte)
		}
		// Overwrites if duplicate op/mode
		opCodes[inst.name][inst.addressMode] = byte(code)
	}
	return opCodes
}

func (c *CPU) initInstructions() {
	c.insts = insts
}

func (c *CPU) setResultFlags(val byte) {
	c.flags.Zero = val == 0
	c.flags.Negative = val>>7 == 1
}

// Ordered as per groupings on this page:
// http://obelisk.me.uk/6502/instructions.html

// Load / store Operations

func lda(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		val := cpu.read8(addr)
		cpu.regs.Accumulator = val
		cpu.setResultFlags(val)
	}
}

func ldx(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		val := cpu.read8(addr)
		cpu.regs.IndexX = val
		cpu.setResultFlags(val)
	}
}

func ldy(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		val := cpu.read8(addr)
		cpu.regs.IndexY = val
		cpu.setResultFlags(val)
	}
}

func sta(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.write8(addr, cpu.regs.Accumulator)
	}
}

func stx(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.write8(addr, cpu.regs.IndexX)
	}
}

func sty(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.write8(addr, cpu.regs.IndexY)
	}
}

// Register Transfers

func tax(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.regs.IndexX = cpu.regs.Accumulator
		cpu.setResultFlags(cpu.regs.Accumulator)
	}
}

func tay(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.regs.IndexY = cpu.regs.Accumulator
		cpu.setResultFlags(cpu.regs.Accumulator)
	}
}

func txa(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.regs.Accumulator = cpu.regs.IndexX
		cpu.setResultFlags(cpu.regs.IndexX)
	}
}

func tya(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.regs.Accumulator = cpu.regs.IndexY
		cpu.setResultFlags(cpu.regs.IndexY)
	}
}

// Stack Operations

func tsx(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.regs.IndexX = cpu.regs.StackPtr
		cpu.setResultFlags(cpu.regs.StackPtr)
	}
}

func txs(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.regs.StackPtr = cpu.regs.IndexX
	}
}

func pha(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.stackPush8(cpu.regs.Accumulator)
	}
}

func php(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.stackPush8(cpu.flags.asByte())
	}
}

func pla(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		val := cpu.stackPull8()
		cpu.setResultFlags(val)
		cpu.regs.Accumulator = val
	}
}

func plp(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.setFlagsFromByte(cpu.stackPull8())
	}
}

// Logical

func and(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		val := cpu.read8(addr)
		cpu.regs.Accumulator &= val
		cpu.setResultFlags(cpu.regs.Accumulator)
	}
}

func eor(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		val := cpu.read8(addr)
		cpu.regs.Accumulator ^= val
		cpu.setResultFlags(cpu.regs.Accumulator)
	}
}

func ora(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		val := cpu.read8(addr)
		cpu.regs.Accumulator |= val
		cpu.setResultFlags(cpu.regs.Accumulator)
	}
}

func bit(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		val := cpu.read8(addr)
		test := val & cpu.regs.Accumulator
		cpu.flags.Zero = test == 0
		cpu.flags.Overflow = (val>>6)&0x1 == 1
		cpu.flags.Negative = (val>>7)&0x1 == 1
	}
}

// Arithmetic

func adc(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		val := cpu.read8(addr)
		accum := cpu.regs.Accumulator
		var carry byte
		if cpu.flags.Carry {
			carry = 1
		}
		cpu.regs.Accumulator = val + carry + cpu.regs.Accumulator
		cpu.setFlagsFromByte(cpu.regs.Accumulator)
		cpu.flags.Carry = uint16(val)+uint16(carry)+uint16(accum) > 0xFF
		cpu.flags.Overflow = (accum^val)&0x80 == 0 && (val^cpu.regs.Accumulator) != 0
	}
}

func sbc(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		val := cpu.read8(addr)
		accum := cpu.regs.Accumulator
		var carry byte
		if cpu.flags.Carry {
			carry = 1
		}
		cpu.regs.Accumulator = cpu.regs.Accumulator - val - (1 - carry)
		cpu.setFlagsFromByte(cpu.regs.Accumulator)
		cpu.flags.Carry = uint16(accum)-uint16(val)-uint16(1-carry) > 0xFF
		cpu.flags.Overflow = (accum^val)&0x80 != 0 && (val^cpu.regs.Accumulator)&0x80 != 0
	}
}

func cmp(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.compare(cpu.regs.Accumulator, cpu.read8(addr))
	}
}

func cpx(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.compare(cpu.regs.IndexX, cpu.read8(addr))
	}
}

func cpy(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.compare(cpu.regs.IndexY, cpu.read8(addr))
	}
}

// Increments and Decrements

func inc(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		val := cpu.read8(addr)
		val++
		cpu.setResultFlags(val)
		cpu.write8(addr, val)
	}
}

func inx(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.regs.IndexX++
		cpu.setResultFlags(cpu.regs.IndexX)
	}
}

func iny(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.regs.IndexY++
		cpu.setResultFlags(cpu.regs.IndexY)
	}
}

func dec(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		val := cpu.read8(addr)
		val--
		cpu.setResultFlags(val)
		cpu.write8(addr, val)
	}
}

func dex(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.regs.IndexX--
		cpu.setResultFlags(cpu.regs.IndexX)
	}
}

func dey(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.regs.IndexY--
		cpu.setResultFlags(cpu.regs.IndexY)
	}
}

// Shifts

func asl(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		if mode == AddressAccumulator {
			cpu.flags.Carry = (cpu.regs.Accumulator>>7)&0x1 == 1
			cpu.regs.Accumulator <<= 1
			cpu.setResultFlags(cpu.regs.Accumulator)
		} else {
			val := cpu.read8(addr)
			cpu.flags.Carry = (val>>7)&0x1 == 1
			val <<= 1
			cpu.write8(addr, val)
			cpu.setResultFlags(val)
		}
	}
}

func lsr(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		if mode == AddressAccumulator {
			cpu.flags.Carry = cpu.regs.Accumulator&0x1 == 1
			cpu.regs.Accumulator >>= 1
			cpu.setResultFlags(cpu.regs.Accumulator)
		} else {
			val := cpu.read8(addr)
			cpu.flags.Carry = val&0x1 == 1
			val >>= 1
			cpu.write8(addr, val)
			cpu.setResultFlags(val)
		}
	}
}

func rol(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		oldCarry := cpu.flags.Carry
		if mode == AddressAccumulator {
			cpu.flags.Carry = (cpu.regs.Accumulator>>7)&0x1 == 1
			cpu.regs.Accumulator <<= 1
			if oldCarry {
				cpu.regs.Accumulator |= 0x1
			}
			cpu.setResultFlags(cpu.regs.Accumulator)
		} else {
			val := cpu.read8(addr)
			cpu.flags.Carry = (val>>7)&0x1 == 1
			val <<= 1
			if oldCarry {
				val |= 0x1
			}
			cpu.write8(addr, val)
			cpu.setResultFlags(val)
		}
	}
}

func ror(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		oldCarry := cpu.flags.Carry
		if mode == AddressAccumulator {
			cpu.flags.Carry = cpu.regs.Accumulator&0x1 == 1
			cpu.regs.Accumulator >>= 1
			if oldCarry {
				cpu.regs.Accumulator |= 0x80
			}
			cpu.setResultFlags(cpu.regs.Accumulator)
		} else {
			val := cpu.read8(addr)
			cpu.flags.Carry = val&0x1 == 1
			val >>= 1
			if oldCarry {
				val |= 0x80
			}
			cpu.write8(addr, val)
			cpu.setResultFlags(val)
		}
	}
}

// Jumps and Calls

func jmp(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.pc = addr
	}
}

func jsr(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.stackPush16(cpu.pc - 1)
		cpu.pc = addr
	}
}

func rts(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.pc = cpu.stackPull16() + 1
	}
}

// Branches

func bcc(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		if !cpu.flags.Carry {
			cpu.branch(addr)
		}
	}
}

func bcs(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		if cpu.flags.Carry {
			cpu.branch(addr)
		}
	}
}

func beq(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		if cpu.flags.Zero {
			cpu.branch(addr)
		}
	}
}

func bmi(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		if cpu.flags.Negative {
			cpu.branch(addr)
		}
	}
}

func bne(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		if !cpu.flags.Zero {
			cpu.branch(addr)
		}
	}
}

func bpl(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		if !cpu.flags.Negative {
			cpu.branch(addr)
		}
	}
}

func bvc(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		if !cpu.flags.Overflow {
			cpu.branch(addr)
		}
	}
}

func bvs(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		if cpu.flags.Overflow {
			cpu.branch(addr)
		}
	}
}

// Status Flag Changes

func clc(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.flags.Carry = false
	}
}

func cld(cpu *CPU, addr uint16, mode AddressMode) func() {
	// We don't track the decimal flag
	return nil
}

func cli(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.flags.InterruptDisable = false
	}
}

func clv(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.flags.Overflow = false
	}
}

func sec(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.flags.Carry = true
	}
}

func sed(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		// We don't track the decimal flag
	}
}

func sei(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.flags.InterruptDisable = true
	}
}

// System Functions

func brk(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.flags.BreakCmd = true
		cpu.irq()
		cpu.flags.BreakCmd = false
	}
}

func rti(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		cpu.setFlagsFromByte(cpu.stackPull8()&0xEF | 0x20)
		cpu.pc = cpu.stackPull16()
	}
}

func nop(cpu *CPU, addr uint16, mode AddressMode) func() {
	return nil
}

// Unofficial / Illegal Opcodes
// http://wiki.nesdev.com/w/index.php/Programming_with_unofficial_opcodes
// http://www.oxyron.de/html/opcodes02.html
// We implement these for compatibility, even though they aren't part of the 6502 spec

func tbd(cpu *CPU, addr uint16, mode AddressMode) func() {
	logrus.Warn("hit unimplemented opcode")
	return nil
}

func alr(cpu *CPU, addr uint16, mode AddressMode) func() {
	return func() {
		and(cpu, addr, mode)
		lsr(cpu, 0, AddressAccumulator)
	}
}

// TODO: Implement more of these

// Helpers

func instHalt() *inst {
	return &inst{
		name:        "KIL",
		addressMode: AddressImplicit,
		op: func(cpu *CPU, addr uint16, mode AddressMode) func() {
			return func() {
				cpu.halted = true
			}
		},
	}
}
