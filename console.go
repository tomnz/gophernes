package gophernes

import (
	"io"

	"github.com/tomnz/gophernes/internal/cpu"
	"github.com/tomnz/gophernes/internal/ppu"
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
	memory := NewMemory(cartridge)

	cpu := cpu.NewCPU(&cpuMemory{memory}, cpu.WithTrace())
	memory.cpu = cpu

	ppu := ppu.NewPPU(&ppuMemory{memory}, ppu.WithTrace())
	memory.ppu = ppu

	cpu.Reset()

	return &Console{
		cpu:    cpu,
		memory: memory,
	}, nil
}

func (c *Console) Run() error {
	_, err := c.cpu.RunTilHalt()
	return err
}

func (c *Console) RunSteps(steps int) error {
	_, err := c.cpu.Run(steps)
	return err
}
