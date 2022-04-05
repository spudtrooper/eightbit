package main

import (
	"flag"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spudtrooper/eightbit/convert"
	"github.com/spudtrooper/eightbit/gitversion"
	"github.com/spudtrooper/goutil/check"
	"github.com/spudtrooper/goutil/slice"
)

var (
	input                 = flag.String("input", "", "input image")
	output                = flag.String("output", "", "output image")
	outputDir             = flag.String("output_dir", "", "output dir")
	pixelateBlockSize     = flag.Int("pixelate_block_size", 16, "blocksize for downsampling")
	blockSize             = flag.Int("block_size", 10, "blocksize overlap and block converters")
	resizeHeight          = flag.Int("resize_height", 0, "height in pixels of the final image; must be used with --resize_width")
	resizeWidth           = flag.Int("resize_width", 0, "width in pixels of the final image; must be used with --resize_height")
	force                 = flag.Bool("force", false, "overwrite existing files")
	converters            = flag.String("converters", "pixelated", "the kinds of converter to use or 'all' for all of them. If you don't specify an output file, the output file will be next to the source file with this tag at the end of the base name.")
	printConverters       = flag.Bool("print_converters", false, "print the names of all the converters and exit")
	colorHist             = flag.Bool("color_hist", false, "print a histogram of web colors from the input image")
	openAll               = flag.Bool("open_all", false, "try to open the output files at the end")
	animateThreads        = flag.Int("animate_threads", 0, "number of threads for producing animations")
	animateBlockSizeStart = flag.Int("animate_block_size_start", 1, "start block size for animations")
	animateBlockSizeEnd   = flag.Int("animate_block_size_end", 150, "end block size for animations")
	animateBlockSizeStep  = flag.Int("animate_block_size_step", 1, "block size step for animations")
	animateReverse        = flag.Bool("animate_reverse", false, "sort the images from higher block size to lower (i.e. reversed)")
	except                = flag.String("except", "", "comma-delimited list of converters to skip; to be used with --converters all --except <foo>")
)

func realMain() error {
	if gitversion.CheckVersionFlag() {
		return nil
	}

	if *printConverters {
		fmt.Println("Printing the names of all the converters...")
		for i, c := range convert.AllConverterNames() {
			fmt.Printf("  [%d] %s\n", i+1, c)
		}
		return nil
	}

	if *input == "" {
		return errors.Errorf("--input required")
	}

	outputs, err := convert.Convert(*input,
		convert.ConvertOutputFile(*output),
		convert.ConvertOutputDir(*outputDir),
		convert.ConvertBlockSize(*blockSize),
		convert.ConvertPixelateBlockSize(*pixelateBlockSize),
		convert.ConvertResizeWidth(uint(*resizeHeight)),
		convert.ConvertResizeHeight(uint(*resizeWidth)),
		convert.ConvertForce(*force),
		convert.ConvertConverters(slice.Strings(*converters, ",")),
		convert.ConvertExcept(slice.Strings(*except, ",")),
		convert.ConvertColorHist(*colorHist),
		convert.ConvertAnimateThreads(*animateThreads),
		convert.ConvertAnimateBlockSizeRange(convert.MakeBlockSizeRange(*animateBlockSizeStart, *animateBlockSizeEnd, *animateBlockSizeStep)),
		convert.ConvertAnimateReverse(*animateReverse),
	)
	if err != nil {
		return err
	}

	if *openAll {
		if err := exec.Command("open", outputs...).Run(); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	flag.Parse()
	check.Err(realMain())
}
