#pragma once

#include <exception>
#include <string>

namespace layrz_protocol
{
/// @brief Exception thrown when a command error occurs.
/// @details This exception is thrown when a command error occurs.
/// @see CrcException
/// @see MalformedException
/// @see ParseException
/// @see ServerException
/// @see UnimplementedException
class CommandException : public std::exception
{
private:
  std::string message_;

public:
  /// @brief Constructor for the CommandException class.
  /// @param message
  explicit CommandException (const std::string &message);
  /// @brief Returns the exception message.
  /// @return The exception message.
  /// @warning The returned pointer is only valid as long as the
  /// CommandException object exists.
  const char *what () const noexcept override;

private:
  /// @brief Returns the exception message.
  /// @return The exception message.
  std::string get_message () const noexcept;
};
}