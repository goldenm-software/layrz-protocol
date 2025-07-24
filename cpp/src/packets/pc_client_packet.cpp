
#include "packets/pc_client_packet.h"

namespace layrz_protocol
{

PcClientPacket::PcClientPacket (std::time_t timestamp, int command_id, const std::string &message)
    : timestamp (timestamp), command_id (command_id), message (message)
{
}

PcClientPacket PcClientPacket::from_packet (const std::string &raw)
{
  if (!helpers::str::starts_with (raw, "<Pc>")
      || !helpers::str::ends_with (raw, "</Pc>"))
    throw MalformedException (
        "Invalid packet definition, should be <Pc>...</Pc>");

  std::string content = raw.substr (4, raw.size () - 9);
  std::vector<std::string> parts;
  std::istringstream ss (content);
  std::string part;
  while (std::getline (ss, part, ';'))
    parts.push_back (part);

  if (parts.size () != 4)
    throw MalformedException (
        "Invalid packet definition, should have 4 parts");

  int received_crc;
  try
    {
      received_crc = std::stoi (parts[3], nullptr, 16);
    }
  catch (const std::invalid_argument &)
    {
      received_crc = 0;
    }

  std::string data_to_crc
      = content.substr (0, content.size () - parts[3].size () - 1);
  std::vector<uint8_t> data_to_crc_bytes (data_to_crc.begin (),
                                          data_to_crc.end ());
  int calculated_crc = helpers::calculate_crc (data_to_crc_bytes);

  if (received_crc != calculated_crc)
    throw CrcException ("Invalid CRC", received_crc, calculated_crc);

  std::time_t timestamp;
  try
    {
      timestamp = std::stoll (parts[0]);
    }
  catch (const std::invalid_argument &)
    {
      throw MalformedException (
          "Invalid timestamp, should be an int or float");
    }

  int command_id;
  try
    {
      command_id = std::stoi (parts[1]);
    }
  catch (const std::invalid_argument &)
    {
      throw MalformedException ("Invalid command_id, should be an int");
    }

  return PcClientPacket (timestamp, command_id, parts[2]);
}

/// @brief Convert the PcClientPacket object to a packet string.
/// @return std::string The packet string.
std::string
PcClientPacket::to_packet () const noexcept
{
  /// @todo Implement the to_packet method
  std::ostringstream raw;
  raw << timestamp << ";" << command_id << ";" << message << ";";
  std::string raw_str = raw.str ();
  std::vector<uint8_t> raw_bytes (raw_str.begin (), raw_str.end ());
  /// @todo Implement the CRC calculation
  uint16_t crc = helpers::calculate_crc (raw_bytes);
  std::ostringstream crc_str;
  crc_str << std::uppercase << std::setfill ('0') << std::setw (4) << std::hex
          << crc;
  return "<Pc>" + raw_str + crc_str.str () + "</Pc>";
}

/// @brief Convert the PcClientPacket object to a string.
/// @return std::string The packet string.
std::string
PcClientPacket::to_string () const noexcept
{
  return to_packet ();
}

} // namespace layrz_protocol