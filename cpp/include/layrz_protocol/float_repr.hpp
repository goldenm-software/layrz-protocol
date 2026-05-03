#pragma once
#ifndef __LAYRZ_PROTOCOL_FLOAT_REPR_HPP__
#define __LAYRZ_PROTOCOL_FLOAT_REPR_HPP__

#include <string>

namespace layrz::protocol {

// Format a double using Python's str(float) semantics:
//   - Shortest round-trip representation
//   - Always includes a decimal point (whole numbers get ".0" suffix)
//   - Examples: 10.0 → "10.0", 1.5 → "1.5", -15.5 → "-15.5", 0.0 → "0.0"
std::string python_repr_float(double v);

} // namespace layrz::protocol

#endif // __LAYRZ_PROTOCOL_FLOAT_REPR_HPP__
