package gophernes

// https://wiki.nesdev.com/w/index.php/CPU_memory_map

const (
	internalRAMSize = 0x7FF
)

type memory struct {
	ram [internalRAMSize]byte
}
