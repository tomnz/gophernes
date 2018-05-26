package gophernes

import (
	"io"

	"github.com/tomnz/gophernes/internal/cpu"
)

// Console implements the main console.
type Console struct {
	cpu    *cpu.CPU
	memory *Memory
}

// NewConsole initializes a new console.
func NewConsole(rom io.Reader) (*Console, error) {
	cartridge, err := loadINES(rom)
	if err != nil {
		return nil, err
	}
	memory := &Memory{
		cartridge: cartridge,
	}

	cpu := cpu.NewCPU(&cpuMemory{memory})
	memory.cpu = cpu

	cpu.Reset()

	return &Console{
		cpu:    cpu,
		memory: memory,
	}, nil
}
