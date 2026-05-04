#include "layrz_protocol/servers/tcp.hpp"

#ifdef LAYRZ_PROTOCOL_SERVERS

#include "layrz_protocol/parser.hpp"
#include <arpa/inet.h>
#include <atomic>
#include <cerrno>
#include <cstring>
#include <iostream>
#include <mutex>
#include <netinet/in.h>
#include <optional>
#include <sys/socket.h>
#include <thread>
#include <unistd.h>
#include <vector>

namespace layrz::protocol::servers {

// ── Pimpl ─────────────────────────────────────────────────────────────────────

struct TcpServer::Impl {
    TcpConfig          cfg;
    int                listen_fd = -1;
    std::atomic<bool>  stop{false};
    std::mutex         threads_mu;
    std::vector<std::thread> conn_threads;
};

// ── Construction ──────────────────────────────────────────────────────────────

TcpServer::TcpServer(std::unique_ptr<Impl> impl) : impl_(std::move(impl)) {}

Result<TcpServer> TcpServer::create(TcpConfig cfg) {
    if (!cfg.on_new_packet)
        return Result<TcpServer>::fail(Error::ParseError);
    if (cfg.port <= 0 || cfg.port >= 65535)
        return Result<TcpServer>::fail(Error::ParseError);

    if (!cfg.on_decode_error) {
        cfg.on_decode_error = [](Error, std::string_view raw, TcpConnection& conn) {
            std::cerr << "[TcpServer] decode error from " << conn.addr
                      << ": " << raw << "\n";
        };
    }

    auto impl = std::make_unique<Impl>();
    impl->cfg = std::move(cfg);
    return Result<TcpServer>::success(TcpServer(std::move(impl)));
}

TcpServer::TcpServer(TcpServer&&) noexcept = default;
TcpServer& TcpServer::operator=(TcpServer&&) noexcept = default;

TcpServer::~TcpServer() {
    if (impl_) close();
}

// ── Lifecycle ─────────────────────────────────────────────────────────────────

Error TcpServer::start() {
    auto& d = *impl_;
    d.listen_fd = ::socket(AF_INET6, SOCK_STREAM, 0);
    bool using_ipv4_fallback = false;

    if (d.listen_fd < 0) {
        d.listen_fd = ::socket(AF_INET, SOCK_STREAM, 0);
        if (d.listen_fd < 0) return Error::ServerError;
        using_ipv4_fallback = true;
    }

    int opt = 1;
    ::setsockopt(d.listen_fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));

    if (!using_ipv4_fallback) {
        int v6only = 0;
        ::setsockopt(d.listen_fd, IPPROTO_IPV6, IPV6_V6ONLY, &v6only, sizeof(v6only));

        struct sockaddr_in6 addr{};
        addr.sin6_family = AF_INET6;
        addr.sin6_port   = htons(static_cast<uint16_t>(d.cfg.port));
        addr.sin6_addr   = in6addr_any;

        if (::bind(d.listen_fd, reinterpret_cast<sockaddr*>(&addr), sizeof(addr)) < 0) {
            ::close(d.listen_fd); d.listen_fd = -1;
            return Error::ServerError;
        }
    } else {
        struct sockaddr_in addr{};
        addr.sin_family      = AF_INET;
        addr.sin_port        = htons(static_cast<uint16_t>(d.cfg.port));
        addr.sin_addr.s_addr = INADDR_ANY;

        if (::bind(d.listen_fd, reinterpret_cast<sockaddr*>(&addr), sizeof(addr)) < 0) {
            ::close(d.listen_fd); d.listen_fd = -1;
            return Error::ServerError;
        }
    }

    if (::listen(d.listen_fd, SOMAXCONN) < 0) {
        ::close(d.listen_fd); d.listen_fd = -1;
        return Error::ServerError;
    }

    d.stop.store(false);

    while (!d.stop.load()) {
        struct sockaddr_storage peer_addr{};
        socklen_t peer_len = sizeof(peer_addr);
        int client_fd = ::accept(d.listen_fd,
                                 reinterpret_cast<sockaddr*>(&peer_addr),
                                 &peer_len);
        if (client_fd < 0) {
            if (d.stop.load()) break;
            if (errno == EINTR) continue;
            break;
        }

        char buf[INET6_ADDRSTRLEN] = {};
        int  remote_port = 0;
        if (peer_addr.ss_family == AF_INET6) {
            auto* s = reinterpret_cast<struct sockaddr_in6*>(&peer_addr);
            ::inet_ntop(AF_INET6, &s->sin6_addr, buf, sizeof(buf));
            remote_port = ntohs(s->sin6_port);
        } else {
            auto* s = reinterpret_cast<struct sockaddr_in*>(&peer_addr);
            ::inet_ntop(AF_INET, &s->sin_addr, buf, sizeof(buf));
            remote_port = ntohs(s->sin_port);
        }

        TcpConnection conn;
        conn.fd   = client_fd;
        conn.addr = std::string(buf) + ":" + std::to_string(remote_port);

        std::lock_guard<std::mutex> lk(d.threads_mu);
        d.conn_threads.emplace_back([this, client_fd, conn = std::move(conn)]() mutable {
            handle_connection(client_fd, conn);
        });
    }

    std::lock_guard<std::mutex> lk(d.threads_mu);
    for (auto& t : d.conn_threads) {
        if (t.joinable()) t.join();
    }
    d.conn_threads.clear();

    return Error::Ok;
}

Error TcpServer::close() {
    if (!impl_) return Error::Ok;
    auto& d = *impl_;
    d.stop.store(true);
    if (d.listen_fd >= 0) {
        ::shutdown(d.listen_fd, SHUT_RDWR);
        ::close(d.listen_fd);
        d.listen_fd = -1;
    }
    return Error::Ok;
}

// ── Per-connection handler ─────────────────────────────────────────────────────

void TcpServer::handle_connection(int client_fd, TcpConnection conn) {
    auto& d = *impl_;
    std::string accumulator;
    accumulator.reserve(4096);
    char buf[1024];

    while (!d.stop.load()) {
        ssize_t n = ::recv(client_fd, buf, sizeof(buf), 0);
        if (n <= 0) break;

        accumulator.append(buf, static_cast<size_t>(n));

        if (accumulator.find('\n') == std::string::npos) continue;

        auto frames = split_client_frames(accumulator);
        accumulator.clear();

        for (const auto& frame : frames) {
            auto decoded = handle_client_input(frame);
            if (!decoded.ok()) {
                if (d.cfg.on_decode_error)
                    d.cfg.on_decode_error(decoded.error, frame, conn);
                continue;
            }

            std::optional<AnyServerPacket> response;
            try {
                response = d.cfg.on_new_packet(decoded.value, conn);
            } catch (...) {
                std::cerr << "[TcpServer] on_new_packet threw from " << conn.addr << "\n";
                continue;
            }

            if (response.has_value()) {
                auto encoded = parse_server_packet_to_string(*response);
                if (encoded.ok()) {
                    ::send(client_fd, encoded.value.data(), encoded.value.size(), MSG_NOSIGNAL);
                }
            }
        }
    }

    ::close(client_fd);
}

} // namespace layrz::protocol::servers

#endif // LAYRZ_PROTOCOL_SERVERS
