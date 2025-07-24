#include "exceptions/unimplemented_exception.h"

/// @brief The namespace of the layrz protocol.
/// @details This namespace contains all the classes and functions of the layrz
/// protocol.
namespace layrz_protocol
{
//------------------------------------------------------------------------------------------------
/// @brief Constructor for the UnimplementedException class.
/// @param message The exception message.
UnimplementedException::UnimplementedException (const std::string &message)
{
  message_ = message;
}
//------------------------------------------------------------------------------------------------
/// @brief Exception thrown when a message is not implemented.
/// @return The exception message.
/// @warning The returned pointer is only valid as long as the
/// UnimplementedException object exists.
//------------------------------------------------------------------------------------------------
const char *
UnimplementedException::what () const noexcept
{
  return message_.c_str ();
}
//------------------------------------------------------------------------------------------------
/// @brief Returns the exception message.
/// @return The exception message.
std::string
UnimplementedException::get_message () const noexcept
{
  return std::string ();
}

} // namespace layrz_protocol