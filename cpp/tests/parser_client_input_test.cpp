#include <gtest/gtest.h>
#include "layrz_protocol/layrz_protocol.hpp"
#include "layrz_protocol/crc.hpp"
#include "fixtures/canonical_frames.hpp"

using namespace layrz::protocol;
using namespace layrz::protocol::packets;
using namespace layrz::protocol::fixtures;

static std::string complete_frame(const std::string& open, const std::string& close,
                                   const std::string& payload) {
    return open + payload + compute_crc_str(payload) + close;
}

// ── handle_client_input dispatch ──────────────────────────────────────────────

TEST(HandleClientInput, Pa) {
    auto r = handle_client_input(FRAME_PA);
    ASSERT_TRUE(r.ok());
    ASSERT_TRUE(std::holds_alternative<PaPacket>(r.value));
    EXPECT_EQ(std::get<PaPacket>(r.value).ident, "123456789012345");
}

TEST(HandleClientInput, Pc) {
    auto r = handle_client_input(FRAME_PC);
    ASSERT_TRUE(r.ok());
    EXPECT_TRUE(std::holds_alternative<PcPacket>(r.value));
}

TEST(HandleClientInput, Pr) {
    auto r = handle_client_input(FRAME_PR);
    ASSERT_TRUE(r.ok());
    EXPECT_TRUE(std::holds_alternative<PrPacket>(r.value));
}

TEST(HandleClientInput, Ts) {
    auto r = handle_client_input(FRAME_TS);
    ASSERT_TRUE(r.ok());
    EXPECT_TRUE(std::holds_alternative<TsPacket>(r.value));
}

TEST(HandleClientInput, Te) {
    auto r = handle_client_input(FRAME_TE);
    ASSERT_TRUE(r.ok());
    EXPECT_TRUE(std::holds_alternative<TePacket>(r.value));
}

TEST(HandleClientInput, Im) {
    auto r = handle_client_input(FRAME_IM);
    ASSERT_TRUE(r.ok());
    EXPECT_TRUE(std::holds_alternative<ImPacket>(r.value));
}

TEST(HandleClientInput, Pi) {
    std::string full = complete_frame("<Pi>", "</Pi>", "testident;1;1;1;1;1;1;1;");
    auto r = handle_client_input(full);
    ASSERT_TRUE(r.ok());
    EXPECT_TRUE(std::holds_alternative<PiPacket>(r.value));
}

TEST(HandleClientInput, Pd) {
    std::string payload =
        "0;10.0;10.0;10.0;10.0;10.0;5;1.0;"
        "test.str:Hola mundo,test.int:1,test.double:1.0,test.bool:true;";
    std::string full = complete_frame("<Pd>", "</Pd>", payload);
    auto r = handle_client_input(full);
    ASSERT_TRUE(r.ok());
    EXPECT_TRUE(std::holds_alternative<PdPacket>(r.value));
}

TEST(HandleClientInput, Ps) {
    std::string payload =
        "0;net_wifi_ssid:AWESOME WIFI,net_wifi_pass:dictadormarico69,"
        "net_wifi_sec:WPA2,static.lat:-15.5,static.lng:15.5;";
    std::string full = complete_frame("<Ps>", "</Ps>", payload);
    auto r = handle_client_input(full);
    ASSERT_TRUE(r.ok());
    EXPECT_TRUE(std::holds_alternative<PsPacket>(r.value));
}

TEST(HandleClientInput, UnknownTag) {
    auto r = handle_client_input("<Xx>garbage</Xx>");
    EXPECT_FALSE(r.ok());
}

TEST(HandleClientInput, MalformedBody) {
    auto r = handle_client_input("<Pa>not-valid-content</Pa>");
    EXPECT_FALSE(r.ok());
}

// ── split_client_frames ───────────────────────────────────────────────────────

TEST(SplitClientFrames, TwoFrames) {
    std::string buf = FRAME_PA + "\n" + FRAME_PR + "\n";
    auto frames = split_client_frames(buf);
    ASSERT_EQ(frames.size(), 2u);
    EXPECT_EQ(frames[0], FRAME_PA);
    EXPECT_EQ(frames[1], FRAME_PR);
}

TEST(SplitClientFrames, FramesWithoutNewline) {
    // Frames may be concatenated without newlines too
    std::string buf = FRAME_PA + FRAME_PR;
    auto frames = split_client_frames(buf);
    ASSERT_EQ(frames.size(), 2u);
}

TEST(SplitClientFrames, SingleFrame) {
    auto frames = split_client_frames(FRAME_PA);
    ASSERT_EQ(frames.size(), 1u);
    EXPECT_EQ(frames[0], FRAME_PA);
}

TEST(SplitClientFrames, TrailingPartial) {
    // Partial frame at end must be dropped silently
    std::string buf = FRAME_PA + "\n<Pa>incomplete";
    auto frames = split_client_frames(buf);
    ASSERT_EQ(frames.size(), 1u);
    EXPECT_EQ(frames[0], FRAME_PA);
}

TEST(SplitClientFrames, Empty) {
    auto frames = split_client_frames("");
    EXPECT_TRUE(frames.empty());
}

TEST(SplitClientFrames, DoesNotSplitServerFrames) {
    // Server output frames (<Ab>, <As>, etc.) must NOT be split out
    std::string buf = FRAME_AS + "\n" + FRAME_PA;
    auto frames = split_client_frames(buf);
    // <As> is not a client frame tag — only <Pa> should be found
    ASSERT_EQ(frames.size(), 1u);
    EXPECT_EQ(frames[0], FRAME_PA);
}

// ── parse_server_packet_to_string ────────────────────────────────────────────

TEST(Parser, ParseServerPacketToString_As) {
    AnyServerPacket pkt = AsPacket{};
    auto r = parse_server_packet_to_string(pkt);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value, FRAME_AS);
}

TEST(Parser, ParseServerPacketToString_Ao) {
    AnyServerPacket pkt = AoPacket{1700000000LL};
    auto r = parse_server_packet_to_string(pkt);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value, FRAME_AO);
}
