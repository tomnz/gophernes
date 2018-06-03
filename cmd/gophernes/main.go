package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tomnz/gophernes"
	"github.com/tomnz/gophernes/internal/cpu"
	"github.com/tomnz/gophernes/internal/ppu"
)

var (
	rom    = flag.String("rom", "", "ROM file to load")
	cycles = flag.Uint64("cycles", 0, "If non-zero, run for a limited number of master clock cycles")

	cputrace = flag.Bool("cputrace", false, "Include the CPU trace")
	pputrace = flag.Bool("pputrace", false, "Include the PPU trace")
)

func main() {
	flag.Parse()
	if *rom == "" {
		logrus.Fatalf("must specify rom file")
	}
	romFile, err := os.Open(*rom)
	if os.IsNotExist(err) {
		logrus.Fatalf("rom file not found: %q", *rom)
	}

	cpuopts := []cpu.Option{
		cpu.WithTrace(*cputrace),
	}
	ppuopts := []ppu.Option{
		ppu.WithTrace(*pputrace),
	}

	logrus.SetLevel(logrus.DebugLevel)

	console, err := gophernes.NewConsole(romFile, cpuopts, ppuopts)
	if err != nil {
		logrus.Fatal(err)
	}
	if *cycles != 0 {
		console.RunCycles(*cycles)
	} else {
		console.Run()
	}
}
