#pragma once
#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/parser.hpp"
#include <atomic>
#include <functional>
#include <string>
#include <thread>

namespace layrz::protocol::transport {

// TCP persistent client — mirrors Go's TcpComm.
// Framing: outbound packets are sent with a trailing "\r\n".
// Inbound: scans the receive buffer for closing "</Xx>" tags and dispatches
// complete frames without discarding unconsumed bytes (safe unlike the Go impl).
// Auth: on Connect() sends PaPacket, waits for AsPacket (timeout configurable).
class TcpComm {
public:
    TcpComm(std::string host, int port,
            std::string ident, std::string password,
            int connect_timeout_secs = 30);
    ~TcpComm();

    // Register callback invoked for every inbound server packet.
    void set_callback(std::function<void(AnyServerPacket)> cb);

    // Dial, authenticate (Pa→As handshake), start listener thread.
    Error connect();

    // Send a device-side packet.
    Error send(const AnyClientPacket& packet);

    // Close the connection and stop the listener thread.
    void close();

private:
    void listen_loop();
    Error write_raw(const std::string& frame);

    std::string host_;
    int         port_;
    std::string ident_;
    std::string password_;
    int         connect_timeout_secs_;

    int                                 sockfd_      = -1;
    std::atomic<bool>                   authenticated_{false};
    std::atomic<bool>                   stop_{false};
    std::thread                         listener_;
    std::function<void(AnyServerPacket)> callback_;
};

} // namespace layrz::protocol::transport
