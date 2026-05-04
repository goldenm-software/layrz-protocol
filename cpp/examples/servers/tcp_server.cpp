#include "layrz_protocol/servers/tcp.hpp"
#include "layrz_protocol/layrz_protocol.hpp"

#include <chrono>
#include <csignal>
#include <ctime>
#include <iostream>
#include <variant>

using namespace layrz::protocol;
using namespace layrz::protocol::packets;
using namespace layrz::protocol::servers;

static TcpServer* g_server = nullptr;

static void on_signal(int) {
    if (g_server) g_server->close();
}

int main() {
    TcpConfig cfg;
    cfg.port = 12345;
    cfg.on_new_packet = [](const AnyClientPacket& pkt, TcpConnection& conn) -> std::optional<AnyServerPacket> {
        std::cout << "Packet from " << conn.addr << "\n";
        if (std::holds_alternative<PaPacket>(pkt)) {
            std::cout << "  Auth: " << std::get<PaPacket>(pkt).ident << "\n";
            return AsPacket{};
        }
        auto now = std::chrono::duration_cast<std::chrono::seconds>(
            std::chrono::system_clock::now().time_since_epoch()).count();
        return AoPacket{static_cast<long long>(now)};
    };

    auto srv_r = TcpServer::create(std::move(cfg));
    if (!srv_r.ok()) {
        std::cerr << "Failed to create TcpServer\n";
        return 1;
    }
    TcpServer srv = std::move(srv_r.value);
    g_server = &srv;

    std::signal(SIGINT,  on_signal);
    std::signal(SIGTERM, on_signal);

    std::cout << "TCP server listening on :12345\n";
    srv.start();
    std::cout << "TCP server stopped\n";
    return 0;
}
