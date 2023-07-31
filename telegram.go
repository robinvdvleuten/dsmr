//go:generate peg -inline -switch grammar.peg

package dsmr

type Telegram struct {
	Header   string
	COSEM    map[string]*COSEM
	Checksum string
}

type COSEM struct {
	OBIS      *OBIS
	Attribute []Attribute
}

type Attribute interface{ value() }

type Measurement struct {
	Value string
	Unit  string
}

func (Measurement) value() {}

type OBIS struct {
	Value string
}

func (OBIS) value() {}

type Text struct {
	Value string
}

func (Text) value() {}

type Timestamp struct {
	Value string
}

func (Timestamp) value() {}
