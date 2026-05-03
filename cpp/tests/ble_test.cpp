#include <gtest/gtest.h>
#include "layrz_protocol/ble/manufacturer_data.hpp"
#include "layrz_protocol/ble/service_data.hpp"
#include "layrz_protocol/ble/advertisement.hpp"
#include "layrz_protocol/crc.hpp"

using namespace layrz::protocol;
using namespace layrz::protocol::ble;

TEST(BleManufacturerData, RoundTrip) {
    ManufacturerData mfr;
    mfr.company_id = 0x004C;
    mfr.data = {0xAA, 0xBB, 0xCC};
    std::string encoded = mfr.to_packet();
    EXPECT_EQ(encoded, "004C:AABBCC");

    auto r = ManufacturerData::from_packet(encoded);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.company_id, 0x004C);
    EXPECT_EQ(r.value.data, (std::vector<uint8_t>{0xAA, 0xBB, 0xCC}));
}

TEST(BleServiceData, RoundTrip) {
    ServiceData sd;
    sd.uuid = 0xFD6F;
    sd.data = {0x01, 0x02};
    std::string encoded = sd.to_packet();
    EXPECT_EQ(encoded, "FD6F:0102");

    auto r = ServiceData::from_packet(encoded);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.uuid, 0xFD6F);
    EXPECT_EQ(r.value.data, (std::vector<uint8_t>{0x01, 0x02}));
}

TEST(BleAdvertisement, RoundTripNoMfrSvc) {
    Advertisement a;
    a.mac_address = "12:34:56:78:90:AB";
    a.timestamp   = 1700000000;
    a.latitude    = 10.0;
    a.longitude   = 20.0;
    a.altitude    = 100.0;
    a.model       = "GENERIC";
    a.device_name = "TestDevice";
    a.rssi        = -70;
    a.tx_power    = std::nullopt;

    std::string encoded = a.to_packet();
    auto r = Advertisement::from_packet(encoded);
    ASSERT_TRUE(r.ok());
    EXPECT_EQ(r.value.mac_address, "12:34:56:78:90:AB");
    EXPECT_EQ(r.value.timestamp,   1700000000LL);
    EXPECT_DOUBLE_EQ(*r.value.latitude,  10.0);
    EXPECT_DOUBLE_EQ(*r.value.longitude, 20.0);
    EXPECT_EQ(r.value.model, "GENERIC");
    EXPECT_EQ(r.value.rssi, -70);
    EXPECT_FALSE(r.value.tx_power.has_value());
    EXPECT_EQ(r.value.to_packet(), encoded);
}

TEST(BleAdvertisement, RoundTripWithMfrSvc) {
    Advertisement a;
    a.mac_address = "AA:BB:CC:DD:EE:FF";
    a.timestamp   = 0;
    a.model       = "X";
    a.device_name = "";
    a.rssi        = -50;
    a.tx_power    = -10;

    ManufacturerData mfr;
    mfr.company_id = 0x0001;
    mfr.data = {0xFF};
    a.manufacturer_data.push_back(mfr);

    ServiceData sd;
    sd.uuid = 0x180D;
    sd.data = {0x00, 0x01};
    a.service_data.push_back(sd);

    std::string encoded = a.to_packet();
    auto r = Advertisement::from_packet(encoded);
    ASSERT_TRUE(r.ok());
    ASSERT_EQ(r.value.manufacturer_data.size(), 1u);
    EXPECT_EQ(r.value.manufacturer_data[0].company_id, 0x0001);
    ASSERT_EQ(r.value.service_data.size(), 1u);
    EXPECT_EQ(r.value.service_data[0].uuid, 0x180D);
    EXPECT_EQ(r.value.to_packet(), encoded);
}

TEST(BleAdvertisement, MacFormatting) {
    Advertisement a;
    a.mac_address = "AA:BB:CC:DD:EE:FF";
    a.timestamp   = 0;
    a.model       = "X";
    a.device_name = "";
    a.rssi        = 0;
    std::string encoded = a.to_packet();
    // Wire MAC must be 12 uppercase hex chars with no colons
    EXPECT_EQ(encoded.substr(0, 12), "AABBCCDDEEFF");
}
