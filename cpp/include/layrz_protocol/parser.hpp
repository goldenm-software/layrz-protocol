#pragma once
#ifndef __LAYRZ_PROTOCOL_PARSER_HPP__
#define __LAYRZ_PROTOCOL_PARSER_HPP__

#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/packets/ab.hpp"
#include "layrz_protocol/packets/ac.hpp"
#include "layrz_protocol/packets/ao.hpp"
#include "layrz_protocol/packets/ar.hpp"
#include "layrz_protocol/packets/as.hpp"
#include "layrz_protocol/packets/au.hpp"
#include "layrz_protocol/packets/im.hpp"
#include "layrz_protocol/packets/pa.hpp"
#include "layrz_protocol/packets/pb.hpp"
#include "layrz_protocol/packets/pc.hpp"
#include "layrz_protocol/packets/pd.hpp"
#include "layrz_protocol/packets/pi.hpp"
#include "layrz_protocol/packets/pm.hpp"
#include "layrz_protocol/packets/pr.hpp"
#include "layrz_protocol/packets/ps.hpp"
#include "layrz_protocol/packets/te.hpp"
#include "layrz_protocol/packets/ts.hpp"
#include <string>
#include <string_view>
#include <variant>

namespace layrz::protocol {

// Inbound: packets sent by the server to the device
using AnyServerPacket = std::variant<
    packets::AbPacket,
    packets::AcPacket,
    packets::AoPacket,
    packets::ArPacket,
    packets::AsPacket,
    packets::AuPacket,
    packets::TsPacket,
    packets::TePacket,
    packets::ImPacket
>;

// Outbound: packets sent by the device to the server
using AnyClientPacket = std::variant<
    packets::PaPacket,
    packets::PbPacket,
    packets::PcPacket,
    packets::PdPacket,
    packets::PiPacket,
    packets::PmPacket,
    packets::PrPacket,
    packets::PsPacket,
    packets::TsPacket,
    packets::TePacket,
    packets::ImPacket
>;

// Dispatch an inbound frame (from the server) to the correct packet type.
// The frame must include the opening <Xx> and closing </Xx> tags.
Result<AnyServerPacket> handle_server_output(std::string_view raw);

// Serialize an outbound packet to its wire frame string.
Result<std::string> parse_packet_to_string(const AnyClientPacket& packet);

} // namespace layrz::protocol

#endif // __LAYRZ_PROTOCOL_PARSER_HPP__
