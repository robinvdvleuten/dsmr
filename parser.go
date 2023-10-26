package dsmr

import (
	"bytes"
	"io"
	"math/big"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Node interface {
	Position() lexer.Position
	children() (children []Node)
}

type Entry interface {
	Key() string
	Node
}

// Data Objects in the root of the AST.
type Data []*Object

// AST for telegram data.
type AST struct {
	Pos lexer.Position `parser:""`

	Header *Header `parser:"@@"`
	Data   Data    `parser:"@@*"`
	Footer *Footer `parser:"@@"`
}

func (a *AST) Position() lexer.Position { return a.Pos }

func (a *AST) entries() (entries []Entry) {
	entries = append(entries, a.Header, a.Footer)
	for _, obj := range a.Data {
		entries = append(entries, obj)
	}

	return
}

func (a *AST) children() (children []Node) {
	for _, entry := range a.entries() {
		children = append(children, entry)
	}

	return
}

type Header struct {
	Pos lexer.Position `parser:""`

	Value *String `parser:"'/' @@ (?=EOL)"`
}

var _ Entry = &Header{}

func (h *Header) Key() string              { return "header" }
func (h *Header) Position() lexer.Position { return h.Pos }
func (h *Header) children() []Node         { return []Node{h.Value} }

type Footer struct {
	Pos lexer.Position `parser:""`

	Value *String `parser:"'!' @@? (?=EOL)"`
}

var _ Entry = &Footer{}

func (f *Footer) Key() string              { return "footer" }
func (f *Footer) Position() lexer.Position { return f.Pos }
func (f *Footer) children() []Node         { return []Node{f.Value} }

// Object is a COSEM object in the Telegram represented by the
// OBIS (Object Identification System) and one or more attributes.
type Object struct {
	Pos lexer.Position `parser:""`

	OBIS  *OBIS `parser:"@@"`
	Value Value `parser:"'(' @@* ')' (?=EOL)"`
}

var _ Entry = &Object{}

func (o *Object) Key() string              { return o.OBIS.Value }
func (o *Object) Position() lexer.Position { return o.Pos }
func (o *Object) children() []Node         { return []Node{o.OBIS, o.Value} }

// Value represents an object value.
type Value interface {
	value()

	Node
}

// EventLog represents a log of events.
type EventLog struct {
	Pos lexer.Position `parser:""`

	Count *Number  `parser:"@@ ')'"`
	OBIS  *OBIS    `parser:"'(' @@ ( ')' (?='(') )?"`
	Value []*Event `parser:"@@*"`
}

var _ Value = &EventLog{}

func (e *EventLog) value()                   {}
func (e *EventLog) Position() lexer.Position { return e.Pos }

func (e *EventLog) children() (children []Node) {
	children = append(children, e.Count, e.OBIS)

	for _, val := range e.Value {
		children = append(children, val)
	}

	return
}

// Event represents a timestamp+duration.
type Event struct {
	Pos lexer.Position `parser:""`

	Timestamp *Timestamp   `parser:"'(' @@ ')'"`
	Value     *Measurement `parser:"'(' @@ ( ')' (?='(') )?"`
}

var _ Value = &Event{}

func (e *Event) value()                      {}
func (e *Event) Position() lexer.Position    { return e.Pos }
func (e *Event) children() (children []Node) { return []Node{e.Timestamp, e.Value} }

// List represents a list of values.
type List struct {
	Pos lexer.Position `parser:""`

	Value []ListValue `parser:"( @@ ')' '(' )+ @@"`
}

var _ Value = &List{}

func (l *List) value()                   {}
func (l *List) Position() lexer.Position { return l.Pos }

func (l *List) children() (children []Node) {
	for _, val := range l.Value {
		children = append(children, val)
	}

	return
}

// ...
type ListValue interface {
	Node
}

// ...
type OBIS struct {
	Pos lexer.Position `parser:""`

	Value string `parser:"@OBIS"`
}

var _ Value = &OBIS{}
var _ ListValue = &OBIS{}

func (o *OBIS) value()                   {}
func (o *OBIS) Position() lexer.Position { return o.Pos }
func (o *OBIS) children() []Node         { return nil }

// Measurement represents a number+unit.
type Measurement struct {
	Pos lexer.Position `parser:""`

	Value *Number `parser:"@@"`
	Unit  *String `parser:"'*' @@"`
}

var _ Value = &Measurement{}
var _ ListValue = &Measurement{}

func (m *Measurement) value()                   {}
func (m *Measurement) Position() lexer.Position { return m.Pos }

func (m *Measurement) children() (children []Node) {
	children = append(children, m.Value, m.Unit)
	return
}

// Timestamp represents a timestamp of a date.
type Timestamp struct {
	Pos lexer.Position `parser:""`

	Value string `parser:"@Timestamp"`
	DST   bool   `parser:"(@'S' | 'W')"`
}

var _ Value = &Timestamp{}
var _ ListValue = &Timestamp{}

func (t *Timestamp) value()                   {}
func (t *Timestamp) Position() lexer.Position { return t.Pos }
func (t *Timestamp) children() []Node         { return nil }

// ...
type Number struct {
	Pos lexer.Position `parser:""`

	Value *big.Float `parser:"@Number"`
}

var _ Value = &Number{}
var _ ListValue = &Number{}

func (n *Number) value()                   {}
func (n *Number) Position() lexer.Position { return n.Pos }
func (n *Number) children() []Node         { return nil }

// String literal.
type String struct {
	Pos lexer.Position `parser:""`

	// Also check for `EOL` token so both Header and Footer
	// can use this Value struct as well.
	Value string `parser:"@(~(')' | EOL)+)"`
}

var _ Value = &String{}
var _ ListValue = &String{}

func (s *String) value()                   {}
func (s *String) Position() lexer.Position { return s.Pos }
func (s *String) children() []Node         { return nil }

var (
	lex = lexer.MustSimple([]lexer.SimpleRule{
		{"OBIS", `\d{1,2}-\d{1,2}:\d{1,2}\.\d{1,2}\.\d{1,2}`},
		{"Timestamp", `\d{12}`},
		{"Number", `\d*\.?\d+`},
		{"Chars", `[[:alnum:]]+`},
		{"Punct", `[-_!*.\\/()]`},
		{"EOL", `\r\n`},
	})

	parser = participle.MustBuild[AST](
		participle.Lexer(lex),
		participle.Elide("EOL"),
		participle.Union[Value](&EventLog{}, &List{}, &OBIS{}, &Measurement{}, &Timestamp{}, &String{}),
		participle.Union[ListValue](&OBIS{}, &Measurement{}, &Timestamp{}, &String{}),
		participle.UseLookahead(4),
	)
)

func Parse(r io.Reader) (*AST, error) {
	ast, err := parser.Parse("", r)
	if err != nil {
		return nil, err
	}

	return ast, ast.VerifyChecksum(r)
}

// ParseString parses telegram from a string.
func ParseString(str string) (*AST, error) {
	ast, err := parser.ParseString("", str)
	if err != nil {
		return nil, err
	}

	return ast, ast.VerifyChecksum(strings.NewReader(str))
}

// ParseBytes parses telegram from bytes.
func ParseBytes(data []byte) (*AST, error) {
	ast, err := parser.ParseBytes("", data)
	if err != nil {
		return nil, err
	}

	return ast, ast.VerifyChecksum(bytes.NewReader(data))
}
