#pragma once

#include <algorithm>
#include <exception>
#include <iomanip>
#include <iterator>
#include <sstream>
#include <sstream>
#include <stdexcept>
#include <string>
#include <vector>

namespace layrz_protocol
{
/// @brief Exception thrown when a malformed error occurs.
/// @details This exception is thrown when a malformed error occurs.
/// @see CommandException
/// @see CrcException
/// @see ParseException
/// @see ServerException
/// @see UnimplementedException
class MalformedException : public std::exception
{
private:
  std::string message_;

public:
  /// @brief Constructor for the MalformedException class.
  /// @param message
  explicit MalformedException (const std::string &message);
  /// @brief Returns the exception message.
  /// @return The exception message.
  /// @warning The returned pointer is only valid as long as the MalformedException object exists.
  const char *what () const noexcept override;
  /// @brief Returns the exception message.
  /// @return The exception message.
  std::string get_message () const noexcept;
};
} // namespace layrz_protocol