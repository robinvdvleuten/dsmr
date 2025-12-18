package dsmr

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Grammar overview (simplified EBNF):
//
//	Telegram   = Header Entry* Footer ;
//	Entry      = Header | Footer | Object | MBusDevice ;
//	Object     = OBIS "(" Value* ")" ;
//	MBusDevice = Object { Object }  // same channel, leading OBIS 0-<channel>:*
//	Value      = EventLog | LastCapture | LegacyLastCapture | Measurement | Timestamp | String ;
//	EventLog   = Number ")" "(" OBIS ")" Event* ;
//	Event      = "(" Timestamp ")" "(" Measurement [")"] ;
//	LastCapture = Timestamp ")" "(" Measurement ;
//	Measurement = Number "*" String ;
//	Timestamp  = <12 digit datetime> ["S" | "W"] ;
//	String     = <any chars except ')' or EOL> ;
//	Number     = ["+" | "-"] digit+ ["." digit+] ;
//	Header/Footer tokens reuse String rules for their payloads.
//
// Node expresses the common behaviour for every AST node in a telegram.
type Node interface {
	Position() lexer.Position
	children() (children []Node)
}

// Entry represents the top-level items that can appear in telegram data.
type Entry interface {
	Key() string
	Node
}

// Value represents the possible values attached to an OBIS object.
type Value interface {
	value()
	Node
}

// Data collects the entries in the body of a telegram.
type Data []Entry

// Header is the opening line of a telegram.
type Header struct {
	Pos lexer.Position `parser:""`

	Value string `parser:"'/' @~EOL+ (?=EOL)"`
}

var _ Entry = &Header{}

func (h *Header) Key() string              { return "header" }
func (h *Header) Position() lexer.Position { return h.Pos }
func (h *Header) children() []Node         { return nil }

// Footer is the closing line of a telegram.
type Footer struct {
	Pos lexer.Position `parser:""`

	Value string `parser:"'!' @~EOL? (?=EOL)"`
}

var _ Entry = &Footer{}

func (f *Footer) Key() string              { return "footer" }
func (f *Footer) Position() lexer.Position { return f.Pos }
func (f *Footer) children() []Node         { return nil }

// OBIS holds the identifier of a COSEM object.
type OBIS struct {
	Pos lexer.Position `parser:""`

	Value string `parser:"@OBIS"`
}

var _ Value = &OBIS{}

func (o *OBIS) value()                   {}
func (o *OBIS) Position() lexer.Position { return o.Pos }
func (o *OBIS) children() []Node         { return nil }

// Object represents a COSEM object with one or more values.
type Object struct {
	Pos lexer.Position `parser:""`

	OBIS  *OBIS `parser:"@@"`
	Value Value `parser:"'(' @@* ')' (?=EOL)"`

	mbusChannel int
	hasMBus     bool
}

var _ Entry = &Object{}

func (o *Object) Key() string              { return o.OBIS.Value }
func (o *Object) Position() lexer.Position { return o.Pos }
func (o *Object) children() []Node         { return []Node{o.OBIS, o.Value} }

func (o *Object) mbusChannelValue() (int, bool) {
	if o == nil || !o.hasMBus {
		return 0, false
	}

	return o.mbusChannel, true
}

// MBusDevice groups OBIS objects that belong to the same M-Bus channel.
type MBusDevice struct {
	Pos     lexer.Position
	Channel int
	Data    []*Object
}

var _ Entry = &MBusDevice{}

func (m *MBusDevice) Key() string              { return fmt.Sprintf("mbus.%d", m.Channel) }
func (m *MBusDevice) Position() lexer.Position { return m.Pos }

func (m *MBusDevice) children() (children []Node) {
	for _, obj := range m.Data {
		if obj == nil {
			continue
		}
		children = append(children, obj)
	}

	return
}

// Telegram is the root node produced for every parsed message.
type Telegram struct {
	Pos lexer.Position `parser:""`

	Header *Header `parser:"@@"`
	Data   Data    `parser:"@@*"`
	Footer *Footer `parser:"@@"`
}

func (t *Telegram) Position() lexer.Position { return t.Pos }

func (t *Telegram) children() (children []Node) {
	children = append(children, t.Header, t.Footer)

	for _, entry := range t.Data {
		if entry == nil {
			continue
		}
		children = append(children, entry)
	}

	return
}

// String represents a literal string segment.
type String struct {
	Pos lexer.Position `parser:""`

	// Also check for `EOL` token so both Header and Footer can use this Value struct as well.
	Value string `parser:"@(~(')' | EOL)+)"`
}

