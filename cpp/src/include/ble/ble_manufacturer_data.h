
#pragma once

#include <iomanip>
#include <sstream>
#include <stdexcept>
#include <string>
#include <vector>

#include "exceptions/malformed_exception.h"

namespace layrz_protocol
{
/// @class BleManufacturerData
/// @brief Represents BLE manufacturer-specific data.
/// This class provides functionality to parse and generate BLE
/// manufacturer-specific data packets.
class BleManufacturerData
{
private:
  int company_id;
  std::vector<int> data;

public:
  /// @brief Construct a new BleManufacturerData object
  /// @return BleManufacturerData
  BleManufacturerData () = default;
  /// @brief Construct a new BleManufacturerData object
  /// @param raw string containing the manufacturer data
  /// @return BleManufacturerData
  BleManufacturerData (const std::string &raw);
  /// @brief Destroy the BleManufacturerData object
  /// @return ~BleManufacturerData
  ~BleManufacturerData () = default;
  /// @brief Construct a new BleManufacturerData object
  /// @param raw string containing the manufacturer data
  /// @return BleManufacturerData
  BleManufacturerData from_packet (const std::string &raw);
  /// @brief Convert the BleManufacturerData object to a packet string
  /// @return std::string The packet string in the format "company_id:data"
  std::string to_packet ();
};

} // namespace layrz_protocol