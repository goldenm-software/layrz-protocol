#include "layrz_protocol/packets/ac.hpp"
#include "packets/helpers.hpp"

namespace layrz::protocol::packets {

using namespace detail;

std::string AcPacket::to_packet() const {
    std::string payload;
    for (auto& cmd : commands) {
        payload += cmd.to_packet() + ";";
    }
    return wrap_packet("<Ac>", "</Ac>", payload);
}

Result<AcPacket> AcPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Ac>", "</Ac>");
    if (!r.ok()) return Result<AcPacket>::fail(r.error);

    auto& parts = r.value;
    // Remove trailing empty parts
    while (!parts.empty() && parts.back().empty()) parts.pop_back();

    // Each command occupies 4 consecutive parts: id;name;args;INNERCRC
    if (parts.size() % 4 != 0)
        return Result<AcPacket>::fail(Error::MalformedFrame);

    AcPacket p;
    for (size_t i = 0; i < parts.size(); i += 4) {
        std::string cmd_raw = parts[i] + ";" + parts[i+1] + ";" + parts[i+2] + ";" + parts[i+3];
        auto c = CommandDefinition::from_packet(cmd_raw);
        if (!c.ok()) return Result<AcPacket>::fail(c.error);
        p.commands.push_back(std::move(c.value));
    }
    return Result<AcPacket>::success(p);
}

} // namespace layrz::protocol::packets
