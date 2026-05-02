#include <gtest/gtest.h>
#include "layrz_protocol/layrz_protocol.hpp"
#include "layrz_protocol/crc.hpp"
#include "fixtures/canonical_frames.hpp"
#include "fixtures/canonical_frames_generated.hpp"

using namespace layrz::protocol;
using namespace layrz::protocol::packets;
using namespace layrz::protocol::fixtures;

// Helper: append CRC to partial payload and wrap with tags
static std::string complete_frame(const std::string& open, const std::string& close,
                                   const std::string& payload) {
    return open + payload + compute_crc_str(payload) + close;
}

// ── Server packets (decode → re-encode) ──────────────────────────────────────

TEST(RoundTrip, As) {
    auto r = AsPacket::from_packet(FRAME_AS);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.to_packet(), FRAME_AS);
}

TEST(RoundTrip, Au) {
    auto r = AuPacket::from_packet(FRAME_AU);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.to_packet(), FRAME_AU);
}

TEST(RoundTrip, Ao) {
    auto r = AoPacket::from_packet(FRAME_AO);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.timestamp, 1700000000LL);
    EXPECT_EQ(r.value.to_packet(), FRAME_AO);
}

TEST(RoundTrip, Ar) {
    auto r = ArPacket::from_packet(FRAME_AR);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.reason, "Unknown reason");
    EXPECT_EQ(r.value.to_packet(), FRAME_AR);
}

TEST(RoundTrip, Ab) {
    auto r = AbPacket::from_packet(FRAME_AB);
    ASSERT_TRUE(r.ok());
    ASSERT_EQ(r.value.devices.size(), 2u);
    EXPECT_EQ(r.value.devices[0].mac_address, "12:34:56:78:90:AB");
    EXPECT_EQ(r.value.devices[0].model,       "GENERIC");
    EXPECT_EQ(r.value.devices[1].mac_address, "BC:09:87:65:43:21");
    EXPECT_EQ(r.value.to_packet(), FRAME_AB);
}

// ── Client packets (encode → decode → re-encode) ─────────────────────────────

TEST(RoundTrip, Pr) {
    auto r = PrPacket::from_packet(FRAME_PR);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.to_packet(), FRAME_PR);
}

TEST(RoundTrip, Pa) {
    auto r = PaPacket::from_packet(FRAME_PA);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.ident,    "123456789012345");
    EXPECT_EQ(r.value.password, "mypassword");
    EXPECT_EQ(r.value.to_packet(), FRAME_PA);
}

TEST(RoundTrip, Pc) {
    auto r = PcPacket::from_packet(FRAME_PC);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.timestamp,   1700000000LL);
    EXPECT_EQ(r.value.command_id,  42);
    EXPECT_EQ(r.value.message,     "ok");
    EXPECT_EQ(r.value.to_packet(), FRAME_PC);
}

TEST(RoundTrip, Ts) {
    auto r = TsPacket::from_packet(FRAME_TS);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.timestamp, 1700000000LL);
    EXPECT_EQ(r.value.trip_id,   "12345678-1234-5678-1234-567812345678");
    EXPECT_EQ(r.value.to_packet(), FRAME_TS);
}

TEST(RoundTrip, Te) {
    auto r = TePacket::from_packet(FRAME_TE);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.timestamp, 1700000000LL);
    EXPECT_EQ(r.value.trip_id,   "12345678-1234-5678-1234-567812345678");
    EXPECT_NEAR(r.value.distance_traveled, 1234.567, 0.0001);
    EXPECT_NEAR(r.value.max_speed,         89.012,   0.0001);
    EXPECT_EQ(r.value.duration, 3600);
    EXPECT_EQ(r.value.to_packet(), FRAME_TE);
}

TEST(RoundTrip, Im) {
    auto r = ImPacket::from_packet(FRAME_IM);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.timestamp, 1700000000LL);
    EXPECT_EQ(r.value.chat_id,   "12345678-1234-5678-1234-567812345678");
    EXPECT_EQ(r.value.message,   "Hello; world"); // '|||' decoded to ';'
    EXPECT_EQ(r.value.to_packet(), FRAME_IM);
}

// ── Python canonical frames ───────────────────────────────────────────────────

TEST(RoundTrip, Pi) {
    // FRAME_PI has the payload without CRC; complete it and test
    std::string full = complete_frame("<Pi>", "</Pi>", "testident;1;1;1;1;1;1;1;");
    auto r = PiPacket::from_packet(full);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.ident,           "testident");
    EXPECT_EQ(std::get<int>(r.value.firmware_id), 1);
    EXPECT_EQ(r.value.firmware_build,  1);
    EXPECT_EQ(r.value.device_id,       1);
    EXPECT_EQ(r.value.hardware_id,     1);
    EXPECT_EQ(r.value.model_id,        1);
    EXPECT_EQ(r.value.firmware_branch, FirmwareBranch::Development);
    EXPECT_TRUE(r.value.fota_enabled);
    EXPECT_EQ(r.value.to_packet(), full);
}

