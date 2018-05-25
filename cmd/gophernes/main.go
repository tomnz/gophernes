package main

import (
	"flag"
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
	_, err := os.Open(*rom)
	if os.IsNotExist(err) {
		log.Fatalf("rom file not found: %q", *rom)
	}
}
