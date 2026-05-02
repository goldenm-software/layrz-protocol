#include <gtest/gtest.h>
#include "layrz_protocol/extras.hpp"

using namespace layrz::protocol;

TEST(Extras, ParseEmpty) {
    auto r = parse_extra("");
    EXPECT_TRUE(r.empty());
}

TEST(Extras, ParseIntFloatBoolString) {
    auto r = parse_extra("a:1,b:1.5,c:true,d:hello");
    ASSERT_EQ(r.size(), 4u);
    EXPECT_EQ(r[0].first, "a");  EXPECT_TRUE(std::holds_alternative<int64_t>(r[0].second));
    EXPECT_EQ(std::get<int64_t>(r[0].second), 1LL);
    EXPECT_EQ(r[1].first, "b");  EXPECT_TRUE(std::holds_alternative<double>(r[1].second));
    EXPECT_DOUBLE_EQ(std::get<double>(r[1].second), 1.5);
    EXPECT_EQ(r[2].first, "c");  EXPECT_TRUE(std::holds_alternative<bool>(r[2].second));
    EXPECT_TRUE(std::get<bool>(r[2].second));
    EXPECT_EQ(r[3].first, "d");  EXPECT_TRUE(std::holds_alternative<std::string>(r[3].second));
    EXPECT_EQ(std::get<std::string>(r[3].second), "hello");
}

TEST(Extras, KeyRemapGpioDigitalInput) {
    auto r = parse_extra("io3.di:1");
    ASSERT_EQ(r.size(), 1u);
    EXPECT_EQ(r[0].first, "gpio.3.digital.input");
}

TEST(Extras, KeyRemapGpioAnalogOutput) {
    auto r = parse_extra("io7.ao:3.3");
    ASSERT_EQ(r.size(), 1u);
    EXPECT_EQ(r[0].first, "gpio.7.analog.output");
}

TEST(Extras, KeyRemapBleTemp) {
    auto r = parse_extra("ble.2.tempc:23.5");
    ASSERT_EQ(r.size(), 1u);
    EXPECT_EQ(r[0].first, "ble.2.temperature.celsius");
}

TEST(Extras, KeyRemapReport) {
    auto r = parse_extra("report:42");
    ASSERT_EQ(r.size(), 1u);
    EXPECT_EQ(r[0].first, "report.code");
    EXPECT_EQ(std::get<int64_t>(r[0].second), 42LL);
}

TEST(Extras, KeyRemapConfiotBle) {
    auto r = parse_extra("confiot_ble:1");
    ASSERT_EQ(r.size(), 1u);
    EXPECT_EQ(r[0].first, "ble.confiot.connection.status");
}

TEST(Extras, CastBasicRoundTrip) {
    ExtrasMap m;
    m.emplace_back("a", ExtraValue{int64_t{1}});
    m.emplace_back("b", ExtraValue{double{1.5}});
    m.emplace_back("c", ExtraValue{true});
    m.emplace_back("d", ExtraValue{std::string{"hello"}});
    std::string out = cast_extra(m);
    EXPECT_EQ(out, "a:1,b:1.5,c:true,d:hello");
}

TEST(Extras, CastBoolFalse) {
    ExtrasMap m;
    m.emplace_back("flag", ExtraValue{false});
    EXPECT_EQ(cast_extra(m), "flag:false");
}

TEST(Extras, CastFloat10) {
    ExtrasMap m;
    m.emplace_back("x", ExtraValue{double{10.0}});
    EXPECT_EQ(cast_extra(m), "x:10.0");
}

TEST(Extras, PsCanonical) {
    // From Python test: params used in Ps canonical frame
    auto r = parse_extra("net_wifi_ssid:AWESOME WIFI,net_wifi_pass:dictadormarico69,"
                         "net_wifi_sec:WPA2,static.lat:-15.5,static.lng:15.5");
    ASSERT_EQ(r.size(), 5u);
    EXPECT_EQ(r[0].first, "net_wifi_ssid");
    EXPECT_EQ(std::get<std::string>(r[0].second), "AWESOME WIFI");
    EXPECT_DOUBLE_EQ(std::get<double>(r[3].second), -15.5);
    EXPECT_DOUBLE_EQ(std::get<double>(r[4].second), 15.5);
}

TEST(Extras, AccentStripping) {
    auto r = parse_extra("msg:caf\xC3\xA9"); // café
    ASSERT_EQ(r.size(), 1u);
    EXPECT_EQ(std::get<std::string>(r[0].second), "cafe");
}