var _ Value = &String{}

func (s *String) value()                   {}
func (s *String) Position() lexer.Position { return s.Pos }
func (s *String) children() []Node         { return nil }

// Number wraps a numeric literal optionally containing a decimal point.
type Number struct {
	Pos lexer.Position `parser:""`

	Value *big.Float `parser:"@Number"`
}

var _ Value = &Number{}

func (n *Number) value()                   {}
func (n *Number) Position() lexer.Position { return n.Pos }
func (n *Number) children() []Node         { return nil }

// Timestamp represents a timestamp of a date.
type Timestamp struct {
	Pos lexer.Position `parser:""`

	Value string       `parser:"@Timestamp"`
	DST   DSTIndicator `parser:"@('S' | 'W')"`
}

var _ Value = &Timestamp{}

func (t *Timestamp) value()                   {}
func (t *Timestamp) Position() lexer.Position { return t.Pos }
func (t *Timestamp) children() []Node         { return nil }

// IsDST reports whether the timestamp is in daylight saving time (summer time).
func (t *Timestamp) IsDST() bool {
	if t == nil {
		return false
	}
	return t.DST.Bool()
}

// DSTIndicator represents daylight saving time indicator.
type DSTIndicator int

const (
	DSTUnknown DSTIndicator = iota
	DSTWinter
	DSTSummer
)

// Capture implements participle capture for DSTIndicator.
func (d *DSTIndicator) Capture(values []string) error {
	if len(values) == 0 {
		*d = DSTUnknown
		return nil
	}

	switch values[0] {
	case "W":
		*d = DSTWinter
	case "S":
		*d = DSTSummer
	default:
		return fmt.Errorf("unexpected DST indicator %q", values[0])
	}

	return nil
}

// Bool reports whether daylight saving time is active (summer time).
func (d DSTIndicator) Bool() bool { return d == DSTSummer }

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

// Event represents a timestamp+duration pair.
type Event struct {
	Pos lexer.Position `parser:""`

	Timestamp *Timestamp   `parser:"'(' @@ ')'"`
	Value     *Measurement `parser:"'(' @@ ( ')' (?='(') )?"`
}

var _ Value = &Event{}

func (e *Event) value()                      {}
func (e *Event) Position() lexer.Position    { return e.Pos }
func (e *Event) children() (children []Node) { return []Node{e.Timestamp, e.Value} }

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

var (
	lex = lexer.MustSimple([]lexer.SimpleRule{
		{"OBIS", `\d-\d:\d{1,2}\.\d{1,2}\.\d{1,2}`},
		{"Timestamp", `\d{12}`},
		{"Number", `[+-]?(?:\d+(?:\.\d+)?|\.\d+)`},
		{"Chars", `[[:alnum:]]+`},
		{"Punct", `[-_!*.\\/()]`},
		{"EOL", `\r\n`},
	})

	obisToken = lex.Symbols()["OBIS"]

	valueUnion = participle.Union[Value](&EventLog{}, &LastCapture{}, &LegacyLastCapture{}, &Measurement{}, &Timestamp{}, &String{})

	parserOnce       sync.Once
	parser           *participle.Parser[Telegram]
	objectParserOnce sync.Once
	objectParser     *participle.Parser[object]
)

// initParser initializes the telegram parser with sync.Once to ensure it's built only once.
func initParser() *participle.Parser[Telegram] {
	parserOnce.Do(func() {
		parser = participle.MustBuild[Telegram](
			participle.Lexer(lex),
			participle.Elide("EOL"),
			participle.ParseTypeWith[Entry](parseEntry),
			valueUnion,
			// We need lookahead to handle legacy last captures correctly.
			participle.UseLookahead(4),
		)
	})
	return parser
}

// initObjectParser initializes the object parser with sync.Once to ensure it's built only once.
func initObjectParser() *participle.Parser[object] {
	objectParserOnce.Do(func() {
		objectParser = participle.MustBuild[object](
			participle.Lexer(lex),
			participle.Elide("EOL"),
			valueUnion,
			participle.UseLookahead(4),
		)
	})
	return objectParser
}

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

	t, err := initParser().ParseString("", str)
	if err != nil {
		return nil, wrapParseError(err)
	}

	return t, verifyChecksum(t, str, &opts)
}

