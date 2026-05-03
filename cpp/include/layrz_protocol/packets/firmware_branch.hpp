#pragma once
#ifndef __LAYRZ_PROTOCOL_PACKETS_FIRMWARE_BRANCH_HPP__
#define __LAYRZ_PROTOCOL_PACKETS_FIRMWARE_BRANCH_HPP__


namespace layrz::protocol::packets {

enum class FirmwareBranch {
    Stable      = 0,
    Development = 1,
};

} // namespace layrz::protocol::packets

#endif // __LAYRZ_PROTOCOL_PACKETS_FIRMWARE_BRANCH_HPP__
