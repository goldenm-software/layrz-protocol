#include "layrz_protocol/layrz_protocol.hpp"
#include "layrz_protocol/transport/tcp_comm.hpp"
#include <iostream>
#include <thread>
#include <chrono>

int main(int argc, char* argv[]) {
    if (argc < 5) {
        std::cerr << "Usage: " << argv[0] << " <host> <port> <ident> <password>\n";
        return 1;
    }
    std::string host     = argv[1];
    int         port     = std::stoi(argv[2]);
    std::string ident    = argv[3];
    std::string password = argv[4];

    using namespace layrz::protocol;
    using namespace layrz::protocol::transport;
    using namespace layrz::protocol::packets;

    TcpComm tcp(host, port, ident, password);

    tcp.set_callback([](AnyServerPacket pkt) {
        std::visit([](auto&& p) {
            using T = std::decay_t<decltype(p)>;
            if constexpr (std::is_same_v<T, ArPacket>) {
                std::cout << "Server error: " << p.reason << "\n";
            } else if constexpr (std::is_same_v<T, AcPacket>) {
                std::cout << "Received " << p.commands.size() << " command(s)\n";
            } else {
                std::cout << "Received server packet\n";
            }
        }, pkt);
    });

    std::cout << "Connecting to " << host << ":" << port << " ...\n";
    if (tcp.connect() != Error::Ok) {
        std::cerr << "Connection/auth failed\n";
        return 1;
    }
    std::cout << "Authenticated.\n";

    // Send a keepalive Pr
    AnyClientPacket pr = PrPacket{};
    if (tcp.send(pr) == Error::Ok) {
        std::cout << "Sent Pr (keepalive)\n";
    }

    std::this_thread::sleep_for(std::chrono::seconds(2));
    tcp.close();
    return 0;
}
