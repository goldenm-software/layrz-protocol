#pragma once
#ifndef __LAYRZ_PROTOCOL_PACKETS_AS_HPP__
#define __LAYRZ_PROTOCOL_PACKETS_AS_HPP__

#include "layrz_protocol/errors.hpp"
#include <string>

namespace layrz::protocol::packets {

// <As> — Server → Device — Auth success (empty body)
struct AsPacket {
    std::string          to_packet() const;
    static Result<AsPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets

#endif // __LAYRZ_PROTOCOL_PACKETS_AS_HPP__
