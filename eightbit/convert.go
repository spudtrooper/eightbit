package eightbit

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"

	"github.com/markdaws/go-effects/pkg/effects"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

func Convert(input, output string, cOpts ...ConvertOption) error {
	opts := MakeConvertOptions(cOpts...)

	pixelated := "pixelated.jpg"

	// First resize the image to 1280,1280 so that we can apply the effects
	inputImage, err := decode(input)
	if err != nil {
		return errors.Errorf("decoding input image: %s", input)
	}
	resizedImage := resize.Resize(1280, 1280, inputImage, resize.Lanczos3)
	resized := "resized" + path.Ext(input)
	if err := encode(resized, resizedImage); err != nil {
		return errors.Errorf("writing resized image to %s", resized)
	}

	img, err := effects.LoadImage(resized)
	if err != nil {
		return errors.Errorf("loading image %s: %v", input, err)
	}
	pixelatedEffectsImg, err := effects.NewPixelate(opts.BlockSize()).Apply(img, 1)
	if err != nil {
		return errors.Errorf("pixelating image %s: %v", input, err)
	}

	saveOpts := effects.SaveOpts{}
	if err := pixelatedEffectsImg.Save(pixelated, saveOpts); err != nil {
		return errors.Errorf("outputting image %s to %s: %v", input, output, err)
	}

	pixelatedImg, err := decode(pixelated)
	if err != nil {
		return errors.Errorf("encoding pixelated image %s: %v", pixelated, err)
	}

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

	var outputImg image.Image
	palette := color.Palette(colorsForPalette())
	palettedImg := image.NewPaletted(image.Rect(0, 0, pixelatedEffectsImg.Width, pixelatedEffectsImg.Height), palette)
	for y := pixelatedImg.Bounds().Min.Y; y < pixelatedImg.Bounds().Max.Y; y++ {
		for x := pixelatedImg.Bounds().Min.X; x < pixelatedImg.Bounds().Max.X; x++ {
			c := pixelatedImg.At(x, y)
			palettedImg.Set(x, y, c)
		}
	}
	outputImg = palettedImg

	if opts.ResizeWidth() != 0 && opts.ResizeHeight() != 0 {
		outputImg = resize.Resize(opts.ResizeWidth(), opts.ResizeHeight(), outputImg, resize.Lanczos3)
	}

	if err := encode(output, outputImg); err != nil {
		return errors.Errorf("encoding image to %s: %v", output, err)
	}

	return nil
}

func encode(output string, outputImg image.Image) error {
	out, err := os.Create(output)
	if err != nil {
		return errors.Errorf("creating output image from %s: %v", output, err)
	}

	switch ext := strings.ToLower(path.Ext(output)); ext {
	case ".png":
		png.Encode(out, outputImg)
	case ".jpg", ".jpeg":
		jpeg.Encode(out, outputImg, &jpeg.Options{})
	default:
		return errors.Errorf("unknown output image format for %s:", output)
	}

	if err := out.Close(); err != nil {
		return errors.Errorf("closing %s: %v", output)
	}

	return nil
}

func decode(input string) (image.Image, error) {
	inputFile, err := os.Open(input)
	if err != nil {
		return nil, errors.Errorf("opening %s: %v", input, err)
	}
	defer inputFile.Close()

	switch ext := strings.ToLower(path.Ext(input)); ext {
	case ".png":
		img, err := png.Decode(inputFile)
		if err != nil {
			return nil, errors.Errorf("decoding png %s: %v", input, err)
		}
		return img, nil
	case ".jpg", ".jpeg":
		img, err := jpeg.Decode(inputFile)
		if err != nil {
			return nil, errors.Errorf("decoding jpg %s: %v", input, err)
		}
		return img, nil
	}

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return nil, errors.Errorf("decoding %s: %v", input, err)
	}

	return img, nil
}
