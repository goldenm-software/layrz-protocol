#pragma once
#include "layrz_protocol/errors.hpp"
#include <cstdint>
#include <string>
#include <vector>

namespace layrz::protocol::packets {

// <Pm> — Device → Server — Media upload
struct PmPacket {
    std::string          filename;
    std::string          content_type;
    std::vector<uint8_t> data;       // raw bytes; serialized as base64

    std::string          to_packet() const;
    static Result<PmPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets
