#include "packets/base_packet.h"

/// @namespace layrz_protocol
/// @brief layrz_protocol namespace
/// @details Contains all the classes and functions for the layrz protocol.
/// @copyright 2025
namespace layrz_protocol
{

/// @brief Create a packet from raw data
/// @param raw The raw data
/// @return The packet
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wunused-parameter"
BasePacket
BasePacket::from_packet (const std::string &raw)
{
  throw std::runtime_error ("Method not implemented");
}
#pragma GCC diagnostic pop
/// @brief Convert packet to raw data
/// @return The raw data
std::string
BasePacket::to_packet () const
{
  throw std::runtime_error ("Method not implemented");
}
/// @brief String representation of the packet
/// @return The string representation
std::string
BasePacket::to_string () const
{
  return to_packet ();
}

} // namespace layrz_protocol