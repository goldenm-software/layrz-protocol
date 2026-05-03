#pragma once
#ifndef __LAYRZ_PROTOCOL_TRANSPORT_TCP_COMM_HPP__
#define __LAYRZ_PROTOCOL_TRANSPORT_TCP_COMM_HPP__

#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/parser.hpp"
#include <functional>
#include <string>

#ifdef LAYRZ_PROTOCOL_CLIENTS
#include <atomic>
#include <thread>
#endif

namespace layrz::protocol::transport {

// TCP persistent client — mirrors Go's TcpComm.
// Requires LAYRZ_PROTOCOL_CLIENTS to be defined (set by CMake when
// LAYRZ_PROTOCOL_BUILD_NET=ON, or manually via -DLAYRZ_PROTOCOL_CLIENTS in
// PlatformIO/ESP-IDF builds).
#ifdef LAYRZ_PROTOCOL_CLIENTS

// Framing: outbound packets are sent with a trailing "\r\n".
// Inbound: scans the receive buffer for closing "</Xx>" tags and dispatches
// complete frames without discarding unconsumed bytes (safe unlike the Go impl).
// Auth: on connect() sends PaPacket, waits for AsPacket (timeout configurable).
class TcpComm {
public:
    TcpComm(std::string host, int port,
            std::string ident, std::string password,
            int connect_timeout_secs = 30);
    ~TcpComm();

    void set_callback(std::function<void(AnyServerPacket)> cb);
    Error connect();
    Error send(const AnyClientPacket& packet);
    void close();

private:
    void listen_loop();
    Error write_raw(const std::string& frame);

    std::string host_;
    int         port_;
    std::string ident_;
    std::string password_;
    int         connect_timeout_secs_;

    int                                  sockfd_        = -1;
    std::atomic<bool>                    authenticated_{false};
    std::atomic<bool>                    stop_{false};
    std::thread                          listener_;
    std::function<void(AnyServerPacket)> callback_;
};

#endif // LAYRZ_PROTOCOL_CLIENTS

} // namespace layrz::protocol::transport

#endif // __LAYRZ_PROTOCOL_TRANSPORT_TCP_COMM_HPP__
