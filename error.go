package dsmr

import "fmt"

type Error interface {
	error
}

type ChecksumError struct {
	Unexpected string
	Expect     string
}

func (e *ChecksumError) Error() string {
	return fmt.Sprintf("unexpected checksum \"%s\" (expected \"%s\")", e.Unexpected, e.Expect)
}
