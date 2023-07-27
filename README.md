# DSMR

A package for parsing Dutch Smart Meter Requirements (DSMR) telegram data.

[![Latest Release](https://img.shields.io/github/release/mijnverbruik/dsmr.svg?style=flat-square)](https://github.com/mijnverbruik/dsmr/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/mijnverbruik/dsmr/test.yml?style=flat-square&branch=main)](https://github.com/mijnverbruik/dsmr/actions?query=workflow%3Atest)
[![MIT license](https://img.shields.io/github/license/mijnverbruik/dsmr.svg?style=flat-square)](https://github.com/mijnverbruik/dsmr/blob/main/LICENSE)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/mijnverbruik/dsmr)](https://pkg.go.dev/github.com/mijnverbruik/dsmr)

## Usage

```go
import "github.com/mijnverbruik/dsmr"

raw := strings.NewReader("" +
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

telegram, err := dsmr.ParseString(raw)
```

## Contributing

Everyone is encouraged to help improve this project. Here are a few ways you can help:

- [Report bugs](https://github.com/mijnverbruik/dsmr/issues)
- Fix bugs and [submit pull requests](https://github.com/mijnverbruik/dsmr/pulls)
- Write, clarify, or fix documentation
- Suggest or add new features

To get started with development:

```
git clone https://github.com/mijnverbruik/dsmr.git
cd dsmr
go test ./...
```

Feel free to open an issue to get feedback on your idea before spending too much time on it.

## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
