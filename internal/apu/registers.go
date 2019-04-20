package apu

const (
	regPulse1_1 byte = iota
	regPulse1_2
	regPulse1_3
	regPulse1_4
	regPulse2_1
	regPulse2_2
	regPulse2_3
	regPulse2_4
	regTriangle_1
	// 09 unused
	_
	regTriangle_2
	regTriangle_3
	regNoise_1
	regNoise_2
	regNoise_3
	regDMC_1
	regDMC_2
	regDMC_3
	regDMC_4
	// 14 unused
	_
	regControl
	// 16 unused
	_
	regFrameCounter
)
