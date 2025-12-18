package dsmr

import "fmt"

// lookupObject finds the first Object with the given OBIS code in main data.
func (t *Telegram) lookupObject(obis string) *Object {
	for _, entry := range t.Data {
		if obj, ok := entry.(*Object); ok && obj.OBIS.Value == obis {
			return obj
		}
	}
	return nil
}

// Type-safe extractors
func asMeasurement(obj *Object) *Measurement {
	if obj == nil {
		return nil
	}
	m, _ := obj.Value.(*Measurement)
	return m
}

func asLastCapture(obj *Object) *LastCapture {
	if obj == nil {
		return nil
	}
	lc, _ := obj.Value.(*LastCapture)
	return lc
}

func asString(obj *Object) *String {
	if obj == nil {
		return nil
	}
	s, _ := obj.Value.(*String)
	return s
}

// MBusDevice returns the M-Bus device on the given channel, or nil.
func (t *Telegram) MBusDevice(channel int) *MBusDevice {
	for _, entry := range t.Data {
		if mbus, ok := entry.(*MBusDevice); ok && mbus.Channel == channel {
			return mbus
		}
	}
	return nil
}

// Version returns the DSMR protocol version.
func (t *Telegram) Version() *String {
	return asString(t.lookupObject("1-3:0.2.8"))
}

// MeasuredAt returns the timestamp of the measurement.
func (t *Telegram) MeasuredAt() *Timestamp {
	obj := t.lookupObject("0-0:1.0.0")
	if obj == nil {
		return nil
	}
	ts, _ := obj.Value.(*Timestamp)
	return ts
}

// EquipmentID returns the meter equipment ID.
func (t *Telegram) EquipmentID() *String {
	return asString(t.lookupObject("0-0:96.1.1"))
}

// ElectricityDelivered returns electricity delivered for the given tariff (1 or 2).
func (t *Telegram) ElectricityDelivered(tariff int) *Measurement {
	return asMeasurement(t.lookupObject(fmt.Sprintf("1-0:1.8.%d", tariff)))
}

// ElectricityReceived returns electricity received for the given tariff (1 or 2).
func (t *Telegram) ElectricityReceived(tariff int) *Measurement {
	return asMeasurement(t.lookupObject(fmt.Sprintf("1-0:2.8.%d", tariff)))
}

// ElectricityTariffIndicator returns the active tariff indicator.
func (t *Telegram) ElectricityTariffIndicator() *String {
	return asString(t.lookupObject("0-0:96.14.0"))
}

// ElectricityCurrentlyDelivered returns current electricity delivered.
func (t *Telegram) ElectricityCurrentlyDelivered() *Measurement {
	return asMeasurement(t.lookupObject("1-0:1.7.0"))
}

// ElectricityCurrentlyReceived returns current electricity received.
func (t *Telegram) ElectricityCurrentlyReceived() *Measurement {
	return asMeasurement(t.lookupObject("1-0:2.7.0"))
}

// PowerFailures returns the number of power failures.
func (t *Telegram) PowerFailures() *String {
	return asString(t.lookupObject("0-0:96.7.21"))
}

// PowerFailuresLong returns the number of long power failures.
func (t *Telegram) PowerFailuresLong() *String {
	return asString(t.lookupObject("0-0:96.7.9"))
}

// VoltageSagsL1 returns the number of voltage sags on phase L1.
func (t *Telegram) VoltageSagsL1() *String {
	return asString(t.lookupObject("1-0:32.32.0"))
}

// VoltageSagsL2 returns the number of voltage sags on phase L2.
func (t *Telegram) VoltageSagsL2() *String {
	return asString(t.lookupObject("1-0:52.32.0"))
}

// VoltageSagsL3 returns the number of voltage sags on phase L3.
func (t *Telegram) VoltageSagsL3() *String {
	return asString(t.lookupObject("1-0:72.32.0"))
}

// VoltageSwellsL1 returns the number of voltage swells on phase L1.
func (t *Telegram) VoltageSwellsL1() *String {
	return asString(t.lookupObject("1-0:32.36.0"))
}

// VoltageSwellsL2 returns the number of voltage swells on phase L2.
func (t *Telegram) VoltageSwellsL2() *String {
	return asString(t.lookupObject("1-0:52.36.0"))
}

// VoltageSwellsL3 returns the number of voltage swells on phase L3.
func (t *Telegram) VoltageSwellsL3() *String {
	return asString(t.lookupObject("1-0:72.36.0"))
}

// CurrentL1 returns current on phase L1.
func (t *Telegram) CurrentL1() *Measurement {
	return asMeasurement(t.lookupObject("1-0:31.7.0"))
}

// CurrentL2 returns current on phase L2.
func (t *Telegram) CurrentL2() *Measurement {
	return asMeasurement(t.lookupObject("1-0:51.7.0"))
}

// CurrentL3 returns current on phase L3.
func (t *Telegram) CurrentL3() *Measurement {
	return asMeasurement(t.lookupObject("1-0:71.7.0"))
}

// VoltageL1 returns voltage on phase L1.
func (t *Telegram) VoltageL1() *Measurement {
	return asMeasurement(t.lookupObject("1-0:32.7.0"))
}

// VoltageL2 returns voltage on phase L2.
func (t *Telegram) VoltageL2() *Measurement {
	return asMeasurement(t.lookupObject("1-0:52.7.0"))
}

// VoltageL3 returns voltage on phase L3.
func (t *Telegram) VoltageL3() *Measurement {
	return asMeasurement(t.lookupObject("1-0:72.7.0"))
}

// PowerDeliveredL1 returns power delivered on phase L1.
func (t *Telegram) PowerDeliveredL1() *Measurement {
	return asMeasurement(t.lookupObject("1-0:21.7.0"))
}

// PowerDeliveredL2 returns power delivered on phase L2.
func (t *Telegram) PowerDeliveredL2() *Measurement {
	return asMeasurement(t.lookupObject("1-0:41.7.0"))
}

// PowerDeliveredL3 returns power delivered on phase L3.
func (t *Telegram) PowerDeliveredL3() *Measurement {
	return asMeasurement(t.lookupObject("1-0:61.7.0"))
}

// PowerReceivedL1 returns power received on phase L1.
func (t *Telegram) PowerReceivedL1() *Measurement {
	return asMeasurement(t.lookupObject("1-0:22.7.0"))
}

// PowerReceivedL2 returns power received on phase L2.
func (t *Telegram) PowerReceivedL2() *Measurement {
	return asMeasurement(t.lookupObject("1-0:42.7.0"))
}

// PowerReceivedL3 returns power received on phase L3.
func (t *Telegram) PowerReceivedL3() *Measurement {
	return asMeasurement(t.lookupObject("1-0:62.7.0"))
}

// GasDelivered returns the last 5-minute gas reading.
func (t *Telegram) GasDelivered() *LastCapture {
	return asLastCapture(t.lookupObject("0-1:24.2.1"))
}
