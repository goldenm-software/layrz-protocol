#pragma once
#ifndef __LAYRZ_PROTOCOL_BLE_ADVERTISEMENT_HPP__
#define __LAYRZ_PROTOCOL_BLE_ADVERTISEMENT_HPP__

#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/ble/manufacturer_data.hpp"
#include "layrz_protocol/ble/service_data.hpp"
#include <cstdint>
#include <ctime>
#include <optional>
#include <string>
#include <vector>

namespace layrz::protocol::ble {

struct Advertisement {
    std::string mac_address;                     // "12:34:56:78:90:AB"
    std::time_t timestamp;                       // Unix seconds
    std::optional<double> latitude;
    std::optional<double> longitude;
    std::optional<double> altitude;
    std::string model;
    std::string device_name;                     // may be empty
    int         rssi;
    std::optional<int> tx_power;                 // nullopt when missing
    std::vector<ManufacturerData> manufacturer_data;
    std::vector<ServiceData>      service_data;

    // Encode to the 12-field wire format ending with inner CRC:
    // MAC;TS;LAT;LON;ALT;MODEL;NAME;RSSI;TXPOW;MFR;SVC;CRC4
    std::string to_packet() const;

    // Parse from the same 12-field format (validates inner CRC).
    static Result<Advertisement> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::ble

#endif // __LAYRZ_PROTOCOL_BLE_ADVERTISEMENT_HPP__
