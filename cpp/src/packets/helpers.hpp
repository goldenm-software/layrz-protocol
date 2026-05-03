#pragma once
// Internal helpers shared across packet implementations — not part of public API.
#include "layrz_protocol/crc.hpp"
#include "layrz_protocol/errors.hpp"
#include <string>
#include <string_view>
#include <vector>

namespace layrz::protocol::detail {

// Split sv on delimiter. Returns all parts including empty ones.
inline std::vector<std::string> split(std::string_view sv, char delim) {
    std::vector<std::string> parts;
    size_t start = 0;
    while (true) {
        size_t pos = sv.find(delim, start);
        if (pos == std::string_view::npos) {
            parts.emplace_back(sv.substr(start));
            break;
        }
        parts.emplace_back(sv.substr(start, pos - start));
        start = pos + 1;
    }
    return parts;
}

// Strip the <Xx> and </Xx> wrapper and validate/strip the CRC suffix.
// Returns the inner field vector (without the CRC element) or error.
// expected_tag: the two-char tag name (e.g. "Ao").
inline Result<std::vector<std::string>>
unwrap_packet(std::string_view raw, const char* open_tag, const char* close_tag) {
    if (raw.size() < 9) return Result<std::vector<std::string>>::fail(Error::MalformedFrame);

    size_t open_len  = std::string_view(open_tag).size();
    size_t close_len = std::string_view(close_tag).size();

    if (raw.substr(0, open_len) != open_tag ||
        raw.substr(raw.size() - close_len) != close_tag) {
        return Result<std::vector<std::string>>::fail(Error::MalformedFrame);
    }

    std::string_view inner = raw.substr(open_len, raw.size() - open_len - close_len);
    auto parts = split(inner, ';');
    if (parts.size() < 2) return Result<std::vector<std::string>>::fail(Error::MalformedFrame);

    // Last part is CRC hex; payload = everything before it + trailing ';'
    std::string received_crc_str = parts.back();
    parts.pop_back();

    // Rebuild payload including trailing ';'
    std::string payload;
    for (auto& p : parts) { payload += p; payload += ';'; }

    uint16_t received_crc = 0;
    try {
        received_crc = static_cast<uint16_t>(std::stoul(received_crc_str, nullptr, 16));
    } catch (...) {
        return Result<std::vector<std::string>>::fail(Error::CrcMismatch);
    }
    uint16_t calculated_crc = crc16_x25(payload);
    if (received_crc != calculated_crc) {
        return Result<std::vector<std::string>>::fail(Error::CrcMismatch);
    }

    return Result<std::vector<std::string>>::success(std::move(parts));
}

// Build a complete packet frame.
// payload must already include the trailing ';' after the last field.
inline std::string wrap_packet(const char* open_tag, const char* close_tag,
                               const std::string& payload) {
    return std::string(open_tag) + payload + compute_crc_str(payload) + close_tag;
}

// Format a MAC: insert ':' every 2 chars → "1234567890AB" → "12:34:56:78:90:AB"
inline std::string mac_with_colons(std::string_view hex12) {
    if (hex12.size() != 12) return std::string(hex12);
    std::string out;
    out.reserve(17);
    for (size_t i = 0; i < 12; i += 2) {
        if (i) out += ':';
        out += hex12[i];
        out += hex12[i+1];
    }
    return out;
}

// Format a MAC for the wire: strip colons and uppercase → "12:34:56:78:90:AB" → "1234567890AB"
inline std::string mac_without_colons(const std::string& mac) {
    std::string out;
    out.reserve(12);
    for (char c : mac) {
        if (c != ':') out += static_cast<char>(std::toupper(static_cast<unsigned char>(c)));
    }
    return out;
}

} // namespace layrz::protocol::detail
