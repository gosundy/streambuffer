package buffer

import "sync"

type Option func(options *Options)
type Options struct {
	blockSize  int
	emptyError error
	fullError  error
	pool       *sync.Pool
}

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}
func WithPool(pool *sync.Pool) Option {
	return func(options *Options) {
		options.pool = pool
	}
}
func WithBlockSize(blockSize int) Option {
	return func(options *Options) {
		options.blockSize = blockSize
	}
}
func WithEmptyError(err error) Option {
	return func(options *Options) {
		options.emptyError = err
	}
}
func WithFullError(err error) Option {
	return func(options *Options) {
		options.fullError = err
	}
}
