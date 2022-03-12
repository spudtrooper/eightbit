package convert

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/spudtrooper/goutil/io"
)

type Converter func(input string, inputImage image.Image, opts ConvertOptions) (image.Image, error)

func Convert(input, output string, cOpts ...ConvertOption) error {
	opts := MakeConvertOptions(cOpts...)

	if !opts.Force() && io.FileExists(output) {
		return errors.Errorf("%s exists. pass --force to write anyway", output)
	}

	inputImage, err := decode(input)
	if err != nil {
		return errors.Errorf("decoding input image: %s", input)
	}

	converter := opts.Converter()
	if converter == nil {
		converter = pixelatedConverter
		// converter = overlapConverter
	}

	outputImg, err := converter(input, inputImage, opts)
	if err != nil {
		return errors.Errorf("converting image: %v", err)
	}
	if outputImg == nil {
		return errors.Errorf("converting image returned nil image")
	}

	if opts.ResizeWidth() != 0 && opts.ResizeHeight() != 0 {
		outputImg = resize.Resize(opts.ResizeWidth(), opts.ResizeHeight(), outputImg, resize.Lanczos3)
	}

	if _, err := io.MkdirAll(path.Dir(output)); err != nil {
		return errors.Errorf("making directory for %s", output)
	}
	if err := encode(output, outputImg); err != nil {
		return errors.Errorf("encoding image to %s: %v", output, err)
	}

	log.Printf("converted %s to %s", input, output)

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
