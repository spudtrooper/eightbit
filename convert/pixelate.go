package convert

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path"
	"strings"

	"github.com/markdaws/go-effects/pkg/effects"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/spudtrooper/goutil/check"
)

type convertPixelatedImageFn func(inputImage, pixelatedImg image.Image, pixelatedWidth, pixelatedHeight int) image.Image

type pixelatedConverter struct {
	convertFn convertPixelatedImageFn
	name      string
}

func (p *pixelatedConverter) Name() string { return p.name }

func (c *pixelatedConverter) OutputFileName(input string, opts ConvertOptions) string {
	ext := path.Ext(input)
	base := strings.Replace(path.Base(input), ext, "", 1)
	return fmt.Sprintf("%s-%s%s", base, c.Name(), ext)
}

func (p *pixelatedConverter) Convert(input string, inputImage image.Image, opts ConvertOptions) (ConvertResult, error) {
	check.Check(p.convertFn != nil, check.CheckMessage(fmt.Sprintf("%s converter has nil convert function", p.Name())))

	pixelated := input + "-pixelated.jpg"
	resized := input + "-resized" + path.Ext(input)
	defer func() {
		if err := os.Remove(pixelated); err != nil {
			log.Printf("trying to delete pixelated: %s: %v", pixelated, err)
		}
		if err := os.Remove(resized); err != nil {
			log.Printf("trying to delete resized: %s: %v", resized, err)
		}
	}()

	// First resize the image to 1280,1280 so that we can apply the effects
	resizedImage := resize.Resize(1280, 1280, inputImage, resize.Lanczos3)
	if err := encodeImage(resized, resizedImage); err != nil {
		return nil, errors.Errorf("writing resized image to %s", resized)
	}

	img, err := effects.LoadImage(resized)
	if err != nil {
		return nil, errors.Errorf("loading image %s: %v", input, err)
	}
	pixelatedEffectsImg, err := effects.NewPixelate(opts.PixelateBlockSize()).Apply(img, 1)
	if err != nil {
		return nil, errors.Errorf("pixelating image %s: %v", input, err)
	}

	saveOpts := effects.SaveOpts{}
	if err := pixelatedEffectsImg.Save(pixelated, saveOpts); err != nil {
		return nil, errors.Errorf("outputting image %s to %s: %v", input, pixelated, err)
	}

	pixelatedImg, err := decode(pixelated)
	if err != nil {
		return nil, errors.Errorf("encoding pixelated image %s: %v", pixelated, err)
	}

	outputImg := p.convertFn(inputImage, pixelatedImg, pixelatedEffectsImg.Width, pixelatedEffectsImg.Height)
	res := makeImageConvertResult(outputImg)

	return res, nil
}

func simpleConvert(inputImage, pixelatedImg image.Image, pixelatedWidth, pixelatedHeight int) image.Image {
	minY, maxY := inputImage.Bounds().Min.Y, inputImage.Bounds().Max.Y
	minX, maxX := inputImage.Bounds().Min.X, inputImage.Bounds().Max.X
	outputImg := image.NewRGBA(image.Rect(minX, minY, maxX, maxY))
	for y := pixelatedImg.Bounds().Min.Y; y < pixelatedImg.Bounds().Max.Y; y++ {
		for x := pixelatedImg.Bounds().Min.X; x < pixelatedImg.Bounds().Max.X; x++ {
			c := pixelatedImg.At(x, y)
			outputImg.Set(x, y, c)
		}
	}

	return outputImg
}

func websafeConvert(inputImage, pixelatedImg image.Image, pixelatedWidth, pixelatedHeight int) image.Image {
	colorsForPalette := func() []color.Color {
		var cols []color.Color
		cols = append(cols,
			color.NRGBA{uint8(0xFF), uint8(0xFF), uint8(0xFF), 255}, // White
			color.NRGBA{uint8(0xC0), uint8(0xC0), uint8(0xC0), 255}, // Silver
			color.NRGBA{uint8(0x80), uint8(0x80), uint8(0x80), 255}, // Gray
			color.NRGBA{uint8(0x00), uint8(0x00), uint8(0x00), 255}, // Black
			color.NRGBA{uint8(0xFF), uint8(0x00), uint8(0x00), 255}, // Red
			color.NRGBA{uint8(0x80), uint8(0x00), uint8(0x00), 255}, // Maroon
			color.NRGBA{uint8(0xFF), uint8(0xFF), uint8(0x00), 255}, // Yellow
			color.NRGBA{uint8(0x80), uint8(0x80), uint8(0x00), 255}, // Olive
			color.NRGBA{uint8(0x00), uint8(0xFF), uint8(0x00), 255}, // Lime
			color.NRGBA{uint8(0x00), uint8(0x80), uint8(0x00), 255}, // Green
			color.NRGBA{uint8(0x00), uint8(0xFF), uint8(0xFF), 255}, // Aqua
			color.NRGBA{uint8(0x00), uint8(0x80), uint8(0x80), 255}, // Teal
			color.NRGBA{uint8(0x00), uint8(0x00), uint8(0xFF), 255}, // Blue
			color.NRGBA{uint8(0x00), uint8(0x00), uint8(0x80), 255}, // Navy
			color.NRGBA{uint8(0xFF), uint8(0x00), uint8(0xFF), 255}, // Fuchsia
			color.NRGBA{uint8(0x80), uint8(0x00), uint8(0x80), 255}, // Purple
		)
		newComps := []int{0x20, 0x60, 0xa0}
		for _, r := range newComps {
			for _, g := range newComps {
				for _, b := range newComps {
					c := color.RGBA{uint8(r), uint8(g), uint8(b), 255}
					cols = append(cols, c)
				}
			}
		}
		return cols
	}
	palette := color.Palette(colorsForPalette())
	palettedImg := image.NewPaletted(image.Rect(0, 0, pixelatedWidth, pixelatedHeight), palette)
	for y := pixelatedImg.Bounds().Min.Y; y < pixelatedImg.Bounds().Max.Y; y++ {
		for x := pixelatedImg.Bounds().Min.X; x < pixelatedImg.Bounds().Max.X; x++ {
			c := pixelatedImg.At(x, y)
			palettedImg.Set(x, y, c)
		}
	}

	return palettedImg
}

func init() {
	globalReg.Register(&pixelatedConverter{
		convertFn: websafeConvert,
		name:      "websafe_pixelated",
	})
	globalReg.Register(&pixelatedConverter{
		convertFn: simpleConvert,
		name:      "pixelated",
	})
}
