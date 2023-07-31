package dsmr

import (
	"io"
	"strings"
)

// Parse a DSRM telegram.
func Parse(r io.Reader) (*Telegram, error) {
	buf := &strings.Builder{}
	if _, err := io.Copy(buf, r); err != nil {
		return nil, err
	}

	return ParseString(buf.String())
}

func ParseString(s string) (*Telegram, error) {
	p := &parser{
		Buffer: s,
		t:      &Telegram{},
	}

	if err := p.Init(); err != nil {
		return nil, err
	}

	if err := p.Parse(); err != nil {
		return nil, err
	}

	p.Execute()

	return p.t, nil
}
