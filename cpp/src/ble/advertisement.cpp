#include "layrz_protocol/ble/advertisement.hpp"
#include "layrz_protocol/crc.hpp"
#include "layrz_protocol/float_repr.hpp"
#include <iomanip>
#include <sstream>
#include <stdexcept>

namespace layrz::protocol::ble {

static std::vector<std::string> split_sv(std::string_view sv, char delim) {
    std::vector<std::string> parts;
    size_t start = 0;
    while (true) {
        size_t pos = sv.find(delim, start);
        if (pos == std::string_view::npos) {
            parts.emplace_back(sv.substr(start));
            break;
        }
        parts.emplace_back(sv.substr(start, pos - start));
        start = pos + 1;
    }
    return parts;
}

static std::string mac_without_colons(const std::string& mac) {
    std::string out;
    out.reserve(12);
    for (char c : mac)
        if (c != ':') out += static_cast<char>(std::toupper(static_cast<unsigned char>(c)));
    return out;
}

static std::string mac_with_colons(const std::string& hex12) {
    if (hex12.size() != 12) return hex12;
    std::string out;
    out.reserve(17);
    for (size_t i = 0; i < 12; i += 2) {
        if (i) out += ':';
        out += hex12[i];
        out += hex12[i+1];
    }
    return out;
}

std::string Advertisement::to_packet() const {
    auto opt_dbl = [](const std::optional<double>& v) -> std::string {
        return v.has_value() ? python_repr_float(*v) : "";
    };
    auto opt_int = [](const std::optional<int>& v) -> std::string {
        return v.has_value() ? std::to_string(*v) : "";
    };

    // Manufacturer data: comma-separated "CCCC:HHHH..."
    std::string mfr_str;
    for (size_t i = 0; i < manufacturer_data.size(); ++i) {
        if (i) mfr_str += ',';
        mfr_str += manufacturer_data[i].to_packet();
    }

    // Service data: comma-separated "UUUU:HHHH..."
    std::string svc_str;
    for (size_t i = 0; i < service_data.size(); ++i) {
        if (i) svc_str += ',';
        svc_str += service_data[i].to_packet();
    }

    std::string raw =
        mac_without_colons(mac_address) + ";"
        + std::to_string(static_cast<long long>(timestamp)) + ";"
        + opt_dbl(latitude) + ";"
        + opt_dbl(longitude) + ";"
        + opt_dbl(altitude) + ";"
        + model + ";"
        + device_name + ";"
        + std::to_string(rssi) + ";"
        + opt_int(tx_power) + ";"
        + mfr_str + ";"
        + svc_str + ";";

    return raw + compute_crc_str(raw);
}

Result<Advertisement> Advertisement::from_packet(std::string_view raw_sv) {
    // 12 fields separated by ';': the last (index 11) is the inner CRC
    auto parts = split_sv(raw_sv, ';');
    // Remove trailing empty tokens
    while (!parts.empty() && parts.back().empty()) parts.pop_back();
    if (parts.size() != 12)
        return Result<Advertisement>::fail(Error::MalformedFrame);

    // Validate inner CRC
    std::string payload;
    for (size_t i = 0; i < 11; ++i) { payload += parts[i]; payload += ';'; }
    uint16_t expected;
    try {
        expected = static_cast<uint16_t>(std::stoul(parts[11], nullptr, 16));
    } catch (...) {
        return Result<Advertisement>::fail(Error::CrcMismatch);
    }
    if (crc16_x25(payload) != expected)
        return Result<Advertisement>::fail(Error::CrcMismatch);

    try {
        Advertisement a;

        // MAC: 12 hex chars → colon-separated
        if (parts[0].size() != 12)
            return Result<Advertisement>::fail(Error::MalformedFrame);
        a.mac_address = mac_with_colons(parts[0]);

        a.timestamp = static_cast<std::time_t>(std::stoll(parts[1]));

        auto opt_dbl = [](const std::string& s) -> std::optional<double> {
            if (s.empty()) return std::nullopt;
            return std::stod(s);
        };
        auto opt_int = [](const std::string& s) -> std::optional<int> {
            if (s.empty()) return std::nullopt;
            return std::stoi(s);
        };

        a.latitude   = opt_dbl(parts[2]);
        a.longitude  = opt_dbl(parts[3]);
        a.altitude   = opt_dbl(parts[4]);
        a.model      = parts[5];
        a.device_name = parts[6];
        a.rssi       = std::stoi(parts[7]);
        a.tx_power   = opt_int(parts[8]);

        // Manufacturer data: comma-separated entries
        if (!parts[9].empty()) {
            for (auto& entry : split_sv(parts[9], ',')) {
                if (entry.empty()) continue;
                auto r = ManufacturerData::from_packet(entry);
                if (!r.ok()) return Result<Advertisement>::fail(r.error);
                a.manufacturer_data.push_back(std::move(r.value));
            }
        }

        // Service data: comma-separated entries
        if (!parts[10].empty()) {
            for (auto& entry : split_sv(parts[10], ',')) {
                if (entry.empty()) continue;
                auto r = ServiceData::from_packet(entry);
                if (!r.ok()) return Result<Advertisement>::fail(r.error);
                a.service_data.push_back(std::move(r.value));
            }
        }

        return Result<Advertisement>::success(a);
    } catch (...) {
        return Result<Advertisement>::fail(Error::ParseError);
    }
}

} // namespace layrz::protocol::ble
