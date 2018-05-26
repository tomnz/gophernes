package cpu

import (
	"fmt"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// testMemory implements basic reads/writes for the purposes of writing low level CPU tests.
// Memory is implemented as a map to detect invalid reads, while allowing writes to any location.
type testMemory struct {
	mem map[uint16]byte
}

// newTestMemory instantiates a new testMemory with the given data.
// The first page is written as contents == position. I.e. address 0x10 contains value 0x10 and so on.
// Data is written past the first page (starting at 0x0100) and program is written starting at 0x8000.
func newTestMemory(data []byte, prg []byte) *testMemory {
	mem := map[uint16]byte{}

	for i := 0; i < 256; i++ {
		mem[uint16(i)] = byte(i)
	}
	for index, val := range data {
		mem[uint16(index+0x100)] = val
	}
	for index, val := range prg {
		mem[uint16(index)+0x8000] = val
	}

	// Initialize reset vector to point to beginning of prg block
	mem[resetVector] = 0
	mem[resetVector+1] = 0x80

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
		data   []byte
		prg    []byte
		numOps int
		// Expected outputs
		regs   *registers
		flags  *flags
		cycles uint64
	}{
		"load accumulator: immediate": {
			prg: []byte{
				// LDA imm
				0xA9,
				2,
			},
			numOps: 1,
			regs: &registers{
				Accumulator: 2,
			},
			cycles: 2,
		},
		"load accumulator: zero page": {
			prg: []byte{
				// LDA zp
				0xA5,
				2,
			},
			numOps: 1,
			regs: &registers{
				Accumulator: 2,
			},
			cycles: 3,
		},
		"load accumulator: zero page x": {
			prg: []byte{
				// LDX imm
				0xA2,
				3,
				// LDA zpx
				0xB5,
				2,
			},
			numOps: 2,
			regs: &registers{
				Accumulator: 5,
				IndexX:      3,
			},
			cycles: 6,
		},
		"load accumulator: absolute": {
			data: []byte{10, 11, 12},
			prg: []byte{
				// LDA abs
				0xAD,
				// Address 0x0102
				0x02,
				0x01,
			},
			numOps: 1,
			regs: &registers{
				Accumulator: 12,
			},
			cycles: 4,
		},
		"load accumulator: absolute x": {
			prg: []byte{
				// LDX imm
				0xA2,
				3,
				// LDA absx
				0xBD,
				// Address 0x00F0
				0xF0,
				0x00,
			},
			numOps: 2,
			regs: &registers{
				Accumulator: 0xF3,
				IndexX:      3,
			},
			cycles: 6,
		},
		"load accumulator: absolute x with page cross": {
			data: []byte{10, 11, 12},
			prg: []byte{
				// LDX imm
				0xA2,
				3,
				// LDA absx
				0xBD,
				// Address 0x00FF
				0xFF,
				0x00,
			},
			numOps: 2,
			regs: &registers{
				Accumulator: 12,
				IndexX:      3,
			},
			// Extra cycle
			cycles: 7,
		},
		"load accumulator: absolute y": {
			prg: []byte{
				// LDY imm
				0xA0,
				3,
				// LDA absy
				0xB9,
				// Address 0x00F0
				0xF0,
				0x00,
			},
			numOps: 2,
			regs: &registers{
				Accumulator: 0xF3,
				IndexY:      3,
			},
			cycles: 6,
		},
		"load accumulator: absolute y with page cross": {
			data: []byte{10, 11, 12},
			prg: []byte{
				// LDY imm
				0xA0,
				3,
				// LDA absy
				0xB9,
				// Address 0x00F0
				0xFF,
				0x00,
			},
			numOps: 2,
			regs: &registers{
				Accumulator: 12,
				IndexY:      3,
			},
			cycles: 7,
		},
		"increment: zero page": {
			prg: []byte{
				// INC zp (5)
				0xE6,
				5,
				// LDA zp (3)
				0xA5,
				5,
			},
			numOps: 2,
			regs: &registers{
				Accumulator: 6,
			},
			cycles: 8,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			log.Printf("Starting test %q", name)
			mem := newTestMemory(tc.data, tc.prg)
			cpu := NewCPU(mem, WithTrace())
			cpu.Reset()
			cpu.Run(tc.numOps)

			if tc.regs != nil {
				if diff := cmp.Diff(*tc.regs, cpu.regs); diff != "" {
					t.Errorf("unexpected registers:\n%s", diff)
				}
			}
			if tc.cycles > 0 && tc.cycles != cpu.cycles {
				t.Errorf("expected %d cycles, got %d", tc.cycles, cpu.cycles)
			}
		})
	}
}
