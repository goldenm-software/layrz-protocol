#pragma once

#include <stdexcept>

/// @namespace layrz_protocol
/// @brief layrz_protocol namespace
/// @details Contains all the classes and functions for the layrz protocol.
/// @copyright 2025
namespace layrz_protocol
{
/// @brief Exception thrown when a parsing error occurs.
class BasePacket
{
public:
  /// @brief Constructor for the BasePacket class.
  /// @details This constructor creates a new BasePacket object.
  BasePacket () = default;
  /// @brief Destructor for the BasePacket class.
  /// @details This destructor destroys a BasePacket object.
  ~BasePacket () = default;
  /// @brief Create a packet from raw data
  /// @param raw The raw data
  /// @return The packet
  static BasePacket from_packet (const std::string &raw);
  /// @brief Convert packet to raw data
  /// @return The raw data
  std::string to_packet () const;
  /// @brief String representation of the packet
  /// @return The string representation
  std::string to_string () const;
};
}
