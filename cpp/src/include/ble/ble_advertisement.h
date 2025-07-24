#pragma once

#include <chrono>
#include <optional>

#include "ble/ble_manufacturer_data.h"
#include "ble/ble_service_data.h"
#include "exceptions/crc_exception.h"
#include "helpers/crc.h"
#include "helpers/strings.h"

/// @brief Namespace for the layrz protocol
/// This namespace contains all the classes and methods related to the layrz
/// protocol
/// @copyright Copyright (c) 2025
namespace layrz_protocol
{
/// @brief Represents a BLE advertisement
/// This class provides functionality to parse and generate BLE advertisement
/// packets. The BLE advertisement packet is a string containing the following
/// fields separated by a semicolon:
/// - MAC address
/// - Timestamp
/// - Latitude
/// - Longitude
/// - Altitude
/// - Model
/// - RSSI
/// - TX Power
/// - Manufacturer data
/// - Service data
/// - CRC
class BleAdvertisement
{
public:
  std::string mac_address;
  std::chrono::system_clock::time_point timestamp;
  std::optional<float> latitude;
  std::optional<float> longitude;
  std::optional<float> altitude;
  std::string model;
  std::optional<int> tx_power;
  std::vector<BleManufacturerData> manufacturer_data;
  std::vector<BleServiceData> service_data;
  int rssi;

public:
  /// @brief Construct a new BleAdvertisement object
  /// @param mac_address MAC address of the device
  /// @param timestamp Timestamp of the advertisement
  /// @param latitude Latitude of the device
  /// @param longitude Longitude of the device
  /// @param altitude Altitude of the device
  /// @param model Model of the device
  /// @param tx_power TX power of the device
  /// @param manufacturer_data Manufacturer data of the device
  /// @param service_data Service data of the device
  /// @param rssi RSSI of the device
  BleAdvertisement (
      const std::string &mac_address,
      const std::chrono::system_clock::time_point &timestamp,
      const std::optional<float> &latitude,
      const std::optional<float> &longitude,
      const std::optional<float> &altitude,
      const std::string &model,
      const std::optional<int> &tx_power,
      const std::vector<BleManufacturerData> &manufacturer_data,
      const std::vector<BleServiceData> &service_data,
      int rssi);
  /// @brief Construct a new BleAdvertisement object
  /// @param raw string containing the advertisement data
  /// @return BleAdvertisement
  BleAdvertisement (const std::string &raw);
  /// @brief Destroy the BleAdvertisement object
  ~BleAdvertisement () = default;
  /// @brief Construct a new BleAdvertisement object
  /// @param raw string containing the advertisement data
  /// @return BleAdvertisement
  BleAdvertisement from_packet (const std::string &raw);
  /// @brief Convert the BleAdvertisement object to a packet string
  /// @return std::string The packet string
  std::string to_packet () const noexcept;
};

} // namespace layrz_protocol