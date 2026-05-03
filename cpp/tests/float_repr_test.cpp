#include <gtest/gtest.h>
#include "layrz_protocol/float_repr.hpp"

using namespace layrz::protocol;

// Golden values verified against CPython: str(float_value)
TEST(FloatRepr, WholeNumbers) {
    EXPECT_EQ(python_repr_float(0.0),   "0.0");
    EXPECT_EQ(python_repr_float(1.0),   "1.0");
    EXPECT_EQ(python_repr_float(10.0),  "10.0");
    EXPECT_EQ(python_repr_float(-1.0),  "-1.0");
    EXPECT_EQ(python_repr_float(100.0), "100.0");
}

TEST(FloatRepr, Fractional) {
    EXPECT_EQ(python_repr_float(1.5),    "1.5");
    EXPECT_EQ(python_repr_float(-15.5),  "-15.5");
    EXPECT_EQ(python_repr_float(10.0),   "10.0");
    EXPECT_EQ(python_repr_float(1.0),    "1.0");
}

TEST(FloatRepr, SmallFractional) {
    EXPECT_EQ(python_repr_float(0.1),  "0.1");
    EXPECT_EQ(python_repr_float(0.5),  "0.5");
}

TEST(FloatRepr, NegativeZero) {
    // Python: str(-0.0) == "-0.0"
    EXPECT_EQ(python_repr_float(-0.0), "-0.0");
}

TEST(FloatRepr, UsedInPdTest) {
    // Pd canonical test uses latitude=10.0
    EXPECT_EQ(python_repr_float(10.0), "10.0");
    // static.lat=-15.5, static.lng=15.5
    EXPECT_EQ(python_repr_float(-15.5), "-15.5");
    EXPECT_EQ(python_repr_float(15.5),  "15.5");
    // hdop=1.0
    EXPECT_EQ(python_repr_float(1.0),  "1.0");
    // test.double=1.0
    EXPECT_EQ(python_repr_float(1.0),  "1.0");
}
