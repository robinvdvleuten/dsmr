package dsmr

type parseOptions struct {
	verifyChecksum bool
}

type Option func(opts *parseOptions)

func VerifyChecksum(v bool) Option {
	return func(opts *parseOptions) {
		opts.verifyChecksum = v
	}
}
