package eightbit

// genopts --opt_type=ConvertOption --prefix=Convert --outfile=eightbit/options.go 'blockSize:int' 'resizeWidth:uint' 'resizeHeight:uint' 'noCleanup'

type ConvertOption func(*convertOptionImpl)

type ConvertOptions interface {
	BlockSize() int
	ResizeWidth() uint
	ResizeHeight() uint
	NoCleanup() bool
}

func ConvertBlockSize(blockSize int) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.blockSize = blockSize
	}
}

func ConvertResizeWidth(resizeWidth uint) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.resizeWidth = resizeWidth
	}
}

func ConvertResizeHeight(resizeHeight uint) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.resizeHeight = resizeHeight
	}
}

func ConvertNoCleanup(noCleanup bool) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.noCleanup = noCleanup
	}
}

type convertOptionImpl struct {
	blockSize    int
	resizeWidth  uint
	resizeHeight uint
	noCleanup    bool
}

func (c *convertOptionImpl) BlockSize() int     { return c.blockSize }
func (c *convertOptionImpl) ResizeWidth() uint  { return c.resizeWidth }
func (c *convertOptionImpl) ResizeHeight() uint { return c.resizeHeight }
func (c *convertOptionImpl) NoCleanup() bool    { return c.noCleanup }

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
