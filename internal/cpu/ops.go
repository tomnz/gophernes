package cpu

type op func(cpu *CPU, addr uint16, mode AddressMode)

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
	instNoop(7),
	instNoop(6),
	instIllegal(),
	instNoop(8),
	instNoop(3),
	instNoop(3),
	instNoop(5),
	instNoop(5),
	{"PHP", 3, php, AddressImplicit, false},
	instNoop(2),
	instNoop(2),
	instNoop(2),
	instNoop(4),
	instNoop(4),
	instNoop(6),
	instNoop(6),

	// 10
	instNoop(2),
	instNoop(5),
	instIllegal(),
	instNoop(8),
	instNoop(4),
	instNoop(4),
	instNoop(6),
	instNoop(6),
	instNoop(2),
	instNoop(4),
	instNoop(2),
	instNoop(7),
	instNoop(4),
	instNoop(4),
	instNoop(7),
	instNoop(7),

	// 20
	instNoop(6),
	instNoop(6),
	instIllegal(),
	instNoop(8),
	instNoop(3),
	instNoop(3),
	instNoop(5),
	instNoop(5),
	{"PLP", 4, plp, AddressImplicit, false},
	instNoop(2),
	instNoop(2),
	instNoop(2),
	instNoop(4),
	instNoop(4),
	instNoop(6),
	instNoop(6),

	// 30
	instNoop(2),
	instNoop(5),
	instIllegal(),
	instNoop(8),
	instNoop(4),
	instNoop(4),
	instNoop(6),
	instNoop(6),
	instNoop(2),
	instNoop(4),
	instNoop(2),
	instNoop(7),
	instNoop(4),
	instNoop(4),
	instNoop(7),
	instNoop(7),

	// 40
	instNoop(6),
	instNoop(6),
	instIllegal(),
	instNoop(8),
	instNoop(3),
	instNoop(3),
	instNoop(5),
	instNoop(5),
	{"PHA", 3, pha, AddressImplicit, false},
	instNoop(2),
	instNoop(2),
	instNoop(2),
	instNoop(3),
	instNoop(4),
	instNoop(6),
	instNoop(6),

	// 50
	instNoop(2),
	instNoop(5),
	instIllegal(),
	instNoop(8),
	instNoop(4),
	instNoop(4),
	instNoop(6),
	instNoop(6),
	instNoop(2),
	instNoop(4),
	instNoop(2),
	instNoop(7),
	instNoop(4),
	instNoop(4),
	instNoop(7),
	instNoop(7),

	// 60
	instNoop(6),
	instNoop(6),
	instIllegal(),
	instNoop(8),
	instNoop(3),
	instNoop(3),
	instNoop(5),
	instNoop(5),
	{"PLA", 4, pla, AddressImplicit, false},
	instNoop(2),
	instNoop(2),
	instNoop(2),
	instNoop(5),
	instNoop(4),
	instNoop(6),
	instNoop(6),

	// 70
	instNoop(2),
	instNoop(5),
	instIllegal(),
	instNoop(8),
	instNoop(4),
	instNoop(4),
	instNoop(6),
	instNoop(6),
	instNoop(2),
	instNoop(4),
	instNoop(2),
	instNoop(7),
	instNoop(4),
	instNoop(4),
	instNoop(7),
	instNoop(7),

	// 80
	instNoop(2),
	{"STA", 6, sta, AddressIndirectX, false},
	instNoop(2),
	instNoop(6),
	{"STY", 3, sty, AddressZeroPage, false},
	{"STA", 3, sta, AddressZeroPage, false},
	{"STX", 3, stx, AddressZeroPage, false},
	instNoop(3),
	instNoop(2),
	instNoop(2),
	{"TXA", 2, txa, AddressImplicit, false},
	instNoop(2),
	{"STY", 4, sty, AddressAbsolute, false},
	{"STA", 4, sta, AddressAbsolute, false},
	{"STX", 4, stx, AddressAbsolute, false},
	instNoop(4),

	// 90
	instNoop(2),
	{"STA", 6, sta, AddressIndirectY, false},
	instIllegal(),
	instNoop(6),
	{"STY", 4, sty, AddressZeroPageX, false},
	{"STA", 4, sta, AddressZeroPageX, false},
	{"STX", 4, stx, AddressZeroPageY, false},
	instNoop(4),
	{"TYA", 2, tya, AddressImplicit, false},
	{"STA", 5, sta, AddressAbsoluteY, false},
	{"TXS", 2, txs, AddressImplicit, false},
	instNoop(5),
	instNoop(5),
	{"STA", 5, sta, AddressAbsoluteX, false},
	instNoop(5),
	instNoop(5),

	// A0
	{"LDY", 2, ldy, AddressImmediate, false},
	instNoop(6),
	{"LDX", 2, ldx, AddressImmediate, false},
	instNoop(6),
	{"LDY", 3, ldy, AddressZeroPage, false},
	{"LDA", 3, lda, AddressZeroPage, false},
	{"LDX", 3, ldx, AddressZeroPage, false},
	instNoop(3),
	{"TAY", 2, tay, AddressImplicit, false},
	{"LDA", 2, lda, AddressImmediate, false},
	{"TAX", 2, tax, AddressImplicit, false},
	instNoop(2),
	{"LDY", 4, ldy, AddressAbsolute, false},
	{"LDA", 4, lda, AddressAbsolute, false},
	{"LDX", 4, ldx, AddressAbsolute, false},
	instNoop(4),

	// B0
	instNoop(2),
	instNoop(5),
	instIllegal(),
	instNoop(5),
	{"LDY", 4, ldy, AddressZeroPageX, false},
	{"LDA", 4, lda, AddressZeroPageX, false},
	{"LDX", 4, ldx, AddressZeroPageY, false},
	instNoop(4),
	instNoop(2),
	{"LDA", 4, lda, AddressAbsoluteY, true},
	{"TSX", 2, tsx, AddressImplicit, false},
	instNoop(4),
	{"LDY", 4, ldy, AddressAbsoluteX, true},
	{"LDA", 4, lda, AddressAbsoluteX, true},
	{"LDX", 4, ldx, AddressAbsoluteY, true},
	instNoop(4),

	// C0
	instNoop(2),
	instNoop(6),
	instNoop(2),
	instNoop(8),
	instNoop(3),
	instNoop(3),
	instNoop(5),
	instNoop(5),
	instNoop(2),
	instNoop(2),
	instNoop(2),
	instNoop(2),
	instNoop(4),
	instNoop(4),
	instNoop(6),
	instNoop(6),

	// D0
	instNoop(2),
	instNoop(5),
	instIllegal(),
	instNoop(8),
	instNoop(3),
	instNoop(3),
	instNoop(5),
	instNoop(5),
	instNoop(2),
	instNoop(2),
	instNoop(2),
	instNoop(2),
	instNoop(4),
	instNoop(4),
	instNoop(6),
	instNoop(6),

	// E0
	instNoop(2),
	instNoop(6),
	instNoop(2),
	instNoop(8),
	instNoop(3),
	instNoop(3),
	{"INC", 5, inc, AddressZeroPage, false},
	instNoop(5),
	instNoop(2),
	instNoop(2),
	instNoop(2),
	instNoop(2),
	instNoop(4),
	instNoop(4),
	{"INC", 6, inc, AddressAbsolute, false},
	instNoop(6),

	// F0
	instNoop(2),
	instNoop(5),
	instIllegal(),
	instNoop(8),
	instNoop(4),
	instNoop(4),
	{"INC", 6, inc, AddressZeroPageX, false},
	instNoop(6),
	instNoop(2),
	instNoop(4),
	instNoop(2),
	instNoop(7),
	instNoop(4),
	instNoop(4),
	{"INC", 7, inc, AddressAbsoluteX, false},
	instNoop(7),
}

