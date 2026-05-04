#pragma once
#ifndef __LAYRZ_PROTOCOL_SERVERS_TCP_HPP__
#define __LAYRZ_PROTOCOL_SERVERS_TCP_HPP__

#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/parser.hpp"
#include <functional>
#include <memory>
#include <optional>
#include <string>

namespace layrz::protocol::servers {

// LAYRZ_PROTOCOL_SERVERS is defined when linking layrz::protocol::servers
// (set by CMake LAYRZ_PROTOCOL_BUILD_SERVERS=ON, or manually).
// Without it, this header is empty — safe to include from embedded targets.
#ifdef LAYRZ_PROTOCOL_SERVERS

// Opaque handle exposing information about a connected TCP client.
struct TcpConnection {
    int         fd   = -1;
    std::string addr;  // "ip:port" string for logging
};

struct TcpConfig {
    int port = 0;

    // Called for every fully-decoded client packet.
    // Return a non-empty optional to write that server packet back to the device.
    // REQUIRED — construction fails if nullptr.
    std::function<std::optional<AnyServerPacket>(const AnyClientPacket&, TcpConnection&)> on_new_packet;

    // Called when a frame cannot be decoded.
    // Defaults to stderr logging if nullptr.
    std::function<void(Error, std::string_view raw, TcpConnection&)> on_decode_error;
};

// TCP server — mirrors Go's TcpServer.
// start() is blocking; call close() from another thread or signal handler to stop.
class TcpServer {
public:
    static Result<TcpServer> create(TcpConfig cfg);

    TcpServer() = default;
    TcpServer(TcpServer&&) noexcept;
    TcpServer& operator=(TcpServer&&) noexcept;
    ~TcpServer();

    // Blocking. Returns when close() is called or a fatal listen error occurs.
    Error start();

    // Unblocks start(); safe to call from a signal handler or another thread.
    Error close();

private:
    struct Impl;
    std::unique_ptr<Impl> impl_;

    explicit TcpServer(std::unique_ptr<Impl> impl);
    void handle_connection(int client_fd, TcpConnection conn);
};

#endif // LAYRZ_PROTOCOL_SERVERS

} // namespace layrz::protocol::servers

#endif // __LAYRZ_PROTOCOL_SERVERS_TCP_HPP__
