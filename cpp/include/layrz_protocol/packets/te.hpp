#pragma once
#include "layrz_protocol/errors.hpp"
#include <ctime>
#include <string>

namespace layrz::protocol::packets {

// <Te> — Service-to-service — Trip end
struct TePacket {
    std::time_t timestamp;
    std::string trip_id;
    double      distance_traveled; // metres, %.3f on wire
    double      max_speed;         // units as sent, %.3f on wire
    int         duration;          // seconds

    std::string          to_packet() const;
    static Result<TePacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets
