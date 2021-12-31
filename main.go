package main

import (
	"flag"

	"github.com/pkg/errors"
	"github.com/spudtrooper/eightbit/eightbit"
	"github.com/spudtrooper/goutil/check"
)

var (
	input     = flag.String("input", "", "input image")
	output    = flag.String("output", "", "output image")
	blockSize = flag.Int("block_size", 16, "blocksize for downsampling")
)

func realMain() error {
	if *input == "" {
		return errors.Errorf("--input required")
	}
	if *output == "" {
		return errors.Errorf("--output required")
	}
	if err := eightbit.Convert(*input, *output, eightbit.ConvertBlockSize(*blockSize)); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	check.Err(realMain())
}
