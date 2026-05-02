#pragma once
#include "layrz_protocol/errors.hpp"
#include <string>

namespace layrz::protocol::packets {

// <Ar> — Server → Device — Error/reject with reason string
struct ArPacket {
    std::string reason;

    std::string          to_packet() const;
    static Result<ArPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets
