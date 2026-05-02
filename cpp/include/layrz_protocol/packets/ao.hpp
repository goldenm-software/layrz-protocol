#pragma once
#include "layrz_protocol/errors.hpp"
#include <ctime>
#include <string>

namespace layrz::protocol::packets {

// <Ao> — Server → Device — ACK with timestamp
struct AoPacket {
    std::time_t timestamp;

    std::string          to_packet() const;
    static Result<AoPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets
