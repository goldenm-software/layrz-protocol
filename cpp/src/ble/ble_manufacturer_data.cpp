#include "ble/ble_manufacturer_data.h"

/// @namespace layrz_protocol
/// @brief Namespace for the layrz protocol
/// @details This namespace contains all the classes and functions related to
/// the layrz protocol.
/// @copyright Copyright (c) 2025
namespace layrz_protocol
{
/// @brief Construct a new BleManufacturerData object
/// @param raw string containing the manufacturer data
/// @return
BleManufacturerData::BleManufacturerData (const std::string &raw)
{
  *this = from_packet(raw);
}
/// @brief Construct a new BleManufacturerData object
/// @param raw string containing the manufacturer data
/// @return BleManufacturerData
BleManufacturerData
BleManufacturerData::from_packet (const std::string &raw)
{
  // Parse the company id
  BleManufacturerData result;
  size_t pos = raw.find (':');
  if (pos == std::string::npos)
    {
      throw MalformedException ("Invalid packet definition");
    }
  try
    {
      result.company_id = std::stoi (raw.substr (0, pos), nullptr, 16);
    }
  catch (const std::invalid_argument &e)
    {
      throw MalformedException ("Invalid company id, should be an int");
    }

  // Parse the data
  std::string data_str = raw.substr (pos + 1);
  while (!data_str.empty ())
    {
      result.data.push_back (std::stoi (data_str.substr (0, 2), nullptr, 16));
      data_str = data_str.substr (2);
    }

  return result;
}
/// @brief Convert the BleManufacturerData object to a packet string
/// @return std::string The packet string in the format "company_id:data"
std::string
BleManufacturerData::to_packet ()
{
  std::ostringstream output;
  output << std::uppercase << std::setfill ('0') << std::setw (4) << std::hex
         << company_id << ':';
  for (int byte : data)
    {
      output << std::setw (2) << byte;
    }
  return output.str ();
}
} // namespace layrz_protocol