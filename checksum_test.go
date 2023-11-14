package dsmr

import (
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestValidChecksum(t *testing.T) {
	raw := "" +
		"/header\r\n" +
		"0-0:0.0.0()\r\n" +
		"!75B7\r\n"

	ast := &AST{Footer: &Footer{Value: "75B7"}}
	err := ast.VerifyChecksum(strings.NewReader(raw))
	assert.NoError(t, err)
}

func TestInvalidChecksum(t *testing.T) {
	raw := "" +
		"/header\r\n" +
		"0-0:0.0.0()\r\n" +
		"!1234\r\n"

	ast := &AST{Footer: &Footer{Value: "1234"}}
	err := ast.VerifyChecksum(strings.NewReader(raw))
	assert.EqualError(t, err, "unexpected checksum \"75B7\" (expected \"1234\")")
}
