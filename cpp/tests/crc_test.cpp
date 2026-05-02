#include <gtest/gtest.h>
#include "layrz_protocol/crc.hpp"

using namespace layrz::protocol;

TEST(CRC, KnownSemicolon) {
    // CRC of ";" is always 0x7F28 — used by As/Au/Pr packets
    EXPECT_EQ(crc16_x25(";"), 0x7F28u);
}

TEST(CRC, HexFormatting) {
    EXPECT_EQ(crc_hex(0x7F28), "7F28");
    EXPECT_EQ(crc_hex(0x00C1), "00C1");
    EXPECT_EQ(crc_hex(0xA058), "A058");
    EXPECT_EQ(crc_hex(0x0000), "0000");
    EXPECT_EQ(crc_hex(0xFFFF), "FFFF");
}

TEST(CRC, AoPayload) {
    // Ao canonical: payload = "1700000000;" → CRC A058
    EXPECT_EQ(crc16_x25("1700000000;"), 0xA058u);
    EXPECT_EQ(compute_crc_str("1700000000;"), "A058");
}

TEST(CRC, ArPayload) {
    // Ar canonical: payload = "Unknown reason;" → CRC 1DFD
    EXPECT_EQ(crc16_x25("Unknown reason;"), 0x1DFDu);
}

TEST(CRC, PaPayload) {
    // Pa canonical: payload = "123456789012345;mypassword;" → CRC 2B64
    EXPECT_EQ(crc16_x25("123456789012345;mypassword;"), 0x2B64u);
}

TEST(CRC, PcPayload) {
    // Pc canonical: payload = "1700000000;42;ok;" → CRC A497
    EXPECT_EQ(crc16_x25("1700000000;42;ok;"), 0xA497u);
}
