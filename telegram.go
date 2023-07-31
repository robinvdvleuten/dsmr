//go:generate peg -inline -switch grammar.peg

package dsmr

import (
	"fmt"
	"time"
)

type Telegram struct {
	Header   string
	COSEM    map[string]*COSEM
	Checksum string
}

type COSEM struct {
	OBIS      *OBIS
	Attribute []Attribute
}

type Attribute interface {
	String() string
}

type Measurement struct {
	Value float64
	Unit  string
}

func (m *Measurement) String() string {
	return fmt.Sprintf("%f%s", m.Value, m.Unit)
}

type OBIS struct {
	Value string
}

func (o *OBIS) String() string {
	return o.Value
}

type Text struct {
	Value string
}

func (t *Text) String() string {
	return t.Value
}

type Timestamp struct {
	Value time.Time
}

func (t *Timestamp) String() string {
	return t.Value.Format(time.RFC3339)
}
