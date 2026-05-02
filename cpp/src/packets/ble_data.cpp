#include "layrz_protocol/packets/ble_data.hpp"
#include "packets/helpers.hpp"
#include <algorithm>
#include <cctype>

namespace layrz::protocol::packets {

using namespace detail;

std::string BleData::to_packet() const {
    return mac_without_colons(mac_address) + ":" + model;
}

Result<BleData> BleData::from_packet(std::string_view raw) {
    auto colon = raw.find(':');
    if (colon == std::string_view::npos)
        return Result<BleData>::fail(Error::MalformedFrame);

    std::string hex = std::string(raw.substr(0, colon));
    if (hex.size() != 12)
        return Result<BleData>::fail(Error::MalformedFrame);

    // Uppercase the hex
    for (char& c : hex) c = static_cast<char>(std::toupper(static_cast<unsigned char>(c)));

    BleData d;
    d.mac_address = mac_with_colons(hex);
    d.model       = std::string(raw.substr(colon + 1));
    return Result<BleData>::success(d);
}

} // namespace layrz::protocol::packets
