#include "layrz_protocol/packets/ab.hpp"
#include "packets/helpers.hpp"

namespace layrz::protocol::packets {

using namespace detail;

std::string AbPacket::to_packet() const {
    std::string payload;
    for (auto& d : devices) { payload += d.to_packet() + ";"; }
    return wrap_packet("<Ab>", "</Ab>", payload);
}

Result<AbPacket> AbPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Ab>", "</Ab>");
    if (!r.ok()) return Result<AbPacket>::fail(r.error);
    // Each device token is "1234567890AB:MODEL"; last token may be empty if trailing ';'
    // unwrap_packet strips trailing ';' implicitly via payload rebuild — but split on ';'
    // means parts may include a trailing empty if the data had trailing ';' before CRC.
    // Remove trailing empty parts.
    auto& parts = r.value;
    while (!parts.empty() && parts.back().empty()) parts.pop_back();
    if (parts.empty()) return Result<AbPacket>::fail(Error::MalformedFrame);

    AbPacket p;
    for (auto& tok : parts) {
        auto d = BleData::from_packet(tok);
        if (!d.ok()) return Result<AbPacket>::fail(d.error);
        p.devices.push_back(std::move(d.value));
    }
    return Result<AbPacket>::success(p);
}

} // namespace layrz::protocol::packets
