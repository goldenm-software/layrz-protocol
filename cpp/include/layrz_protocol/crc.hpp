#pragma once
#include <cstdint>
#include <string>
#include <string_view>

namespace layrz::protocol {

// CRC-16/X-25: init=0xFFFF, poly=0x1021 (reflected 0x8408), refin=true,
// refout=true, xorout=0xFFFF. Used by every Layrz Protocol packet.
uint16_t crc16_x25(std::string_view data);

// Render a CRC value as 4 uppercase hex digits (e.g. 0x7F28 → "7F28").
std::string crc_hex(uint16_t crc);

// Compute CRC over payload (including trailing ';') and return the 4-char hex string.
inline std::string compute_crc_str(std::string_view payload) {
    return crc_hex(crc16_x25(payload));
}

} // namespace layrz::protocol
