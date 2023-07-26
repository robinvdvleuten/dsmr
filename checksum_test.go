package dsmr

import (
	"strings"
	"testing"

	require "github.com/alecthomas/assert/v2"
)

func TestValidChecksum(t *testing.T) {
	raw := "" +
		"/header\r\n" +
		"0-0:0.0.0()\r\n" +
		"!75B7\r\n"

	telegram := &Telegram{Checksum: "75B7"}

	err := telegram.check(strings.NewReader(raw))
	require.NoError(t, err)
}

func TestInvalidChecksum(t *testing.T) {
	raw := "" +
		"/header\r\n" +
		"0-0:0.0.0()\r\n" +
		"!1234\r\n"

	telegram := &Telegram{Checksum: "1234"}

	err := telegram.check(strings.NewReader(raw))
	require.EqualError(t, err, "unexpected checksum \"75B7\" (expected \"1234\")")
}
