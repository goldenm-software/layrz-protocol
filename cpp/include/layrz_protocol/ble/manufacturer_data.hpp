#pragma once
#ifndef __LAYRZ_PROTOCOL_BLE_MANUFACTURER_DATA_HPP__
#define __LAYRZ_PROTOCOL_BLE_MANUFACTURER_DATA_HPP__

#include "layrz_protocol/errors.hpp"
#include <cstdint>
#include <string>
#include <vector>

namespace layrz::protocol::ble {

struct ManufacturerData {
    int                  company_id; // e.g. 0x004C
    std::vector<uint8_t> data;

    // Encode to "004C:AABBCC..." (company_id %04X, bytes %02X)
    std::string to_packet() const;

    // Parse from "004C:AABBCC..."
    static Result<ManufacturerData> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::ble

#endif // __LAYRZ_PROTOCOL_BLE_MANUFACTURER_DATA_HPP__
