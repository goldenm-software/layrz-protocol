#pragma once
#ifndef __LAYRZ_PROTOCOL_PACKETS_IM_HPP__
#define __LAYRZ_PROTOCOL_PACKETS_IM_HPP__

#include "layrz_protocol/errors.hpp"
#include <ctime>
#include <string>

namespace layrz::protocol::packets {

// <Im> — Service-to-service — AI chat message
// Wire escape: ';' in message body is replaced with "|||" on encode and reversed on decode.
struct ImPacket {
    std::time_t timestamp;
    std::string chat_id;  // UUID string
    std::string message;  // stored unescaped (';' is literal)

    std::string          to_packet() const;
    static Result<ImPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets

#endif // __LAYRZ_PROTOCOL_PACKETS_IM_HPP__