func (c *CPU) initInstructions() {
	c.insts = insts
}

func (c *CPU) setResultFlags(val byte) {
	c.flags.Zero = val == 0
	c.flags.Negative = val>>7 == 1
}

func inc(cpu *CPU, addr uint16, mode AddressMode) {
	val := cpu.read8(addr)
	val++
	cpu.setResultFlags(val)
	cpu.write8(addr, val)
}

// Ordered as per groupings on this page:
// http://obelisk.me.uk/6502/instructions.html

// Load / store Operations

func lda(cpu *CPU, addr uint16, mode AddressMode) {
	val := cpu.read8(addr)
	cpu.regs.Accumulator = val
	cpu.setResultFlags(val)
}

func ldx(cpu *CPU, addr uint16, mode AddressMode) {
	val := cpu.read8(addr)
	cpu.regs.IndexX = val
	cpu.setResultFlags(val)
}

func ldy(cpu *CPU, addr uint16, mode AddressMode) {
	val := cpu.read8(addr)
	cpu.regs.IndexY = val
	cpu.setResultFlags(val)
}

func sta(cpu *CPU, addr uint16, mode AddressMode) {
	cpu.write8(addr, cpu.regs.Accumulator)
}

func stx(cpu *CPU, addr uint16, mode AddressMode) {
	cpu.write8(addr, cpu.regs.Accumulator)
}

