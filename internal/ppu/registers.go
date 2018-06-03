package ppu

const (
	regController byte = iota
	regMask
	regStatus
	regOAMAddress
	regOAMData
	regScroll
	regAddress
	regData
)

const (
	addrNametable0 uint16 = 0x2000
	addrNametable1 uint16 = 0x2400
	addrNametable2 uint16 = 0x2800
	addrNametable3 uint16 = 0x2C00
)

const (
	addrPatternTable0 uint16 = 0
	addrPatternTable1 uint16 = 0x1000
)

type Registers struct {
	// Control
	BaseNametableAddress uint16
	VRAMAddressIncrement uint16
	SpritePatternTableAddress,
	BackgroundPatternTableAddress uint16
	TallSprites,
	NMIGenerate bool
	// Status
	NMIOccurred bool
	// Mask
	Grayscale,
	ShowLeftBackground,
	ShowLeftSprites,
	ShowBackground,
	ShowSprites,
	EmphasizeRed,
	EmphasizeGreen,
	EmphasizeBlue bool
	// OAM Address
	OAMAddr byte
	// Scroll
	ScrollX,
	ScrollY byte
	// VRAM Address
	VRAMAddr uint16
}
