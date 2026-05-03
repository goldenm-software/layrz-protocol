#pragma once
#ifndef __LAYRZ_PROTOCOL_PACKETS_PC_HPP__
#define __LAYRZ_PROTOCOL_PACKETS_PC_HPP__

#include "layrz_protocol/errors.hpp"
#include <ctime>
#include <string>

namespace layrz::protocol::packets {

// <Pc> — Device → Server — Command response
struct PcPacket {
    std::time_t timestamp;
    int         command_id;
    std::string message;

    std::string          to_packet() const;
    static Result<PcPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets

#endif // __LAYRZ_PROTOCOL_PACKETS_PC_HPP__
