#include "layrz_protocol/servers/http.hpp"
#include "layrz_protocol/layrz_protocol.hpp"

#include <chrono>
#include <csignal>
#include <iostream>
#include <variant>

using namespace layrz::protocol;
using namespace layrz::protocol::packets;
using namespace layrz::protocol::servers;

static HttpServer* g_server = nullptr;

static void on_signal(int) {
    if (g_server) g_server->close();
}

int main() {
    HttpConfig cfg;
    cfg.port = 8080;

    cfg.on_authenticate = [](std::string_view ident, std::string_view passwd, const HttpRequest&) -> bool {
        return ident == "device001" && passwd == "secret";
    };

    cfg.on_new_packet = [](const AnyClientPacket& pkt, const HttpRequest& req) -> std::optional<AnyServerPacket> {
        std::cout << "Packet from " << req.remote_addr << "\n";
        if (std::holds_alternative<PaPacket>(pkt)) {
            std::cout << "  Auth: " << std::get<PaPacket>(pkt).ident << "\n";
            return AsPacket{};
        }
        auto now = std::chrono::duration_cast<std::chrono::seconds>(
            std::chrono::system_clock::now().time_since_epoch()).count();
        return AoPacket{static_cast<long long>(now)};
    };

    cfg.on_pull_commands = [](std::string_view ident, std::string_view, const HttpRequest&) -> std::optional<AnyServerPacket> {
        std::cout << "Pull commands from " << ident << "\n";
        return std::nullopt;
    };

    auto srv_r = HttpServer::create(std::move(cfg));
    if (!srv_r.ok()) {
        std::cerr << "Failed to create HttpServer\n";
        return 1;
    }
    HttpServer srv = std::move(srv_r.value);
    g_server = &srv;

    std::signal(SIGINT,  on_signal);
    std::signal(SIGTERM, on_signal);

    std::cout << "HTTP server listening on :8080\n";
    srv.start();
    std::cout << "HTTP server stopped\n";
    return 0;
}
