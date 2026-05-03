# layrz-protocol — C++ (v3)

C++17 implementation of the [Layrz Link Protocol](https://developers.layrz.com/protocol/).

## What's included

| Target | Description |
|--------|-------------|
| `layrz::protocol::core` | Packet codec — works everywhere (desktop, Arduino, ESP-IDF) |
| `layrz::protocol::net`  | TCP + HTTP transport — desktop/POSIX only by default |

The transport layer (`layrz::protocol::net`) is desktop/POSIX-only. Embedded targets (PlatformIO, ESP-IDF, Arduino) receive the codec only and handle connectivity through their own platform mechanisms.

## Desktop CMake

### Build

```bash
cmake -S . -B build
cmake --build build
ctest --test-dir build
```

Options:

| Option | Default | Description |
|--------|---------|-------------|
| `LAYRZ_PROTOCOL_BUILD_NET` | `ON` | Build TCP/HTTP transport |
| `LAYRZ_PROTOCOL_BUILD_TESTS` | `ON` | Build GoogleTest suite |

### Consume with `find_package`

```cmake
find_package(LayrzProtocol REQUIRED)
target_link_libraries(my_target PRIVATE layrz::protocol::core)

# Or, to also pull in the transport layer:
find_package(LayrzProtocol REQUIRED COMPONENTS net)
target_link_libraries(my_target PRIVATE layrz::protocol::net)
```

`LAYRZ_PROTOCOL_CLIENTS` is automatically defined when linking against `layrz::protocol::net`.

### Install

```bash
cmake --install build --prefix /usr/local
```

## PlatformIO

Add to `platformio.ini`:

```ini
lib_deps =
    https://github.com/goldenm-software/layrz-protocol#v3.1.0

build_flags = -std=gnu++17
```

Only the packet codec is compiled. Connectivity is left to the application (ESP-IDF sockets, Arduino `WiFiClient`, etc.).

## Packet reference

See [protocol documentation](https://developers.layrz.com/protocol/) for the full packet list and wire format.

Server → Device: `AB`, `AR`, `AU`, `AS`, `AO`, `AC`  
Device → Server: `PB`, `PM`, `PD`, `PC`, `PI`, `PR`, `PS`, `PA`  
Special: `AI`, `Ts`, `Te`

## Quick start

```cpp
#include <layrz_protocol/layrz_protocol.hpp>
using namespace layrz::protocol;

// Encode
packets::PaPacket pa;
pa.ident    = "MY_DEVICE";
pa.password = "secret";
std::string frame = pa.to_packet(); // "<Pa>...</Pa>"

// Decode server frame
auto result = handle_server_output("<As><imei>MY_DEVICE</imei></As>");
if (result.ok()) {
    std::visit([](auto&& pkt) { /* handle */ }, result.value);
}
```
