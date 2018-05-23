package gophernes

// CPU implements the main NES CPU.
type CPU struct {
	pc     uint16
	regs   registers
	flags  flags
	cycles uint64
}

type registers struct {
	stackPtr,
	accumulator,
	indexX,
	indexY byte
}

type flags struct {
	negative,
	overflow,
	breakCmd,
	interruptDisable,
	zero,
	carry bool
}