func sty(cpu *CPU, addr uint16, mode AddressMode) {
	cpu.write8(addr, cpu.regs.Accumulator)
}

// Register Transfers

func tax(cpu *CPU, addr uint16, mode AddressMode) {
	cpu.regs.IndexX = cpu.regs.Accumulator
	cpu.setResultFlags(cpu.regs.Accumulator)
}

func tay(cpu *CPU, addr uint16, mode AddressMode) {
	cpu.regs.IndexY = cpu.regs.Accumulator
	cpu.setResultFlags(cpu.regs.Accumulator)
}

func txa(cpu *CPU, addr uint16, mode AddressMode) {
	cpu.regs.Accumulator = cpu.regs.IndexX
	cpu.setResultFlags(cpu.regs.IndexX)
}

func tya(cpu *CPU, addr uint16, mode AddressMode) {
	cpu.regs.Accumulator = cpu.regs.IndexY
	cpu.setResultFlags(cpu.regs.IndexY)
}

// Stack Operations

func tsx(cpu *CPU, addr uint16, mode AddressMode) {
	cpu.regs.IndexX = cpu.regs.StackPtr
	cpu.setResultFlags(cpu.regs.StackPtr)
}

func txs(cpu *CPU, addr uint16, mode AddressMode) {
	cpu.regs.StackPtr = cpu.regs.IndexX
}

func pha(cpu *CPU, addr uint16, mode AddressMode) {
	stackAddr := uint16(cpu.regs.StackPtr)
	stackAddr |= 0x100
	cpu.write8(stackAddr, cpu.regs.Accumulator)
	cpu.regs.StackPtr++
}

func php(cpu *CPU, addr uint16, mode AddressMode) {
	stackAddr := uint16(cpu.regs.StackPtr)
	stackAddr |= 0x100
	cpu.write8(stackAddr, cpu.flags.asByte())
	cpu.regs.StackPtr++
}

func pla(cpu *CPU, addr uint16, mode AddressMode) {
	cpu.regs.StackPtr--
	stackAddr := uint16(cpu.regs.StackPtr)
	stackAddr |= 0x100
	cpu.regs.Accumulator = cpu.read8(stackAddr)
	cpu.setResultFlags(cpu.regs.Accumulator)
}

func plp(cpu *CPU, addr uint16, mode AddressMode) {
	cpu.regs.StackPtr--
	stackAddr := uint16(cpu.regs.StackPtr)
	stackAddr |= 0x100
	cpu.setFlagsFromByte(cpu.read8(stackAddr))
}

// Helpers

func instNoop(cycles uint64) *inst {
	return &inst{
		name:   "NOP",
		cycles: cycles,
		op: func(cpu *CPU, addr uint16, mode AddressMode) {
		},
	}
}

func instIllegal() *inst {
	return &inst{
		name: "ILLEGAL",
		op: func(cpu *CPU, addr uint16, mode AddressMode) {
			panic("illegal operation")
		},
	}
}
