package buffer

type Option func(options *Options)
type Options struct {
	blockSize int
}

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}
func WithBlockSize(blockSize int) Option {
	if blockSize <= 0 {
		blockSize = 4096
	}
	return func(options *Options) {
		options.blockSize = blockSize
	}
}
