#pragma once
#ifndef __LAYRZ_PROTOCOL_SERVERS_HTTP_HPP__
#define __LAYRZ_PROTOCOL_SERVERS_HTTP_HPP__

#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/parser.hpp"
#include <functional>
#include <string>
#include <string_view>

#ifdef LAYRZ_PROTOCOL_SERVERS
#include <memory>
#endif

namespace layrz::protocol::servers {

#ifdef LAYRZ_PROTOCOL_SERVERS

// Opaque wrapper around an incoming HTTP request, exposing only what server
// callbacks need without leaking the httplib header into consumer code.
struct HttpRequest {
    std::string method;
    std::string path;
    std::string remote_addr;
    // Returns the value of a header, or empty string if absent.
    std::function<std::string(std::string_view)> get_header;
};

struct HttpConfig {
    int port = 0;

    // Called for every POST /v2/message.
    // Return non-empty optional to write the encoded server packet as the response body (200).
    // Return empty optional for 204.
    // Throw or return an error-carrying state for 500.
    // REQUIRED — construction fails if nullptr.
    std::function<std::optional<AnyServerPacket>(const AnyClientPacket&, const HttpRequest&)> on_new_packet;

    // Called for every GET /v2/commands.
    // ident and passwd are extracted from the LayrzAuth header.
    // Return empty optional for 204; non-empty for 200 with encoded packet.
    // If nullptr, all /v2/commands requests receive 204.
    std::function<std::optional<AnyServerPacket>(std::string_view ident, std::string_view passwd, const HttpRequest&)> on_pull_commands;

    // Called to authenticate every request.
    // Return true to allow, false to deny (401).
    // If nullptr, all requests are allowed.
    std::function<bool(std::string_view ident, std::string_view passwd, const HttpRequest&)> on_authenticate;

    // Called when a client packet cannot be decoded.
    // Defaults to stderr logging if nullptr.
    std::function<void(Error, std::string_view raw, const HttpRequest&)> on_decode_error;
};

// HTTP server — mirrors Go's HttpServer.
// Registers POST /v2/message and GET /v2/commands using cpp-httplib.
// start() is blocking; close() unblocks it gracefully.
class HttpServer {
public:
    static Result<HttpServer> create(HttpConfig cfg);

    HttpServer() = default;
    HttpServer(HttpServer&&) noexcept;
    HttpServer& operator=(HttpServer&&) noexcept;
    ~HttpServer();

    // Blocking. Returns when close() is called or listen fails.
    Error start();

    // Graceful shutdown (waits for in-flight requests to finish).
    Error close();

private:
    explicit HttpServer(HttpConfig cfg);

    struct Impl;
    std::unique_ptr<Impl> impl_;
};

#endif // LAYRZ_PROTOCOL_SERVERS

} // namespace layrz::protocol::servers

#endif // __LAYRZ_PROTOCOL_SERVERS_HTTP_HPP__
