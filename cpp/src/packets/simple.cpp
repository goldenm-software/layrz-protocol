// Implementations for the simple (empty-body or single-field) packets:
// As, Au, Pr  (empty body → CRC over ";")
// Ao          (timestamp)
// Ar          (reason string)
// Pa          (ident + password)
// Pc          (timestamp + command_id + message)
// Ts          (timestamp + trip_id)
// Te          (timestamp + trip_id + distance + max_speed + duration)
// Im          (timestamp + chat_id + message, with ';' ↔ "|||" escape)

#include "packets/helpers.hpp"
#include "layrz_protocol/packets/as.hpp"
#include "layrz_protocol/packets/au.hpp"
#include "layrz_protocol/packets/pr.hpp"
#include "layrz_protocol/packets/ao.hpp"
#include "layrz_protocol/packets/ar.hpp"
#include "layrz_protocol/packets/pa.hpp"
#include "layrz_protocol/packets/pc.hpp"
#include "layrz_protocol/packets/ts.hpp"
#include "layrz_protocol/packets/te.hpp"
#include "layrz_protocol/packets/im.hpp"
#include <cstdio>
#include <stdexcept>

namespace layrz::protocol::packets {

using namespace detail;

// ─── As ──────────────────────────────────────────────────────────────────────
std::string AsPacket::to_packet() const {
    return wrap_packet("<As>", "</As>", ";");
}
Result<AsPacket> AsPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<As>", "</As>");
    if (!r.ok()) return Result<AsPacket>::fail(r.error);
    if (r.value.size() != 1 || !r.value[0].empty())
        return Result<AsPacket>::fail(Error::MalformedFrame);
    return Result<AsPacket>::success({});
}

// ─── Au ──────────────────────────────────────────────────────────────────────
std::string AuPacket::to_packet() const {
    return wrap_packet("<Au>", "</Au>", ";");
}
Result<AuPacket> AuPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Au>", "</Au>");
    if (!r.ok()) return Result<AuPacket>::fail(r.error);
    if (r.value.size() != 1 || !r.value[0].empty())
        return Result<AuPacket>::fail(Error::MalformedFrame);
    return Result<AuPacket>::success({});
}

// ─── Pr ──────────────────────────────────────────────────────────────────────
std::string PrPacket::to_packet() const {
    return wrap_packet("<Pr>", "</Pr>", ";");
}
Result<PrPacket> PrPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Pr>", "</Pr>");
    if (!r.ok()) return Result<PrPacket>::fail(r.error);
    if (r.value.size() != 1 || !r.value[0].empty())
        return Result<PrPacket>::fail(Error::MalformedFrame);
    return Result<PrPacket>::success({});
}

// ─── Ao ──────────────────────────────────────────────────────────────────────
std::string AoPacket::to_packet() const {
    std::string payload = std::to_string(static_cast<long long>(timestamp)) + ";";
    return wrap_packet("<Ao>", "</Ao>", payload);
}
Result<AoPacket> AoPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Ao>", "</Ao>");
    if (!r.ok()) return Result<AoPacket>::fail(r.error);
    if (r.value.size() != 1) return Result<AoPacket>::fail(Error::MalformedFrame);
    try {
        AoPacket p;
        p.timestamp = static_cast<std::time_t>(std::stoll(r.value[0]));
        return Result<AoPacket>::success(p);
    } catch (...) {
        return Result<AoPacket>::fail(Error::ParseError);
    }
}

// ─── Ar ──────────────────────────────────────────────────────────────────────
std::string ArPacket::to_packet() const {
    return wrap_packet("<Ar>", "</Ar>", reason + ";");
}
Result<ArPacket> ArPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Ar>", "</Ar>");
    if (!r.ok()) return Result<ArPacket>::fail(r.error);
    if (r.value.size() != 1) return Result<ArPacket>::fail(Error::MalformedFrame);
    ArPacket p;
    p.reason = r.value[0];
    return Result<ArPacket>::success(p);
}

// ─── Pa ──────────────────────────────────────────────────────────────────────
std::string PaPacket::to_packet() const {
    return wrap_packet("<Pa>", "</Pa>", ident + ";" + password + ";");
}
Result<PaPacket> PaPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Pa>", "</Pa>");
    if (!r.ok()) return Result<PaPacket>::fail(r.error);
    if (r.value.size() != 2) return Result<PaPacket>::fail(Error::MalformedFrame);
    PaPacket p;
    p.ident    = r.value[0];
    p.password = r.value[1];
    return Result<PaPacket>::success(p);
}

