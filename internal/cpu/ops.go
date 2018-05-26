package cpu

type op func(addr uint16)

type inst struct {
	name           string
	cycles         uint64
	op             op
	addressMode    addressMode
	pageCrossCycle bool
}

func (i *inst) fullName() string {
	return i.name + " " + i.addressMode.name()
}

func (c *CPU) initInstructions() {
	c.insts = [256]*inst{
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
		{"LDY", 2, c.ldy, addressImmediate, false},
		instNoop(6),
		{"LDX", 2, c.ldx, addressImmediate, false},
		instNoop(6),
		{"LDY", 3, c.ldy, addressZeroPage, false},
		{"LDA", 3, c.lda, addressZeroPage, false},
		{"LDX", 3, c.ldx, addressZeroPage, false},
		instNoop(3),
		instNoop(2),
		{"LDA", 2, c.lda, addressImmediate, false},
		instNoop(2),
		instNoop(2),
		{"LDY", 4, c.ldy, addressAbsolute, false},
		{"LDA", 4, c.lda, addressAbsolute, false},
		{"LDX", 4, c.ldx, addressAbsolute, false},
		instNoop(4),

		// B0
		instNoop(2),
		instNoop(5),
		instIllegal(),
		instNoop(5),
		{"LDY", 4, c.ldy, addressZeroPageX, false},
		{"LDA", 4, c.lda, addressZeroPageX, false},
		{"LDX", 4, c.ldx, addressZeroPageY, false},
		instNoop(4),
		instNoop(2),
		{"LDA", 4, c.lda, addressAbsoluteY, true},
		instNoop(2),
		instNoop(4),
		{"LDY", 4, c.ldy, addressAbsoluteX, true},
		{"LDA", 4, c.lda, addressAbsoluteX, true},
		{"LDX", 4, c.ldx, addressAbsoluteY, true},
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
		{"INC", 5, c.inc, addressZeroPage, false},
		instNoop(5),
		instNoop(2),
		instNoop(2),
		instNoop(2),
		instNoop(2),
		instNoop(4),
		instNoop(4),
		{"INC", 6, c.inc, addressAbsolute, false},
		instNoop(6),

		// F0
		instNoop(2),
		instNoop(5),
		instIllegal(),
		instNoop(8),
		instNoop(4),
		instNoop(4),
		{"INC", 6, c.inc, addressZeroPageX, false},
		instNoop(6),
		instNoop(2),
		instNoop(4),
		instNoop(2),
		instNoop(7),
		instNoop(4),
		instNoop(4),
		{"INC", 7, c.inc, addressAbsoluteX, false},
		instNoop(7),
	}
}

func (c *CPU) setResultFlags(val byte) {
	c.flags.Zero = val == 0
	c.flags.Negative = val>>7 == 1
}

func (c *CPU) inc(addr uint16) {
	val := c.read8(addr)
	val++
	c.setResultFlags(val)
	c.write8(addr, val)
}

func (c *CPU) lda(addr uint16) {
	val := c.read8(addr)
	c.regs.Accumulator = val
	c.setResultFlags(val)
}

func (c *CPU) ldx(addr uint16) {
	val := c.read8(addr)
	c.regs.IndexX = val
	c.setResultFlags(val)
}

func (c *CPU) ldy(addr uint16) {
	val := c.read8(addr)
	c.regs.IndexY = val
	c.setResultFlags(val)
}

func instNoop(cycles uint64) *inst {
	return &inst{
		name:   "NOP",
		cycles: cycles,
		op: func(addr uint16) {
		},
	}
}

func instIllegal() *inst {
	return &inst{
		name: "ILLEGAL",
		op: func(addr uint16) {
			panic("illegal operation")
		},
	}
}
