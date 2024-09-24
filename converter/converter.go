package asciiart

import (
	"image"
)

type Converter struct {
	charset []rune
	scale   int
}

type Options struct {
	scale int
}

func NewConverter(options ...func(*Converter)) *Converter {
	c := &Converter{
		scale:   8,
		charset: []rune(" .:-=+*#%@"),
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func WithScale(scale int) func(*Converter) {
	return func(c *Converter) {
		c.scale = scale
	}
}

func (c *Converter) Convert(img image.Image) string {
	// TODO

	return ""
}
