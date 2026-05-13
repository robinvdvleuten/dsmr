# DSMR

A package for parsing Dutch Smart Meter Requirements (DSMR) telegram data.

[![Latest Release](https://img.shields.io/github/release/robinvdvleuten/dsmr.svg?style=flat-square)](https://github.com/robinvdvleuten/dsmr/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/robinvdvleuten/dsmr/test.yml?style=flat-square&branch=main)](https://github.com/robinvdvleuten/dsmr/actions?query=workflow%3Atest)
[![MIT license](https://img.shields.io/github/license/robinvdvleuten/dsmr.svg?style=flat-square)](https://github.com/robinvdvleuten/dsmr/blob/main/LICENSE)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/robinvdvleuten/dsmr)](https://pkg.go.dev/github.com/robinvdvleuten/dsmr)

The package focuses on turning raw telegram strings into strongly typed Go
structures, making it easier to work with smart meter measurements such as
energy consumption, production, gas readings, and meter metadata. It supports
common DSMR versions out of the box and hides the boilerplate of dealing with
checksums, optional fields, and multi-line records.

## Usage

```go
import "github.com/robinvdvleuten/dsmr"

raw := "" +
    "/ISk5\\2MT382-1000\r\n" +
    "\r\n" +
    "1-3:0.2.8(50)\r\n" +
    "0-0:1.0.0(170102192002W)\r\n" +
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
    "0-1:24.2.1(170102161005W)(00000.107*m3)\r\n" +
    "0-2:24.1.0(003)\r\n" +
    "0-2:96.1.0()\r\n" +
    "!6EEE\r\n"

telegram, err := dsmr.Parse(raw)
if err != nil {
    // Handle checksum mismatches or invalid telegrams.
    log.Fatal(err)
}

// Use convenience accessors for common readings
if delivery := telegram.ElectricityDelivered(1); delivery != nil {
    fmt.Printf("Electricity delivered (tariff 1): %s %s\n", 
        delivery.Value.Value, delivery.Unit.Value)
}

if gas := telegram.GasDelivered(); gas != nil {
    fmt.Printf("Gas delivered: %s %s\n", 
        gas.Value.Value.Value, gas.Value.Unit.Value)
}

// Access 3-phase measurements
if voltageL1 := telegram.VoltageL1(); voltageL1 != nil {
    fmt.Printf("Voltage L1: %s %s\n", 
        voltageL1.Value.Value, voltageL1.Unit.Value)
}

if currentL1 := telegram.CurrentL1(); currentL1 != nil {
    fmt.Printf("Current L1: %s %s\n", 
        currentL1.Value.Value, currentL1.Unit.Value)
}

// Iterate over M-Bus devices
if mbus := telegram.MBusDevice(1); mbus != nil {
    fmt.Printf("M-Bus channel %d has %d objects\n", mbus.Channel, len(mbus.Data))
}
```

`dsmr.Parse` accepts raw telegram data as a string. See the [`_examples`](./_examples)
directory for additional usage patterns.

## Serialization

Telegrams can also be constructed and serialized back to DSMR text. The checksum
is recalculated when marshaling.

```go
amount, err := dsmr.NewNumber("000004.426")
if err != nil {
    log.Fatal(err)
}

telegram := &dsmr.Telegram{
    Header: &dsmr.Header{Value: "ISk5\\2MT382-1000"},
    Data: dsmr.Data{
        &dsmr.Object{
            OBIS:  &dsmr.OBIS{Value: "1-0:1.8.1"},
            Value: &dsmr.Measurement{Value: amount, Unit: &dsmr.String{Value: "kWh"}},
        },
    },
}

raw, err := telegram.MarshalText()
if err != nil {
    log.Fatal(err)
}
```

`telegram.String()` is also available when the telegram is already known to be
valid.

## Contributing

Everyone is encouraged to help improve this project. Here are a few ways you can help:

- [Report bugs](https://github.com/robinvdvleuten/dsmr/issues)
- Fix bugs and [submit pull requests](https://github.com/robinvdvleuten/dsmr/pulls)
- Write, clarify, or fix documentation
- Suggest or add new features

To get started with development:

```
git clone https://github.com/robinvdvleuten/dsmr.git
cd dsmr
go test ./...
```

Before submitting a pull request, please make sure to run
`go fmt` on any Go source files you touched so the code stays consistent.

Feel free to open an issue to get feedback on your idea before spending too much time on it.

## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
