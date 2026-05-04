#include <gtest/gtest.h>
#include "layrz_protocol/servers/tcp.hpp"
#include "layrz_protocol/layrz_protocol.hpp"
#include "fixtures/canonical_frames.hpp"

#include <arpa/inet.h>
#include <atomic>
#include <chrono>
#include <cstring>
#include <netinet/in.h>
#include <sys/socket.h>
#include <thread>
#include <unistd.h>

using namespace layrz::protocol;
using namespace layrz::protocol::packets;
using namespace layrz::protocol::servers;
using namespace layrz::protocol::fixtures;

// ── Helpers ───────────────────────────────────────────────────────────────────

static int find_free_port() {
    int s = ::socket(AF_INET, SOCK_STREAM, 0);
    struct sockaddr_in addr{};
    addr.sin_family = AF_INET;
    addr.sin_port   = 0;
    addr.sin_addr.s_addr = INADDR_ANY;
    ::bind(s, reinterpret_cast<sockaddr*>(&addr), sizeof(addr));
    socklen_t len = sizeof(addr);
    ::getsockname(s, reinterpret_cast<sockaddr*>(&addr), &len);
    int port = ntohs(addr.sin_port);
    ::close(s);
    return port;
}

static int connect_to(int port) {
    int s = ::socket(AF_INET, SOCK_STREAM, 0);
    struct sockaddr_in addr{};
    addr.sin_family      = AF_INET;
    addr.sin_port        = htons(static_cast<uint16_t>(port));
    addr.sin_addr.s_addr = htonl(INADDR_LOOPBACK);
    if (::connect(s, reinterpret_cast<sockaddr*>(&addr), sizeof(addr)) < 0) {
        ::close(s);
        return -1;
    }
    return s;
}

static std::string recv_all(int s, size_t expected_min, int timeout_ms = 2000) {
    std::string buf;
    buf.reserve(expected_min);
    auto deadline = std::chrono::steady_clock::now() + std::chrono::milliseconds(timeout_ms);
    char tmp[1024];
    while (buf.size() < expected_min &&
           std::chrono::steady_clock::now() < deadline) {
        fd_set fds; FD_ZERO(&fds); FD_SET(s, &fds);
        struct timeval tv{0, 20000}; // 20ms
        if (::select(s + 1, &fds, nullptr, nullptr, &tv) <= 0) continue;
        ssize_t n = ::recv(s, tmp, sizeof(tmp), 0);
        if (n <= 0) break;
        buf.append(tmp, static_cast<size_t>(n));
    }
    return buf;
}

// ── Tests ─────────────────────────────────────────────────────────────────────

TEST(TcpServer, Create_RequiresOnNewPacket) {
    TcpConfig cfg;
    cfg.port = 30001;
    cfg.on_new_packet = nullptr;
    auto r = TcpServer::create(std::move(cfg));
    EXPECT_FALSE(r.ok());
}

TEST(TcpServer, Create_InvalidPort) {
    TcpConfig cfg;
    cfg.port = 0;
    cfg.on_new_packet = [](const AnyClientPacket&, TcpConnection&) -> std::optional<AnyServerPacket> {
        return std::nullopt;
    };
    EXPECT_FALSE(TcpServer::create(std::move(cfg)).ok());

    cfg.port = 65535;
    cfg.on_new_packet = [](const AnyClientPacket&, TcpConnection&) -> std::optional<AnyServerPacket> {
        return std::nullopt;
    };
    EXPECT_FALSE(TcpServer::create(std::move(cfg)).ok());
}