TEST(RoundTrip, Pd) {
    std::string payload =
        "0;10.0;10.0;10.0;10.0;10.0;5;1.0;"
        "test.str:Hola mundo,test.int:1,test.double:1.0,test.bool:true;";
    std::string full = complete_frame("<Pd>", "</Pd>", payload);
    auto r = PdPacket::from_packet(full);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.timestamp, 0LL);
    ASSERT_TRUE(r.value.position.has_value());
    EXPECT_DOUBLE_EQ(*r.value.position->latitude,  10.0);
    EXPECT_DOUBLE_EQ(*r.value.position->longitude, 10.0);
    EXPECT_EQ(r.value.extra_data.size(), 4u);
    // Re-encode must match
    EXPECT_EQ(r.value.to_packet(), full);
}

TEST(RoundTrip, Ps) {
    std::string payload =
        "0;net_wifi_ssid:AWESOME WIFI,net_wifi_pass:dictadormarico69,"
        "net_wifi_sec:WPA2,static.lat:-15.5,static.lng:15.5;";
    std::string full = complete_frame("<Ps>", "</Ps>", payload);
    auto r = PsPacket::from_packet(full);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.timestamp, 0LL);
    EXPECT_EQ(r.value.params.size(), 5u);
    EXPECT_EQ(r.value.to_packet(), full);
}

// ── Parser dispatch ──────────────────────────────────────────────────────────

TEST(Parser, HandleServerOutput_As) {
    auto r = handle_server_output(FRAME_AS);
    ASSERT_TRUE(r.ok());
    EXPECT_TRUE(std::holds_alternative<AsPacket>(r.value));
}

TEST(Parser, HandleServerOutput_Ar) {
    auto r = handle_server_output(FRAME_AR);
    ASSERT_TRUE(r.ok());
    ASSERT_TRUE(std::holds_alternative<ArPacket>(r.value));
    EXPECT_EQ(std::get<ArPacket>(r.value).reason, "Unknown reason");
}

TEST(Parser, HandleServerOutput_Unknown) {
    auto r = handle_server_output("<Xx>garbage</Xx>");
    EXPECT_FALSE(r.ok());
}

TEST(Parser, ParsePacketToString_Pa) {
    AnyClientPacket pkt = PaPacket{"123456789012345", "mypassword"};
    auto r = parse_packet_to_string(pkt);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value, FRAME_PA);
}

TEST(Parser, ParsePacketToString_Pr) {
    AnyClientPacket pkt = PrPacket{};
    auto r = parse_packet_to_string(pkt);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value, FRAME_PR);
}

// ── Generated fixtures (Pb / Pm) ─────────────────────────────────────────────

TEST(RoundTrip, Pb) {
    auto r = PbPacket::from_packet(FRAME_PB);
    ASSERT_TRUE(r.ok()) << "Pb parse failed";
    ASSERT_EQ(r.value.advertisements.size(), 2u);

    auto& adv1 = r.value.advertisements[0];
    EXPECT_EQ(adv1.mac_address, "12:34:56:78:90:AB");
    EXPECT_EQ(adv1.timestamp, 1700000000LL);
    EXPECT_DOUBLE_EQ(*adv1.latitude,  10.0);
    EXPECT_DOUBLE_EQ(*adv1.longitude, 20.0);
    EXPECT_DOUBLE_EQ(*adv1.altitude,  100.0);
    EXPECT_EQ(adv1.model, "GENERIC");
    EXPECT_EQ(adv1.rssi, -70);
    ASSERT_EQ(adv1.manufacturer_data.size(), 1u);
    EXPECT_EQ(adv1.manufacturer_data[0].company_id, 0x004C);
    ASSERT_EQ(adv1.service_data.size(), 1u);
    EXPECT_EQ(adv1.service_data[0].uuid, 0xFD6F);

    auto& adv2 = r.value.advertisements[1];
    EXPECT_EQ(adv2.mac_address, "BC:09:87:65:43:21");
    EXPECT_EQ(adv2.rssi, -80);
    EXPECT_EQ(*adv2.tx_power, -10);

    EXPECT_EQ(r.value.to_packet(), FRAME_PB);
}

TEST(RoundTrip, Pm) {
    auto r = PmPacket::from_packet(FRAME_PM);
    ASSERT_TRUE(r.ok()) << "Pm parse failed";
    EXPECT_EQ(r.value.filename,     "test.jpg");
    EXPECT_EQ(r.value.content_type, "image/jpeg");
    // data = {0xFF, 0xD8, 0xFF, 0xE0}
    ASSERT_EQ(r.value.data.size(), 4u);
    EXPECT_EQ(r.value.data[0], 0xFF);
    EXPECT_EQ(r.value.data[1], 0xD8);
    EXPECT_EQ(r.value.data[2], 0xFF);
    EXPECT_EQ(r.value.data[3], 0xE0);
    EXPECT_EQ(r.value.to_packet(), FRAME_PM);
}
