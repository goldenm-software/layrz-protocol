#include "layrz_protocol/packets/pb.hpp"
#include "packets/helpers.hpp"

namespace layrz::protocol::packets {

using namespace detail;

std::string PbPacket::to_packet() const {
    std::string payload;
    for (auto& adv : advertisements) {
        payload += adv.to_packet() + ";";
    }
    return wrap_packet("<Pb>", "</Pb>", payload);
}

Result<PbPacket> PbPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Pb>", "</Pb>");
    if (!r.ok()) return Result<PbPacket>::fail(r.error);

    auto& parts = r.value;
    while (!parts.empty() && parts.back().empty()) parts.pop_back();

    // Each advertisement is 12 ';'-separated fields (11 data + 1 inner CRC).
    // After unwrap_packet split, each adv occupies 12 consecutive parts.
    if (parts.size() % 12 != 0)
        return Result<PbPacket>::fail(Error::MalformedFrame);

    PbPacket p;
    for (size_t i = 0; i < parts.size(); i += 12) {
        // Reconstruct the adv raw string (12 fields joined by ';')
        std::string adv_raw;
        for (size_t j = 0; j < 12; ++j) {
            adv_raw += parts[i + j];
            adv_raw += ';';
        }
        auto adv = ble::Advertisement::from_packet(adv_raw);
        if (!adv.ok()) return Result<PbPacket>::fail(adv.error);
        p.advertisements.push_back(std::move(adv.value));
    }
    return Result<PbPacket>::success(p);
}

} // namespace layrz::protocol::packets
