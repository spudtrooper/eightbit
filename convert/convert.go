package convert

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"
	"time"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/spudtrooper/goutil/hist"
	"github.com/spudtrooper/goutil/io"
	"github.com/spudtrooper/goutil/or"
)

var AllConverters = []string{"pixelate", "overlap_mean", "overlap_median"}

type Converter func(input string, inputImage image.Image, opts ConvertOptions) (image.Image, error)

func Convert(input string, cOpts ...ConvertOption) error {
	opts := MakeConvertOptions(cOpts...)

	switch ext := path.Ext(input); ext {
	case ".png", ".jpg", ".jpeg":
	default:
		return errors.Errorf("invalid input image type: %s", input)
	}

	inputImage, err := decode(input)
	if err != nil {
		return errors.Errorf("decoding input image: %s", input)
	}

	if opts.ColorHist() {
		colorHist := hist.MakeHistogram()
		for y := inputImage.Bounds().Min.Y; y < inputImage.Bounds().Max.Y; y++ {
			for x := inputImage.Bounds().Min.X; x < inputImage.Bounds().Max.X; x++ {
				c := inputImage.At(x, y)
				colorHist.Add(colorName(c), 1)
			}
		}
		fmt.Println("Printing color histogram...")
		fmt.Println(hist.HistString(colorHist))
	}

	converters := opts.Converters()
	if len(converters) == 1 && converters[0] == "all" {
		converters = AllConverters
	}
	if len(converters) == 0 {
		if opts.ColorHist() {
			return nil
		}
		return errors.Errorf("you must specify at least one converter")
	}
	if len(converters) > 1 && opts.OutputFile() != "" {
		return errors.Errorf("you cannot specify an output with >1 converter")
	}

	for _, c := range converters {
		output := or.String(opts.OutputFile(), makeOutput(input, opts.OutputDir(), c))
		if !opts.Force() && io.FileExists(output) {
			return errors.Errorf("%s exists. pass --force to write anyway", output)
		}
		if err := convertOne(inputImage, input, output, c, opts); err != nil {
			return errors.Errorf("converting %s to %s: %v", input, output)
		}
	}

	return nil
}

func convertOne(inputImage image.Image, input, output string, converterName string, opts ConvertOptions) error {
	start := time.Now()

	conv := makeConverter(converterName)
	if conv == nil {
		return errors.Errorf("invalid converter string: %s", converterName)
	}

	outputImg, err := conv(input, inputImage, opts)
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

	log.Printf("converted %s to %s in %v", input, output, time.Since(start))

	return nil
}

func makeOutput(input, outputDir, converter string) string {
	dir := or.String(outputDir, path.Dir(input))
	ext := path.Ext(input)
	base := strings.Replace(path.Base(input), ext, "", 1)
	return path.Join(dir, base+"-"+converter+ext)
}

func makeConverter(s string) Converter {
	switch s {
	case "pixelate":
		return pixelatedConverter
	case "overlap_mean":
		return overlapMeanConverter
	case "overlap_median":
		return overlapMedianConverter
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
		return errors.Errorf("closing %s: %v", output, err)
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
