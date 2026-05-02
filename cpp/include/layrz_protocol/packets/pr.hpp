#pragma once
#include "layrz_protocol/errors.hpp"
#include <string>

namespace layrz::protocol::packets {

// <Pr> — Device → Server — Sync/keepalive request (empty body)
struct PrPacket {
    std::string          to_packet() const;
    static Result<PrPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets
