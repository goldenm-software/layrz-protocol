#pragma once
#include "layrz_protocol/errors.hpp"
#include <ctime>
#include <string>

namespace layrz::protocol::packets {

// <Ts> — Service-to-service — Trip start
struct TsPacket {
    std::time_t timestamp;
    std::string trip_id;   // UUID string

    std::string          to_packet() const;
    static Result<TsPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets
