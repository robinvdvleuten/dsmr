package dsmr

import (
	"math/big"

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

// Struct for telegram data.
type Telegram struct {
	Pos lexer.Position `parser:""`

	Header *Header `parser:"@@"`
	Data   Data    `parser:"@@*"`
	Footer *Footer `parser:"@@"`
}

func (t *Telegram) Position() lexer.Position { return t.Pos }

func (t *Telegram) children() (children []Node) {
	children = append(children, t.Header, t.Footer)

	for _, obj := range t.Data {
		children = append(children, obj)
	}

	return
}

type Header struct {
	Pos lexer.Position `parser:""`

	Value string `parser:"'/' @~EOL+ (?=EOL)"`
}

var _ Entry = &Header{}

func (h *Header) Key() string              { return "header" }
func (h *Header) Position() lexer.Position { return h.Pos }
func (h *Header) children() []Node         { return nil }

type Footer struct {
	Pos lexer.Position `parser:""`

	Value string `parser:"'!' @~EOL? (?=EOL)"`
}

var _ Entry = &Footer{}

func (f *Footer) Key() string              { return "footer" }
func (f *Footer) Position() lexer.Position { return f.Pos }
func (f *Footer) children() []Node         { return nil }

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

// LastCapture represents the last 5-minute capture of a MBus device.
type LastCapture struct {
	Pos lexer.Position `parser:""`

	Timestamp *Timestamp   `parser:"@@ ')'"`
	Value     *Measurement `parser:"'(' @@ ( ?!')' '(' )"`
}

var _ Value = &LastCapture{}

func (l *LastCapture) value()                      {}
func (l *LastCapture) Position() lexer.Position    { return l.Pos }
func (l *LastCapture) children() (children []Node) { return []Node{l.Timestamp, l.Value} }

// LegacyLastCapture represents the last 5-minute capture of an older MBus device (DSMR v2.2 or v3.0).
type LegacyLastCapture struct {
	Pos lexer.Position `parser:""`

	// We ignore any extraneous values between timestamp and OBIS as specs are unclear about their purpose.
	Timestamp *String            `parser:"@@ ')' ( '(' ~(')' | OBIS) ')' (?='(') )+ '('"`
	OBIS      *OBIS              `parser:"@@ ')' '('"`
	Value     *LegacyMeasurement `parser:"@@"`
}

var _ Value = &LegacyLastCapture{}

func (l *LegacyLastCapture) value()                      {}
func (l *LegacyLastCapture) Position() lexer.Position    { return l.Pos }
func (l *LegacyLastCapture) children() (children []Node) { return []Node{l.Timestamp, l.OBIS, l.Value} }

// ...
type OBIS struct {
	Pos lexer.Position `parser:""`

	Value string `parser:"@OBIS"`
}

var _ Value = &OBIS{}

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

func (m *Measurement) value()                   {}
func (m *Measurement) Position() lexer.Position { return m.Pos }

func (m *Measurement) children() (children []Node) {
	children = append(children, m.Value, m.Unit)
	return
}

// LegacyMeasurement represents a number+unit of a [LegacyLastCapture].
type LegacyMeasurement struct {
	Pos lexer.Position `parser:""`

	Unit  *String `parser:"@@ ')' '('"`
	Value *Number `parser:"@@"`
}

var _ Value = &LegacyMeasurement{}

func (m *LegacyMeasurement) value()                   {}
func (m *LegacyMeasurement) Position() lexer.Position { return m.Pos }

func (m *LegacyMeasurement) children() (children []Node) {
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

func (t *Timestamp) value()                   {}
func (t *Timestamp) Position() lexer.Position { return t.Pos }
func (t *Timestamp) children() []Node         { return nil }

// ...
type Number struct {
	Pos lexer.Position `parser:""`

	Value *big.Float `parser:"@Number"`
}

var _ Value = &Number{}

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

func (s *String) value()                   {}
func (s *String) Position() lexer.Position { return s.Pos }
func (s *String) children() []Node         { return nil }

var (
	lex = lexer.MustSimple([]lexer.SimpleRule{
		{"OBIS", `\d-\d:\d{1,2}\.\d{1,2}\.\d{1,2}`},
		{"Timestamp", `\d{12}`},
		{"Number", `\d*\.?\d+`},
		{"Chars", `[[:alnum:]]+`},
		{"Punct", `[-_!*.\\/()]`},
		{"EOL", `\r\n`},
	})

	parser = participle.MustBuild[Telegram](
		participle.Lexer(lex),
		participle.Elide("EOL"),
		participle.Union[Value](&EventLog{}, &LastCapture{}, &LegacyLastCapture{}, &Measurement{}, &Timestamp{}, &String{}),
		// We need lookahead to handle legacy last captures correctly.
		participle.UseLookahead(4),
	)
)

// Parse parses telegram from a string.
func Parse(str string, options ...Option) (*Telegram, error) {
	opts := parseOptions{
		verifyChecksum: true,
	}

	for _, option := range options {
		if err := option(&opts); err != nil {
			return nil, err
		}
	}

	t, err := parser.ParseString("", str)
	if err != nil {
		return nil, err
	}

	return t, verifyChecksum(t, str, &opts)
}
