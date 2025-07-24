#include "exceptions/command_exception.h"

namespace layrz_protocol
{
//------------------------------------------------------------------------------------------------
/// @brief Constructor for the CommandException class.
/// @param message The exception message.
CommandException::CommandException (const std::string &message)
{
  message_ = message;
}
//------------------------------------------------------------------------------------------------
/// @brief Returns the exception message.
/// @return The exception message.
/// @warning The returned pointer is only valid as long as the CommandException
/// object exists.
//------------------------------------------------------------------------------------------------
const char *
CommandException::what () const noexcept
{
  return message_.c_str ();
}
//------------------------------------------------------------------------------------------------
/// @brief Returns the exception message.
/// @return The exception message.
std::string
CommandException::get_message () const noexcept
{
  return message_;
}
}