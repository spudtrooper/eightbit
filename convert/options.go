package convert

//go:generate genopts --opt_type=ConvertOption --prefix=Convert --outfile=convert/options.go "blockSize:int" "resizeWidth:uint" "resizeHeight:uint"

type ConvertOption func(*convertOptionImpl)

type ConvertOptions interface {
	BlockSize() int
	ResizeWidth() uint
	ResizeHeight() uint
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

type convertOptionImpl struct {
	blockSize    int
	resizeWidth  uint
	resizeHeight uint
}

func (c *convertOptionImpl) BlockSize() int     { return c.blockSize }
func (c *convertOptionImpl) ResizeWidth() uint  { return c.resizeWidth }
func (c *convertOptionImpl) ResizeHeight() uint { return c.resizeHeight }

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
