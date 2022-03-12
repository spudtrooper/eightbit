package main

import (
	"flag"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spudtrooper/eightbit/convert"
	"github.com/spudtrooper/eightbit/gitversion"
	"github.com/spudtrooper/goutil/check"
	"github.com/spudtrooper/goutil/slice"
)

var (
	input           = flag.String("input", "", "input image")
	output          = flag.String("output", "", "output image")
	blockSize       = flag.Int("block_size", 16, "blocksize for downsampling")
	resizeHeight    = flag.Int("resize_height", 0, "height in pixels of the final image; must be used with --resize_width")
	resizeWidth     = flag.Int("resize_width", 0, "width in pixels of the final image; must be used with --resize_height")
	force           = flag.Bool("force", false, "overwrite existing files")
	converters      = flag.String("converters", "pixelate", "the kinds of converter to use or 'all' for all of them. If you don't specify an output file, the output file will be next to the source file with this tag at the end of the base name.")
	printConverters = flag.Bool("print_converters", false, "print the names of all the converters and exit")
	colorHist       = flag.Bool("color_hist", false, "print a histogram of web colors from the input image")
)

func realMain() error {
	if gitversion.CheckVersionFlag() {
		return nil
	}

	if *printConverters {
		fmt.Println("Printing the names of all the converters...")
		for i, c := range convert.AllConverters {
			fmt.Printf("  [%d] %s\n", i+1, c)
		}
		return nil
	}

	if *input == "" {
		return errors.Errorf("--input required")
	}

	if err := convert.Convert(*input,
		convert.ConvertOutputFile(*output),
		convert.ConvertBlockSize(*blockSize),
		convert.ConvertResizeWidth(uint(*resizeHeight)),
		convert.ConvertResizeHeight(uint(*resizeWidth)),
		convert.ConvertForce(*force),
		convert.ConvertConverters(slice.Strings(*converters, ",")),
		convert.ConvertColorHist(*colorHist),
	); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	check.Err(realMain())
}
