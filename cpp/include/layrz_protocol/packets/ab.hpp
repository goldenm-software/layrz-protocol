#pragma once
#ifndef __LAYRZ_PROTOCOL_PACKETS_AB_HPP__
#define __LAYRZ_PROTOCOL_PACKETS_AB_HPP__

#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/packets/ble_data.hpp"
#include <string>
#include <vector>

namespace layrz::protocol::packets {

// <Ab> — Server → Device — BLE allow-list push
struct AbPacket {
    std::vector<BleData> devices;

    std::string          to_packet() const;
    static Result<AbPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets

#endif // __LAYRZ_PROTOCOL_PACKETS_AB_HPP__
