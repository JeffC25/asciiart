package asciiart

import "image"

type Converter struct {
	img   image.Image
	scale int
}

type Options struct {
	scale int
}

func NewConverter(img image.Image, opts ...func(*Converter)) *Converter {
	c := &Converter{
		img:   img,
		scale: 8,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithScale(scale int) func(*Converter) {
	return func(c *Converter) {
		c.scale = scale
	}
}
