#pragma once

#include <exception>
#include <string>

namespace layrz_protocol
{
/// @brief Exception thrown when a server error occurs.
/// @details This exception is thrown when a server error occurs.
/// @see CommandException
/// @see CrcException
/// @see MalformedException
/// @see ParseException
/// @see UnimplementedException
class ServerException : public std::exception
{
private:
  std::string message_;

public:
  /// @brief Constructor for the ServerException class.
  /// @param message
  explicit ServerException (const std::string &message) : message_ (message) {}
  /// @brief Returns the exception message.
  /// @return The exception message.
  /// @warning The returned pointer is only valid as long as the
  /// ServerException object exists.
  const char *what () const noexcept override;

private:
  /// @brief Returns the exception message.
  /// @return The exception message.
  std::string get_message () const noexcept;
};

}
