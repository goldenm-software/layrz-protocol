#pragma once
#ifndef __LAYRZ_PROTOCOL_BLE_SERVICE_DATA_HPP__
#define __LAYRZ_PROTOCOL_BLE_SERVICE_DATA_HPP__

#include "layrz_protocol/errors.hpp"
#include <cstdint>
#include <string>
#include <vector>

namespace layrz::protocol::ble {

struct ServiceData {
    int                  uuid; // e.g. 0xFD6F
    std::vector<uint8_t> data;

    // Encode to "FD6F:AABBCC..." (uuid %04X, bytes %02X)
    std::string to_packet() const;

    // Parse from "FD6F:AABBCC..."
    static Result<ServiceData> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::ble

#endif // __LAYRZ_PROTOCOL_BLE_SERVICE_DATA_HPP__
