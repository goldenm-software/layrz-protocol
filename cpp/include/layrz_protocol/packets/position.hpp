#pragma once
#include <optional>

namespace layrz::protocol::packets {

struct Position {
    std::optional<double> latitude;
    std::optional<double> longitude;
    std::optional<double> altitude;
    std::optional<double> speed;
    std::optional<double> direction;
    std::optional<int>    satellites;
    std::optional<double> hdop;
};

} // namespace layrz::protocol::packets
