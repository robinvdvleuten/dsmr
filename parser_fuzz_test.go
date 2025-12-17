package dsmr

import "testing"

func FuzzParse(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		telegram, err := Parse(string(data), VerifyChecksum(false))
		if err != nil {
			return
		}

		if telegram == nil || telegram.Header == nil || telegram.Footer == nil {
			t.Fatalf("parsed telegram missing header or footer")
		}

		normalizeTelegram(telegram)
	})
}

func FuzzParseObject(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		parsed, err := objectParser.ParseString("", string(data))
		if err != nil {
			return
		}

		if parsed == nil {
			t.Fatalf("parsed object is nil")
		}

		normalizeObjectForTest(&Object{OBIS: parsed.OBIS, Value: parsed.Value})
	})
}
