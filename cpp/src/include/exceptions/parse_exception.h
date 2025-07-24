#pragma once

#include <exception>
#include <string>

namespace layrz_protocol
{
/// @brief Exception thrown when a parsing error occurs.
/// @details This exception is thrown when a parsing error occurs.
/// @see CommandException
/// @see CrcException
/// @see MalformedException
/// @see ServerException
/// @see UnimplementedException
class ParseException : public std::exception
{
private:
  std::string message_;

public:
  /// @brief Constructor for the ParseException class.
  /// @param message
  explicit ParseException (const std::string &message);
  /// @brief Returns the exception message.
  /// @return The exception message.
  /// @warning The returned pointer is only valid as long as the ParseException
  /// object exists.
  const char *what () const noexcept override;

private:
  /// @brief Returns the exception message.
  /// @return The exception message.
  std::string get_message () const noexcept;
};

}