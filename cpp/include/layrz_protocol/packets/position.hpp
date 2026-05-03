#pragma once
#ifndef __LAYRZ_PROTOCOL_PACKETS_POSITION_HPP__
#define __LAYRZ_PROTOCOL_PACKETS_POSITION_HPP__

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

#endif // __LAYRZ_PROTOCOL_PACKETS_POSITION_HPP__
