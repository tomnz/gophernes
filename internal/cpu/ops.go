package cpu

type op func(cpu *CPU, addr uint16)

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
	instNoop(3),
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
	instNoop(4),
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
	instNoop(3),
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
	instNoop(5),
	instNoop(5),
	instNoop(4),
	instNoop(2),
	instNoop(2),
	instNoop(2),
	instNoop(5),
	instNoop(4),
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
	instNoop(6),
	instNoop(2),
	instNoop(6),
	instNoop(3),
	instNoop(3),
	instNoop(3),
	instNoop(3),
	instNoop(2),
	instNoop(2),
	instNoop(2),
	instNoop(2),
	instNoop(4),
	instNoop(4),
	instNoop(4),
	instNoop(4),

	// 90
	instNoop(2),
	instNoop(6),
	instIllegal(),
	instNoop(6),
	instNoop(4),
	instNoop(4),
	instNoop(4),
	instNoop(4),
	instNoop(2),
	instNoop(5),
	instNoop(2),
	instNoop(5),
	instNoop(5),
	instNoop(5),
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
	instNoop(2),
	{"LDA", 2, lda, AddressImmediate, false},
	instNoop(2),
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
	instNoop(2),
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

func inc(cpu *CPU, addr uint16) {
	val := cpu.read8(addr)
	val++
	cpu.setResultFlags(val)
	cpu.write8(addr, val)
}

func lda(cpu *CPU, addr uint16) {
	val := cpu.read8(addr)
	cpu.regs.Accumulator = val
	cpu.setResultFlags(val)
}

func ldx(cpu *CPU, addr uint16) {
	val := cpu.read8(addr)
	cpu.regs.IndexX = val
	cpu.setResultFlags(val)
}

func ldy(cpu *CPU, addr uint16) {
	val := cpu.read8(addr)
	cpu.regs.IndexY = val
	cpu.setResultFlags(val)
}

func instNoop(cycles uint64) *inst {
	return &inst{
		name:   "NOP",
		cycles: cycles,
		op: func(cpu *CPU, addr uint16) {
		},
	}
}

func instIllegal() *inst {
	return &inst{
		name: "ILLEGAL",
		op: func(cpu *CPU, addr uint16) {
			panic("illegal operation")
		},
	}
}
