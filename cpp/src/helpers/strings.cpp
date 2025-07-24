#include "helpers/strings.h"

#include <sstream>

namespace layers_protocol
{
namespace helpers
{
namespace str
{
/// @brief Check if a string starts with a prefix
/// @param value
/// @param prefix
/// @return true if the string starts with the prefix
bool
starts_with (const std::string &value, const std::string &prefix)
{
  return value.size () >= prefix.size ()
         && value.compare (0, prefix.size (), prefix) == 0;
}
/// @brief Check if a string ends with a suffix
/// @param value
/// @param suffix
/// @return true if the string ends with the suffix
bool
ends_with (const std::string &value, const std::string &suffix)
{
  return value.size () >= suffix.size ()
         && value.compare (value.size () - suffix.size (), suffix.size (),
                           suffix)
                == 0;
}
/// @brief Split a string into a vector of strings
/// @param value
/// @param delimiter
/// @return a vector of strings
std::vector<std::string>
split (const std::string &value, char delimiter)
{
  std::vector<std::string> tokens;
  std::string token;
  std::istringstream tokenStream (value);
  while (std::getline (tokenStream, token, delimiter))
    {
      tokens.push_back (token);
    }
  return tokens;
}
/// @brief Join a vector of strings into a single string
/// @param begin iterator to the beginning of the vector
/// @param end iterator to the end of the vector
/// @param delimiter character to separate the strings
/// @return a single string
template <typename Iterator>
std::string
join (Iterator begin, Iterator end, char delimiter)
{
  std::ostringstream oss;
  if (begin != end)
    {
      oss << *begin++;
    }
  while (begin != end)
    {
      oss << delimiter << *begin++;
    }
  return oss.str ();
}
template <typename T>
std::string
join (const std::vector<T> &vec, const std::string &delimiter)
{
  std::ostringstream result;
  for (size_t i = 0; i < vec.size (); ++i)
    {
      if (i != 0)
        {
          result << delimiter;
        }
      result << vec[i].to_packet ();
    }
  return result.str ();
}
} // namespace str
} // namespace helpers
} // namespace layers_protocol