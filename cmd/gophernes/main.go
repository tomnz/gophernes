package main

import (
	"flag"
	"image"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"

	"github.com/hajimehoshi/ebiten"
	"github.com/sirupsen/logrus"
	"github.com/tomnz/gophernes"
	"github.com/tomnz/gophernes/internal/cpu"
	"github.com/tomnz/gophernes/internal/ppu"
)

var (
	rom    = flag.String("rom", "", "ROM file to load")
	cycles = flag.Uint64("cycles", 0, "If non-zero, run for a limited number of master clock cycles")
	frames = flag.Uint64("frames", 0, "If non-zero, run for a limited number of frames")
	rate   = flag.Float64("rate", 1.0, "Emulation rate - 1.0 runs at normal speed or slower, 0 runs without any delays")

	cputrace = flag.Bool("cputrace", false, "Include the CPU trace")
	pputrace = flag.Bool("pputrace", false, "Include the PPU trace")

	cpuprofile = flag.String("cpuprofile", "", "Write host CPU profile to this file")
	memprofile = flag.String("memprofile", "", "Write host memory profile to this file")
)

var (
	lastFrame   *image.RGBA
	lastFrameMu sync.Mutex
)

func update(screen *ebiten.Image) error {
	lastFrameMu.Lock()
	defer lastFrameMu.Unlock()
	if lastFrame != nil {
		return screen.ReplacePixels(lastFrame.Pix)
	}
	return nil
}

func draw(frame *image.RGBA) {
	lastFrameMu.Lock()
	defer lastFrameMu.Unlock()
	lastFrame = frame
}

func main() {
	flag.Parse()
	if *rom == "" {
		logrus.Fatalf("Must specify rom file!")
	}
	romFile, err := os.Open(*rom)
	if os.IsNotExist(err) {
		logrus.Fatalf("ROM file not found: %q", *rom)
	}

	if *cpuprofile != "" {
		cpuFile, err := os.Create(*cpuprofile)
		if err != nil {
			logrus.Fatalf("Could not create host CPU profile file: %q", *cpuprofile)
		}
		if err := pprof.StartCPUProfile(cpuFile); err != nil {
			logrus.Fatalf("Could not start CPU profile: %s", err)
		}
		defer pprof.StopCPUProfile()
	}

	cpuopts := []cpu.Option{
		cpu.WithTrace(*cputrace),
	}
	ppuopts := []ppu.Option{
		ppu.WithTrace(*pputrace),
	}

	logrus.SetLevel(logrus.DebugLevel)

	console, err := gophernes.NewConsole(
		romFile,
		cpuopts,
		ppuopts,
		gophernes.WithRate(*rate),
		gophernes.WithDraw(draw),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	go func(console *gophernes.Console) {
		if *frames != 0 {
			console.RunFrames(*frames)
		} else if *cycles != 0 {
			console.RunCycles(*cycles)
		} else {
			console.Run()
		}
	}(console)

	if err := ebiten.Run(update, ppu.DisplayWidth, ppu.DisplayHeight, 1, "NES"); err != nil {
		logrus.Fatal(err)
	}

	if *memprofile != "" {
		memFile, err := os.Create(*memprofile)
		if err != nil {
			logrus.Fatalf("Could not create host memory profile file: %q", *memprofile)
		}
		runtime.GC()
		if err := pprof.WriteHeapProfile(memFile); err != nil {
			logrus.Fatalf("Could not start memory profile: %s", err)
		}
		memFile.Close()
	}
}
