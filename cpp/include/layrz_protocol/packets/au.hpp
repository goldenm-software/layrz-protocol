#pragma once
#ifndef __LAYRZ_PROTOCOL_PACKETS_AU_HPP__
#define __LAYRZ_PROTOCOL_PACKETS_AU_HPP__

#include "layrz_protocol/errors.hpp"
#include <string>

namespace layrz::protocol::packets {

// <Au> — Server → Device — Auth required (deprecated, same shape as As)
struct AuPacket {
    std::string          to_packet() const;
    static Result<AuPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets

#endif // __LAYRZ_PROTOCOL_PACKETS_AU_HPP__
