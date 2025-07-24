#pragma once

#include <exception>
#include <sstream>
#include <string>

/// @namespace layrz_protocol
/// @brief Namespace for the layrz protocol.
/// @details This namespace contains all the classes and functions related to the layrz protocol.
/// @copyright 2025
namespace layrz_protocol
{
/// @brief Exception thrown when a unimplemented error occurs.
/// @details This exception is thrown when a unimplemented error occurs.
/// @see CommandException
/// @see CrcException
/// @see MalformedException
/// @see ParseException
/// @see ServerException
class UnimplementedException : public std::exception
{
private:
  std::string message_;

public:

  /// @brief Constructor for the UnimplementedException class.
  /// @param message The exception message.  
  explicit UnimplementedException (const std::string &message);
  /// @brief Returns the exception message.
  /// @return  The exception message.
  /// @warning The returned pointer is only valid as long as the UnimplementedException object exists.
  const char *what () const noexcept override;
  /// @brief Returns the exception message.
  /// @return The exception message.
  std::string get_message () const noexcept;
};

} // namespace layrz_protocol
