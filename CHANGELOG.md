# Changelog

## 3.3.1

- Fixed extra-arg colon escaping to also cover the **key**, not just the value. 3.3.0 escaped only values, but a Zigbee entry carries its colons in the key (`{mac}.{expose}`, e.g. `a4:c1:38:5c:02:f6:b4:53.energy`); `:` is now escaped as `___` on both key and value when serializing and reversed on parse, so such keys round-trip intact. Applied across Python, Dart and Go.

## 3.3.0

- Fixed extra-arg parsing so values containing a colon (e.g. a MAC-like identifier `a4:c1:38:5c:02:f6:b4:53`) round-trip intact: string values now escape `:` as `___` on serialize and the parser reverses it, keeping the `key:value` split unambiguous. Applied across Python, Dart and Go.

## 3.2.0

- Added C++ `TcpServer` and `HttpServer` (Tier-3 `layrz_protocol_servers` CMake target) mirroring Go's server API; gated by `LAYRZ_PROTOCOL_SERVERS` macro
- Added Go packet subpackages (`packets/client`, `packets/server`, `packets/ai`, `packets/trips`, `packets/helpers`) with decoder, union, and helper types
- Added Python asyncio `TcpServer` and `HttpServer` with `OnNewPacket`, `OnDecodeError`, `OnAuthenticate`, and `OnPullCommands` callbacks
- Added Dart pure-Dart `TcpServer` and `HttpServer` under `lib/servers/`; removed Flutter/drift database dependency
- Refactored Dart package: removed `lib/src/` nesting, flattened all packet paths to `lib/packets/`, `lib/clients/`, `lib/servers/`, `lib/utils/`
- Added `split_client_frames` and `handle_client_input` to C++ core parser for server-side frame decoding
- Added server examples for Go, Python, Dart, and C++ under each language's `examples/servers/` directory
- Added global `make coverage` target aggregating per-language coverage reports with 80% threshold enforcement

## 3.1.2

- Fixed Dart TCP client `<AuPacket>` handling using `return` instead of `continue`, which caused `<AsPacket>` to be missed when both arrived in the same TCP chunk
- Added Dart TCP client IPv6 preference: hostname resolution now picks an AAAA record when available, falling back to A

## 3.1.1

- Fixed Go `BleAdvertisement.Latitude`, `Longitude`, and `Altitude` fields changed from `float64` to `*float64` so devices without GPS fix correctly report absent coordinates as `nil` instead of `0.0`
- Fixed Go `<Pb>` packet parser to return an error (rather than silently zero-ing) when coordinate fields are present but non-numeric
- Fixed Go `<Pb>` parser to normalize MAC addresses to uppercase colon-separated format (`AA:BB:CC:DD:EE:FF`) after CRC validation

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
