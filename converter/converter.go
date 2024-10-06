package asciiart

import (
	"image"
)

type Converter struct {
	img        image.Image
	charset    []rune
	width      int
	dOpts      DoGOptions
	sThreshold float32
	eThreshold float32
	squash     float32
}

type Options struct {
	width int
}

func NewConverter(img image.Image, options ...func(*Converter)) *Converter {
	c := &Converter{
		img:        img,
		width:      100,
		charset:    []rune(" .:-=+*#%@"),
		dOpts:      DoGOptions{Sigma1: 1, Sigma2: 1.5, Epsilon: 0.65, Tau: 0.8, Phi: 25},
		sThreshold: 0.2,
		eThreshold: 0.05,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func WithWidth(width int) func(*Converter) {
	return func(c *Converter) {
		c.width = width
	}
}

func WithCharset(charset []rune) func(*Converter) {
	return func(c *Converter) {
		c.charset = charset
	}
}

func WithDSigma1(sigma float32) func(*Converter) {
	return func(c *Converter) {
		c.dOpts.Sigma1 = sigma
	}
}

func WithDSigma2(sigma float32) func(*Converter) {
	return func(c *Converter) {
		c.dOpts.Sigma2 = sigma
	}
}

func WithDEpsilon(epsilon float32) func(*Converter) {
	return func(c *Converter) {
		c.dOpts.Epsilon = epsilon
	}
}

func WithDTau(tau float32) func(*Converter) {
	return func(c *Converter) {
		c.dOpts.Tau = tau
	}
}

func WithDPhi(phi float32) func(*Converter) {
	return func(c *Converter) {
		c.dOpts.Phi = phi
	}
}

func WithSThreshold(threshold float32) func(*Converter) {
	return func(c *Converter) {
		c.sThreshold = threshold
	}
}

func WithEThreshold(threshold float32) func(*Converter) {
	return func(c *Converter) {
		c.sThreshold = threshold
	}
}

func WithSquash(squash float32) func(*Converter) {
	return func(c *Converter) {
		c.squash = squash
	}
}

func (c *Converter) Convert(img image.Image) string {
	// TODO

	return ""
}
