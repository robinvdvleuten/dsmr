package dsmr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

func TestTelegramMarshalTextRoundTrip(t *testing.T) {
	for _, test := range marshalRoundTripTelegrams() {
		t.Run(test.name, func(t *testing.T) {
			telegram, err := Parse(test.telegram, VerifyChecksum(false))
			assert.NoError(t, err)

			text, err := telegram.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, test.telegram, string(text))
			assert.Equal(t, test.telegram, telegram.String())
		})
	}
}

func TestTelegramMarshalTextCustomTelegram(t *testing.T) {
	amount, err := NewNumber("000004.426")
	assert.NoError(t, err)

	telegram := &Telegram{
		Header: &Header{Value: "ISk5\\2MT382-1000"},
		Data: Data{
			obj("1-3:0.2.8", str("50")),
			obj("1-0:1.8.1", &Measurement{Value: amount, Unit: str("kWh")}),
		},
		Footer: footer("BEEF"),
	}

	payload := "" +
		"/ISk5\\2MT382-1000\r\n" +
		"\r\n" +
		"1-3:0.2.8(50)\r\n" +
		"1-0:1.8.1(000004.426*kWh)\r\n" +
		"!"
	checksum, err := telegram.checksum()
	assert.NoError(t, err)
	expected := payload + checksum + "\r\n"

	text, err := telegram.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, expected, string(text))
}

func TestTelegramMarshalTextValueRendering(t *testing.T) {
	tests := []struct {
		name string
		data Data
		line string
	}{
		{
			name: "empty string",
			data: Data{obj("0-0:96.13.0", nil)},
			line: "0-0:96.13.0()\r\n",
		},
		{
			name: "timestamp",
			data: Data{obj("0-0:1.0.0", ts("161030020000", DSTSummer))},
			line: "0-0:1.0.0(161030020000S)\r\n",
		},
		{
			name: "measurement",
			data: Data{obj("1-0:1.8.1", mm("000004.426", "kWh"))},
			line: "1-0:1.8.1(000004.426*kWh)\r\n",
		},
		{
			name: "last capture",
			data: Data{obj("0-1:24.2.1", lc(ts("161030020000", DSTSummer), mm("00000.107", "m3")))},
			line: "0-1:24.2.1(161030020000S)(00000.107*m3)\r\n",
		},
		{
			name: "event log",
			data: Data{obj("1-0:99.97.0", events("1", "0-0:96.7.19", event(ts("161107190000", DSTSummer), "00015")))},
			line: "1-0:99.97.0(1)(0-0:96.7.19)(161107190000S)(00015*s)\r\n",
		},
		{
			name: "m-bus device",
			data: Data{mbus(1,
				obj("0-1:24.1.0", str("003")),
				obj("0-1:96.1.0", str("3232323241424344313233343536373839")),
			)},
			line: "0-1:24.1.0(003)\r\n0-1:96.1.0(3232323241424344313233343536373839)\r\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			telegram := &Telegram{Header: header("TEST"), Data: test.data}

			text, err := telegram.MarshalText()
			assert.NoError(t, err)
			assert.True(t, strings.Contains(string(text), test.line), string(text))
		})
	}
}

func TestNewNumber(t *testing.T) {
	for _, input := range []string{"000004.426", "00.000", "+1.23", "-1.23", ".5"} {
		t.Run(fmt.Sprintf("valid/%s", input), func(t *testing.T) {
			number, err := NewNumber(input)
			assert.NoError(t, err)
			assert.Equal(t, input, number.Text)
			assert.NotZero(t, number.Value)
		})
	}

	for _, input := range []string{"", "abc", "1.2.3", "1*kWh"} {
		t.Run(fmt.Sprintf("invalid/%s", input), func(t *testing.T) {
			number, err := NewNumber(input)
			assert.Error(t, err)
			assert.Zero(t, number)
		})
	}
}

func TestNumberParseKeepsText(t *testing.T) {
	parsed, err := initObjectParser().ParseString("", "1-0:1.8.1(000004.426*kWh)\r\n")
	assert.NoError(t, err)

	measurement, ok := parsed.Value.(*Measurement)
	assert.True(t, ok)
	assert.Equal(t, "000004.426", measurement.Value.Text)
	assert.NotZero(t, measurement.Value.Value)
}

