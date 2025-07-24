#pragma once

#include <ctime>
#include <iomanip>
#include <sstream>
#include <stdexcept>
#include <string>
#include <vector>

#include "exceptions/crc_exception.h"
#include "exceptions/malformed_exception.h"
#include "helpers/crc.h"
#include "helpers/strings.h"

namespace layrz_protocol
{
class PcClientPacket
{
public:
  std::time_t timestamp;
  int command_id;
  std::string message;

public:
  PcClientPacket (std::time_t timestamp, int command_id, const std::string &message)ñ
  // ~PcClientPacket () = default;
  PcClientPacket from_packet (const std::string &raw);
  std::string to_packet () const noexcept;
  std::string to_string () const noexcept;
};

}
