#include "ble/ble_advertisement.h"

/// @brief Namespace for the layrz protocol
/// This namespace contains all the classes and methods related to the layrz
/// protocol
/// @copyright Copyright (c) 2025
namespace layrz_protocol
{
/// @brief Construct a new BleAdvertisement object
/// @param raw string containing the advertisement data
/// @return BleAdvertisement
BleAdvertisement::BleAdvertisement (
    const std::string &mac_address,
    const std::chrono::system_clock::time_point &timestamp,
    const std::optional<float> &latitude,
    const std::optional<float> &longitude,
    const std::optional<float> &altitude, const std::string &model,
    const std::optional<int> &tx_power,
    const std::vector<BleManufacturerData> &manufacturer_data,
    const std::vector<BleServiceData> &service_data, int rssi)
{
  this->mac_address = mac_address;
  this->timestamp = timestamp;
  this->latitude = latitude;
  this->longitude = longitude;
  this->altitude = altitude;
  this->model = model;
  this->tx_power = tx_power;
  this->manufacturer_data = manufacturer_data;
  this->service_data = service_data;
  this->rssi = rssi;
}
/// @brief Construct a new BleAdvertisement object
/// @param raw string containing the advertisement data
/// @return
BleAdvertisement::BleAdvertisement (const std::string &raw)
{
  *this = from_packet (raw);
}
/// @brief Construct a new BleAdvertisement object
/// @param raw string containing the advertisement data
/// @return BleAdvertisement
BleAdvertisement
BleAdvertisement::from_packet (const std::string &raw)
{
  std::vector<std::string> parts;
  std::stringstream ss (raw);
  std::string item;
  // Split the string by the semicolon
  while (std::getline (ss, item, ';'))
    parts.push_back (item);
  // Check if the packet has the correct number of fields
  if (parts.size () != 11)
    throw MalformedException ("Invalid packet definition");
  // Extract the fields
  auto raw_mac_address = parts[0];
  auto raw_timestamp = parts[1];
  auto raw_latitude = parts[2];
  auto raw_longitude = parts[3];
  auto raw_altitude = parts[4];
  auto raw_model = parts[5];
  auto raw_rssi = parts[6];
  auto raw_tx_power = parts[7];
  auto raw_manufacturer_data = parts[8];
  auto raw_service_data = parts[9];
  auto raw_crc = parts[10];
  // Convert the CRC to an integer
  int received_crc;
  try
    {
      received_crc = std::stoi (raw_crc, nullptr, 16);
    }
  catch (const std::invalid_argument &)
    {
      received_crc = 0;
    }
  // Check the CRC
  std::string sub_str = raw.substr (0, raw.size () - raw_crc.size () - 1);
  std::vector<uint8_t> bytes (sub_str.begin (), sub_str.end ());
  int calculated_crc = helpers::calculate_crc (bytes);
  // If the CRC is invalid, throw an exception
  if (received_crc != calculated_crc)
    throw CrcException ("Invalid CRC", received_crc, calculated_crc);
  // Check the MAC address
  if (raw_mac_address.size () != 12)
    throw MalformedException ("Invalid MAC Address");
  // Convert the fields to the correct types
  std::string mac_address;
  for (size_t i = 0; i < raw_mac_address.size (); i += 2)
    {
      if (i > 0)
        mac_address += ":";
      mac_address += raw_mac_address.substr (i, 2);
    }
  // Convert the timestamp to a time_point
  std::chrono::system_clock::time_point timestamp;
  try
    {
      timestamp = std::chrono::system_clock::from_time_t (
          std::stoll (raw_timestamp));
    }
  catch (const std::invalid_argument &)
    {
      throw MalformedException (
          "Invalid timestamp, should be an int or float");
    }
  // Convert the latitude to a float
  std::optional<float> latitude;
  if (!raw_latitude.empty ())
    {
      try
        {
          latitude = std::stof (raw_latitude);
        }
      catch (const std::invalid_argument &)
        {
          throw MalformedException ("Invalid latitude, should be a float");
        }
    }
  // Convert the longitude to a float
  std::optional<float> longitude;
  if (!raw_longitude.empty ())
    {
      try
        {
          longitude = std::stof (raw_longitude);
        }
      catch (const std::invalid_argument &)
        {
          throw MalformedException ("Invalid longitude, should be a float");
        }
    }
  // Convert the altitude to a float
  std::optional<float> altitude;
  if (!raw_altitude.empty ())
    {
      try
        {
          altitude = std::stof (raw_altitude);
        }
      catch (const std::invalid_argument &)
        {
          throw MalformedException ("Invalid altitude, should be a float");
        }
    }
  // Convert the RSSI to an int
  int rssi;
  try
    {
      rssi = std::stoi (raw_rssi);
    }
  catch (const std::invalid_argument &)
    {
      throw MalformedException ("Invalid RSSI, should be an int");
    }
  // Convert the TX Power to an int
  std::optional<int> tx_power;
  if (!raw_tx_power.empty ())
    {
      try
        {
          tx_power = std::stoi (raw_tx_power);
        }
      catch (const std::invalid_argument &)
        {
          throw MalformedException ("Invalid TX Power, should be an int");
        }
    }
  // Convert the manufacturer data to a vector of BleManufacturerData
  std::vector<BleManufacturerData> manufacturer_data;
  std::stringstream ss_manufacturer (raw_manufacturer_data);
  while (std::getline (ss_manufacturer, item, ','))
    {
      BleManufacturerData ble = BleManufacturerData (item);
      // ble.from_packet (item);
      manufacturer_data.push_back (ble);
    }

  std::vector<BleServiceData> service_data;
  std::stringstream ss_service (raw_service_data);
  while (std::getline (ss_service, item, ','))
    {
      service_data.push_back (BleServiceData::from_packet (item));
    }

  BleAdvertisement ble = BleAdvertisement (
      mac_address, timestamp, latitude, longitude, altitude, raw_model,
      tx_power, manufacturer_data, service_data, rssi);
  return ble;
}

/// @brief Convert the BleAdvertisement object to a packet string
/// @return std::string The packet string
std::string
BleAdvertisement::to_packet () const noexcept
{
  std::stringstream raw;
  raw << mac_address;
  raw << ";"
      << std::chrono::duration_cast<std::chrono::seconds> (
             timestamp.time_since_epoch ())
             .count ();
  raw << ";" << (latitude ? std::to_string (*latitude) : "");
  raw << ";" << (longitude ? std::to_string (*longitude) : "");
  raw << ";" << (altitude ? std::to_string (*altitude) : "");
  raw << ";" << model;
  raw << ";" << rssi;
  raw << ";" << (tx_power ? std::to_string (*tx_power) : "");
  raw << ";" << helpers::str::join (manufacturer_data, ",");
  raw << ";" << helpers::str::join (service_data, ",");
  raw << ";";

  std::vector<uint8_t> raw_bytes (raw.str ().begin (), raw.str ().end ());
  int crc = helpers::calculate_crc (raw_bytes);
  raw << std::hex << std::uppercase << std::setw (4) << std::setfill ('0')
      << crc;

  return raw.str ();
}
} // namespace layrz_protocol
