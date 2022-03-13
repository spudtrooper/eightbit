package convert

import (
	"image"
	"image/gif"
)

type ConvertResult interface {
	Image() image.Image
	GIF() gif.GIF
}

type convertResult struct {
	image image.Image
	gif   gif.GIF
}

func (r *convertResult) Image() image.Image { return r.image }
func (r *convertResult) GIF() gif.GIF       { return r.gif }

func makeImageConvertResult(image image.Image) ConvertResult {
	return &convertResult{image: image}
}

func makeGIFConvertResult(gif gif.GIF) ConvertResult {
	return &convertResult{gif: gif}
}

type Converter interface {
	Convert(input string, inputImage image.Image, opts ConvertOptions) (ConvertResult, error)
	Name() string
	OutputFileName(input string, opts ConvertOptions) string
}

type baseConverter struct {
	name string
	conv func(input string, inputImage image.Image, opts ConvertOptions) (ConvertResult, error)
}

func (c *baseConverter) Name() string { return c.name }
func (c *baseConverter) Convert(input string, inputImage image.Image, opts ConvertOptions) (ConvertResult, error) {
	return c.conv(input, inputImage, opts)
}
