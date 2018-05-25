package gophernes

type mapper uint16

const (
	mapperNROM mapper = 0
)

type cartridge struct {
	prg    []byte
	chr    []byte
	mapper mapper
}
