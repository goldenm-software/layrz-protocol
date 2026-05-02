#include "layrz_protocol/ble/manufacturer_data.hpp"
#include "layrz_protocol/errors.hpp"
#include <cctype>
#include <iomanip>
#include <sstream>

namespace layrz::protocol::ble {

std::string ManufacturerData::to_packet() const {
    std::ostringstream ss;
    ss << std::uppercase << std::hex << std::setw(4) << std::setfill('0') << company_id << ':';
    for (uint8_t b : data) ss << std::setw(2) << std::setfill('0') << static_cast<int>(b);
    return ss.str();
}

Result<ManufacturerData> ManufacturerData::from_packet(std::string_view raw) {
    auto colon = raw.find(':');
    if (colon == std::string_view::npos)
        return Result<ManufacturerData>::fail(Error::MalformedFrame);

    ManufacturerData md;
    try {
        md.company_id = static_cast<int>(std::stoul(std::string(raw.substr(0, colon)), nullptr, 16));
    } catch (...) {
        return Result<ManufacturerData>::fail(Error::ParseError);
    }

    std::string_view hex = raw.substr(colon + 1);
    if (hex.size() % 2 != 0) return Result<ManufacturerData>::fail(Error::MalformedFrame);
    for (size_t i = 0; i < hex.size(); i += 2) {
        try {
            md.data.push_back(static_cast<uint8_t>(
                std::stoul(std::string(hex.substr(i, 2)), nullptr, 16)));
        } catch (...) {
            return Result<ManufacturerData>::fail(Error::ParseError);
        }
    }
    return Result<ManufacturerData>::success(md);
}

} // namespace layrz::protocol::ble
