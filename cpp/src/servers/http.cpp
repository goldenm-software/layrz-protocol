#include "layrz_protocol/servers/http.hpp"

#ifdef LAYRZ_PROTOCOL_SERVERS

#include "layrz_protocol/parser.hpp"
#include <httplib.h>
#include <iostream>
#include <optional>

namespace layrz::protocol::servers {

// ── Pimpl ─────────────────────────────────────────────────────────────────────

struct HttpServer::Impl {
    HttpConfig      cfg;
    httplib::Server srv;
};

// ── Helpers ───────────────────────────────────────────────────────────────────

static bool parse_layrz_auth(const std::string& header,
                              std::string& ident, std::string& passwd) {
    const std::string prefix = "LayrzAuth ";
    if (header.size() <= prefix.size()) return false;
    if (header.substr(0, prefix.size()) != prefix) return false;
    std::string rest = header.substr(prefix.size());
    auto sep = rest.find(';');
    if (sep == std::string::npos) return false;
    ident  = rest.substr(0, sep);
    passwd = rest.substr(sep + 1);
    return !ident.empty();
}

static HttpRequest make_request(const httplib::Request& req) {
    HttpRequest r;
    r.method      = req.method;
    r.path        = req.path;
    r.remote_addr = req.remote_addr;
    r.get_header  = [&req](std::string_view name) -> std::string {
        auto it = req.headers.find(std::string(name));
        if (it == req.headers.end()) return {};
        return it->second;
    };
    return r;
}

static void send_packet_response(httplib::Response& res, const AnyServerPacket& pkt) {
    auto encoded = parse_server_packet_to_string(pkt);
    if (!encoded.ok()) {
        res.status = 500;
        res.set_content("internal server error", "text/plain");
        return;
    }
    res.status = 200;
    res.set_content(encoded.value, "text/plain; charset=utf-8");
}

// ── Construction ──────────────────────────────────────────────────────────────

HttpServer::HttpServer(HttpConfig cfg) : impl_(std::make_unique<Impl>()) {
    impl_->cfg = std::move(cfg);
}

Result<HttpServer> HttpServer::create(HttpConfig cfg) {
    if (!cfg.on_new_packet)
        return Result<HttpServer>::fail(Error::ParseError);
    if (cfg.port <= 0 || cfg.port >= 65535)
        return Result<HttpServer>::fail(Error::ParseError);

    if (!cfg.on_decode_error) {
        cfg.on_decode_error = [](Error, std::string_view raw, const HttpRequest& req) {
            std::cerr << "[HttpServer] decode error from " << req.remote_addr
                      << ": " << raw << "\n";
        };
    }

    return Result<HttpServer>::success(HttpServer(std::move(cfg)));
}

HttpServer::HttpServer(HttpServer&&) noexcept = default;
HttpServer& HttpServer::operator=(HttpServer&&) noexcept = default;
HttpServer::~HttpServer() { close(); }

// ── Lifecycle ─────────────────────────────────────────────────────────────────

Error HttpServer::start() {
    auto& srv = impl_->srv;
    auto& cfg = impl_->cfg;

    srv.set_payload_max_length(1 << 20); // 1 MiB

    // POST /v2/message
    srv.Post("/v2/message", [&cfg](const httplib::Request& req, httplib::Response& res) {
        std::string ident, passwd;
        if (!parse_layrz_auth(req.get_header_value("Authorization"), ident, passwd)) {
            res.status = 401;
            res.set_content("unauthorized", "text/plain");
            return;
        }

        auto hr = make_request(req);
        if (cfg.on_authenticate && !cfg.on_authenticate(ident, passwd, hr)) {
            res.status = 401;
            res.set_content("unauthorized", "text/plain");
            return;
        }

        auto decoded = handle_client_input(req.body);
        if (!decoded.ok()) {
            if (cfg.on_decode_error) cfg.on_decode_error(decoded.error, req.body, hr);
            res.status = 400;
            res.set_content("invalid packet", "text/plain");
            return;
        }

        std::optional<AnyServerPacket> response;
        try {
            response = cfg.on_new_packet(decoded.value, hr);
        } catch (...) {
            res.status = 500;
            res.set_content("internal server error", "text/plain");
            return;
        }

        if (!response.has_value()) {
            res.status = 204;
            return;
        }
        send_packet_response(res, *response);
    });

    // GET /v2/commands
    srv.Get("/v2/commands", [&cfg](const httplib::Request& req, httplib::Response& res) {
        std::string ident, passwd;
        if (!parse_layrz_auth(req.get_header_value("Authorization"), ident, passwd)) {
            res.status = 401;
            res.set_content("unauthorized", "text/plain");
            return;
        }

        auto hr = make_request(req);
        if (cfg.on_authenticate && !cfg.on_authenticate(ident, passwd, hr)) {
            res.status = 401;
            res.set_content("unauthorized", "text/plain");
            return;
        }

        if (!cfg.on_pull_commands) {
            res.status = 204;
            return;
        }

        std::optional<AnyServerPacket> response;
        try {
            response = cfg.on_pull_commands(ident, passwd, hr);
        } catch (...) {
            res.status = 500;
            res.set_content("internal server error", "text/plain");
            return;
        }

        if (!response.has_value()) {
            res.status = 204;
            return;
        }
        send_packet_response(res, *response);
    });

    // Method-not-allowed for the complementary methods
    auto method_not_allowed = [](const httplib::Request&, httplib::Response& res) {
        res.status = 405;
        res.set_content("method not allowed", "text/plain");
    };
    srv.Get("/v2/message",  method_not_allowed);
    srv.Put("/v2/message",  method_not_allowed);
    srv.Delete("/v2/message", method_not_allowed);
    srv.Post("/v2/commands", method_not_allowed);
    srv.Put("/v2/commands",  method_not_allowed);
    srv.Delete("/v2/commands", method_not_allowed);

    if (!srv.listen("0.0.0.0", cfg.port)) {
        return Error::ServerError;
    }
    return Error::Ok;
}

Error HttpServer::close() {
    if (impl_) impl_->srv.stop();
    return Error::Ok;
}

} // namespace layrz::protocol::servers

#endif // LAYRZ_PROTOCOL_SERVERS
