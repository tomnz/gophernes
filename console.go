package gophernes

import (
	"io"
)

// Console implements the main console.
type Console struct {
	cpu       *cpu
	memory    *memory
	cartridge *cartridge
}

// NewConsole initializes a new console.
func NewConsole() *Console {
	cpu := &cpu{}
	cpu.init()

	memory := &memory{}
	memory.init()

	return &Console{
		cpu:    cpu,
		memory: memory,
	}
}

// LoadROM loads the given ROM data into the console.
func (c *Console) LoadROM(rom io.Reader) error {
	// TODO: Formats other than iNES?
	cartridge, err := loadINES(rom)
	if err != nil {
		return err
	}

	c.cartridge = cartridge
	return c.loadCartridge()
}

func (c *Console) loadCartridge() error {

}
