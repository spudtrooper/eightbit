package convert

//go:generate genopts --prefix=Convert --outfile=convert/options.go "blockSize:int" "resizeWidth:uint" "resizeHeight:uint" "converter:Converter" "force:bool"

type ConvertOption func(*convertOptionImpl)

type ConvertOptions interface {
	BlockSize() int
	ResizeWidth() uint
	ResizeHeight() uint
	Converter() Converter
	Force() bool
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

func ConvertConverter(converter Converter) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.converter = converter
	}
}
func ConvertConverterFlag(converter *Converter) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.converter = *converter
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

type convertOptionImpl struct {
	blockSize    int
	resizeWidth  uint
	resizeHeight uint
	converter    Converter
	force        bool
}

func (c *convertOptionImpl) BlockSize() int       { return c.blockSize }
func (c *convertOptionImpl) ResizeWidth() uint    { return c.resizeWidth }
func (c *convertOptionImpl) ResizeHeight() uint   { return c.resizeHeight }
func (c *convertOptionImpl) Converter() Converter { return c.converter }
func (c *convertOptionImpl) Force() bool          { return c.force }

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
