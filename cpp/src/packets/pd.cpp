#include "layrz_protocol/packets/pd.hpp"
#include "layrz_protocol/float_repr.hpp"
#include "packets/helpers.hpp"

namespace layrz::protocol::packets {

using namespace detail;

static std::string opt_float(const std::optional<double>& v) {
    return v.has_value() ? python_repr_float(*v) : "";
}
static std::string opt_int(const std::optional<int>& v) {
    return v.has_value() ? std::to_string(*v) : "";
}

std::string PdPacket::to_packet() const {
    const Position* pos = position.has_value() ? &*position : nullptr;
    std::string payload =
        std::to_string(static_cast<long long>(timestamp)) + ";"
        + (pos ? opt_float(pos->latitude)  : "") + ";"
        + (pos ? opt_float(pos->longitude) : "") + ";"
        + (pos ? opt_float(pos->altitude)  : "") + ";"
        + (pos ? opt_float(pos->speed)     : "") + ";"
        + (pos ? opt_float(pos->direction) : "") + ";"
        + (pos ? opt_int(pos->satellites)  : "") + ";"
        + (pos ? opt_float(pos->hdop)      : "") + ";"
        + cast_extra(extra_data) + ";";
    return wrap_packet("<Pd>", "</Pd>", payload);
}

Result<PdPacket> PdPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Pd>", "</Pd>");
    if (!r.ok()) return Result<PdPacket>::fail(r.error);
    if (r.value.size() != 9) return Result<PdPacket>::fail(Error::MalformedFrame);

    auto& p = r.value;
    auto opt_dbl = [](const std::string& s) -> std::optional<double> {
        if (s.empty()) return std::nullopt;
        try { return std::stod(s); } catch (...) { return std::nullopt; }
    };
    auto opt_int_f = [](const std::string& s) -> std::optional<int> {
        if (s.empty()) return std::nullopt;
        try { return std::stoi(s); } catch (...) { return std::nullopt; }
    };

    try {
        PdPacket pkt;
        pkt.timestamp = static_cast<std::time_t>(std::stoll(p[0]));

        Position pos;
        pos.latitude   = opt_dbl(p[1]);
        pos.longitude  = opt_dbl(p[2]);
        pos.altitude   = opt_dbl(p[3]);
        pos.speed      = opt_dbl(p[4]);
        pos.direction  = opt_dbl(p[5]);
        pos.satellites = opt_int_f(p[6]);
        pos.hdop       = opt_dbl(p[7]);

        bool any = pos.latitude || pos.longitude || pos.altitude ||
                   pos.speed    || pos.direction || pos.satellites || pos.hdop;
        if (any) pkt.position = pos;

        pkt.extra_data = parse_extra(p[8]);
        return Result<PdPacket>::success(pkt);
    } catch (...) {
        return Result<PdPacket>::fail(Error::ParseError);
    }
}

} // namespace layrz::protocol::packets
