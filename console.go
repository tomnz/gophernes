package gophernes

import (
	"io"

	"context"
	"github.com/tomnz/gophernes/internal/cpu"
	"github.com/tomnz/gophernes/internal/ppu"
)

// Console implements the main console.
type Console struct {
	cpu    *cpu.CPU
	ppu    *ppu.PPU
	memory *Memory
}

// NewConsole initializes a new console.
func NewConsole(rom io.Reader, cpuopts []cpu.Option, ppuopts []ppu.Option) (*Console, error) {
	cartridge, err := loadINES(rom)
	if err != nil {
		return nil, err
	}
	memory := NewMemory(cartridge)

	cpu := cpu.NewCPU(&cpuMemory{memory}, cpuopts...)
	memory.cpu = cpu

	ppu := ppu.NewPPU(&ppuMemory{memory}, ppuopts...)
	memory.ppu = ppu

	cpu.Reset()

	return &Console{
		cpu:    cpu,
		ppu:    ppu,
		memory: memory,
	}, nil
}

func (c *Console) Run() {
	ctx := context.Background()

	// Single-length to keep in sync, but allow parallelism
	cpuClock := make(chan struct{}, 1)
	ppuClock := make(chan struct{}, 1)

	go c.cpu.Run(ctx, cpuClock)
	go c.ppu.Run(ctx, ppuClock)

	for {
		// TODO: Timing
		cpuClock <- struct{}{}
		ppuClock <- struct{}{}
	}
}

func (c *Console) RunCycles(cycles uint64) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Single-length to keep in sync, but allow parallelism
	cpuClock := make(chan struct{}, 1)
	ppuClock := make(chan struct{}, 1)

	go c.cpu.Run(ctx, cpuClock)
	go c.ppu.Run(ctx, ppuClock)

	for i := uint64(0); i < cycles; i++ {
		// TODO: Timing
		cpuClock <- struct{}{}
		ppuClock <- struct{}{}
	}
}
