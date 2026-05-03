#include "layrz_protocol/packets/ps.hpp"
#include "packets/helpers.hpp"

namespace layrz::protocol::packets {

using namespace detail;

std::string PsPacket::to_packet() const {
    std::string payload = std::to_string(static_cast<long long>(timestamp)) + ";"
                        + cast_extra(params) + ";";
    return wrap_packet("<Ps>", "</Ps>", payload);
}

Result<PsPacket> PsPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Ps>", "</Ps>");
    if (!r.ok()) return Result<PsPacket>::fail(r.error);
    if (r.value.size() != 2) return Result<PsPacket>::fail(Error::MalformedFrame);
    try {
        PsPacket p;
        p.timestamp = static_cast<std::time_t>(std::stoll(r.value[0]));
        p.params    = parse_extra(r.value[1]);
        return Result<PsPacket>::success(p);
    } catch (...) {
        return Result<PsPacket>::fail(Error::ParseError);
    }
}

} // namespace layrz::protocol::packets
