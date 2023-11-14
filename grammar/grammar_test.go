package grammar

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/alecthomas/repr"
)

func TestTelegramParsing(t *testing.T) {
	type testCase struct {
		telegram string
		fail     string
		expected interface{}
	}

	tests := map[string]testCase{
		"v2.2": {
			telegram: "" +
				"/ISk5\\2MT382-1004\r\n" +
				"\r\n" +
				"0-0:96.1.1(00000000000000)\r\n" +
				// "1-0:1.8.1(00001.001*kWh)\r\n" +
				// "1-0:1.8.2(00001.001*kWh)\r\n" +
				// "1-0:2.8.1(00001.001*kWh)\r\n" +
				// "1-0:2.8.2(00001.001*kWh)\r\n" +
				// "0-0:96.14.0(0001)\r\n" +
				// "1-0:1.7.0(0001.01*kW)\r\n" +
				// "1-0:2.7.0(0000.00*kW)\r\n" +
				// "0-0:17.0.0(0999.00*kW)\r\n" +
				// "0-0:96.3.10(1)\r\n" +
				// "0-0:96.13.1()\r\n" +
				// "0-0:96.13.0()\r\n" +
				// "0-1:24.1.0(3)\r\n" +
				// "0-1:96.1.0(000000000000)\r\n" +
				// "0-1:24.3.0(161107190000)(00)(60)(1)(0-1:24.2.1)(m3)\r\n" +
				// "(00001.001)\r\n" +
				// "0-1:24.4.0(1)\r\n" +
				"!\r\n",
			expected: &Telegram{
				Header: header("ISk5\\2MT382-1004"),
				Data: []*Object{
					obj(obis(0, 0, 96, 1, 1), "00000000000000"),
					// obj(obis(1, 0, 1, 8, 1), "00001.001*kWh"),
					// obj(obis(1, 0, 1, 8, 2), "00001.001*kWh"),
					// obj(obis(1, 0, 2, 8, 1), "00001.001*kWh"),
					// obj(obis(1, 0, 2, 8, 2), "00001.001*kWh"),
				},
				Footer: footer(""),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			telegram, err := Parse("", []byte(test.telegram))
			if test.fail != "" {
				assert.EqualError(t, err, test.fail)
			} else {
				assert.NoError(t, err)
				assert.Equal(t,
					repr.String(test.expected, repr.Indent("  ")),
					repr.String(telegram, repr.Indent("  ")))
			}
		})
	}
}

func header(v string) *Header {
	return &Header{Value: v}
}

func footer(v string) *Footer {
	return &Footer{Value: v}
}

func obj(o *OBIS, v string) *Object {
	return &Object{Id: o, Value: v}
}

func obis(a, b, c, d, e int) *OBIS {
	return &OBIS{a, b, c, d, e}
}
