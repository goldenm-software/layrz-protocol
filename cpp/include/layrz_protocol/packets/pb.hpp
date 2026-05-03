#pragma once
#ifndef __LAYRZ_PROTOCOL_PACKETS_PB_HPP__
#define __LAYRZ_PROTOCOL_PACKETS_PB_HPP__

#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/ble/advertisement.hpp"
#include <string>
#include <vector>

namespace layrz::protocol::packets {

// <Pb> — Device → Server — BLE advertisements report
struct PbPacket {
    std::vector<ble::Advertisement> advertisements;

    std::string          to_packet() const;
    static Result<PbPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets

#endif // __LAYRZ_PROTOCOL_PACKETS_PB_HPP__
