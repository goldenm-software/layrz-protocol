#include "layrz_protocol/ble/service_data.hpp"
#include "layrz_protocol/errors.hpp"
#include <iomanip>
#include <sstream>

namespace layrz::protocol::ble {

std::string ServiceData::to_packet() const {
    std::ostringstream ss;
    ss << std::uppercase << std::hex << std::setw(4) << std::setfill('0') << uuid << ':';
    for (uint8_t b : data) ss << std::setw(2) << std::setfill('0') << static_cast<int>(b);
    return ss.str();
}

Result<ServiceData> ServiceData::from_packet(std::string_view raw) {
    auto colon = raw.find(':');
    if (colon == std::string_view::npos)
        return Result<ServiceData>::fail(Error::MalformedFrame);

    ServiceData sd;
    try {
        sd.uuid = static_cast<int>(std::stoul(std::string(raw.substr(0, colon)), nullptr, 16));
    } catch (...) {
        return Result<ServiceData>::fail(Error::ParseError);
    }

    std::string_view hex = raw.substr(colon + 1);
    if (hex.size() % 2 != 0) return Result<ServiceData>::fail(Error::MalformedFrame);
    for (size_t i = 0; i < hex.size(); i += 2) {
        try {
            sd.data.push_back(static_cast<uint8_t>(
                std::stoul(std::string(hex.substr(i, 2)), nullptr, 16)));
        } catch (...) {
            return Result<ServiceData>::fail(Error::ParseError);
        }
    }
    return Result<ServiceData>::success(sd);
}

} // namespace layrz::protocol::ble
