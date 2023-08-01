package dsmr

import (
	"fmt"
	"io"
	"strings"

	"github.com/snksoft/crc"
)

func (t *Telegram) check(r io.Reader) error {
	// Only check checksum if we found one while parsing.
	if t.checksum == "" {
		return nil
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	// Compute expected checksum from original message (including the "!" character).
	msg := strings.Split(string(b), "!")[0] + "!"
	checksum := fmt.Sprintf("%04X", crc.CalculateCRC(crc.CRC16, []byte(msg)))

	if t.checksum != checksum {
		return &ChecksumError{Unexpected: checksum, Expect: t.checksum}
	}

	return nil
}
