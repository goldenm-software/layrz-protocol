#pragma once

#include <string>
#include <vector>

/// @brief layrz_protocol namespace
/// @details Namespace for the layrz protocol
/// @namespace layrz_protocol
/// @copyright Copyright (c) 2025
namespace layrz_protocol
{
/// @namespace helpers
/// @brief helpers namespace
/// @details Namespace for helper functions
namespace helpers
{
/// @namespace str
/// @brief str namespace
/// @details Namespace for string helper functions
namespace str
{
/// @brief Check if a string ends with a suffix
/// @param str string to check
/// @param suffix suffix to check
/// @return true if the string ends with the suffix, false otherwise
bool ends_with (const std::string &value, const std::string &suffix);
/// @brief Check if a string starts with a prefix
/// @param str string to check
/// @param prefix prefix to check
/// @return true if the string starts with the prefix, false otherwise
bool starts_with (const std::string &value, const std::string &prefix);
/// @brief Split a string into a vector of strings
/// @param str string to split
/// @param delimiter delimiter to split the string
/// @return vector of strings
std::vector<std::string> split (const std::string &value, char delimiter);
template <typename Iterator>
std::string join (Iterator begin, Iterator end, char delimiter);
template <typename T>
std::string
join (const std::vector<T> &vec, const std::string &delimiter);

} // namespace str
} // namespace helpers
} // namespace layrz_protocol
