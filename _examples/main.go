package main

import (
	"fmt"
	"log"

	"github.com/robinvdvleuten/dsmr"
)

func main() {
	raw := "" +
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
		"!6796\r\n"

	telegram, err := dsmr.Parse(raw)
	if err != nil {
		log.Fatal(err)
	}

	// Header and timestamp
	fmt.Printf("Header: %s\n", telegram.Header.Value)
	if ts := telegram.MeasuredAt(); ts != nil {
		fmt.Printf("Measured at: %s (DST: %v)\n", ts.Value, ts.IsDST())
	}

	// Electricity consumption
	fmt.Println("\nElectricity Consumption:")
	if d1 := telegram.ElectricityDelivered(1); d1 != nil {
		fmt.Printf("  Tariff 1: %v %s\n", d1.Value.Value, d1.Unit.Value)
	}
	if d2 := telegram.ElectricityDelivered(2); d2 != nil {
		fmt.Printf("  Tariff 2: %v %s\n", d2.Value.Value, d2.Unit.Value)
	}

	// Current consumption
	if current := telegram.ElectricityCurrentlyDelivered(); current != nil {
		fmt.Printf("  Current: %v %s\n", current.Value.Value, current.Unit.Value)
	}

	// 3-phase measurements
	fmt.Println("\n3-Phase Measurements:")
	phases := []struct {
		name    string
		voltage func(*dsmr.Telegram) *dsmr.Measurement
		current func(*dsmr.Telegram) *dsmr.Measurement
		power   func(*dsmr.Telegram) *dsmr.Measurement
	}{
		{
			name: "L1",
			voltage: (*dsmr.Telegram).VoltageL1,
			current: (*dsmr.Telegram).CurrentL1,
			power:   (*dsmr.Telegram).PowerDeliveredL1,
		},
		{
			name: "L2",
			voltage: (*dsmr.Telegram).VoltageL2,
			current: (*dsmr.Telegram).CurrentL2,
			power:   (*dsmr.Telegram).PowerDeliveredL2,
		},
		{
			name: "L3",
			voltage: (*dsmr.Telegram).VoltageL3,
			current: (*dsmr.Telegram).CurrentL3,
			power:   (*dsmr.Telegram).PowerDeliveredL3,
		},
	}

	for _, p := range phases {
		fmt.Printf("  %s:\n", p.name)
		if v := p.voltage(telegram); v != nil {
			fmt.Printf("    Voltage: %v %s\n", v.Value.Value, v.Unit.Value)
		}
		if c := p.current(telegram); c != nil {
			fmt.Printf("    Current: %v %s\n", c.Value.Value, c.Unit.Value)
		}
		if pw := p.power(telegram); pw != nil {
			fmt.Printf("    Power: %v %s\n", pw.Value.Value, pw.Unit.Value)
		}
	}

	// Gas reading
	if gas := telegram.GasDelivered(); gas != nil {
		fmt.Printf("\nGas Delivered: %v %s (measured at %s)\n",
			gas.Value.Value.Value, gas.Value.Unit.Value, gas.Timestamp.Value)
	}

	// Power failures
	fmt.Println("\nPower Quality:")
	if failures := telegram.PowerFailures(); failures != nil {
		fmt.Printf("  Power failures: %s\n", failures.Value)
	}
	if failuresLong := telegram.PowerFailuresLong(); failuresLong != nil {
		fmt.Printf("  Long power failures: %s\n", failuresLong.Value)
	}

	// M-Bus devices
	for i := 1; i <= 4; i++ {
		if mbus := telegram.MBusDevice(i); mbus != nil {
			fmt.Printf("\nM-Bus Device %d: %d objects\n", i, len(mbus.Data))
		}
	}
}
