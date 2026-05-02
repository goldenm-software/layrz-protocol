#include "layrz_protocol/packets/pi.hpp"
#include "packets/helpers.hpp"

namespace layrz::protocol::packets {

using namespace detail;

std::string PiPacket::to_packet() const {
    std::string fw_id_str = std::visit([](auto&& v) -> std::string {
        using T = std::decay_t<decltype(v)>;
        if constexpr (std::is_same_v<T, int>) return std::to_string(v);
        else return v;
    }, firmware_id);

    std::string payload =
        ident + ";"
        + fw_id_str + ";"
        + std::to_string(firmware_build) + ";"
        + std::to_string(device_id) + ";"
        + std::to_string(hardware_id) + ";"
        + std::to_string(model_id) + ";"
        + std::to_string(static_cast<int>(firmware_branch)) + ";"
        + (fota_enabled ? "1" : "0") + ";";
    return wrap_packet("<Pi>", "</Pi>", payload);
}

Result<PiPacket> PiPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Pi>", "</Pi>");
    if (!r.ok()) return Result<PiPacket>::fail(r.error);
    if (r.value.size() != 8) return Result<PiPacket>::fail(Error::MalformedFrame);

    auto& p = r.value;
    try {
        PiPacket pkt;
        pkt.ident = p[0];

        // firmware_id: try int first, fall back to string
        try {
            pkt.firmware_id = std::stoi(p[1]);
        } catch (...) {
            pkt.firmware_id = p[1];
        }

        pkt.firmware_build = std::stoi(p[2]);
        pkt.device_id      = std::stoi(p[3]);
        pkt.hardware_id    = std::stoi(p[4]);
        pkt.model_id       = std::stoi(p[5]);

        try {
            int branch = std::stoi(p[6]);
            pkt.firmware_branch = (branch == 1) ? FirmwareBranch::Development : FirmwareBranch::Stable;
        } catch (...) {
            pkt.firmware_branch = FirmwareBranch::Stable;
        }

        std::string fota_low = p[7];
        for (char& c : fota_low) c = static_cast<char>(std::tolower(static_cast<unsigned char>(c)));
        pkt.fota_enabled = (fota_low == "true" || fota_low == "1");

        return Result<PiPacket>::success(pkt);
    } catch (...) {
        return Result<PiPacket>::fail(Error::ParseError);
    }
}

} // namespace layrz::protocol::packets
