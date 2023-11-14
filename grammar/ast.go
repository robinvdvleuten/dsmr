package grammar

type Node interface{}

type Telegram struct {
	Header *Header
	Data   []*Object
	Footer *Footer
}

type Header struct {
	Value string
}

type Footer struct {
	Value string
}

// Object is a COSEM object in the Telegram represented by the
// OBIS (Object Identification System) and one or more attributes.
type Object struct {
	Id    *OBIS
	Value interface{}
}

// ...
type OBIS struct {
	A int
	B int
	C int
	D int
	E int
}
