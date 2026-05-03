#include "layrz_protocol/float_repr.hpp"
#include <charconv>
#include <cmath>
#include <cstring>
#include <string>

namespace layrz::protocol {

// Matches Python's str(float):
//   - Uses the shortest decimal representation that round-trips exactly.
//   - Always includes a '.' (whole-valued floats get a ".0" suffix).
std::string python_repr_float(double v) {
    // Handle special values
    if (std::isnan(v))  return "nan";
    if (std::isinf(v))  return v > 0 ? "inf" : "-inf";

    char buf[64];
    auto [ptr, ec] = std::to_chars(buf, buf + sizeof(buf), v, std::chars_format::general);
    if (ec != std::errc{}) {
        // Fallback (should never happen for finite doubles)
        return std::to_string(v);
    }
    std::string s(buf, ptr);

    // If the representation has no '.', no 'e', and no 'E' — it's a whole number
    // stored without a decimal point (e.g. "10"). Add ".0" to match Python.
    bool has_dot = s.find('.') != std::string::npos;
    bool has_exp = s.find('e') != std::string::npos || s.find('E') != std::string::npos;
    if (!has_dot && !has_exp) {
        s += ".0";
    }

    return s;
}

} // namespace layrz::protocol
