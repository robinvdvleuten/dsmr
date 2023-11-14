package dsmr

import (
	"fmt"
	"strings"

	"github.com/snksoft/crc"
)

func (t *Telegram) VerifyChecksum(str string) error {
	// Only check footer if we found one while parsing.
	if t.Footer.Value == "" {
		return nil
	}

	// Compute expected checksum from original message (including the "!" character).
	msg, _, _ := strings.Cut(str, "!")
	checksum := fmt.Sprintf("%04X", crc.CalculateCRC(crc.CRC16, []byte(msg+"!")))

	if t.Footer.Value != checksum {
		return &ChecksumError{Unexpected: checksum, Expect: t.Footer.Value}
	}

	return nil
}
