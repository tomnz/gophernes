package gophernes

import (
	"io"

	"context"
	"github.com/sirupsen/logrus"
	"github.com/tomnz/gophernes/internal/cartridge"
	"github.com/tomnz/gophernes/internal/cpu"
	"github.com/tomnz/gophernes/internal/ppu"
	"time"
)

// Console implements the main console.
type Console struct {
	config    *config
	ram       []byte
	cpu       *cpu.CPU
	ppu       *ppu.PPU
	cartridge cartridge.Cartridge
}

const (
	internalRAMSize uint16 = 0x800
	frameTime              = 1.0 / 60
)

// NewConsole initializes a new console.
func NewConsole(rom io.Reader, cpuopts []cpu.Option, ppuopts []ppu.Option, opts ...Option) (*Console, error) {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	cartridge, err := loadINES(rom)
	if err != nil {
		return nil, err
	}

	console := &Console{
		config:    config,
		ram:       make([]byte, internalRAMSize),
		cartridge: cartridge,
	}

	cpu := cpu.NewCPU(&cpuMemory{console}, cpuopts...)
	ppu := ppu.NewPPU(&ppuMemory{console}, ppuopts...)

	cpu.Reset()
	ppu.Reset()

	console.cpu = cpu
	console.ppu = ppu

	return console, nil
}

func (c *Console) Run() {
	ctx := context.Background()

	// Single-length to keep in sync, but allow parallelism
	cpuClock := make(chan struct{}, 1)
	ppuClock := make(chan struct{}, 1)

	go c.cpu.Run(ctx, cpuClock)
	go c.ppu.Run(ctx, ppuClock)

	startTime := time.Now()
	var frames uint64

	for {
		// TODO: Timing
		cpuClock <- struct{}{}
		ppuClock <- struct{}{}
		if frames != c.ppu.Frames() {
			frames = c.ppu.Frames()
			c.frameWait(startTime, frames)
		}
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

	startTime := time.Now()
	var frames uint64

	for i := uint64(0); i < cycles; i++ {
		// TODO: Timing
		cpuClock <- struct{}{}
		ppuClock <- struct{}{}
		if frames != c.ppu.Frames() {
			frames = c.ppu.Frames()
			c.frameWait(startTime, frames)
		}
	}
}

func (c *Console) frameWait(startTime time.Time, frames uint64) {
	expected := startTime.Add(time.Duration(float64(time.Second) * frameTime * float64(frames)))
	sleepDuration := expected.Sub(time.Now())
	logrus.Debugf("%s", sleepDuration)
	time.Sleep(sleepDuration)
}
