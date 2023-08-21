package dsmr

import (
	"reflect"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/alecthomas/repr"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		telegram string
		fail     string
		dest     interface{}
	}{
		{
			name: "ScalarAttributes",
			telegram: "" +
				"/ISk5\\2MT382-1004\r\n" +
				"\r\n" +
				"0-0:96.1.1(4B384547303034303436333935353037)\r\n" +
				"1-0:1.8.1(12345.678*kWh)\r\n" +
				"1-0:32.32.0(00000)\r\n" +
				"!\r\n",
			dest: struct {
				Str   string  `dsmr:"0-0:96.1.1"`
				Float float64 `dsmr:"1-0:1.8.1"`
				Int   int32   `dsmr:"1-0:32.32.0"`
			}{
				Str:   "4B384547303034303436333935353037",
				Float: 12345.678,
				Int:   0,
			},
		},
		{
			name: "TimestampAsTime",
			telegram: "" +
				"/ISk5\\2MT382-1004\r\n" +
				"\r\n" +
				"0-0:1.0.0(161113205757W)\r\n" +
				"!\r\n",
			dest: struct {
				Time time.Time `dsmr:"0-0:1.0.0"`
			}{
				Time: time.Now(),
			},
		},
		{
			name: "MeasurementAsStringListAttribute",
			telegram: "" +
				"/ISk5\\2MT382-1004\r\n" +
				"\r\n" +
				"1-0:1.8.1(12345.678*kWh)\r\n" +
				"!\r\n",
			dest: struct {
				List []string `dsmr:"1-0:1.8.1"`
			}{
				List: []string{"12345.678", "kWh"},
			},
		},
		{
			name: "MeasurementAsAnyListAttribute",
			telegram: "" +
				"/ISk5\\2MT382-1004\r\n" +
				"\r\n" +
				"1-0:1.8.1(12345.678*kWh)\r\n" +
				"!\r\n",
			dest: struct {
				List []any `dsmr:"1-0:1.8.1"`
			}{
				List: append(make([]any, 0), 12345.678, "kWh"),
			},
		},
		{
			name: "ListAsStringListAttribute",
			telegram: "" +
				"/ISk5\\2MT382-1004\r\n" +
				"\r\n" +
				"0-1:24.2.1(161129200000W)(00981.443*m3)\r\n" +
				"!\r\n",
			dest: struct {
				List []string `dsmr:"0-1:24.2.1"`
			}{
				List: []string{"161129200000W", "981.443"},
			},
		},
		{
			name: "ListAsAnyListAttribute",
			telegram: "" +
				"/ISk5\\2MT382-1004\r\n" +
				"\r\n" +
				"0-1:24.2.1(161129200000W)(00981.443*m3)\r\n" +
				"!\r\n",
			dest: struct {
				List []any `dsmr:"0-1:24.2.1"`
			}{
				List: append(make([]any, 0), "161129200000W", 981.443),
			},
		},
		{
			name: "PointerScalars",
			telegram: "" +
				"/ISk5\\2MT382-1004\r\n" +
				"\r\n" +
				"0-0:96.1.1(4B384547303034303436333935353037)\r\n" +
				"!\r\n",
			dest: struct {
				Ptr *string `dsmr:"0-0:96.1.1"`
			}{Ptr: strp("4B384547303034303436333935353037")},
		},
		{
			name: "PointerScalarsNil",
			telegram: "" +
				"/ISk5\\2MT382-1004\r\n" +
				"\r\n" +
				"!\r\n",
			dest: struct {
				Ptr *string `dsmr:"0-0:96.1.1,optional"`
			}{Ptr: nil},
		},
		{
			name: "MissingRequired",
			telegram: "" +
				"/ISk5\\2MT382-1004\r\n" +
				"\r\n" +
				"!\r\n",
			dest: struct {
				Str string `dsmr:"0-0:96.1.1"`
			}{},
			fail: `missing required attribute "0-0:96.1.1"`,
		},
		{
			name: "MissingWithDefault",
			telegram: "" +
				"/ISk5\\2MT382-1004\r\n" +
				"\r\n" +
				"!\r\n",
			dest: struct {
				Str   string  `dsmr:"0-0:96.1.1,optional" default:"foo"`
				Float float64 `dsmr:"1-0:52.36.0,optional" default:"09.84"`
				Int   int32   `dsmr:"1-0:52.36.0,optional" default:"3"`
			}{
				Str:   "foo",
				Float: 9.84,
				Int:   3,
			},
		},
		{
			name: "WrongDefault",
			telegram: "" +
				"/ISk5\\2MT382-1004\r\n" +
				"\r\n" +
				"!\r\n",
			dest: struct {
				Int int32 `dsmr:"0-0:96.1.1,optional" default:"foo"`
			}{},
			fail: `error parsing default value: error converting "foo" to int`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Helper()
			rv := reflect.New(reflect.TypeOf(test.dest))
			actual := rv.Interface()
			err := Unmarshal([]byte(test.telegram), actual)
			if test.fail != "" {
				assert.EqualError(t, err, test.fail)
			} else {
				assert.NoError(t, err)
				assert.Equal(t,
					repr.String(test.dest, repr.Indent("  ")),
					repr.String(rv.Elem().Interface(), repr.Indent("  ")))
			}
		})
	}
}

func strp(s string) *string { return &s }
