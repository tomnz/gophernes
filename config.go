package gophernes

type config struct {
	rate float64
}

func defaultConfig() *config {
	return &config{
		rate: 1.0,
	}
}

type Option func(*config)

func WithRate(rate float64) Option {
	return func(config *config) {
		config.rate = rate
	}
}
