package dsmr

import (
	"bytes"
	"io"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	dsmrLexer = lexer.MustSimple([]lexer.SimpleRule{
		{`OBIS`, `\d+-\d+:\d+\.\d+\.\d{1,2}`},
		{`Timestamp`, `\d{12}[SW]`},
		{`Number`, `\d+\.?\d*`},
		{`Punct`, `[!\\().]`},
		{`Char`, `[^\r]`},
		{"Separator", `\r\n`},
	})

	dsmrParser = participle.MustBuild[Telegram](
		participle.Lexer(dsmrLexer),
		participle.Elide("Separator"),
		participle.Union[Attribute](
			OBIS{},
			Measurement{},
			Timestamp{},
			Text{},
		),
	)
)

// Parse a DSRM telegram.
func Parse(r io.Reader) (*Telegram, error) {
	buf := &bytes.Buffer{}
	raw := io.TeeReader(r, buf)

	telegram, err := dsmrParser.Parse("", raw)
	if err != nil {
		return nil, err
	}

	err = telegram.check(buf)
	if err != nil {
		return nil, err
	}

	return telegram, nil
}

func ParseString(s string) (*Telegram, error) {
	return Parse(strings.NewReader(s))
}
