package main

import (
	"flag"

	"github.com/pkg/errors"
	"github.com/spudtrooper/eightbit/convert"
	"github.com/spudtrooper/eightbit/gitversion"
	"github.com/spudtrooper/goutil/check"
)

var (
	input        = flag.String("input", "", "input image")
	output       = flag.String("output", "", "output image")
	blockSize    = flag.Int("block_size", 16, "blocksize for downsampling")
	resizeHeight = flag.Int("resize_height", 0, "height in pixels of the final image; must be used with --resize_width")
	resizeWidth  = flag.Int("resize_width", 0, "width in pixels of the final image; must be used with --resize_height")
	force        = flag.Bool("force", false, "overwrite existing files")
)

func realMain() error {
	if gitversion.CheckVersionFlag() {
		return nil
	}
	if *input == "" {
		return errors.Errorf("--input required")
	}
	if *output == "" {
		return errors.Errorf("--output required")
	}
	if err := convert.Convert(*input, *output,
		convert.ConvertBlockSize(*blockSize),
		convert.ConvertResizeWidth(uint(*resizeHeight)),
		convert.ConvertResizeHeight(uint(*resizeWidth)),
		convert.ConvertForce(*force),
	); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	check.Err(realMain())
}
