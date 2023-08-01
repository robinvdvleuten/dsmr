//go:generate peg -inline -switch grammar.peg

package dsmr

import (
	"time"
)

type Telegram struct {
	header   string
	cosem    map[string]*COSEM
	checksum string
}

func (t *Telegram) Header() string {
	return t.header
}

func (t *Telegram) Checksum() string {
	return t.checksum
}

func (t *Telegram) COSEM(obis string) []Attribute {
	return t.cosem[obis].attribute
}

type COSEM struct {
	obis      *OBIS
	attribute []Attribute
}

func (c *COSEM) OBIS() *OBIS {
	return c.obis
}

func (c *COSEM) Attribute() []Attribute {
	return c.attribute
}

type Attribute interface {
	attribute()
}

type Measurement struct {
	value float64
	unit  string
}

func (m *Measurement) attribute() {}

func (m *Measurement) Value() float64 {
	return m.value
}

func (m *Measurement) Unit() string {
	return m.unit
}

type OBIS struct {
	value string
}

func (o *OBIS) attribute() {}

func (o *OBIS) Value() string {
	return o.value
}

type Text struct {
	value string
}

func (t *Text) attribute() {}

func (t *Text) Value() string {
	return t.value
}

type Timestamp struct {
	value time.Time
}

func (t *Timestamp) attribute() {}

func (t *Timestamp) Value() time.Time {
	return t.value
}
