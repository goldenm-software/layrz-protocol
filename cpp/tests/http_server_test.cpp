#include <gtest/gtest.h>
#include "layrz_protocol/servers/http.hpp"
#include "layrz_protocol/layrz_protocol.hpp"
#include "fixtures/canonical_frames.hpp"

#include <httplib.h>

#include <atomic>
#include <chrono>
#include <thread>

using namespace layrz::protocol;
using namespace layrz::protocol::packets;
using namespace layrz::protocol::servers;
using namespace layrz::protocol::fixtures;

// ── Helpers ───────────────────────────────────────────────────────────────────

static int find_free_port() {
    int s = ::socket(AF_INET, SOCK_STREAM, 0);
    struct sockaddr_in addr{};
    addr.sin_family      = AF_INET;
    addr.sin_port        = 0;
    addr.sin_addr.s_addr = INADDR_ANY;
    ::bind(s, reinterpret_cast<sockaddr*>(&addr), sizeof(addr));
    socklen_t len = sizeof(addr);
    ::getsockname(s, reinterpret_cast<sockaddr*>(&addr), &len);
    int port = ntohs(addr.sin_port);
    ::close(s);
    return port;
}

static const std::string VALID_AUTH = "LayrzAuth device001;secret";

// ── Tests ─────────────────────────────────────────────────────────────────────

TEST(HttpServer, Create_RequiresOnNewPacket) {
    HttpConfig cfg;
    cfg.port = 30100;
    cfg.on_new_packet = nullptr;
    EXPECT_FALSE(HttpServer::create(std::move(cfg)).ok());
}

TEST(HttpServer, Create_InvalidPort) {
    HttpConfig cfg;
    cfg.on_new_packet = [](const AnyClientPacket&, const HttpRequest&) -> std::optional<AnyServerPacket> {
        return std::nullopt;
    };
    cfg.port = 0;
    EXPECT_FALSE(HttpServer::create(std::move(cfg)).ok());
}

TEST(HttpServer, PostMessage_ReturnsResponse) {
    int port = find_free_port();

    std::atomic<bool> got_packet{false};

    HttpConfig cfg;
    cfg.port = port;
    cfg.on_new_packet = [&](const AnyClientPacket& pkt, const HttpRequest&) -> std::optional<AnyServerPacket> {
        if (std::holds_alternative<PaPacket>(pkt)) {
            got_packet.store(true);
            return AsPacket{};
        }
        return std::nullopt;
    };

    auto srv_r = HttpServer::create(std::move(cfg));
    ASSERT_TRUE(srv_r.ok());
    HttpServer srv = std::move(srv_r.value);

    std::thread srv_thread([&srv]() { srv.start(); });
    std::this_thread::sleep_for(std::chrono::milliseconds(80));

    httplib::Client cli("127.0.0.1", port);
    httplib::Headers headers = {{"Authorization", VALID_AUTH}};
    auto res = cli.Post("/v2/message", headers, FRAME_PA, "text/plain");

    srv.close();
    srv_thread.join();

    ASSERT_TRUE(res);
    EXPECT_EQ(res->status, 200);
    EXPECT_EQ(res->body, FRAME_AS);
    EXPECT_TRUE(got_packet.load());
}

TEST(HttpServer, PostMessage_NoResponse_Returns204) {
    int port = find_free_port();

    HttpConfig cfg;
    cfg.port = port;
    cfg.on_new_packet = [](const AnyClientPacket&, const HttpRequest&) -> std::optional<AnyServerPacket> {
        return std::nullopt;
    };

    auto srv_r = HttpServer::create(std::move(cfg));
    ASSERT_TRUE(srv_r.ok());
    HttpServer srv = std::move(srv_r.value);

    std::thread srv_thread([&srv]() { srv.start(); });
    std::this_thread::sleep_for(std::chrono::milliseconds(80));

    httplib::Client cli("127.0.0.1", port);
    httplib::Headers headers = {{"Authorization", VALID_AUTH}};
    auto res = cli.Post("/v2/message", headers, FRAME_PA, "text/plain");

    srv.close();
    srv_thread.join();

    ASSERT_TRUE(res);
    EXPECT_EQ(res->status, 204);
}

TEST(HttpServer, MissingAuth_Returns401) {
    int port = find_free_port();

    HttpConfig cfg;
    cfg.port = port;
    cfg.on_new_packet = [](const AnyClientPacket&, const HttpRequest&) -> std::optional<AnyServerPacket> {
        return std::nullopt;
    };

    auto srv_r = HttpServer::create(std::move(cfg));
    ASSERT_TRUE(srv_r.ok());
    HttpServer srv = std::move(srv_r.value);

    std::thread srv_thread([&srv]() { srv.start(); });
    std::this_thread::sleep_for(std::chrono::milliseconds(80));

    httplib::Client cli("127.0.0.1", port);
    auto res = cli.Post("/v2/message", FRAME_PA, "text/plain"); // no auth

    srv.close();
    srv_thread.join();

    ASSERT_TRUE(res);
    EXPECT_EQ(res->status, 401);
}

