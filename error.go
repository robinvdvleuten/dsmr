package dsmr

import "fmt"

type ChecksumError struct {
	Got  string
	Want string
}

func (e *ChecksumError) Error() string {
	return fmt.Sprintf("unexpected checksum \"%s\" (expected \"%s\")", e.Got, e.Want)
}
