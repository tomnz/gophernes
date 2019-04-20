package apu

import "fmt"

func NewAPU(irqer IRQer, opts ...Option) *APU {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	return &APU{
		config:   config,
		irqer:    irqer,
		buf:      newBuffer(int(config.sampleRate)),
		pulse1:   newPulseChannel(),
		pulse2:   newPulseChannel(),
		triangle: newTriangleChannel(),
		noise:    newNoiseChannel(),
		dmc:      newDMCChannel(),
	}
}

// APU is the main APU of the NES system.
type APU struct {
	config *config
	cycles uint64

	flags Flags
	irqer IRQer
	buf   *buffer

	// Channels
	pulse1,
	pulse2 *pulseChannel
	triangle *triangleChannel
	noise    *noiseChannel
	dmc      *dmcChannel
}

type IRQer interface {
	IRQ()
}

func newPulseChannel() *pulseChannel {
	return &pulseChannel{}
}

type pulseChannel struct {
	// Registers
	duty byte
	counterDisable,
	constantVol bool
	envelopePeriod  byte
	sweepEnabled    bool
	sweepPeriod     byte
	sweepNegative   bool
	sweepShiftCount byte
	timerLoad       uint16
	dividerLoad     byte

	// State
	divider byte
	timer   uint16
}

func newTriangleChannel() *triangleChannel {
	return &triangleChannel{}
}

type triangleChannel struct {
	counterDisable   bool
	counterReloadVal byte
	timer            uint16
	dividerLoad      byte
}

func newNoiseChannel() *noiseChannel {
	return &noiseChannel{}
}

type noiseChannel struct {
	counterDisable,
	constantVol bool
	envelopePeriod byte
	loopNoise      bool
	period         byte
	dividerLoad    byte
}

func newDMCChannel() *dmcChannel {
	return &dmcChannel{}
}

type dmcChannel struct {
	irqEnable     bool
	loopSample    bool
	freqIndex     byte
	directLoad    byte
	sampleAddress byte
	sampleLength  byte
}

func (p *pulseChannel) step() {
	//if p.timer == 0 {
	//	p.timer =
	//}
	//
	//if p.divider == 0 {
	//	tick = true
	//	p.divider = p.dividerLoad
	//} else {
	//	p.divider--
	//}
	//
	//if !tick {
	//	return
	//}

}

type Flags struct {
	dmcEnable,
	noiseEnable,
	triangleEnable,
	pulse2Enable,
	pulse1Enable bool
}

func (a *APU) Flags() Flags {
	return a.flags
}

func (a *APU) Reset() {
}

func (a *APU) Step() {
}

func (a *APU) Close() {
}

func (a *APU) Buffer() *buffer {
	return a.buf
}

func (a *APU) WriteReg(reg byte, val byte) {
	switch reg {
	case regPulse1_1:
		a.pulse1.duty = val >> 6
		a.pulse1.counterDisable = val>>5&1 == 1
		a.pulse1.constantVol = val>>4&1 == 1
		a.pulse1.envelopePeriod = val & 0xF

	case regPulse1_2:
		a.pulse1.sweepEnabled = val>>7&1 == 1
		a.pulse1.sweepPeriod = val >> 4 & 0x7
		a.pulse1.sweepNegative = val>>3&1 == 1
		a.pulse1.sweepShiftCount = val & 0x7

	case regPulse1_3:
		a.pulse1.timer &= 0xFF00
		a.pulse1.timer |= uint16(val)

	case regPulse1_4:
		a.pulse1.timer &= 0x00FF
		a.pulse1.timer |= uint16(val&0x7) << 8
		a.pulse1.dividerLoad = val >> 3
		// TODO: Reset duty and start envelope?

	case regPulse2_1:
		a.pulse2.duty = val >> 6
		a.pulse2.counterDisable = val>>5&1 == 1
		a.pulse2.constantVol = val>>4&1 == 1
		a.pulse2.envelopePeriod = val & 0xF

	case regPulse2_2:
		a.pulse2.sweepEnabled = val>>7&1 == 1
		a.pulse2.sweepPeriod = val >> 4 & 0x7
		a.pulse2.sweepNegative = val>>3&1 == 1
		a.pulse2.sweepShiftCount = val & 0x7

	case regPulse2_3:
		a.pulse2.timer &= 0xFF00
		a.pulse2.timer |= uint16(val)

	case regPulse2_4:
		a.pulse2.timer &= 0x00FF
		a.pulse2.timer |= uint16(val&0x7) << 8
		a.pulse2.dividerLoad = val >> 3
		// TODO: Reset duty and start envelope?

	case regTriangle_1:
		a.triangle.counterDisable = val>>7&1 == 1
		a.triangle.counterReloadVal = val & 0x7F

	case regTriangle_2:
		a.triangle.timer &= 0xFF00
		a.triangle.timer |= uint16(val)

	case regTriangle_3:
		a.triangle.timer &= 0x00FF
		a.triangle.timer |= uint16(val&0x7) << 8
		a.triangle.dividerLoad = val >> 3
		// TODO: Reload linear counter?

	case regNoise_1:
		a.noise.counterDisable = val>>5&1 == 1
		a.noise.constantVol = val>>4&1 == 1
		a.noise.envelopePeriod = val & 0xF

	case regNoise_2:
		a.noise.loopNoise = val>>7&1 == 1
		a.noise.period = val & 0xF

	case regNoise_3:
		a.noise.dividerLoad = val >> 3

	case regDMC_1:
		a.dmc.irqEnable = val>>7&1 == 1
		a.dmc.loopSample = val>>6&1 == 1
		a.dmc.freqIndex = val & 0xF

	case regDMC_2:
		a.dmc.directLoad = val & 0x7F

	case regDMC_3:
		a.dmc.sampleAddress = val

	case regDMC_4:
		a.dmc.sampleLength = val

	case regControl:
		a.flags.dmcEnable = val>>4&1 == 1
		a.flags.noiseEnable = val>>3&1 == 1
		a.flags.triangleEnable = val>>2&1 == 1
		a.flags.pulse2Enable = val>>1&1 == 1
		a.flags.pulse1Enable = val&1 == 1

	case regFrameCounter:
		// TODO: What is this I don't even

	default:
		// panic(fmt.Sprintf("write to unknown APU register %#x", reg))

	}
}

func (a *APU) ReadReg(reg byte) byte {
	switch reg {
	case regControl:
		// TODO: Do this shit
		return 0

	}
	panic(fmt.Sprintf("read from unknown APU register %#x", reg))
}
