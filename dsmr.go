package dsmr

import (
	"io"
	"strings"
	"time"
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

// WithLocation sets the location used for parsing timestamps.
// By default the "Europe/Amsterdam" location is used, but can
// be overriden to handle meters in different timezones.
func WithLocation(name string) Option {
	return func(p *parser) {
		tz, _ := time.LoadLocation(name)
		p.tz = tz
	}
}

// Parse a DSRM telegram.
func ParseString(s string, options ...Option) (*Telegram, error) {
	t := &Telegram{
		cosem: map[string]*COSEM{},
	}

	tz, _ := time.LoadLocation("Europe/Amsterdam")

	p := &parser{
		Buffer: s,
		t:      t,
		tz:     tz,
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
