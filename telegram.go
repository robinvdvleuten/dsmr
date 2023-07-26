package dsmr

type Telegram struct {
	Header   string   `"/" @~Separator+`
	COSEM    []*COSEM `@@+`
	Checksum string   `"!" @~Separator*`
}

type COSEM struct {
	OBIS      *OBIS       `@@`
	Attribute []Attribute `("(" @@* ")")+`
}

type Attribute interface{ value() }

type Measurement struct {
	Value string `@Number "*"`
	Unit  string `@~")"+`
}

func (Measurement) value() {}

type OBIS struct {
	Value string `@OBIS`
}

func (OBIS) value() {}

type Text struct {
	Value string `@~")"+`
}

func (Text) value() {}

type Timestamp struct {
	Value string `@Timestamp`
}

func (Timestamp) value() {}