func wrapParseError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(*ParseError); ok {
		return err
	}

	if perr, ok := err.(participle.Error); ok {
		wrapped := &ParseError{
			Pos:     perr.Position(),
			Message: perr.Message(),
			Err:     err,
		}
		if unexpectedProvider, ok := perr.(interface{ Unexpected() string }); ok {
			wrapped.Unexpected = unexpectedProvider.Unexpected()
		}
		return wrapped
	}

	if lerr, ok := err.(*lexer.Error); ok {
		return &ParseError{
			Pos:     lerr.Pos,
			Message: lerr.Message(),
			Err:     err,
		}
	}

	return err
}

func parseEntry(lex *lexer.PeekingLexer) (Entry, error) {
	tok := lex.Peek()
	if tok.Type != obisToken {
		return nil, participle.NextMatch
	}

	first, err := parseObject(lex)
	if err != nil {
		return nil, err
	}

	channel, ok := first.mbusChannelValue()
	if !ok {
		return first, nil
	}

	device := &MBusDevice{
		Pos:     first.Position(),
		Channel: channel,
		Data:    []*Object{first},
	}

	for {
		nextTok := lex.Peek()
		if nextTok.Type != obisToken {
			break
		}

		checkpoint := lex.MakeCheckpoint()
		nextObj, err := parseObject(lex)
		if err != nil {
			return nil, err
		}

		nextChannel, ok := nextObj.mbusChannelValue()
		if !ok || nextChannel != channel {
			lex.LoadCheckpoint(checkpoint)
			break
		}

		device.Data = append(device.Data, nextObj)
	}

	return device, nil
}

func parseObject(lex *lexer.PeekingLexer) (*Object, error) {
	start := lex.Peek()

	parsed, err := initObjectParser().ParseFromLexer(lex, participle.AllowTrailing(true))
	if err != nil {
		return nil, participle.Wrapf(start.Pos, err, "unable to parse OBIS object")
	}

	obj := &Object{
		Pos:   parsed.Pos,
		OBIS:  parsed.OBIS,
		Value: parsed.Value,
	}

	if channel, ok := mbusChannel(obj.OBIS.Value); ok {
		obj.mbusChannel = channel
		obj.hasMBus = true
	}

	return obj, nil
}

type object struct {
	Pos   lexer.Position `parser:""`
	OBIS  *OBIS          `parser:"@@"`
	Value Value          `parser:"'(' @@* ')'"`
}

type obisCapture struct {
	Value     string
	Device    int
	hasDevice bool
	Channel   int
}

func (o *obisCapture) Capture(values []string) error {
	if len(values) == 0 {
		return fmt.Errorf("missing OBIS value")
	}

	o.Value = values[0]
	o.Device = 0
	o.Channel = 0
	o.hasDevice = false

	parts := strings.SplitN(o.Value, ":", 2)
	if len(parts) != 2 {
		return nil
	}

	segments := strings.SplitN(parts[0], "-", 2)
	if len(segments) != 2 {
		return nil
	}

	if device, err := strconv.Atoi(segments[0]); err == nil {
		o.Device = device
		o.hasDevice = true
	}

	if channel, err := strconv.Atoi(segments[1]); err == nil {
		o.Channel = channel
	}

	return nil
}

func mbusChannel(obis string) (int, bool) {
	var captured obisCapture
	if err := captured.Capture([]string{obis}); err != nil {
		return 0, false
	}

	if !captured.hasDevice || captured.Device != 0 {
		return 0, false
	}

	if captured.Channel <= 0 {
		return 0, false
	}

	return captured.Channel, true
}

// ParseError wraps errors returned by participle with position information and context.
type ParseError struct {
	Pos        lexer.Position
	Message    string
	Unexpected string
	Err        error
}

func (e *ParseError) Error() string {
	if e == nil {
		return "<nil>"
	}

	pos := e.Pos
	msg := "parse error"
	if e.Message != "" {
		msg = fmt.Sprintf("%s: %s", msg, e.Message)
	}

	if e.Message == "" && e.Err != nil {
		msg = fmt.Sprintf("%s: %s", msg, e.Err.Error())
	}

	if e.Unexpected != "" && !strings.Contains(msg, "unexpected") {
		msg = fmt.Sprintf("%s: unexpected %s", msg, e.Unexpected)
	}

	if pos.Line > 0 && pos.Column > 0 {
		return fmt.Sprintf("%d:%d: %s", pos.Line, pos.Column, msg)
	}

	return msg
}

func (e *ParseError) Unwrap() error { return e.Err }
