#include "exceptions/server_exception.h"

namespace layrz_protocol
{
//------------------------------------------------------------------------------------------------
/// @brief Exception thrown when a server error occurs.
/// @details This exception is thrown when a server error occurs.
/// @warning The returned pointer is only valid as long as the ServerException
/// object exists.
//------------------------------------------------------------------------------------------------
const char *
ServerException::what () const noexcept
{
  return message_.c_str ();
}
//------------------------------------------------------------------------------------------------
/// @brief Returns the exception message.
/// @return The exception message.
std::string ServerException::get_message () const noexcept
{ 
  return message_;
}
} // namespace layrz_protocol