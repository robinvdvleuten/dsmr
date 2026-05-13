package dsmr

import (
	"fmt"
	"strings"
)

// MarshalText renders a telegram as DSMR text.
func (t *Telegram) MarshalText() ([]byte, error) {
	if t == nil {
		return nil, nil
	}

	var b strings.Builder
	if err := t.appendTo(&b); err != nil {
		return nil, err
	}

	return []byte(b.String()), nil
}

// String renders a telegram as DSMR text.
func (t *Telegram) String() string {
	text, err := t.MarshalText()
	if err != nil {
		return ""
	}
	return string(text)
}

func (t *Telegram) appendTo(b *strings.Builder) error {
	if err := t.appendPayload(b); err != nil {
		return err
	}
	if t.shouldWriteChecksum() {
		checksum, err := t.checksum()
		if err != nil {
			return err
		}
		b.WriteString(checksum)
	}
	b.WriteString("\r\n")
	return nil
}

func (t *Telegram) appendPayload(b *strings.Builder) error {
	if t.Header == nil {
		return fmt.Errorf("telegram header is required")
	}

	b.WriteByte('/')
	b.WriteString(t.Header.Value)
	b.WriteString("\r\n\r\n")

	for _, entry := range t.Data {
		if err := appendEntry(b, entry); err != nil {
			return err
		}
	}

	b.WriteByte('!')
	return nil
}

func (t *Telegram) shouldWriteChecksum() bool {
	return t.Footer == nil || t.Footer.Value != ""
}

func appendEntry(b *strings.Builder, entry Entry) error {
	switch entry := entry.(type) {
	case nil:
		return nil
	case *Object:
		return appendObject(b, entry)
	case *MBusDevice:
		for _, obj := range entry.Data {
			if err := appendObject(b, obj); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported entry %T", entry)
	}
}

func appendObject(b *strings.Builder, obj *Object) error {
	if obj == nil {
		return nil
	}
	if obj.OBIS == nil {
		return fmt.Errorf("object OBIS is required")
	}
	b.WriteString(obj.OBIS.Value)
	b.WriteByte('(')
	if err := appendValue(b, obj.Value); err != nil {
		return err
	}
	b.WriteString(")\r\n")
	return nil
}

func appendValue(b *strings.Builder, value Value) error {
	switch value := value.(type) {
	case nil:
		return nil
	case *String:
		b.WriteString(value.Value)
	case *Timestamp:
		appendTimestamp(b, value)
	case *Measurement:
		return appendMeasurement(b, value)
	case *LegacyMeasurement:
		return appendLegacyMeasurement(b, value)
	case *Event:
		return appendEvent(b, value, false)
	case *EventLog:
		return appendEventLog(b, value)
	case *LastCapture:
		return appendLastCapture(b, value)
	case *LegacyLastCapture:
		return appendLegacyLastCapture(b, value)
	default:
		return fmt.Errorf("unsupported value %T", value)
	}
	return nil
}

func appendTimestamp(b *strings.Builder, ts *Timestamp) {
	if ts == nil {
		return
	}
	b.WriteString(ts.Value)
	switch ts.DST {
	case DSTWinter:
		b.WriteByte('W')
	case DSTSummer:
		b.WriteByte('S')
	}
}

func appendMeasurement(b *strings.Builder, measurement *Measurement) error {
	if measurement == nil {
		return nil
	}
	if measurement.Value == nil {
		return fmt.Errorf("measurement value is required")
	}
	if measurement.Unit == nil {
		return fmt.Errorf("measurement unit is required")
	}
	b.WriteString(measurement.Value.String())
	b.WriteByte('*')
	b.WriteString(measurement.Unit.Value)
	return nil
}

func appendLegacyMeasurement(b *strings.Builder, measurement *LegacyMeasurement) error {
	if measurement == nil {
		return nil
	}
	if measurement.Unit == nil {
		return fmt.Errorf("legacy measurement unit is required")
	}
	if measurement.Value == nil {
		return fmt.Errorf("legacy measurement value is required")
	}
	b.WriteString(measurement.Unit.Value)
	b.WriteString(")\r\n(")
	b.WriteString(measurement.Value.String())
	return nil
}

func appendEventLog(b *strings.Builder, log *EventLog) error {
	if log == nil {
		return nil
	}
	if log.Count == nil {
		return fmt.Errorf("event log count is required")
	}
	if log.OBIS == nil {
		return fmt.Errorf("event log OBIS is required")
	}

	b.WriteString(log.Count.String())
	b.WriteString(")(")
	b.WriteString(log.OBIS.Value)
	for i, event := range log.Value {
		if i == 0 {
			b.WriteByte(')')
		}
		if err := appendEvent(b, event, false); err != nil {
			return err
		}
		if i < len(log.Value)-1 {
			b.WriteByte(')')
		}
	}
	return nil
}

func appendEvent(b *strings.Builder, event *Event, closeMeasurement bool) error {
	if event == nil {
		return nil
	}
	if event.Timestamp == nil {
		return fmt.Errorf("event timestamp is required")
	}
	if event.Value == nil {
		return fmt.Errorf("event value is required")
	}

	b.WriteByte('(')
	appendTimestamp(b, event.Timestamp)
	b.WriteString(")(")
	if err := appendMeasurement(b, event.Value); err != nil {
		return err
	}
	if closeMeasurement {
		b.WriteByte(')')
	}
	return nil
}

func appendLastCapture(b *strings.Builder, capture *LastCapture) error {
	if capture == nil {
		return nil
	}
	if capture.Timestamp == nil {
		return fmt.Errorf("last capture timestamp is required")
	}
	if capture.Value == nil {
		return fmt.Errorf("last capture value is required")
	}

	appendTimestamp(b, capture.Timestamp)
	b.WriteString(")(")
	return appendMeasurement(b, capture.Value)
}

func appendLegacyLastCapture(b *strings.Builder, capture *LegacyLastCapture) error {
	if capture == nil {
		return nil
	}
	if capture.Timestamp == nil {
		return fmt.Errorf("legacy last capture timestamp is required")
	}
	for _, extra := range capture.Extra {
		if extra == nil {
			return fmt.Errorf("legacy last capture extra value is required")
		}
	}
	if capture.OBIS == nil {
		return fmt.Errorf("legacy last capture OBIS is required")
	}
	if capture.Value == nil {
		return fmt.Errorf("legacy last capture value is required")
	}

	b.WriteString(capture.Timestamp.Value)
	for _, extra := range capture.Extra {
		b.WriteString(")(")
		b.WriteString(extra.Value)
	}
	b.WriteString(")(")
	b.WriteString(capture.OBIS.Value)
	b.WriteString(")(")
	return appendLegacyMeasurement(b, capture.Value)
}
