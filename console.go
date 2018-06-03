package gophernes

import (
	"io"

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

const (
	cpuClockDivisor = 12
	ppuClockDivisor = 12
)

func (c *Console) Run() {
	startTime := time.Now()
	var clock, frames uint64

	for {
		if clock%cpuClockDivisor == 0 {
			c.cpu.Step()
		}
		if clock%ppuClockDivisor == 0 {
			c.ppu.Step()
			currFrames := c.ppu.Frames()
			if frames != currFrames {
				frames = currFrames
				c.frameWait(startTime, frames)
			}
		}
		clock++
	}
}

func (c *Console) RunCycles(cycles uint64) {
	startTime := time.Now()
	var clock, frames uint64

	for clock < cycles {
		if clock%cpuClockDivisor == 0 {
			c.cpu.Step()
		}
		if clock%ppuClockDivisor == 0 {
			c.ppu.Step()
			currFrames := c.ppu.Frames()
			if frames != currFrames {
				frames = currFrames
				c.frameWait(startTime, frames)
			}
		}
		clock++
	}
}

func (c *Console) frameWait(startTime time.Time, frames uint64) {
	if c.config.rate <= 0 {
		return
	}
	expected := startTime.Add(time.Duration(
		float64(time.Second) * frameTime * float64(frames) / c.config.rate))
	sleepDuration := expected.Sub(time.Now())
	logrus.Debugf("%s", sleepDuration)
	time.Sleep(sleepDuration)
}
