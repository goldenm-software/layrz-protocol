#include "layrz_protocol/packets/pm.hpp"
#include "packets/helpers.hpp"
#include <cstdint>

namespace layrz::protocol::packets {

using namespace detail;

// Minimal standard-base64 encoder (RFC 4648, no line wrapping)
static const char B64_CHARS[] =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

static std::string base64_encode(const std::vector<uint8_t>& data) {
    std::string out;
    size_t n = data.size();
    out.reserve(((n + 2) / 3) * 4);
    for (size_t i = 0; i < n; i += 3) {
        uint32_t a = data[i];
        uint32_t b = (i + 1 < n) ? data[i + 1] : 0;
        uint32_t c = (i + 2 < n) ? data[i + 2] : 0;
        uint32_t t = (a << 16) | (b << 8) | c;
        out += B64_CHARS[(t >> 18) & 0x3F];
        out += B64_CHARS[(t >> 12) & 0x3F];
        out += (i + 1 < n) ? B64_CHARS[(t >> 6) & 0x3F] : '=';
        out += (i + 2 < n) ? B64_CHARS[ t       & 0x3F] : '=';
    }
    return out;
}

static std::vector<uint8_t> base64_decode(const std::string& s) {
    static const int8_t D[256] = {
        -1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
        -1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
        -1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,62,-1,-1,-1,63,
        52,53,54,55,56,57,58,59,60,61,-1,-1,-1,-1,-1,-1,
        -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9,10,11,12,13,14,
        15,16,17,18,19,20,21,22,23,24,25,-1,-1,-1,-1,-1,
        -1,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,
        41,42,43,44,45,46,47,48,49,50,51,-1,-1,-1,-1,-1,
        -1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
        -1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
        -1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
        -1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
        -1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
        -1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
        -1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
        -1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,
    };
    std::vector<uint8_t> out;
    out.reserve((s.size() * 3) / 4);
    uint32_t val = 0; int bits = -8;
    for (unsigned char c : s) {
        if (D[c] == -1) continue; // skip padding '=' or invalid
        val = (val << 6) | static_cast<uint32_t>(D[c]);
        bits += 6;
        if (bits >= 0) {
            out.push_back(static_cast<uint8_t>((val >> bits) & 0xFF));
            bits -= 8;
        }
    }
    return out;
}

std::string PmPacket::to_packet() const {
    std::string payload = filename + ";" + content_type + ";" + base64_encode(data) + ";";
    return wrap_packet("<Pm>", "</Pm>", payload);
}

Result<PmPacket> PmPacket::from_packet(std::string_view raw) {
    auto r = unwrap_packet(raw, "<Pm>", "</Pm>");
    if (!r.ok()) return Result<PmPacket>::fail(r.error);
    if (r.value.size() != 3) return Result<PmPacket>::fail(Error::MalformedFrame);

    PmPacket p;
    p.filename     = r.value[0];
    p.content_type = r.value[1];
    p.data         = base64_decode(r.value[2]);
    return Result<PmPacket>::success(p);
}

} // namespace layrz::protocol::packets
