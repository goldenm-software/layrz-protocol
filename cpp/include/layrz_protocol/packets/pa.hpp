#pragma once
#include "layrz_protocol/errors.hpp"
#include <string>

namespace layrz::protocol::packets {

// <Pa> — Device → Server — Authentication
struct PaPacket {
    std::string ident;
    std::string password;

    std::string          to_packet() const;
    static Result<PaPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets
