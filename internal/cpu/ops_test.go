package cpu_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tomnz/gophernes/internal/cpu"
)

// testMemory implements basic reads/writes for the purposes of writing low level CPU tests.
// Memory is implemented as a map to detect invalid reads, while allowing writes to any location.
type testMemory struct {
	mem map[uint16]byte
}

var ops = cpu.OpCodes()

// newTestMemory instantiates a new testMemory with the given data.
// The first page is written as contents == position. I.e. address 0x10 contains value 0x10 and so on.
// Data is written past the first page (starting at 0x0100) and program is written starting at 0x8000.
func newTestMemory(data []byte, prg []byte) *testMemory {
	mem := map[uint16]byte{}

	log.Printf("%#v", prg)
	for i := 0; i < 256; i++ {
		mem[uint16(i)] = byte(i)
	}
	for index, val := range data {
		mem[uint16(index+0x100)] = val
	}
	for index, val := range prg {
		mem[uint16(index)+0x8000] = val
	}
	// Add halt code to end test
	mem[uint16(len(prg))+0x8000] = ops["KIL"][cpu.AddressImplicit]

	// Initialize reset vector to point to beginning of prg block
	mem[0xFFFC] = 0
	mem[0xFFFD] = 0x80

	return &testMemory{
		mem: mem,
	}
}

func (t *testMemory) Read(addr uint16) byte {
	val, ok := t.mem[addr]
	if !ok {
		panic(fmt.Sprintf("access to uninitialized memory at address %#x", addr))
	}
	return val
}

func (t *testMemory) Write(addr uint16, val byte) {
	t.mem[addr] = val
}

func TestOps(t *testing.T) {
	testCases := map[string]struct {
		// Inputs
		data []byte
		prg  []byte
		// Expected outputs
		regs   *cpu.Registers
		flags  *cpu.Flags
		cycles uint64
	}{
		"load accumulator: immediate": {
			prg: []byte{
				ops["LDA"][cpu.AddressImmediate],
				2,
			},
			regs: &cpu.Registers{
				Accumulator: 2,
			},
			cycles: 2,
		},
		"load accumulator: zero page": {
			prg: []byte{
				ops["LDA"][cpu.AddressZeroPage],
				2,
			},
			regs: &cpu.Registers{
				Accumulator: 2,
			},
			cycles: 3,
		},
		"load accumulator: zero page x": {
			prg: []byte{
				ops["LDX"][cpu.AddressImmediate],
				3,
				ops["LDA"][cpu.AddressZeroPageX],
				2,
			},
			regs: &cpu.Registers{
				Accumulator: 5,
				IndexX:      3,
			},
			cycles: 6,
		},
		"load accumulator: absolute": {
			data: []byte{10, 11, 12},
			prg: []byte{
				ops["LDA"][cpu.AddressAbsolute],
				// Address 0x0102
				0x02,
				0x01,
			},
			regs: &cpu.Registers{
				Accumulator: 12,
			},
			cycles: 4,
		},
		"load accumulator: absolute x": {
			prg: []byte{
				ops["LDX"][cpu.AddressImmediate],
				3,
				ops["LDA"][cpu.AddressAbsoluteX],
				// Address 0x00F0
				0xF0,
				0x00,
			},
			regs: &cpu.Registers{
				Accumulator: 0xF3,
				IndexX:      3,
			},
			cycles: 6,
		},
		"load accumulator: absolute x with page cross": {
			data: []byte{10, 11, 12},
			prg: []byte{
				ops["LDX"][cpu.AddressImmediate],
				3,
				ops["LDA"][cpu.AddressAbsoluteX],
				// Address 0x00FF
				0xFF,
				0x00,
			},
			regs: &cpu.Registers{
				Accumulator: 12,
				IndexX:      3,
			},
			// Extra cycle
			cycles: 7,
		},
		"load accumulator: absolute y": {
			prg: []byte{
				ops["LDY"][cpu.AddressImmediate],
				3,
				ops["LDA"][cpu.AddressAbsoluteY],
				// Address 0x00F0
				0xF0,
				0x00,
			},
			regs: &cpu.Registers{
				Accumulator: 0xF3,
				IndexY:      3,
			},
			cycles: 6,
		},
		"load accumulator: absolute y with page cross": {
			data: []byte{10, 11, 12},
			prg: []byte{
				ops["LDY"][cpu.AddressImmediate],
				3,
				ops["LDA"][cpu.AddressAbsoluteY],
				// Address 0x00F0
				0xFF,
				0x00,
			},
			regs: &cpu.Registers{
				Accumulator: 12,
				IndexY:      3,
			},
			cycles: 7,
		},
		"increment: zero page": {
			prg: []byte{
				ops["INC"][cpu.AddressZeroPage],
				5,
				ops["LDA"][cpu.AddressZeroPage],
				5,
			},
			regs: &cpu.Registers{
				Accumulator: 6,
			},
			cycles: 8,
		},
		"load accumulator: zero flag": {
			prg: []byte{
				ops["LDA"][cpu.AddressImmediate],
				0,
			},
			flags: &cpu.Flags{
				BreakCmd:         true,
				InterruptDisable: true,
				Zero:             true,
			},
			cycles: 2,
		},
		"load accumulator: nonzero flag": {
			prg: []byte{
				ops["LDA"][cpu.AddressImmediate],
				1,
			},
			flags: &cpu.Flags{
				BreakCmd:         true,
				InterruptDisable: true,
				Zero:             false,
			},
			cycles: 2,
		},
		"load accumulator: negative": {
			prg: []byte{
				ops["LDA"][cpu.AddressImmediate],
				0xF0,
			},
			flags: &cpu.Flags{
				BreakCmd:         true,
				InterruptDisable: true,
				Negative:         true,
			},
			cycles: 2,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mem := newTestMemory(tc.data, tc.prg)
			cpu := cpu.NewCPU(mem, cpu.WithTrace())
			cpu.Reset()
			cycles, err := cpu.RunTilHalt()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if tc.regs != nil {
				if diff := cmp.Diff(*tc.regs, cpu.Registers()); diff != "" {
					t.Errorf("unexpected registers:\n%s", diff)
				}
			}
			if tc.flags != nil {
				if diff := cmp.Diff(*tc.flags, cpu.Flags()); diff != "" {
					t.Errorf("unexpected flags:\n%s", diff)
				}
			}
			if tc.cycles > 0 && tc.cycles != cycles {
				t.Errorf("expected %d cycles, got %d", tc.cycles, cycles)
			}
		})
	}
}
