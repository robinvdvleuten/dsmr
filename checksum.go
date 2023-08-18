package dsmr

import (
	"fmt"
	"io"
	"strings"

	"github.com/snksoft/crc"
)

func (a *AST) VerifyChecksum(r io.Reader) error {
	// Only check footer if we found one while parsing.
	if a.Footer.Value == nil {
		return nil
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	// Compute expected checksum from original message (including the "!" character).
	msg, _, _ := strings.Cut(string(b), "!")
	checksum := fmt.Sprintf("%04X", crc.CalculateCRC(crc.CRC16, []byte(msg+"!")))

	if a.Footer.Value.Value != checksum {
		return &ChecksumError{Unexpected: checksum, Expect: a.Footer.Value.Value}
	}

	return nil
}
