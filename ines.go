package gophernes

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/tomnz/gophernes/internal/cartridge"
)

const (
	inesMagic        = 0x1a53454e
	prgLenMultiplier = 16384
	chrLenMultiplier = 8192
)

type inesHeader struct {
	Magic uint32
	PrgLen,
	ChrLen,
	Flags6,
	Flags7,
	Flags8,
	Flags9,
	Flags10,
	Flags11,
	Flags12,
	Flags13 byte
	_ [2]byte
}

func loadINES(file io.Reader) (cartridge.Cartridge, error) {
	header := inesHeader{}
	if err := binary.Read(file, binary.LittleEndian, &header); err != nil {
		return nil, err
	}
	if header.Magic != inesMagic {
		return nil, errors.New("does not appear to be an iNES file: invalid header")
	}

	// TODO: Detect iNES version
	// https://wiki.nesdev.com/w/index.php/INES#Variant_comparison
	// For now, assume iNES 1.0

	mapper := uint16(header.Flags6>>4 | header.Flags7&0xf0)

	if (header.Flags6>>3)&1 == 1 {
		// Trainer is present in the ROM - ignore
		if _, err := io.ReadFull(file, make([]byte, 512)); err != nil {
			return nil, err
		}
	}

	prg := make([]byte, prgLenMultiplier*int(header.PrgLen))
	if _, err := io.ReadFull(file, prg); err != nil {
		return nil, err
	}

	var chr []byte
	if header.ChrLen == 0 {
		// Special case - provide an empty block
		chr = make([]byte, 8192)
	} else {
		chr = make([]byte, chrLenMultiplier*int(header.ChrLen))
		if _, err := io.ReadFull(file, chr); err != nil {
			return nil, err
		}
	}

	return cartridge.NewCartridge(mapper, prg, chr)
}
