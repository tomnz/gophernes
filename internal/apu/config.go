package apu

type config struct {
	trace      bool
	sampleRate uint32
}

const DefaultSampleRate = 44100

func defaultConfig() *config {
	return &config{
		trace:      false,
		sampleRate: DefaultSampleRate,
	}
}

type Option func(*config)

func WithTrace(trace bool) Option {
	return func(config *config) {
		config.trace = trace
	}
}

func WithSampleRate(sampleRate uint32) Option {
	return func(config *config) {
		config.sampleRate = sampleRate
	}
}
