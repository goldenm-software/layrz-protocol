#pragma once

#include <iomanip>
#include <sstream>
#include <stdexcept>
#include <string>
#include <string>
#include <vector>

#include "exceptions/malformed_exception.h"

namespace layrz_protocol
{

/// @brief Class to represent a BLE service data
/// @details This class is used to represent a BLE service data
/// @details It contains the UUID of the service and the data to be sent
class BleServiceData
{
private:
  int uuid;
  std::vector<int> data;

public:
  /// @brief Construct a new Ble Service Data object
  /// @param uuid unique identifier for the service
  /// @param data data to be sent
  BleServiceData (int uuid, const std::vector<int> &data);
  /// @brief Construct a new Ble Service Data object from a packet
  /// @param raw packet to be converted
  /// @return BleServiceData
  static BleServiceData from_packet (const std::string &raw);
  /// @brief Convert the object to a packet
  /// @return std::string
  std::string to_packet () const noexcept;
};

};
