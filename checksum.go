package dsmr

import (
	"fmt"
	"io"
	"strings"

	"github.com/snksoft/crc"
)

func (a *AST) check(r io.Reader) error {
	// Only check footer if we found one while parsing.
	if a.Footer.Value.Value == "" {
		return nil
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	// Compute expected checksum from original message (including the "!" character).
	msg := strings.Split(string(b), "!")[0] + "!"
	checksum := fmt.Sprintf("%04X", crc.CalculateCRC(crc.CRC16, []byte(msg)))

	if a.Footer.Value.Value != checksum {
		return &ChecksumError{Unexpected: checksum, Expect: a.Footer.Value.Value}
	}

	return nil
}
