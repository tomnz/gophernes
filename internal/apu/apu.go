package apu

import "fmt"

func NewAPU(irqer IRQer, opts ...Option) *APU {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	return &APU{
		config: config,
		irqer:  irqer,
	}
}

// APU is the main APU of the NES system.
type APU struct {
	config *config
	regs   Registers
	flags  Flags
	irqer  IRQer
}

type IRQer interface {
	IRQ()
}

type Registers struct {
	pulse1,
	pulse2 PulseChannel
	triangle TriangleChannel
	noise    NoiseChannel
	dmc      DMCChannel
}

type PulseChannel struct {
	duty byte
	counterDisable,
	constantVol bool
	envelopePeriod    byte
	sweepEnabled      bool
	sweepPeriod       byte
	sweepNegative     bool
	sweepShiftCount   byte
	timer             uint16
	lengthCounterLoad byte
}

type TriangleChannel struct {
	counterDisable    bool
	counterReloadVal  byte
	timer             uint16
	lengthCounterLoad byte
}

type NoiseChannel struct {
	counterDisable,
	constantVol bool
	envelopePeriod    byte
	loopNoise         bool
	period            byte
	lengthCounterLoad byte
}

type DMCChannel struct {
	irqEnable     bool
	loopSample    bool
	freqIndex     byte
	directLoad    byte
	sampleAddress byte
	sampleLength  byte
}

type Flags struct {
	dmcEnable,
	noiseEnable,
	triangleEnable,
	pulse2Enable,
	pulse1Enable bool
}

func (a *APU) Registers() Registers {
	return a.regs
}

func (a *APU) Flags() Flags {
	return a.flags
}

func (a *APU) Reset() {
}

func (a *APU) Close() {

}

func (a *APU) WriteReg(reg byte, val byte) {
	switch reg {
	case regPulse1_1:
		a.regs.pulse1.duty = val >> 6
		a.regs.pulse1.counterDisable = val>>5&1 == 1
		a.regs.pulse1.constantVol = val>>4&1 == 1
		a.regs.pulse1.envelopePeriod = val & 0xF

	case regPulse1_2:
		a.regs.pulse1.sweepEnabled = val>>7&1 == 1
		a.regs.pulse1.sweepPeriod = val >> 4 & 0x7
		a.regs.pulse1.sweepNegative = val>>3&1 == 1
		a.regs.pulse1.sweepShiftCount = val & 0x7

	case regPulse1_3:
		a.regs.pulse1.timer &= 0xFF00
		a.regs.pulse1.timer |= uint16(val)

	case regPulse1_4:
		a.regs.pulse1.timer &= 0x00FF
		a.regs.pulse1.timer |= uint16(val&0x7) << 8
		a.regs.pulse1.lengthCounterLoad = val >> 3
		// TODO: Reset duty and start envelope?

	case regPulse2_1:
		a.regs.pulse2.duty = val >> 6
		a.regs.pulse2.counterDisable = val>>5&1 == 1
		a.regs.pulse2.constantVol = val>>4&1 == 1
		a.regs.pulse2.envelopePeriod = val & 0xF

	case regPulse2_2:
		a.regs.pulse2.sweepEnabled = val>>7&1 == 1
		a.regs.pulse2.sweepPeriod = val >> 4 & 0x7
		a.regs.pulse2.sweepNegative = val>>3&1 == 1
		a.regs.pulse2.sweepShiftCount = val & 0x7

	case regPulse2_3:
		a.regs.pulse2.timer &= 0xFF00
		a.regs.pulse2.timer |= uint16(val)

	case regPulse2_4:
		a.regs.pulse2.timer &= 0x00FF
		a.regs.pulse2.timer |= uint16(val&0x7) << 8
		a.regs.pulse2.lengthCounterLoad = val >> 3
		// TODO: Reset duty and start envelope?

	case regTriangle_1:
		a.regs.triangle.counterDisable = val>>7&1 == 1
		a.regs.triangle.counterReloadVal = val & 0x7F

	case regTriangle_2:
		a.regs.triangle.timer &= 0xFF00
		a.regs.triangle.timer |= uint16(val)

	case regTriangle_3:
		a.regs.triangle.timer &= 0x00FF
		a.regs.triangle.timer |= uint16(val&0x7) << 8
		a.regs.triangle.lengthCounterLoad = val >> 3
		// TODO: Reload linear counter?

	case regNoise_1:
		a.regs.noise.counterDisable = val>>5&1 == 1
		a.regs.noise.constantVol = val>>4&1 == 1
		a.regs.noise.envelopePeriod = val & 0xF

	case regNoise_2:
		a.regs.noise.loopNoise = val>>7&1 == 1
		a.regs.noise.period = val & 0xF

	case regNoise_3:
		a.regs.noise.lengthCounterLoad = val >> 3

	case regDMC_1:
		a.regs.dmc.irqEnable = val>>7&1 == 1
		a.regs.dmc.loopSample = val>>6&1 == 1
		a.regs.dmc.freqIndex = val & 0xF

	case regDMC_2:
		a.regs.dmc.directLoad = val & 0x7F

	case regDMC_3:
		a.regs.dmc.sampleAddress = val

	case regDMC_4:
		a.regs.dmc.sampleLength = val

	case regControl:
		a.flags.dmcEnable = val>>4&1 == 1
		a.flags.noiseEnable = val>>3&1 == 1
		a.flags.triangleEnable = val>>2&1 == 1
		a.flags.pulse2Enable = val>>1&1 == 1
		a.flags.pulse1Enable = val&1 == 1

	case regFrameCounter:
		// TODO: What is this I don't even

	default:
		//panic(fmt.Sprintf("write to unknown PPU register %#x", reg))

	}
}

func (a *APU) ReadReg(reg byte) byte {
	switch reg {
	case regControl:
		// TODO: Do this shit

	default:
		panic(fmt.Sprintf("read from unknown APU register %#x", reg))
	}
	return 0
}
