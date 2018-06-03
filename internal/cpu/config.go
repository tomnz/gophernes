package cpu

type config struct {
	trace bool
}

func defaultConfig() *config {
	return &config{
		trace: false,
	}
}

type Option func(*config)

func WithTrace(trace bool) Option {
	return func(config *config) {
		config.trace = trace
	}
}
