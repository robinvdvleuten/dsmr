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

	ast := &AST{Footer: &Footer{Value: &String{Value: "75B7"}}}
	err := ast.check(strings.NewReader(raw))
	require.NoError(t, err)
}

func TestInvalidChecksum(t *testing.T) {
	raw := "" +
		"/header\r\n" +
		"0-0:0.0.0()\r\n" +
		"!1234\r\n"

	ast := &AST{Footer: &Footer{Value: &String{Value: "1234"}}}
	err := ast.check(strings.NewReader(raw))
	require.EqualError(t, err, "unexpected checksum \"75B7\" (expected \"1234\")")
}
