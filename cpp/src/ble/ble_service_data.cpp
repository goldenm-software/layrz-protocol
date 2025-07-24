#include "ble/ble_service_data.h"

namespace layrz_protocol
{
/// @brief Construct a new BleServiceData object
/// @brief Construct a new BleServiceData object
/// @param uuid unique identifier of the service
/// @param data data of the service
BleServiceData::BleServiceData (int uuid, const std::vector<int> &data)
    : uuid (uuid), data (data)
{
}
/// @brief Construct a new BleServiceData object
/// @param raw string containing the service data
/// @return BleServiceData
BleServiceData
BleServiceData::from_packet (const std::string &raw)
{
  size_t pos = raw.find (':');
  if (pos == std::string::npos)
    {
      throw MalformedException ("Invalid packet definition");
    }

  std::string uuid_str = raw.substr (0, pos);
  std::string data_str = raw.substr (pos + 1);

  int service_uuid;
  try
    {
      service_uuid = std::stoi (uuid_str, nullptr, 16);
    }
  catch (const std::invalid_argument &e)
    {
      throw MalformedException ("Invalid service uuid, should be an int");
    }

  std::vector<int> data;
  while (!data_str.empty ())
    {
      data.push_back (std::stoi (data_str.substr (0, 2), nullptr, 16));
      data_str = data_str.substr (2);
    }
  return BleServiceData (service_uuid, data);
}
/// @brief Get the UUID of the service
/// @return int
std::string
BleServiceData::to_packet () const noexcept
{
  std::ostringstream output;
  output << std::uppercase << std::setfill ('0') << std::setw (4) << std::hex
         << uuid << ":";
  for (int byte : data)
    {
      output << std::setw (2) << byte;
    }
  return output.str ();
}
} // namespace layrz_protocol