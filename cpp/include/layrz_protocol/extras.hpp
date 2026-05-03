#pragma once
#ifndef __LAYRZ_PROTOCOL_EXTRAS_HPP__
#define __LAYRZ_PROTOCOL_EXTRAS_HPP__

#include <map>
#include <string>
#include <variant>
#include <vector>

namespace layrz::protocol {

// Wire value type: int | double | bool | string  (matching Python's parse_extra typing)
using ExtraValue = std::variant<int64_t, double, bool, std::string>;

// Ordered extras map (insertion-order preserved, matching Python's dict ordering)
using ExtrasMap = std::vector<std::pair<std::string, ExtraValue>>;

// Parse a "k1:v1,k2:v2,..." extras field from the wire.
// Keys are remapped per the Layrz one-way decode-side remap table (e.g. io1.di →
// gpio.1.digital.input). String values have accent-stripping applied (ASCII_MAP).
ExtrasMap parse_extra(std::string_view raw);

// Encode an extras map back to "k1:v1,k2:v2,..." wire form.
// Keys are emitted verbatim (no reverse remap). Booleans emit as "true"/"false".
// Floats use python_repr_float. Ints emit as plain decimal.
std::string cast_extra(const ExtrasMap& extras);

} // namespace layrz::protocol

#endif // __LAYRZ_PROTOCOL_EXTRAS_HPP__
