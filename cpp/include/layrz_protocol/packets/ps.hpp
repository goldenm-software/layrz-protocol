#pragma once
#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/extras.hpp"
#include <ctime>
#include <string>

namespace layrz::protocol::packets {

// <Ps> — Device → Server — Settings/status report
struct PsPacket {
    std::time_t timestamp;
    ExtrasMap   params;

    std::string          to_packet() const;
    static Result<PsPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets
