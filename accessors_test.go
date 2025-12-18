package dsmr

import (
	"math/big"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestAccessors(t *testing.T) {
	telegram := &Telegram{
		Data: Data{
			&Object{
				OBIS:  &OBIS{Value: "1-3:0.2.8"},
				Value: &String{Value: "50"},
			},
			&Object{
				OBIS: &OBIS{Value: "0-0:1.0.0"},
				Value: &Timestamp{
					Value: "170102192002",
					DST:   DSTWinter,
				},
			},
			&Object{
				OBIS:  &OBIS{Value: "0-0:96.1.1"},
				Value: &String{Value: "4B384547303034303436333935353037"},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:1.8.1"},
				Value: &Measurement{
					Value: numFloat("4.426"),
					Unit:  &String{Value: "kWh"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:1.8.2"},
				Value: &Measurement{
					Value: numFloat("2.399"),
					Unit:  &String{Value: "kWh"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:2.8.1"},
				Value: &Measurement{
					Value: numFloat("2.444"),
					Unit:  &String{Value: "kWh"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:2.8.2"},
				Value: &Measurement{
					Value: numFloat("0.000"),
					Unit:  &String{Value: "kWh"},
				},
			},
			&Object{
				OBIS:  &OBIS{Value: "0-0:96.14.0"},
				Value: &String{Value: "0002"},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:1.7.0"},
				Value: &Measurement{
					Value: numFloat("0.244"),
					Unit:  &String{Value: "kW"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:2.7.0"},
				Value: &Measurement{
					Value: numFloat("0.000"),
					Unit:  &String{Value: "kW"},
				},
			},
			&Object{
				OBIS:  &OBIS{Value: "0-0:96.7.21"},
				Value: &String{Value: "00013"},
			},
			&Object{
				OBIS:  &OBIS{Value: "0-0:96.7.9"},
				Value: &String{Value: "00000"},
			},
			&Object{
				OBIS:  &OBIS{Value: "1-0:32.32.0"},
				Value: &String{Value: "00000"},
			},
			&Object{
				OBIS:  &OBIS{Value: "1-0:52.32.0"},
				Value: &String{Value: "00000"},
			},
			&Object{
				OBIS:  &OBIS{Value: "1-0:72.32.0"},
				Value: &String{Value: "00000"},
			},
			&Object{
				OBIS:  &OBIS{Value: "1-0:32.36.0"},
				Value: &String{Value: "00000"},
			},
			&Object{
				OBIS:  &OBIS{Value: "1-0:52.36.0"},
				Value: &String{Value: "00000"},
			},
			&Object{
				OBIS:  &OBIS{Value: "1-0:72.36.0"},
				Value: &String{Value: "00000"},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:31.7.0"},
				Value: &Measurement{
					Value: numFloat("0.48"),
					Unit:  &String{Value: "A"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:51.7.0"},
				Value: &Measurement{
					Value: numFloat("0.44"),
					Unit:  &String{Value: "A"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:71.7.0"},
				Value: &Measurement{
					Value: numFloat("0.86"),
					Unit:  &String{Value: "A"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:32.7.0"},
				Value: &Measurement{
					Value: numFloat("230.0"),
					Unit:  &String{Value: "V"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:52.7.0"},
				Value: &Measurement{
					Value: numFloat("230.0"),
					Unit:  &String{Value: "V"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:72.7.0"},
				Value: &Measurement{
					Value: numFloat("229.0"),
					Unit:  &String{Value: "V"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:21.7.0"},
				Value: &Measurement{
					Value: numFloat("0.070"),
					Unit:  &String{Value: "kW"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:41.7.0"},
				Value: &Measurement{
					Value: numFloat("0.032"),
					Unit:  &String{Value: "kW"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:61.7.0"},
				Value: &Measurement{
					Value: numFloat("0.142"),
					Unit:  &String{Value: "kW"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:22.7.0"},
				Value: &Measurement{
					Value: numFloat("0.000"),
					Unit:  &String{Value: "kW"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:42.7.0"},
				Value: &Measurement{
					Value: numFloat("0.000"),
					Unit:  &String{Value: "kW"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "1-0:62.7.0"},
				Value: &Measurement{
					Value: numFloat("0.000"),
					Unit:  &String{Value: "kW"},
				},
			},
			&Object{
				OBIS: &OBIS{Value: "0-1:24.2.1"},
				Value: &LastCapture{
					Timestamp: &Timestamp{Value: "170102161005", DST: DSTWinter},
					Value: &Measurement{
						Value: numFloat("0.107"),
						Unit:  &String{Value: "m3"},
					},
				},
			},
		},
	}

	t.Run("Version", func(t *testing.T) {
		v := telegram.Version()
		assert.True(t, v != nil)
		assert.Equal(t, "50", v.Value)
	})

	t.Run("MeasuredAt", func(t *testing.T) {
		ts := telegram.MeasuredAt()
		assert.True(t, ts != nil)
		assert.Equal(t, "170102192002", ts.Value)
		assert.Equal(t, DSTWinter, ts.DST)
	})

	t.Run("EquipmentID", func(t *testing.T) {
		id := telegram.EquipmentID()
		assert.True(t, id != nil)
		assert.Equal(t, "4B384547303034303436333935353037", id.Value)
	})

	t.Run("ElectricityDelivered", func(t *testing.T) {
		m := telegram.ElectricityDelivered(1)
		assert.True(t, m != nil)
		assert.Equal(t, "4.426", m.Value.Value.String())
		assert.Equal(t, "kWh", m.Unit.Value)
	})

	t.Run("ElectricityReceived", func(t *testing.T) {
		m := telegram.ElectricityReceived(1)
		assert.True(t, m != nil)
		assert.Equal(t, "2.444", m.Value.Value.String())
		assert.Equal(t, "kWh", m.Unit.Value)
	})

	t.Run("ElectricityTariffIndicator", func(t *testing.T) {
		s := telegram.ElectricityTariffIndicator()
		assert.True(t, s != nil)
		assert.Equal(t, "0002", s.Value)
	})

	t.Run("ElectricityCurrentlyDelivered", func(t *testing.T) {
		m := telegram.ElectricityCurrentlyDelivered()
		assert.True(t, m != nil)
		assert.Equal(t, "0.244", m.Value.Value.String())
		assert.Equal(t, "kW", m.Unit.Value)
	})

	t.Run("PowerFailures", func(t *testing.T) {
		s := telegram.PowerFailures()
		assert.True(t, s != nil)
		assert.Equal(t, "00013", s.Value)
	})

	t.Run("CurrentL1", func(t *testing.T) {
		m := telegram.CurrentL1()
		assert.True(t, m != nil)
		assert.Equal(t, "0.48", m.Value.Value.String())
		assert.Equal(t, "A", m.Unit.Value)
	})

	t.Run("VoltageL1", func(t *testing.T) {
		m := telegram.VoltageL1()
		assert.True(t, m != nil)
		assert.Equal(t, "230", m.Value.Value.String())
		assert.Equal(t, "V", m.Unit.Value)
	})

	t.Run("PowerDeliveredL1", func(t *testing.T) {
		m := telegram.PowerDeliveredL1()
		assert.True(t, m != nil)
		assert.Equal(t, "0.07", m.Value.Value.String())
		assert.Equal(t, "kW", m.Unit.Value)
	})

	t.Run("GasDelivered", func(t *testing.T) {
		lc := telegram.GasDelivered()
		assert.True(t, lc != nil)
		assert.Equal(t, "170102161005", lc.Timestamp.Value)
		assert.Equal(t, "0.107", lc.Value.Value.Value.String())
		assert.Equal(t, "m3", lc.Value.Unit.Value)
	})

	t.Run("NonexistentAccessor", func(t *testing.T) {
		m := telegram.lookupObject("9-9:99.99.99")
		assert.True(t, m == nil)
	})
}

func TestMBusDevice(t *testing.T) {
	telegram := &Telegram{
		Data: Data{
			&MBusDevice{
				Channel: 1,
				Data: []*Object{
					{
						OBIS:  &OBIS{Value: "0-1:24.1.0"},
						Value: &String{Value: "003"},
					},
				},
			},
			&MBusDevice{
				Channel: 2,
				Data: []*Object{
					{
						OBIS:  &OBIS{Value: "0-2:24.1.0"},
						Value: &String{Value: "003"},
					},
				},
			},
		},
	}

	t.Run("MBusDeviceExists", func(t *testing.T) {
		mbus := telegram.MBusDevice(1)
		assert.True(t, mbus != nil)
		assert.Equal(t, 1, mbus.Channel)
	})

	t.Run("MBusDeviceNotExists", func(t *testing.T) {
		mbus := telegram.MBusDevice(3)
		assert.True(t, mbus == nil)
	})
}

func numFloat(s string) *Number {
	f := new(big.Float)
	f.SetString(s)
	return &Number{Value: f}
}
