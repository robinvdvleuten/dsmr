package dsmr

import (
	"io"
	"strings"
)

// Option function.
type Option func(*parser)

func Parse(r io.Reader, options ...Option) (*Telegram, error) {
	buf := &strings.Builder{}
	if _, err := io.Copy(buf, r); err != nil {
		return nil, err
	}

	return ParseString(buf.String(), options...)
}

// Parse a DSRM telegram.
func ParseString(s string, options ...Option) (*Telegram, error) {
	t := &Telegram{
		COSEM: map[string]*COSEM{},
	}

	p := &parser{
		Buffer: s,
		t:      t,
	}

	for _, o := range options {
		o(p)
	}

	if err := p.Init(); err != nil {
		return nil, err
	}

	if err := p.Parse(); err != nil {
		return nil, err
	}

	p.Execute()

	return t, nil
}
