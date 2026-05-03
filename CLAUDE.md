# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

Layrz Protocol is a **multi-language monorepo** implementing the Layrz Link Protocol. Each language lives in its own subdirectory with an independent package/module structure.

| Directory | Language | Package |
|-----------|----------|---------|
| `python/` | Python 3.10+ | `layrz_protocol` (PyPI) |
| `dart/` | Dart/Flutter | `layrz_protocol` (pub.dev) |
| `go/` | Go 1.26 | `github.com/goldenm-software/layrz-protocol/go/v3` |
| `cpp/` | C++ | Pending |

Docs: https://developers.layrz.com/protocol/

## Commands

### Python (`cd python/`)
```bash
uv sync --only-group dev        # Install dev dependencies
uv run ruff check               # Lint
uv run mypy .                   # Type check
uv run coverage run -m unittest discover -s tests -v  # Run all tests
uv run coverage run -m unittest tests.test_foo -v     # Run single test module
uv run coverage report -m       # Coverage report
```

### Dart (`cd dart/`)
```bash
flutter analyze                 # Static analysis
flutter test --reporter expanded --coverage  # Run all tests
flutter test test/ab_test.dart  # Run single test file
dart run build_runner build     # Regenerate drift DB code
```

### Go (`cd golang/v3/` or `cd golang/v2/`)
```bash
go vet ./...
go build ./...
go test ./...
```

## Architecture

### Protocol Packets

The protocol defines two directions of packets:

- **Server → Client** (device receives): `AB`, `AR`, `AU`, `AS`, `AO`, `AC`
- **Client → Server** (device sends): `PB`, `PM`, `PD`, `PC`, `PI`, `PR`, `PS`, `PA`
- **Special**: `AI` (AI packets), `Ts`/`Te` (trip start/end)

Each language implements the same packet types. In Dart they live under `dart/lib/src/packets/src/{server,client,ai,trips}/`. In Python under `python/layrz_protocol/packets/`. In Go, each packet is a top-level `.go` file (`pb.go`, `ab.go`, etc.).

### Dart Package

- **Entry point**: `dart/lib/layrz_protocol.dart` — exports all packets, CRC util, errors; `part`s the two client implementations.
- **Clients**: `src/clients/http.dart` and `src/clients/tcp.dart` — implement `LayrzProtocolMode` (HTTP polling vs. TCP persistent connection).
- **Database**: `src/database/` — drift ORM, requires `build_runner` to regenerate generated files after schema changes.
- **Utils**: `src/utils/crc.dart` (CRC calculation), `src/utils/errors.dart`.

### Python Package

- **`client.py`**: Main protocol client.
- **`constants.py`**: Shared protocol constants.
- **`packets/`**: One module per packet type.
- **`utils/`**: CRC and helpers.
- Uses `pydantic` v2 for data models, `requests` for HTTP transport.
- Ruff config: 2-space indent, 120-char line length, single quotes.
- Mypy: strict mode.

### Go Packages

- v2 and v3 share the same file layout. Differences are in the protocol version they implement.
- `parser.go`: top-level packet parser.
- `http.go` / `tcp.go`: transport implementations.
- `args.go`: shared argument/option types.
- `crc.go`: CRC implementation.
- Uses `github.com/iancoleman/orderedmap` for ordered JSON maps.

## Versioning & Publishing

Version is bumped via `deploy.py` in each language directory before publishing. Publishing targets:
- Python → PyPI via `twine`
- Dart → pub.dev via `flutter pub publish`
- Go → GitHub tags (module proxy resolves automatically)
