package main

import (
	"flag"
	"github.com/tomnz/gophernes"
	"log"
	"os"
)

var (
	rom   = flag.String("rom", "", "ROM file to load")
	steps = flag.Int("steps", 0, "If non-zero, run for a limited number of CPU operations")
)

func main() {
	flag.Parse()
	if *rom == "" {
		log.Fatalf("must specify rom file")
	}
	romFile, err := os.Open(*rom)
	if os.IsNotExist(err) {
		log.Fatalf("rom file not found: %q", *rom)
	}
	console, err := gophernes.NewConsole(romFile)
	if err != nil {
		log.Fatal(err)
	}
	if *steps != 0 {
		err = console.RunSteps(*steps)
	} else {
		err = console.Run()
	}
	if err != nil {
		log.Fatal(err)
	}
}