// ─── Pc ──────────────────────────────────────────────────────────────────────
std::string PcPacket::to_packet() const {
    std::string payload = std::to_string(static_cast<long long>(timestamp)) + ";"
                        + std::to_string(command_id) + ";"
                        + message + ";";
    return wrap_packet("<Pc>", "</Pc>", payload);
}
Result<PcPacket> PcPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Pc>", "</Pc>");
    if (!r.ok()) return Result<PcPacket>::fail(r.error);
    if (r.value.size() != 3) return Result<PcPacket>::fail(Error::MalformedFrame);
    try {
        PcPacket p;
        p.timestamp  = static_cast<std::time_t>(std::stoll(r.value[0]));
        p.command_id = std::stoi(r.value[1]);
        p.message    = r.value[2];
        return Result<PcPacket>::success(p);
    } catch (...) {
        return Result<PcPacket>::fail(Error::ParseError);
    }
}

// ─── Ts ──────────────────────────────────────────────────────────────────────
std::string TsPacket::to_packet() const {
    std::string payload = std::to_string(static_cast<long long>(timestamp)) + ";"
                        + trip_id + ";";
    return wrap_packet("<Ts>", "</Ts>", payload);
}
Result<TsPacket> TsPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Ts>", "</Ts>");
    if (!r.ok()) return Result<TsPacket>::fail(r.error);
    if (r.value.size() != 2) return Result<TsPacket>::fail(Error::MalformedFrame);
    try {
        TsPacket p;
        p.timestamp = static_cast<std::time_t>(std::stoll(r.value[0]));
        p.trip_id   = r.value[1];
        return Result<TsPacket>::success(p);
    } catch (...) {
        return Result<TsPacket>::fail(Error::ParseError);
    }
}

// ─── Te ──────────────────────────────────────────────────────────────────────
std::string TePacket::to_packet() const {
    char dist_buf[32], speed_buf[32];
    std::snprintf(dist_buf,  sizeof(dist_buf),  "%.3f", distance_traveled);
    std::snprintf(speed_buf, sizeof(speed_buf), "%.3f", max_speed);
    std::string payload = std::to_string(static_cast<long long>(timestamp)) + ";"
                        + trip_id + ";"
                        + dist_buf + ";"
                        + speed_buf + ";"
                        + std::to_string(duration) + ";";
    return wrap_packet("<Te>", "</Te>", payload);
}
Result<TePacket> TePacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Te>", "</Te>");
    if (!r.ok()) return Result<TePacket>::fail(r.error);
    if (r.value.size() != 5) return Result<TePacket>::fail(Error::MalformedFrame);
    try {
        TePacket p;
        p.timestamp          = static_cast<std::time_t>(std::stoll(r.value[0]));
        p.trip_id            = r.value[1];
        p.distance_traveled  = std::stod(r.value[2]);
        p.max_speed          = std::stod(r.value[3]);
        p.duration           = std::stoi(r.value[4]);
        return Result<TePacket>::success(p);
    } catch (...) {
        return Result<TePacket>::fail(Error::ParseError);
    }
}

// ─── Im ──────────────────────────────────────────────────────────────────────
static std::string escape_im(const std::string& s) {
    std::string out;
    out.reserve(s.size());
    for (size_t i = 0; i < s.size(); ++i) {
        if (s[i] == ';') out += "|||";
        else             out += s[i];
    }
    return out;
}

static std::string unescape_im(const std::string& s) {
    std::string out;
    out.reserve(s.size());
    for (size_t i = 0; i < s.size(); ) {
        if (s.compare(i, 3, "|||") == 0) { out += ';'; i += 3; }
        else                              { out += s[i++]; }
    }
    return out;
}

std::string ImPacket::to_packet() const {
    std::string payload = std::to_string(static_cast<long long>(timestamp)) + ";"
                        + chat_id + ";"
                        + escape_im(message) + ";";
    return wrap_packet("<Im>", "</Im>", payload);
}
Result<ImPacket> ImPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Im>", "</Im>");
    if (!r.ok()) return Result<ImPacket>::fail(r.error);
    if (r.value.size() != 3) return Result<ImPacket>::fail(Error::MalformedFrame);
    try {
        ImPacket p;
        p.timestamp = static_cast<std::time_t>(std::stoll(r.value[0]));
        p.chat_id   = r.value[1];
        p.message   = unescape_im(r.value[2]);
        return Result<ImPacket>::success(p);
    } catch (...) {
        return Result<ImPacket>::fail(Error::ParseError);
    }
}

} // namespace layrz::protocol::packets
