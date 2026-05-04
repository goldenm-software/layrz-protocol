#include "layrz_protocol/parser.hpp"
#include <cctype>
#include <string>
#include <string_view>
#include <vector>

namespace layrz::protocol {

static bool starts_with(std::string_view s, const char* prefix) {
    std::string_view p(prefix);
    return s.size() >= p.size() && s.substr(0, p.size()) == p;
}

Result<AnyServerPacket> handle_server_output(std::string_view raw) {
    if (starts_with(raw, "<Ab>")) {
        auto r = packets::AbPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Ac>")) {
        auto r = packets::AcPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Ao>")) {
        auto r = packets::AoPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Ar>")) {
        auto r = packets::ArPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<As>")) {
        auto r = packets::AsPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Au>")) {
        auto r = packets::AuPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Ts>")) {
        auto r = packets::TsPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Te>")) {
        auto r = packets::TePacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Im>")) {
        auto r = packets::ImPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyServerPacket>::fail(r.error);
        return Result<AnyServerPacket>::success(std::move(r.value));
    }
    return Result<AnyServerPacket>::fail(Error::MalformedFrame);
}

Result<AnyClientPacket> handle_client_input(std::string_view raw) {
    if (starts_with(raw, "<Pa>")) {
        auto r = packets::PaPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyClientPacket>::fail(r.error);
        return Result<AnyClientPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Pb>")) {
        auto r = packets::PbPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyClientPacket>::fail(r.error);
        return Result<AnyClientPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Pc>")) {
        auto r = packets::PcPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyClientPacket>::fail(r.error);
        return Result<AnyClientPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Pd>")) {
        auto r = packets::PdPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyClientPacket>::fail(r.error);
        return Result<AnyClientPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Pi>")) {
        auto r = packets::PiPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyClientPacket>::fail(r.error);
        return Result<AnyClientPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Pm>")) {
        auto r = packets::PmPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyClientPacket>::fail(r.error);
        return Result<AnyClientPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Pr>")) {
        auto r = packets::PrPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyClientPacket>::fail(r.error);
        return Result<AnyClientPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Ps>")) {
        auto r = packets::PsPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyClientPacket>::fail(r.error);
        return Result<AnyClientPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Ts>")) {
        auto r = packets::TsPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyClientPacket>::fail(r.error);
        return Result<AnyClientPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Te>")) {
        auto r = packets::TePacket::from_packet(raw);
        if (!r.ok()) return Result<AnyClientPacket>::fail(r.error);
        return Result<AnyClientPacket>::success(std::move(r.value));
    }
    if (starts_with(raw, "<Im>")) {
        auto r = packets::ImPacket::from_packet(raw);
        if (!r.ok()) return Result<AnyClientPacket>::fail(r.error);
        return Result<AnyClientPacket>::success(std::move(r.value));
    }
    return Result<AnyClientPacket>::fail(Error::MalformedFrame);
}

// Finds positions of opening client-packet tags in a buffer.
// Tags are exactly 4 bytes: '<', uppercase, lowercase, '>'.
// Only recognised client tags are included.
static bool is_client_open_tag(std::string_view s) {
    if (s.size() < 4) return false;
    if (s[0] != '<' || s[3] != '>') return false;
    if (!std::isupper(static_cast<unsigned char>(s[1]))) return false;
    if (!std::islower(static_cast<unsigned char>(s[2]))) return false;
    // Only allow known client packet prefixes (server tags like <Ab> must not split)
    char hi = s[1], lo = s[2];
    if (hi == 'P') return lo == 'a' || lo == 'b' || lo == 'c' || lo == 'd' ||
                          lo == 'i' || lo == 'm' || lo == 'r' || lo == 's';
    if (hi == 'T') return lo == 's' || lo == 'e';
    if (hi == 'I') return lo == 'm';
    return false;
}

std::vector<std::string> split_client_frames(std::string_view buffer) {
    std::vector<std::string> frames;

    // Collect positions of all recognised open tags
    std::vector<size_t> tag_positions;
    for (size_t i = 0; i + 3 < buffer.size(); ++i) {
        if (is_client_open_tag(buffer.substr(i, 4))) {
            tag_positions.push_back(i);
        }
    }

    for (size_t idx = 0; idx < tag_positions.size(); ++idx) {
        size_t start = tag_positions[idx];
        size_t end   = (idx + 1 < tag_positions.size())
                           ? tag_positions[idx + 1]
                           : buffer.size();

        // Build the expected closing tag </Xx>
        char close_tag[6] = {'<', '/', buffer[start + 1], buffer[start + 2], '>', '\0'};
        std::string_view slice = buffer.substr(start, end - start);
        size_t close_pos = slice.rfind(close_tag);
        if (close_pos == std::string_view::npos) continue; // incomplete frame

        std::string frame(slice.substr(0, close_pos + 5));
        // Trim trailing whitespace/CRLF
        while (!frame.empty() && (frame.back() == '\r' || frame.back() == '\n' || frame.back() == ' '))
            frame.pop_back();
        if (!frame.empty()) frames.push_back(std::move(frame));
    }

    return frames;
}

Result<std::string> parse_packet_to_string(const AnyClientPacket& packet) {
    return std::visit([](auto&& p) -> Result<std::string> {
        return Result<std::string>::success(p.to_packet());
    }, packet);
}

Result<std::string> parse_server_packet_to_string(const AnyServerPacket& packet) {
    return std::visit([](auto&& p) -> Result<std::string> {
        return Result<std::string>::success(p.to_packet());
    }, packet);
}

} // namespace layrz::protocol
