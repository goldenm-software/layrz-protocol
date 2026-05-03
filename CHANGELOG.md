# Changelog

## [Unreleased]

### Breaking changes (Go)
- `CommandDefinition.Args`, `PdPacket.ExtraData`, and `PsPacket.Params` changed from `*map[string]any` to `map[string]any`. Callers must drop `&` on assignment and `*` on read. Ranging over a nil map is safe (zero iterations).

## 3.1.0

- Added full C++ implementation (`cpp/`) with CMake C++17 build, covering all packet types (`Pa`, `Pb`, `Pc`, `Pd`, `Pi`, `Pm`, `Pr`, `Ps`, `Ab`, `Ac`, `Ao`, `Ar`, `As`, `Au`, `Ts`, `Te`, `Im`), CRC-16/X-25, BLE advertisement codec, extras parser, and HTTP/TCP transports
- Added Go v3 implementation with full protocol parity against the Python source of truth
- Fixed Dart parity issues against Python wire format (field order, CRC scope, `<Te>` and `<Ts>` packets)
- Migrated Python type checker from `mypy` to `ty`; bumped `requires-python` to `>=3.13`

## 3.0.7

- Changes on Apache-2.0 license format (Keeping only the appendix with copyright holder)

## 3.0.6

- Crated `<Im>` packet for AI streaming messages

## 3.0.5

- Fixes on Dart package for `<Te>` packet parsing (CRC and field order)

## 3.0.4

- Updated `<Te>` packet to include additional metadata fields

## 3.0.3

- Renamed `<Pt>` to `<Ts>` and `<Pe>` to `<Te>` for Trip Start and Trip End packets

## 3.0.2

- Added `<Pt>` and `<Pe>` packets

## 3.0.1

- Fixes on Golang AuPacket deprecation warning

## 3.0.0

- Major Protocol update

## 2.7.7

- Sync of Golang

## 2.7.6

- Changes on documentation

## 2.7.5

- Sync of Dart

## 2.7.4

- Initial public release
