#include "layrz_protocol/parser.hpp"
#include <string_view>

namespace layrz::protocol {

static bool starts_with(std::string_view s, const char* prefix) {
    std::string_view p(prefix);
    return s.size() >= p.size() && s.substr(0, p.size()) == p;
}

Result<AnyServerPacket> handle_server_output(std::string_view raw) {
    if (starts_with(raw, "<Ab>")) {
        auto r = packets::AbPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Ac>")) {
        auto r = packets::AcPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Ao>")) {
        auto r = packets::AoPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Ar>")) {
        auto r = packets::ArPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<As>")) {
        auto r = packets::AsPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Au>")) {
        auto r = packets::AuPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Ts>")) {
        auto r = packets::TsPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Te>")) {
        auto r = packets::TePacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Im>")) {
        auto r = packets::ImPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    return Result<AnyServerPacket>::fail(Error::MalformedFrame);
}

Result<std::string> parse_packet_to_string(const AnyClientPacket& packet) {
    return std::visit([](auto&& p) -> Result<std::string> {
        return Result<std::string>::success(p.to_packet());
    }, packet);
}

} // namespace layrz::protocol