TEST(TcpServer, BasicRoundtrip) {
    int port = find_free_port();

    std::atomic<bool> got_packet{false};

    TcpConfig cfg;
    cfg.port = port;
    cfg.on_new_packet = [&](const AnyClientPacket& pkt, TcpConnection&) -> std::optional<AnyServerPacket> {
        if (std::holds_alternative<PaPacket>(pkt)) {
            got_packet.store(true);
            return AsPacket{};
        }
        return std::nullopt;
    };

    auto srv_r = TcpServer::create(std::move(cfg));
    ASSERT_TRUE(srv_r.ok());
    TcpServer srv = std::move(srv_r.value);

    std::thread srv_thread([&srv]() { srv.start(); });

    // Give the server a moment to bind
    std::this_thread::sleep_for(std::chrono::milliseconds(50));

    int s = connect_to(port);
    ASSERT_GE(s, 0) << "Failed to connect to server";

    // Send a Pa frame with newline terminator (mimics TcpComm client)
    std::string frame = FRAME_PA + "\n";
    ::send(s, frame.data(), frame.size(), MSG_NOSIGNAL);

    std::string response = recv_all(s, FRAME_AS.size());
    ::close(s);

    srv.close();
    srv_thread.join();

    EXPECT_TRUE(got_packet.load());
    EXPECT_EQ(response, FRAME_AS);
}

TEST(TcpServer, MalformedFrame_CallsDecodeError) {
    int port = find_free_port();

    std::atomic<bool> decode_error_called{false};

    TcpConfig cfg;
    cfg.port = port;
    cfg.on_new_packet = [](const AnyClientPacket&, TcpConnection&) -> std::optional<AnyServerPacket> {
        return std::nullopt;
    };
    cfg.on_decode_error = [&](Error, std::string_view, TcpConnection&) {
        decode_error_called.store(true);
    };

    auto srv_r = TcpServer::create(std::move(cfg));
    ASSERT_TRUE(srv_r.ok());
    TcpServer srv = std::move(srv_r.value);

    std::thread srv_thread([&srv]() { srv.start(); });
    std::this_thread::sleep_for(std::chrono::milliseconds(50));

    int s = connect_to(port);
    ASSERT_GE(s, 0);

    // Send a Pa frame with bad CRC so from_packet returns error
    std::string bad = "<Pa>ident;badcrc;0000</Pa>\n";
    ::send(s, bad.data(), bad.size(), MSG_NOSIGNAL);

    std::this_thread::sleep_for(std::chrono::milliseconds(100));
    ::close(s);

    srv.close();
    srv_thread.join();

    EXPECT_TRUE(decode_error_called.load());
}

TEST(TcpServer, TwoFramesInOneWrite) {
    int port = find_free_port();

    std::atomic<int> packet_count{0};

    TcpConfig cfg;
    cfg.port = port;
    cfg.on_new_packet = [&](const AnyClientPacket&, TcpConnection&) -> std::optional<AnyServerPacket> {
        packet_count.fetch_add(1);
        return std::nullopt;
    };

    auto srv_r = TcpServer::create(std::move(cfg));
    ASSERT_TRUE(srv_r.ok());
    TcpServer srv = std::move(srv_r.value);

    std::thread srv_thread([&srv]() { srv.start(); });
    std::this_thread::sleep_for(std::chrono::milliseconds(50));

    int s = connect_to(port);
    ASSERT_GE(s, 0);

    // Send two complete frames concatenated with a single write + newline
    std::string two = FRAME_PA + "\n" + FRAME_PR + "\n";
    ::send(s, two.data(), two.size(), MSG_NOSIGNAL);

    std::this_thread::sleep_for(std::chrono::milliseconds(100));
    ::close(s);

    srv.close();
    srv_thread.join();

    EXPECT_EQ(packet_count.load(), 2);
}

TEST(TcpServer, CloseUnblocksStart) {
    int port = find_free_port();

    TcpConfig cfg;
    cfg.port = port;
    cfg.on_new_packet = [](const AnyClientPacket&, TcpConnection&) -> std::optional<AnyServerPacket> {
        return std::nullopt;
    };

    auto srv_r = TcpServer::create(std::move(cfg));
    ASSERT_TRUE(srv_r.ok());
    TcpServer srv = std::move(srv_r.value);

    std::thread srv_thread([&srv]() { srv.start(); });
    std::this_thread::sleep_for(std::chrono::milliseconds(50));

    srv.close();
    srv_thread.join(); // must not hang
    SUCCEED();
}
