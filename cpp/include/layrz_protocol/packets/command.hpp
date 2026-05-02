#pragma once
#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/extras.hpp"
#include <string>

namespace layrz::protocol::packets {

struct CommandDefinition {
    int         command_id;
    std::string command_name;
    ExtrasMap   args;  // may be empty

    // Encode to "id;name;args;INNERCRС" (inner CRC over "id;name;args;")
    std::string to_packet() const;

    // Parse from "id;name;args;INNERCRC" fragment (4 ';'-separated parts)
    static Result<CommandDefinition> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets
