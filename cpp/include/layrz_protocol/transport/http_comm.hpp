#pragma once
#ifndef __LAYRZ_PROTOCOL_TRANSPORT_HTTP_COMM_HPP__
#define __LAYRZ_PROTOCOL_TRANSPORT_HTTP_COMM_HPP__

#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/parser.hpp"
#include <string>

namespace layrz::protocol::transport {

// HTTP client — mirrors Go's HttpComm.
// POST /v2/message   (raw packet text body)
// GET  /v2/commands
// Authorization: LayrzAuth <ident>;<password>
//
// Requires LAYRZ_PROTOCOL_CLIENTS to be defined (set by CMake when
// LAYRZ_PROTOCOL_BUILD_NET=ON, or manually via -DLAYRZ_PROTOCOL_CLIENTS in
// PlatformIO/ESP-IDF builds).
#ifdef LAYRZ_PROTOCOL_CLIENTS

enum class HttpScheme { Http, Https };

class HttpComm {
public:
    HttpComm(HttpScheme scheme, std::string host,
             std::string ident, std::string password);

    Result<AnyServerPacket> send(const AnyClientPacket& packet);
    Result<AnyServerPacket> get_commands();

private:
    std::string make_url(const std::string& path) const;
    std::string auth_header() const;

    HttpScheme  scheme_;
    std::string host_;
    std::string ident_;
    std::string password_;
};

#endif // LAYRZ_PROTOCOL_CLIENTS

} // namespace layrz::protocol::transport

#endif // __LAYRZ_PROTOCOL_TRANSPORT_HTTP_COMM_HPP__
