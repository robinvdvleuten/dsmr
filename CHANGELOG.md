# Changelog

All notable changes to this project will be documented in this file.

## [0.6.0](https://github.com/robinvdvleuten/dsmr/compare/v0.5.0...v0.6.0) (2026-05-13)


### Features

* add telegram text marshaling ([#51](https://github.com/robinvdvleuten/dsmr/issues/51)) ([e9f628a](https://github.com/robinvdvleuten/dsmr/commit/e9f628a0c92608249bc458805ff792ab3076d68a))
* **api:** add convenience accessors for common DSMR fields ([a508f34](https://github.com/robinvdvleuten/dsmr/commit/a508f34fd5985a12369af878487343a43b1ac3f5))
* **deps:** Bump github.com/alecthomas/repr from 0.4.0 to 0.5.2 ([#37](https://github.com/robinvdvleuten/dsmr/issues/37)) ([4d7dbc9](https://github.com/robinvdvleuten/dsmr/commit/4d7dbc902d73558b975510ee73efe229d5e0f038))


### Performance Improvements

* implement lazy initialization of participle parsers with sync.Once ([f0e2253](https://github.com/robinvdvleuten/dsmr/commit/f0e225301c3c89a30fda843ecd828bb42f56a709))

## [0.5.0](https://github.com/robinvdvleuten/dsmr/compare/v0.4.0..v0.5.0) - 2025-09-22

### Features

- parseObject now caches the computed MBus channel on each Object - ([e6f1a73](https://github.com/robinvdvleuten/dsmr/commit/e6f1a731a096989f5ff9460e87b61ac36cb63d85))
- Optimised the MBus grouping logic - ([051b827](https://github.com/robinvdvleuten/dsmr/commit/051b8275f955cbbdc7b06100664862682e8853e3))
- Expose formatted errors via participle.Wrapf - ([e6fc9bb](https://github.com/robinvdvleuten/dsmr/commit/e6fc9bba9e53781cb36a4a165dd07ab62b93d774))
- Lean on Capture more broadly - ([0fc6a7f](https://github.com/robinvdvleuten/dsmr/commit/0fc6a7f077d5a9e4b8f8df91b3913734d5c5895d))
- Prefer typed enums for repetitive literals - ([d1de80c](https://github.com/robinvdvleuten/dsmr/commit/d1de80c2cc540b22d1df12fcaa09bf633e4e08cb))
- Expose richer error context - ([f108fb8](https://github.com/robinvdvleuten/dsmr/commit/f108fb8620159ddf104f6f2ed314045a700bba38))
- MBus devices are grouped per channel - ([5bc1b51](https://github.com/robinvdvleuten/dsmr/commit/5bc1b515fe94098768eb039935916da55ce7cccc))
- Do not accept OBIS identifier in value union - ([52756dd](https://github.com/robinvdvleuten/dsmr/commit/52756dd3abf6f79c707cd9418c5730160a52c0ea))
- Only allow single digits for medium and channel in OBIS identifer - ([5276322](https://github.com/robinvdvleuten/dsmr/commit/5276322c9a6ccbf86490df24db02eb7ae02c73d5))
- Replace list with more specific structs - ([448c09f](https://github.com/robinvdvleuten/dsmr/commit/448c09ff138240777f678a50d2ec030fa280ff21))
- Make verifying checksum optional - ([f55d971](https://github.com/robinvdvleuten/dsmr/commit/f55d971f673ed9a6f9840f566e9e63ededfeeb6a))
- Only allow telegrams as strings to be parsed - ([188bb84](https://github.com/robinvdvleuten/dsmr/commit/188bb844c251f4c479a7f28c910dff1204960def))

### Documentation

- Document the grammar. - ([27b4457](https://github.com/robinvdvleuten/dsmr/commit/27b4457caaf71bf68364374653994ec36e0d14a6))

### Feet

- Keep lexer rules symmetric with parser needs - ([38c68a1](https://github.com/robinvdvleuten/dsmr/commit/38c68a192bd57465a2d982306cefd95e1d3aaecc))

## [0.4.0](https://github.com/robinvdvleuten/dsmr/compare/v0.3.0..v0.4.0) - 2023-11-14

### Features

- Remove nested `Value` within `Header` and `Footer` - ([8f2872b](https://github.com/robinvdvleuten/dsmr/commit/8f2872b87c3cdbd49b156199dc1562b6c61df160))
- Extract DST indicator seperate from timestamp - ([28304b8](https://github.com/robinvdvleuten/dsmr/commit/28304b811e3c38328a3a294bac7ab887e1c79d25))
- Verify checksum after parsing - ([0e12763](https://github.com/robinvdvleuten/dsmr/commit/0e127635fc1d6dbdaffbc03677b4a6cf743bed50))

## [0.3.0](https://github.com/robinvdvleuten/dsmr/compare/v0.2.0..v0.3.0) - 2023-08-14

### Features

- Make value interface more explicit and private - ([98166d1](https://github.com/robinvdvleuten/dsmr/commit/98166d1bc3cb310eeff19f6be14eb791d20b345f))
- Parse event log into corresponding struct - ([5caaea6](https://github.com/robinvdvleuten/dsmr/commit/5caaea6b4d5130d390928dec8519fad311c6312a))
- Replace PEG based parser back to Participle - ([c3d1320](https://github.com/robinvdvleuten/dsmr/commit/c3d13208d208b2a32b458da66abd383915ca2f87))

### Bug Fixes

- Do not include `/` when capturing header value - ([99fd4c9](https://github.com/robinvdvleuten/dsmr/commit/99fd4c9e5fcf806ec3be0bfd7e764ad343e43b39))

### Documentation

- Fix example in README [skip ci] - ([7978619](https://github.com/robinvdvleuten/dsmr/commit/7978619931a04d4d8a1f8631c966e20c55af368c))

## [0.2.0](https://github.com/robinvdvleuten/dsmr/compare/v0.1.0..v0.2.0) - 2023-08-01

### Features

- Return error when decimal cannot be parsed - ([7f42222](https://github.com/robinvdvleuten/dsmr/commit/7f42222c6835b5bb875aef79b14c833d3179c398))
- Use decimal package to represent measurement values - ([06061ea](https://github.com/robinvdvleuten/dsmr/commit/06061eab031e1066051ba6453cc1e0eec67ca4b2))
- Only allow access to properties through getters - ([ffc4f3c](https://github.com/robinvdvleuten/dsmr/commit/ffc4f3ce4d080b04a6268d3c3dc317459b353dc1))
- Convert attribute values to their actual types - ([c8b7884](https://github.com/robinvdvleuten/dsmr/commit/c8b78842aaf58bb6e79f751d7a1cf91d3a51a58d))
- Allow options to be passed to parser - ([f0a12bc](https://github.com/robinvdvleuten/dsmr/commit/f0a12bcfa61d802ec5dfd7c38b6533aefafaaea4))
- Make COSEM a lookup map - ([c78208e](https://github.com/robinvdvleuten/dsmr/commit/c78208eca3c90573804ae17f4fc333e9e1bfdd23))
- Generate PEG based parser ([#3](https://github.com/orhun/git-cliff/issues/3)) - ([79c778d](https://github.com/robinvdvleuten/dsmr/commit/79c778dfea79d9d69834ef7b07b16aab1679cc85))

### Bug Fixes

- Correctly represent OBIS as OBIS attribute - ([2f5a39a](https://github.com/robinvdvleuten/dsmr/commit/2f5a39a859b708875b81d408b1e9d0ae40db0d19))
- Obsolete closing bracket in Footer rule - ([98da80a](https://github.com/robinvdvleuten/dsmr/commit/98da80a559c976dad29936c4b9484c82e1018e90))

## [0.1.0] - 2023-07-27

- Initial commit - ([2ed3459](https://github.com/robinvdvleuten/dsmr/commit/2ed3459a254ee803c4a3c2c7726bb8f26715c93a))
