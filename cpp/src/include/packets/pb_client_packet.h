#pragma once

#include <algorithm>
#include <iomanip>
#include <sstream>
#include <stdexcept>
#include <string>
#include <vector>

#include "ble/ble_advertisement.h"
#include "exceptions/crc_exception.h"
#include "exceptions/malformed_exception.h"
#include "helpers/crc.h"
#include "helpers/strings.h"

namespace layrz_protocol
{

class PbClientPacket
{
public:
  std::vector<BleAdvertisement> advertisements;

public:
  PbClientPacket (const std::vector<BleAdvertisement> &advertisements);
  PbClientPacket from_packet (const std::string &raw);
  std::string to_packet () const noexcept;
  std::string to_string () const noexcept;
};

} // namespace layrz_protocol
