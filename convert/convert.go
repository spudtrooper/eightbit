package convert

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"
	"time"

	"github.com/nfnt/resize"
	"github.com/noelyahan/impexp"
	"github.com/noelyahan/mergi"
	"github.com/pkg/errors"
	"github.com/spudtrooper/goutil/hist"
	"github.com/spudtrooper/goutil/io"
	"github.com/spudtrooper/goutil/or"
)

func Convert(input string, cOpts ...ConvertOption) ([]string, error) {
	opts := MakeConvertOptions(cOpts...)

	switch ext := path.Ext(input); ext {
	case ".png", ".jpg", ".jpeg":
	default:
		return nil, errors.Errorf("invalid input image type: %s", input)
	}

	inputImage, err := decode(input)
	if err != nil {
		return nil, errors.Errorf("decoding input image: %s", input)
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
		converters = globalReg.AllConverterNames()
	}
	if len(converters) == 0 {
		if opts.ColorHist() {
			return nil, nil
		}
		return nil, errors.Errorf("you must specify at least one converter")
	}
	if len(converters) > 1 && opts.OutputFile() != "" {
		return nil, errors.Errorf("you cannot specify an output with >1 converter")
	}

	var outputs []string
	for _, convName := range converters {
		conv := globalReg.Get(convName)
		if conv == nil {
			return nil, errors.Errorf("invalid converter string: %s", convName)
		}
		output := or.String(opts.OutputFile(), makeOutput(conv, input, opts.OutputDir(), opts))
		if !opts.Force() && io.FileExists(output) {
			return nil, errors.Errorf("%s exists. pass --force to write anyway", output)
		}
		if err := convertOne(inputImage, input, output, conv, opts); err != nil {
			return nil, errors.Errorf("converting %s to %s: %v", input, output, err)
		}
		outputs = append(outputs, output)
	}

	return outputs, nil
}

func convertOne(inputImage image.Image, input, output string, conv Converter, opts ConvertOptions) error {
	start := time.Now()

	outputImgRes, err := conv.Convert(input, inputImage, opts)
	if err != nil {
		return errors.Errorf("converting image: %v", err)
	}
	if outputImgRes == nil {
		return errors.Errorf("converting image returned nil image")
	}

	if opts.ResizeWidth() != 0 && opts.ResizeHeight() != 0 {
		outputImg := resize.Resize(opts.ResizeWidth(), opts.ResizeHeight(), outputImgRes.Image(), resize.Lanczos3)
		outputImgRes = makeImageConvertResult(outputImg)
	}

	if !opts.Force() && io.FileExists(output) {
		return errors.Errorf("%s exists. pass --force to write anyway", output)
	}
	if _, err := io.MkdirAll(path.Dir(output)); err != nil {
		return errors.Errorf("making directory for %s", output)
	}
	if err := encode(output, outputImgRes); err != nil {
		return errors.Errorf("encoding image to %s: %v", output, err)
	}

	log.Printf("converted %s to %s in %v", input, output, time.Since(start))

	return nil
}

func makeOutput(c Converter, input, outputDir string, opts ConvertOptions) string {
	dir := or.String(outputDir, path.Dir(input))
	output := c.OutputFileName(input, opts)
	return path.Join(dir, output)
}

func encode(output string, res ConvertResult) error {
	if res.Image() != nil {
		return encodeImage(output, res.Image())
	}
	if len(res.GIF().Image) > 0 {
		return encodeGIF(output, res.GIF())
	}
	return errors.Errorf("no image in result")
}

func encodeGIF(output string, gif gif.GIF) error {
	return mergi.Export(impexp.NewAnimationExporter(gif, output))
}

func encodeImage(output string, outputImg image.Image) error {
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
