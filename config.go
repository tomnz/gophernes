package gophernes

import "image"

type config struct {
	rate    float64
	palette Palette
	draw    func(*image.RGBA)
}

func defaultConfig() *config {
	return &config{
		rate:    1.0,
		palette: defaultPalette(),
	}
}

type Option func(*config)

func WithRate(rate float64) Option {
	return func(config *config) {
		config.rate = rate
	}
}

func WithPalette(palette Palette) Option {
	return func(config *config) {
		config.palette = palette
	}
}

func WithDraw(draw func(rgba *image.RGBA)) Option {
	return func(config *config) {
		config.draw = draw
	}
}