func TestTelegramMarshalTextErrors(t *testing.T) {
	tests := []struct {
		name     string
		telegram *Telegram
	}{
		{
			name:     "nil header",
			telegram: &Telegram{},
		},
		{
			name: "nil object OBIS",
			telegram: &Telegram{
				Header: header("TEST"),
				Data:   Data{&Object{Value: str("value")}},
			},
		},
		{
			name: "unsupported value",
			telegram: &Telegram{
				Header: header("TEST"),
				Data:   Data{obj("0-0:96.13.0", unknownValue{})},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			text, err := test.telegram.MarshalText()
			assert.Error(t, err)
			assert.Zero(t, text)
		})
	}
}

type unknownValue struct{}

func (unknownValue) value() {}

func (unknownValue) Position() lexer.Position { return lexer.Position{} }

func (unknownValue) children() []Node { return nil }

func marshalRoundTripTelegrams() []struct {
	name     string
	telegram string
} {
	return []struct {
		name     string
		telegram string
	}{
		{
			name: "v2.2",
			telegram: "" +
				"/ISk5\\2MT382-1004\r\n" +
				"\r\n" +
				"0-0:96.1.1(00000000000000)\r\n" +
				"1-0:1.8.1(00001.001*kWh)\r\n" +
				"1-0:1.8.2(00001.001*kWh)\r\n" +
				"1-0:2.8.1(00001.001*kWh)\r\n" +
				"1-0:2.8.2(00001.001*kWh)\r\n" +
				"0-0:96.14.0(0001)\r\n" +
				"1-0:1.7.0(0001.01*kW)\r\n" +
				"1-0:2.7.0(0000.00*kW)\r\n" +
				"0-0:17.0.0(0999.00*kW)\r\n" +
				"0-0:96.3.10(1)\r\n" +
				"0-0:96.13.1()\r\n" +
				"0-0:96.13.0()\r\n" +
				"0-1:24.1.0(3)\r\n" +
				"0-1:96.1.0(000000000000)\r\n" +
				"0-1:24.3.0(161107190000)(00)(60)(1)(0-1:24.2.1)(m3)\r\n" +
				"(00001.001)\r\n" +
				"0-1:24.4.0(1)\r\n" +
				"!\r\n",
		},
		{
			name: "v3.0",
			telegram: "" +
				"/ISk5\\2MT382-1000\r\n" +
				"\r\n" +
				"0-0:96.1.1(4B384547303034303436333935353037)\r\n" +
				"1-0:1.8.1(12345.678*kWh)\r\n" +
				"1-0:1.8.2(12345.678*kWh)\r\n" +
				"1-0:2.8.1(12345.678*kWh)\r\n" +
				"1-0:2.8.2(12345.678*kWh)\r\n" +
				"0-0:96.14.0(0002)\r\n" +
				"1-0:1.7.0(001.19*kW)\r\n" +
				"1-0:2.7.0(000.00*kW)\r\n" +
				"0-0:17.0.0(016*A)\r\n" +
				"0-0:96.3.10(1)\r\n" +
				"0-0:96.13.1(303132333435363738)\r\n" +
				"0-0:96.13.0(303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E" +
				"3F303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F30313233" +
				"3435363738393A3B3C3D3E3F)\r\n" +
				"0-1:96.1.0(3232323241424344313233343536373839)\r\n" +
				"0-1:24.1.0(03)\r\n" +
				"0-1:24.3.0(090212160000)(00)(60)(1)(0-1:24.2.1)(m3)\r\n" +
				"(00001.001)\r\n" +
				"0-1:24.4.0(1)\r\n" +
				"!\r\n",
		},
		{
			name: "v4.2",
			telegram: "" +
				"/KFM5KAIFA-METER\r\n" +
				"\r\n" +
				"1-3:0.2.8(42)\r\n" +
				"0-0:1.0.0(161113205757W)\r\n" +
				"0-0:96.1.1(3960221976967177082151037881335713)\r\n" +
				"1-0:1.8.1(001581.123*kWh)\r\n" +
				"1-0:1.8.2(001435.706*kWh)\r\n" +
				"1-0:2.8.1(000000.000*kWh)\r\n" +
				"1-0:2.8.2(000000.000*kWh)\r\n" +
				"0-0:96.14.0(0002)\r\n" +
				"1-0:1.7.0(02.027*kW)\r\n" +
				"1-0:2.7.0(00.000*kW)\r\n" +
				"0-0:96.7.21(00015)\r\n" +
				"0-0:96.7.9(00007)\r\n" +
				"1-0:99.97.0(3)(0-0:96.7.19)(000104180320W)(0000237126*s)(000101000001W)" +
				"(2147583646*s)(000102000003W)(2317482647*s)\r\n" +
				"1-0:32.32.0(00000)\r\n" +
				"1-0:52.32.0(00000)\r\n" +
				"1-0:72.32.0(00000)\r\n" +
				"1-0:32.36.0(00000)\r\n" +
				"1-0:52.36.0(00000)\r\n" +
				"1-0:72.36.0(00000)\r\n" +
				"0-0:96.13.1()\r\n" +
				"0-0:96.13.0()\r\n" +
				"1-0:31.7.0(000*A)\r\n" +
				"1-0:51.7.0(006*A)\r\n" +
				"1-0:71.7.0(002*A)\r\n" +
				"1-0:21.7.0(00.170*kW)\r\n" +
				"1-0:22.7.0(00.000*kW)\r\n" +
				"1-0:41.7.0(01.247*kW)\r\n" +
				"1-0:42.7.0(00.000*kW)\r\n" +
				"1-0:61.7.0(00.209*kW)\r\n" +
				"1-0:62.7.0(00.000*kW)\r\n" +
				"0-1:24.1.0(003)\r\n" +
				"0-1:96.1.0(4819243993373755377509728609491464)\r\n" +
				"0-1:24.2.1(161129200000W)(00981.443*m3)\r\n" +
				"!6796\r\n",
		},
		{
			name: "v5.0",
			telegram: "" +
				"/ISk5\\2MT382-1000\r\n" +
				"\r\n" +
				"1-3:0.2.8(50)\r\n" +
				"0-0:1.0.0(161030020000S)\r\n" +
				"0-0:96.1.1(4B384547303034303436333935353037)\r\n" +
				"1-0:1.8.1(000004.426*kWh)\r\n" +
				"1-0:1.8.2(000002.399*kWh)\r\n" +
				"1-0:2.8.1(000002.444*kWh)\r\n" +
				"1-0:2.8.2(000000.000*kWh)\r\n" +
				"0-0:96.14.0(0002)\r\n" +
				"1-0:1.7.0(00.244*kW)\r\n" +
				"1-0:2.7.0(00.000*kW)\r\n" +
				"0-0:96.7.21(00013)\r\n" +
				"0-0:96.7.9(00000)\r\n" +
				"1-0:99.97.0(0)(0-0:96.7.19)\r\n" +
				"1-0:32.32.0(00000)\r\n" +
				"1-0:52.32.0(00000)\r\n" +
				"1-0:72.32.0(00000)\r\n" +
				"1-0:32.36.0(00000)\r\n" +
				"1-0:52.36.0(00000)\r\n" +
				"1-0:72.36.0(00000)\r\n" +
				"0-0:96.13.0()\r\n" +
				"1-0:32.7.0(0230.0*V)\r\n" +
				"1-0:52.7.0(0230.0*V)\r\n" +
				"1-0:72.7.0(0229.0*V)\r\n" +
				"1-0:31.7.0(0.48*A)\r\n" +
				"1-0:51.7.0(0.44*A)\r\n" +
				"1-0:71.7.0(0.86*A)\r\n" +
				"1-0:21.7.0(00.070*kW)\r\n" +
				"1-0:41.7.0(00.032*kW)\r\n" +
				"1-0:61.7.0(00.142*kW)\r\n" +
				"1-0:22.7.0(00.000*kW)\r\n" +
				"1-0:42.7.0(00.000*kW)\r\n" +
				"1-0:62.7.0(00.000*kW)\r\n" +
				"0-1:24.1.0(003)\r\n" +
				"0-1:96.1.0(3232323241424344313233343536373839)\r\n" +
				"0-1:24.2.1(161030020000S)(00000.107*m3)\r\n" +
				"0-2:24.1.0(003)\r\n" +
				"0-2:96.1.0()\r\n" +
				"!8397\r\n",
		},
	}
}
