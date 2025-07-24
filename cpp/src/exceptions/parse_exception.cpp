#include "exceptions/parse_exception.h"

namespace layrz_protocol
{
/// @brief Returns the exception message.
/// @return The exception message.
/// @warning The returned pointer is only valid as long as the ParseException
/// object exists.
const char *
ParseException::what () const noexcept
{
  return message_.c_str ();
}
/// @brief Returns the exception message.
/// @return The exception message.
std::string
ParseException::get_message () const noexcept
{
  return message_;
}

} // namespace layrz_protocol