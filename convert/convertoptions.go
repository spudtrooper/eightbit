package convert

//go:generate genopts --prefix=Convert --outfile=convert/convertoptions.go "blockSize:int" "animateBlockSizeRange:blockSizeRange" "pixelateBlockSize:int" "resizeWidth:uint" "resizeHeight:uint" "force:bool" "converters:[]string" "outputDir:string" "outputFile:string" "colorHist:bool" "animateThreads:int" "animateReverse"

type ConvertOption func(*convertOptionImpl)

type ConvertOptions interface {
	BlockSize() int
	AnimateBlockSizeRange() blockSizeRange
	PixelateBlockSize() int
	ResizeWidth() uint
	ResizeHeight() uint
	Force() bool
	Converters() []string
	OutputDir() string
	OutputFile() string
	ColorHist() bool
	AnimateThreads() int
	AnimateReverse() bool
}

func ConvertBlockSize(blockSize int) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.blockSize = blockSize
	}
}
func ConvertBlockSizeFlag(blockSize *int) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.blockSize = *blockSize
	}
}

func ConvertAnimateBlockSizeRange(animateBlockSizeRange blockSizeRange) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.animateBlockSizeRange = animateBlockSizeRange
	}
}
func ConvertAnimateBlockSizeRangeFlag(animateBlockSizeRange *blockSizeRange) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.animateBlockSizeRange = *animateBlockSizeRange
	}
}

func ConvertPixelateBlockSize(pixelateBlockSize int) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.pixelateBlockSize = pixelateBlockSize
	}
}
func ConvertPixelateBlockSizeFlag(pixelateBlockSize *int) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.pixelateBlockSize = *pixelateBlockSize
	}
}

func ConvertResizeWidth(resizeWidth uint) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.resizeWidth = resizeWidth
	}
}
func ConvertResizeWidthFlag(resizeWidth *uint) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.resizeWidth = *resizeWidth
	}
}

func ConvertResizeHeight(resizeHeight uint) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.resizeHeight = resizeHeight
	}
}
func ConvertResizeHeightFlag(resizeHeight *uint) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.resizeHeight = *resizeHeight
	}
}

func ConvertForce(force bool) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.force = force
	}
}
func ConvertForceFlag(force *bool) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.force = *force
	}
}

func ConvertConverters(converters []string) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.converters = converters
	}
}
func ConvertConvertersFlag(converters *[]string) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.converters = *converters
	}
}

func ConvertOutputDir(outputDir string) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.outputDir = outputDir
	}
}
func ConvertOutputDirFlag(outputDir *string) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.outputDir = *outputDir
	}
}

func ConvertOutputFile(outputFile string) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.outputFile = outputFile
	}
}
func ConvertOutputFileFlag(outputFile *string) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.outputFile = *outputFile
	}
}

func ConvertColorHist(colorHist bool) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.colorHist = colorHist
	}
}
func ConvertColorHistFlag(colorHist *bool) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.colorHist = *colorHist
	}
}

func ConvertAnimateThreads(animateThreads int) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.animateThreads = animateThreads
	}
}
func ConvertAnimateThreadsFlag(animateThreads *int) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.animateThreads = *animateThreads
	}
}

func ConvertAnimateReverse(animateReverse bool) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.animateReverse = animateReverse
	}
}
func ConvertAnimateReverseFlag(animateReverse *bool) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.animateReverse = *animateReverse
	}
}

type convertOptionImpl struct {
	blockSize             int
	animateBlockSizeRange blockSizeRange
	pixelateBlockSize     int
	resizeWidth           uint
	resizeHeight          uint
	force                 bool
	converters            []string
	outputDir             string
	outputFile            string
	colorHist             bool
	animateThreads        int
	animateReverse        bool
}

func (c *convertOptionImpl) BlockSize() int                        { return c.blockSize }
func (c *convertOptionImpl) AnimateBlockSizeRange() blockSizeRange { return c.animateBlockSizeRange }
func (c *convertOptionImpl) PixelateBlockSize() int                { return c.pixelateBlockSize }
func (c *convertOptionImpl) ResizeWidth() uint                     { return c.resizeWidth }
func (c *convertOptionImpl) ResizeHeight() uint                    { return c.resizeHeight }
func (c *convertOptionImpl) Force() bool                           { return c.force }
func (c *convertOptionImpl) Converters() []string                  { return c.converters }
func (c *convertOptionImpl) OutputDir() string                     { return c.outputDir }
func (c *convertOptionImpl) OutputFile() string                    { return c.outputFile }
func (c *convertOptionImpl) ColorHist() bool                       { return c.colorHist }
func (c *convertOptionImpl) AnimateThreads() int                   { return c.animateThreads }
func (c *convertOptionImpl) AnimateReverse() bool                  { return c.animateReverse }

func makeConvertOptionImpl(opts ...ConvertOption) *convertOptionImpl {
	res := &convertOptionImpl{}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func MakeConvertOptions(opts ...ConvertOption) ConvertOptions {
	return makeConvertOptionImpl(opts...)
}