TEST(HttpServer, OnAuthenticate_Denied_Returns401) {
    int port = find_free_port();

    HttpConfig cfg;
    cfg.port = port;
    cfg.on_new_packet = [](const AnyClientPacket&, const HttpRequest&) -> std::optional<AnyServerPacket> {
        return std::nullopt;
    };
    cfg.on_authenticate = [](std::string_view, std::string_view, const HttpRequest&) {
        return false; // deny all
    };

    auto srv_r = HttpServer::create(std::move(cfg));
    ASSERT_TRUE(srv_r.ok());
    HttpServer srv = std::move(srv_r.value);

    std::thread srv_thread([&srv]() { srv.start(); });
    std::this_thread::sleep_for(std::chrono::milliseconds(80));

    httplib::Client cli("127.0.0.1", port);
    httplib::Headers headers = {{"Authorization", VALID_AUTH}};
    auto res = cli.Post("/v2/message", headers, FRAME_PA, "text/plain");

    srv.close();
    srv_thread.join();

    ASSERT_TRUE(res);
    EXPECT_EQ(res->status, 401);
}

TEST(HttpServer, MalformedBody_Returns400_CallsDecodeError) {
    int port = find_free_port();

    std::atomic<bool> decode_error_called{false};

    HttpConfig cfg;
    cfg.port = port;
    cfg.on_new_packet = [](const AnyClientPacket&, const HttpRequest&) -> std::optional<AnyServerPacket> {
        return std::nullopt;
    };
    cfg.on_decode_error = [&](Error, std::string_view, const HttpRequest&) {
        decode_error_called.store(true);
    };

    auto srv_r = HttpServer::create(std::move(cfg));
    ASSERT_TRUE(srv_r.ok());
    HttpServer srv = std::move(srv_r.value);

    std::thread srv_thread([&srv]() { srv.start(); });
    std::this_thread::sleep_for(std::chrono::milliseconds(80));

    httplib::Client cli("127.0.0.1", port);
    httplib::Headers headers = {{"Authorization", VALID_AUTH}};
    auto res = cli.Post("/v2/message", headers, "not a valid packet", "text/plain");

    srv.close();
    srv_thread.join();

    ASSERT_TRUE(res);
    EXPECT_EQ(res->status, 400);
    EXPECT_TRUE(decode_error_called.load());
}

TEST(HttpServer, GetCommands_NullHandler_Returns204) {
    int port = find_free_port();

    HttpConfig cfg;
    cfg.port = port;
    cfg.on_new_packet = [](const AnyClientPacket&, const HttpRequest&) -> std::optional<AnyServerPacket> {
        return std::nullopt;
    };
    // on_pull_commands not set

    auto srv_r = HttpServer::create(std::move(cfg));
    ASSERT_TRUE(srv_r.ok());
    HttpServer srv = std::move(srv_r.value);

    std::thread srv_thread([&srv]() { srv.start(); });
    std::this_thread::sleep_for(std::chrono::milliseconds(80));

    httplib::Client cli("127.0.0.1", port);
    httplib::Headers headers = {{"Authorization", VALID_AUTH}};
    auto res = cli.Get("/v2/commands", headers);

    srv.close();
    srv_thread.join();

    ASSERT_TRUE(res);
    EXPECT_EQ(res->status, 204);
}

TEST(HttpServer, GetCommands_WithHandler_Returns200) {
    int port = find_free_port();

    HttpConfig cfg;
    cfg.port = port;
    cfg.on_new_packet = [](const AnyClientPacket&, const HttpRequest&) -> std::optional<AnyServerPacket> {
        return std::nullopt;
    };
    cfg.on_pull_commands = [](std::string_view, std::string_view, const HttpRequest&) -> std::optional<AnyServerPacket> {
        return AoPacket{1700000000LL};
    };

    auto srv_r = HttpServer::create(std::move(cfg));
    ASSERT_TRUE(srv_r.ok());
    HttpServer srv = std::move(srv_r.value);

    std::thread srv_thread([&srv]() { srv.start(); });
    std::this_thread::sleep_for(std::chrono::milliseconds(80));

    httplib::Client cli("127.0.0.1", port);
    httplib::Headers headers = {{"Authorization", VALID_AUTH}};
    auto res = cli.Get("/v2/commands", headers);

    srv.close();
    srv_thread.join();

    ASSERT_TRUE(res);
    EXPECT_EQ(res->status, 200);
    EXPECT_EQ(res->body, FRAME_AO);
}

TEST(HttpServer, WrongMethod_Returns405) {
    int port = find_free_port();

    HttpConfig cfg;
    cfg.port = port;
    cfg.on_new_packet = [](const AnyClientPacket&, const HttpRequest&) -> std::optional<AnyServerPacket> {
        return std::nullopt;
    };

    auto srv_r = HttpServer::create(std::move(cfg));
    ASSERT_TRUE(srv_r.ok());
    HttpServer srv = std::move(srv_r.value);

    std::thread srv_thread([&srv]() { srv.start(); });
    std::this_thread::sleep_for(std::chrono::milliseconds(80));

    httplib::Client cli("127.0.0.1", port);
    httplib::Headers headers = {{"Authorization", VALID_AUTH}};
    // GET on /v2/message (should be POST)
    auto res = cli.Get("/v2/message", headers);

    srv.close();
    srv_thread.join();

    ASSERT_TRUE(res);
    EXPECT_EQ(res->status, 405);
}
