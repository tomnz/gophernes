package ppu

type config struct {
	trace bool
}

func defaultConfig() *config {
	return &config{
		trace: false,
	}
}

type Option func(*config)

func WithTrace() Option {
	return func(config *config) {
		config.trace = true
	}
}
