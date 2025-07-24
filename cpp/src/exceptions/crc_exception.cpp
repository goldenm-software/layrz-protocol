#include "exceptions/crc_exception.h"

namespace layrz_protocol
{
/// @brief Constructor for the CrcException class.
/// @param message The exception message.
/// @param received The received CRC value.
/// @param calculated The calculated CRC value.
CrcException::CrcException (const std::string &message, int received,
                            int calculated)

{
  message_ = message;
  received_ = received;
  calculated_  = to_hex (calculated);
}
/// @brief Returns the exception message.
/// @return The exception message.
/// @warning The returned pointer is only valid as long as the CrcException
const char *
CrcException::what () const noexcept
{
  return message_.c_str ();
}
/// @brief Returns the exception message.
/// @return The exception message.
std::string
CrcException::get_message () const noexcept
{
  return message_;
}
/// @brief Converts an integer to a hexadecimal string.
/// @param value The integer to convert.
/// @return The hexadecimal string.
/// @warning The returned pointer is only valid as long as the CrcException
/// object exists.
std::string
CrcException::to_hex (int value) const noexcept
{
  std::stringstream ss;
  ss << std::hex << value;
  return ss.str ();
}

} // namespace layrz_protocol