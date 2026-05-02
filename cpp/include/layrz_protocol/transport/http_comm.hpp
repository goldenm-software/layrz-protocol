#pragma once
#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/parser.hpp"
#include <string>

namespace layrz::protocol::transport {

enum class HttpScheme { Http, Https };

// HTTP client — mirrors Go's HttpComm.
// POST /v2/message   (raw packet text body)
// GET  /v2/commands
// Authorization: LayrzAuth <ident>;<password>
class HttpComm {
public:
    HttpComm(HttpScheme scheme, std::string host,
             std::string ident, std::string password);

    // Send a device-side packet; returns the server response packet.
    Result<AnyServerPacket> send(const AnyClientPacket& packet);

    // Poll for pending commands from the server.
    Result<AnyServerPacket> get_commands();

private:
    std::string make_url(const std::string& path) const;
    std::string auth_header() const;

    HttpScheme  scheme_;
    std::string host_;
    std::string ident_;
    std::string password_;
};

} // namespace layrz::protocol::transport
