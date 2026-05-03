#pragma once
#ifndef __LAYRZ_PROTOCOL_PACKETS_PI_HPP__
#define __LAYRZ_PROTOCOL_PACKETS_PI_HPP__

#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/packets/firmware_branch.hpp"
#include <string>
#include <variant>

namespace layrz::protocol::packets {

// firmware_id may be int or string on the wire
using FirmwareId = std::variant<int, std::string>;

// <Pi> — Device → Server — Identification
struct PiPacket {
    std::string    ident;
    FirmwareId     firmware_id;
    int            firmware_build;
    int            device_id;
    int            hardware_id;
    int            model_id;
    FirmwareBranch firmware_branch = FirmwareBranch::Stable;
    bool           fota_enabled    = false;

    std::string          to_packet() const;
    static Result<PiPacket> from_packet(std::string_view raw);
};

} // namespace layrz::protocol::packets

#endif // __LAYRZ_PROTOCOL_PACKETS_PI_HPP__
