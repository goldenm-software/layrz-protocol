#pragma once
#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/extras.hpp"
#include "layrz_protocol/packets/position.hpp"
#include <ctime>
#include <optional>
#include <string>

namespace layrz::protocol::packets {

// <Pd> — Device → Server — Position + extra data
struct PdPacket {
    std::time_t          timestamp;
    std::optional<Position> position;
    ExtrasMap            extra_data;

    std::string          to_packet() const;
    static Result<PdPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets
