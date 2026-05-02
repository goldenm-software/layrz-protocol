#include "layrz_protocol/packets/command.hpp"
#include "layrz_protocol/crc.hpp"
#include "layrz_protocol/extras.hpp"
#include "packets/helpers.hpp"

namespace layrz::protocol::packets {

using namespace detail;

std::string CommandDefinition::to_packet() const {
    std::string payload = std::to_string(command_id) + ";"
                        + command_name + ";"
                        + cast_extra(args) + ";";
    std::string inner_crc = compute_crc_str(payload);
    return payload + inner_crc;
}

Result<CommandDefinition> CommandDefinition::from_packet(std::string_view raw) {
    // raw = "id;name;args;INNERCRC"  (4 parts when split on ';')
    auto parts = split(raw, ';');
    // split("id;name;args;CRC") gives ["id","name","args","CRC"] — exactly 4
    if (parts.size() != 4)
        return Result<CommandDefinition>::fail(Error::MalformedFrame);

    // Validate inner CRC (scope = "id;name;args;")
    std::string scope = parts[0] + ";" + parts[1] + ";" + parts[2] + ";";
    uint16_t expected;
    try {
        expected = static_cast<uint16_t>(std::stoul(parts[3], nullptr, 16));
    } catch (...) {
        return Result<CommandDefinition>::fail(Error::CrcMismatch);
    }
    if (crc16_x25(scope) != expected)
        return Result<CommandDefinition>::fail(Error::CrcMismatch);

    try {
        CommandDefinition c;
        c.command_id   = std::stoi(parts[0]);
        c.command_name = parts[1];
        c.args         = parse_extra(parts[2]);
        return Result<CommandDefinition>::success(c);
    } catch (...) {
        return Result<CommandDefinition>::fail(Error::ParseError);
    }
}

} // namespace layrz::protocol::packets
