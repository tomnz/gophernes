package ppu

// Adapted from: https://github.com/fogleman/nes/blob/master/nes/ppu.go
// Over time, this will be refactored and reworked into the core files. Trying to get unblocked to actually show
// some pixels on the screen in the meantime.

// Copyright (C) 2015 Michael Fogleman
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and
// to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of
// the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO
// THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
// CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

func (p *PPU) Step() {
	// Perform the vertical blank handling as needed
	p.stepNMI()

	renderEnabled := p.regs.ShowBackground || p.regs.ShowSprites
	preLine := p.scanLine == 261
	visibleLine := p.scanLine < 240
	// postLine := ppu.ScanLine == 240
	renderLine := preLine || visibleLine
	preFetchCycle := p.lineCycle >= 321 && p.lineCycle <= 336
	visibleCycle := p.lineCycle >= 1 && p.lineCycle <= 256
	fetchCycle := preFetchCycle || visibleCycle

	if renderEnabled {
		if visibleLine && visibleCycle {
			p.renderPixel()
		}
		if renderLine && fetchCycle {
			p.tileData <<= 4
			switch p.lineCycle % 8 {
			case 1:
				p.fetchNameTableByte()
			case 3:
				p.fetchAttributeTableByte()
			case 5:
				p.fetchLowTileByte()
			case 7:
				p.fetchHighTileByte()
			case 0:
				p.storeTileData()
			}
		}

		if renderLine {
			if fetchCycle && p.lineCycle%8 == 0 {
				p.incrementX()
			}
			if p.lineCycle == 256 {
				p.incrementY()
			}
			if p.lineCycle == 257 {
				p.resetX()
			}
		}

		if preLine && p.lineCycle >= 280 && p.lineCycle <= 304 {
			p.resetY()
		}
	}

	if renderEnabled {
		if p.lineCycle == 257 {
			if visibleLine {
				p.evaluateSprites()
			} else {
				p.spriteCount = 0
			}
		}
	}

	// Move to the next cell/scanline as needed
	p.stepScan()

	p.cycles++
}

func (p *PPU) fetchNameTableByte() {
	v := p.vramAddr
	address := 0x2000 | (v & 0x0FFF)
	p.nameTableByte = p.read8(address)
}

func (p *PPU) readPalette(address uint16) byte {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	return p.paletteData[address]
}

func (p *PPU) writePalette(addr uint16, val byte) {
	if addr >= 16 && addr%4 == 0 {
		addr -= 16
	}
	p.paletteData[addr] = val
}

func (p *PPU) fetchAttributeTableByte() {
	v := p.vramAddr
	address := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
	shift := ((v >> 4) & 4) | (v & 2)
	p.attributeTableByte = ((p.read8(address) >> shift) & 3) << 2
}

func (p *PPU) fetchLowTileByte() {
	fineY := (p.vramAddr >> 12) & 7
	tile := p.nameTableByte
	address := p.regs.BackgroundPatternTableAddress + uint16(tile)*16 + fineY
	p.lowTileByte = p.read8(address)
}

func (p *PPU) fetchHighTileByte() {
	fineY := (p.vramAddr >> 12) & 7
	tile := p.nameTableByte
	address := p.regs.BackgroundPatternTableAddress + uint16(tile)*16 + fineY
	p.highTileByte = p.read8(address + 8)
}

func (p *PPU) storeTileData() {
	var data uint32
	for i := 0; i < 8; i++ {
		a := p.attributeTableByte
		p1 := (p.lowTileByte & 0x80) >> 7
		p2 := (p.highTileByte & 0x80) >> 6
		p.lowTileByte <<= 1
		p.highTileByte <<= 1
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	p.tileData |= uint64(data)
}

func (p *PPU) fetchTileData() uint32 {
	return uint32(p.tileData >> 32)
}

func (p *PPU) backgroundPixel() byte {
	if !p.regs.ShowBackground {
		return 0
	}
	data := p.fetchTileData() >> ((7 - p.fineX) * 4)
	return byte(data & 0x0F)
}

func (p *PPU) spritePixel() (byte, byte) {
	if !p.regs.ShowSprites {
		return 0, 0
	}
	for i := 0; i < p.spriteCount; i++ {
		offset := (p.lineCycle - 1) - int(p.spritePositions[i])
		if offset < 0 || offset > 7 {
			continue
		}
		offset = 7 - offset
		color := byte((p.spritePatterns[i] >> byte(offset*4)) & 0x0F)
		if color%4 == 0 {
			continue
		}
		return byte(i), color
	}
	return 0, 0
}

func (p *PPU) renderPixel() {
	x := p.lineCycle - 1
	y := p.scanLine
	background := p.backgroundPixel()
	i, sprite := p.spritePixel()
	if x < 8 && !p.regs.ShowLeftBackground {
		background = 0
	}
	if x < 8 && !p.regs.ShowLeftSprites {
		sprite = 0
	}
	b := background%4 != 0
	s := sprite%4 != 0
	var color byte
	if !b && !s {
		color = 0
	} else if !b && s {
		color = sprite | 0x10
	} else if b && !s {
		color = background
	} else {
		if p.spriteIndexes[i] == 0 && x < 255 {
			p.sprite0Hit = true
		}
		if p.spritePriorities[i] == 0 {
			color = sprite | 0x10
		} else {
			color = background
		}
	}
	p.backBuffer[y][x] = p.readPalette(uint16(color)) % 64
}

func (p *PPU) fetchSpritePattern(i, row int) uint32 {
	tile := p.oam[i*4+1]
	attributes := p.oam[i*4+2]
	var address uint16
	if !p.regs.TallSprites {
		if attributes&0x80 == 0x80 {
			row = 7 - row
		}
		address = p.regs.SpritePatternTableAddress + uint16(tile)*16 + uint16(row)
	} else {
		if attributes&0x80 == 0x80 {
			row = 15 - row
		}
		table := tile & 1
		tile &= 0xFE
		if row > 7 {
			tile++
			row -= 8
		}
		address = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
	}
	a := (attributes & 3) << 2
	lowTileByte := p.read8(address)
	highTileByte := p.read8(address + 8)
	var data uint32
	for i := 0; i < 8; i++ {
		var p1, p2 byte
		if attributes&0x40 == 0x40 {
			p1 = (lowTileByte & 1) << 0
			p2 = (highTileByte & 1) << 1
			lowTileByte >>= 1
			highTileByte >>= 1
		} else {
			p1 = (lowTileByte & 0x80) >> 7
			p2 = (highTileByte & 0x80) >> 6
			lowTileByte <<= 1
			highTileByte <<= 1
		}
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	return data
}

func (p *PPU) evaluateSprites() {
	var h int
	if !p.regs.TallSprites {
		h = 8
	} else {
		h = 16
	}
	count := 0
	for i := 0; i < 64; i++ {
		y := p.oam[i*4+0]
		a := p.oam[i*4+2]
		x := p.oam[i*4+3]
		row := p.scanLine - int(y)
		if row < 0 || row >= h {
			continue
		}
		if count < 8 {
			p.spritePatterns[count] = p.fetchSpritePattern(i, row)
			p.spritePositions[count] = x
			p.spritePriorities[count] = (a >> 5) & 1
			p.spriteIndexes[count] = byte(i)
		}
		count++
	}
	if count > 8 {
		count = 8
		p.spriteOverflow = true
	}
	p.spriteCount = count
}
