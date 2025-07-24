#pragma once

#include <exception>
#include <sstream>
#include <string>

/// @brief The layrz_protocol namespace.
/// @details The layrz_protocol namespace contains all classes, functions, and
/// data types of the layrz_protocol library.
/// @namespace layrz_protocol
/// @copyright 2025
namespace layrz_protocol
{
/// @brief Exception thrown when a command error occurs.
/// @details This exception is thrown when a command error occurs.
/// @see CommandException
/// @see MalformedException
/// @see ParseException
/// @see ServerException
/// @see UnimplementedException
class CrcException : public std::exception
{
private:
  std::string message_;
  int received_;
  std::string calculated_;

public:
  /// @brief Constructor for the CrcException class.
  /// @param message
  /// @param received
  /// @param calculated
  CrcException (const std::string &message, int received, int calculated);
  /// @brief Returns the exception message.
  /// @return The exception message.
  /// @warning The returned pointer is only valid as long as the CrcException
  /// object exists.
  const char *what () const noexcept override;
  /// @brief Returns the exception message.
  /// @return The exception message.
  std::string get_message () const noexcept;

private:
  /// @brief Converts an integer to a hexadecimal string.
  /// @param value
  /// @return The hexadecimal string.
  std::string to_hex (int value) const noexcept;
};
}