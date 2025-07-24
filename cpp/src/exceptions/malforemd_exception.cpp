#include "exceptions/malformed_exception.h"

/// @namespace layrz_protocol
/// @brief The layrz_protocol namespace.
/// @details The layrz_protocol namespace contains all classes, functions, and
/// data types of the layrz_protocol library.
/// @copyright 2025
namespace layrz_protocol
{
/// @brief Constructor for the MalformedException class.
/// @param message The exception message.
MalformedException::MalformedException (const std::string &message)
{
  message_ = message;
}
/// @brief Exception thrown when a message is malformed.
/// @return The exception message.
/// @warning The returned pointer is only valid as long as the
/// MalformedException object exists.
const char *
MalformedException::what () const noexcept
{
  return message_.c_str ();
}
/// @brief Returns the exception message.
/// @return The exception message.
std::string
MalformedException::get_message () const noexcept
{
  return message_;
}

} // namespace layrz_protocol