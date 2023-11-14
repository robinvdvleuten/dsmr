package dsmr

type parseOptions struct {
	verifyChecksum bool
}

type Option func(opts *parseOptions) error

func VerifyChecksum(v bool) Option {
	return func(opts *parseOptions) error {
		opts.verifyChecksum = v
		return nil
	}
}
