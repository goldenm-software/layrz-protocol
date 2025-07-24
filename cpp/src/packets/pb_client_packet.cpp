#include "packets/pb_client_packet.h"

/// @namespace layrz_protocol
/// @brief Namespace for the layrz protocol
/// @details This namespace contains all the classes and methods related to the
/// layrz protocol
/// @copyright Copyright (c) 2025
namespace layrz_protocol
{
/// @brief Construct a new PbClientPacket object
/// @param advertisements
/// @return PbClientPacket The PbClientPacket object. 
PbClientPacket::PbClientPacket (const std::vector<BleAdvertisement> &advertisements)
    : advertisements (advertisements)
{
}
/// @brief Construct a new PbClientPacket object
/// @param raw The packet string.
/// @return PbClientPacket The PbClientPacket object.
PbClientPacket
PbClientPacket::from_packet (const std::string &raw)
{
  if (!helpers::str::starts_with (raw, "<Pb>") || !helpers::str::ends_with (raw, "</Pb>"))
    {
      throw MalformedException (
          "Invalid packet definition, should be <Pb>...</Pb>");
    }

  std::string content = raw.substr (4, raw.size () - 9);
  std::vector<std::string> parts = helpers::str::split (content, ';');
  if (parts.empty () || (parts.size () - 1) % 11 != 0)
    {
      throw MalformedException ("Invalid packet definition");
    }

  int received_crc = 0;
  try
    {
      received_crc = std::stoi (parts.back (), nullptr, 16);
    }
  catch (const std::invalid_argument &)
    {
      received_crc = 0;
    }

  std::string data = helpers::str::join (parts.begin (), parts.end () - 1, ';') + ";";
  std::vector<uint8_t> data_bytes (data.begin (), data.end ());
  int calculated_crc = helpers::calculate_crc (data_bytes);

  if (received_crc != calculated_crc)
    {
      throw CrcException ("Invalid CRC", received_crc, calculated_crc);
    }

  std::vector<BleAdvertisement> advertisements;
  for (size_t i = 0; i < parts.size () - 1; i += 11)
    {
      PbClientPacket packet = from_packet(helpers::str::join(parts.begin() + i, parts.begin() + i + 11, ';'));
      advertisements.insert(advertisements.end(), packet.advertisements.begin(), packet.advertisements.end());
    }

  return PbClientPacket{ advertisements };
}
/// @brief Convert the PbClientPacket object to a string.
/// @details This method converts the PbClientPacket object to a string.
/// @return std::string The packet string.
std::string PbClientPacket::to_packet () const noexcept
{
  std::string raw;
  for (const auto &ad : advertisements)
    {
      raw += ad.to_packet () + ";";
    }
  std::vector<uint8_t> raw_bytes (raw.begin (), raw.end ());
  int crc = helpers::calculate_crc (raw_bytes);
  std::stringstream ss;
  ss << std::uppercase << std::setfill ('0') << std::setw (4) << std::hex
     << crc;
  return "<Pb>" + raw + ss.str () + "</Pb>";
}
/// @brief Convert the PbClientPacket object to a string.
/// @return std::string The packet string.
std::string
PbClientPacket::to_string () const noexcept
{
  return to_packet ();
}

} // namespace layrz_protocol
