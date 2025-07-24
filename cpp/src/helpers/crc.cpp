#include "helpers/crc.h"

namespace layrz_protocol
{
namespace helpers
{
/// @brief Calculate the CRC of a byte array
/// @param data
/// @return
uint16_t
calculate_crc (const std::vector<uint8_t> &data)
{
  uint16_t fcs = 0xFFFF;
  for (auto b : data)
    {
      uint8_t index = (fcs ^ b) & 0xFF;
      fcs = (fcs >> 8) ^ crc_tab[index];
    }
  return fcs ^ 0xFFFF;
}

} // namespace helpers
} // namespace layrz_protocol
