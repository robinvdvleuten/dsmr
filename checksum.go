package dsmr

import (
	"fmt"
	"strings"

	"github.com/snksoft/crc"
)

func verifyChecksum(t *Telegram, raw string, opts *parseOptions) error {
	// Only check footer if verifying is enabled and we found one while parsing
	if !opts.verifyChecksum || t.Footer.Value == "" {
		return nil
	}

	// Compute expected checksum from original message (including the "!" character).
	msg, _, _ := strings.Cut(raw, "!")
	want := crc16(msg + "!")

	if t.Footer.Value != want {
		return &ChecksumError{Got: t.Footer.Value, Want: want}
	}

	return nil
}

func (t *Telegram) checksum() (string, error) {
	var b strings.Builder
	if err := t.appendPayload(&b); err != nil {
		return "", err
	}
	return crc16(b.String()), nil
}

func crc16(msg string) string {
	return fmt.Sprintf("%04X", crc.CalculateCRC(crc.CRC16, []byte(msg)))
}
