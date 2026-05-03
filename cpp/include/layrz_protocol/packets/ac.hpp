#pragma once
#ifndef __LAYRZ_PROTOCOL_PACKETS_AC_HPP__
#define __LAYRZ_PROTOCOL_PACKETS_AC_HPP__

#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/packets/command.hpp"
#include <string>
#include <vector>

namespace layrz::protocol::packets {

// <Ac> — Server → Device — Command queue push
struct AcPacket {
    std::vector<CommandDefinition> commands;

    std::string          to_packet() const;
    static Result<AcPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets

#endif // __LAYRZ_PROTOCOL_PACKETS_AC_HPP__
