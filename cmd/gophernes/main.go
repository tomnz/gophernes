package main

import (
	"flag"
	"github.com/tomnz/gophernes"
	"log"
	"os"
)

var (
	rom = flag.String("rom", "", "ROM file to load")
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
	console.Run()
}
