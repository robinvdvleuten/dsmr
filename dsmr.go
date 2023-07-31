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
	t := &Telegram{
		COSEM: map[string]*COSEM{},
	}

	p := &parser{
		Buffer: s,
		t:      t,
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
