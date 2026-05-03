#pragma once
#ifndef __LAYRZ_PROTOCOL_PACKETS_BLE_DATA_HPP__
#define __LAYRZ_PROTOCOL_PACKETS_BLE_DATA_HPP__

#include "layrz_protocol/errors.hpp"
#include <string>

namespace layrz::protocol::packets {

struct BleData {
    std::string mac_address; // stored colon-separated: "12:34:56:78:90:AB"
    std::string model;

    // Encode to "1234567890AB:MODEL" (MAC without colons, uppercase)
    std::string to_packet() const;

    // Parse from "1234567890AB:MODEL"
    static Result<BleData> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets

#endif // __LAYRZ_PROTOCOL_PACKETS_BLE_DATA_HPP__
