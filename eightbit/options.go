package eightbit

// genopts --opt_type=ConvertOption --prefix=Convert --outfile=eightbit/options.go 'blockSize:int'

type ConvertOption func(*convertOptionImpl)

type ConvertOptions interface {
	BlockSize() int
}

func ConvertBlockSize(blockSize int) ConvertOption {
	return func(opts *convertOptionImpl) {
		opts.blockSize = blockSize
	}
}

type convertOptionImpl struct {
	blockSize int
}

func (c *convertOptionImpl) BlockSize() int { return c.blockSize }

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
