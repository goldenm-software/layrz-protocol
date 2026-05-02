#include "layrz_protocol/extras.hpp"
#include "layrz_protocol/float_repr.hpp"
#include <cctype>
#include <regex>
#include <sstream>
#include <string>

namespace layrz::protocol {

// ASCII accent-stripping map  (matches Python's ASCII_MAP in constants.py)
// Keys are multi-byte UTF-8 sequences; stored as string pairs for lookup.
static const std::pair<const char*, char> ASCII_MAP[] = {
    {"\xC3\xA1", 'a'}, // á
    {"\xC3\xA9", 'e'}, // é
    {"\xC3\xAD", 'i'}, // í
    {"\xC3\xB3", 'o'}, // ó
    {"\xC3\xBA", 'u'}, // ú
    {"\xC3\xB1", 'n'}, // ñ
    {"\xC3\xBC", 'u'}, // ü
    {"\xC3\xA0", 'a'}, // à
    {"\xC3\xA8", 'e'}, // è
    {"\xC3\xAC", 'i'}, // ì
    {"\xC3\xB2", 'o'}, // ò
    {"\xC3\xB9", 'u'}, // ù
    {"\xC3\xA2", 'a'}, // â
    {"\xC3\xAA", 'e'}, // ê
    {"\xC3\xAE", 'i'}, // î
    {"\xC3\xB4", 'o'}, // ô
    {"\xC3\xBB", 'u'}, // û
    {"\xC3\xA4", 'a'}, // ä
    {"\xC3\xAB", 'e'}, // ë
    {"\xC3\xAF", 'i'}, // ï
    {"\xC3\xB6", 'o'}, // ö
};

static std::string strip_accents(const std::string& s) {
    std::string out;
    out.reserve(s.size());
    size_t i = 0;
    while (i < s.size()) {
        bool replaced = false;
        // All accented chars in our map are 2-byte UTF-8 sequences
        if (i + 1 < s.size()) {
            char buf[3] = {s[i], s[i+1], '\0'};
            for (auto& [seq, rep] : ASCII_MAP) {
                if (seq[0] == buf[0] && seq[1] == buf[1]) {
                    out += rep;
                    i += 2;
                    replaced = true;
                    break;
                }
            }
        }
        if (!replaced) {
            out += s[i++];
        }
    }
    return out;
}

// Remap short BLE/GPIO wire keys to their canonical dot-case names.
// This is a one-way decode-side operation (matches Python's parse_extra).
static std::string remap_key(const std::string& key) {
    // GPIO: io<N>.(di|do|ai|ao|counter)
    static const std::regex re_io_di(R"(^io(\d+)\.di$)");
    static const std::regex re_io_do(R"(^io(\d+)\.do$)");
    static const std::regex re_io_ai(R"(^io(\d+)\.ai$)");
    static const std::regex re_io_ao(R"(^io(\d+)\.ao$)");
    static const std::regex re_io_cnt(R"(^io(\d+)\.counter$)");

    // BLE: ble.<N>.<suffix>
    static const std::regex re_ble_id(R"(^ble\.(\d+)\.id$)");
    static const std::regex re_ble_hum(R"(^ble\.(\d+)\.hum$)");
    static const std::regex re_ble_tempc(R"(^ble\.(\d+)\.tempc$)");
    static const std::regex re_ble_tempf(R"(^ble\.(\d+)\.tempf$)");
    static const std::regex re_ble_model(R"(^ble\.(\d+)\.model_id$)");
    static const std::regex re_ble_batt(R"(^ble\.(\d+)\.batt$)");
    static const std::regex re_ble_lux(R"(^ble\.(\d+)\.lux$)");
    static const std::regex re_ble_volt(R"(^ble\.(\d+)\.volt$)");
    static const std::regex re_ble_rpm(R"(^ble\.(\d+)\.rpm$)");
    static const std::regex re_ble_press(R"(^ble\.(\d+)\.press$)");
    static const std::regex re_ble_cnt(R"(^ble\.(\d+)\.counter$)");
    static const std::regex re_ble_xacc(R"(^ble\.(\d+)\.x_acc$)");
    static const std::regex re_ble_yacc(R"(^ble\.(\d+)\.y_acc$)");
    static const std::regex re_ble_zacc(R"(^ble\.(\d+)\.z_acc$)");
    static const std::regex re_ble_msgcnt(R"(^ble\.(\d+)\.msg_count$)");
    static const std::regex re_ble_msg(R"(^ble\.(\d+)\.msg$)");
    static const std::regex re_ble_magcnt(R"(^ble\.(\d+)\.mag_counter)");
    static const std::regex re_ble_magdata(R"(^ble\.(\d+)\.mag_data)");
    static const std::regex re_ble_rssi(R"(^ble\.(\d+)\.rssi)");

    std::smatch m;

    if (std::regex_match(key, m, re_io_di))    return "gpio." + m[1].str() + ".digital.input";
    if (std::regex_match(key, m, re_io_do))    return "gpio." + m[1].str() + ".digital.output";
    if (std::regex_match(key, m, re_io_ai))    return "gpio." + m[1].str() + ".analog.input";
    if (std::regex_match(key, m, re_io_ao))    return "gpio." + m[1].str() + ".analog.output";
    if (std::regex_match(key, m, re_io_cnt))   return "gpio." + m[1].str() + ".event.count";

    if (std::regex_match(key, m, re_ble_id))    return "ble." + m[1].str() + ".mac.address";
    if (std::regex_match(key, m, re_ble_hum))   return "ble." + m[1].str() + ".humidity";
    if (std::regex_match(key, m, re_ble_tempc)) return "ble." + m[1].str() + ".temperature.celsius";
    if (std::regex_match(key, m, re_ble_tempf)) return "ble." + m[1].str() + ".temperature.fahrenheit";
    if (std::regex_match(key, m, re_ble_model)) return "ble." + m[1].str() + ".model.id";
    if (std::regex_match(key, m, re_ble_batt))  return "ble." + m[1].str() + ".battery.level";
    if (std::regex_match(key, m, re_ble_lux))   return "ble." + m[1].str() + ".light.level.lux";
    if (std::regex_match(key, m, re_ble_volt))  return "ble." + m[1].str() + ".voltage";
    if (std::regex_match(key, m, re_ble_rpm))   return "ble." + m[1].str() + ".rpm";
    if (std::regex_match(key, m, re_ble_press)) return "ble." + m[1].str() + ".pressure";
    if (std::regex_match(key, m, re_ble_cnt))   return "ble." + m[1].str() + ".event.count";
    if (std::regex_match(key, m, re_ble_xacc))  return "ble." + m[1].str() + ".acceleration.x";
    if (std::regex_match(key, m, re_ble_yacc))  return "ble." + m[1].str() + ".acceleration.y";
    if (std::regex_match(key, m, re_ble_zacc))  return "ble." + m[1].str() + ".acceleration.z";
    if (std::regex_match(key, m, re_ble_msgcnt)) return "ble." + m[1].str() + ".message.count";
    if (std::regex_match(key, m, re_ble_msg))   return "ble." + m[1].str() + ".message";
    if (std::regex_search(key, m, re_ble_magcnt)) return "ble." + m[1].str() + ".magnetic.event.count";
    if (std::regex_search(key, m, re_ble_magdata)) return "ble." + m[1].str() + ".magnetic.data";
    if (std::regex_search(key, m, re_ble_rssi)) return "ble." + m[1].str() + ".rssi.dbm";

    if (key == "report")        return "report.code";
    if (key == "confiot_ble")   return "ble.confiot.connection.status";
    if (key == "confiot_serial") return "serial.confiot.connection.status";

    return key;
}

static ExtraValue parse_value(const std::string& v) {
    // Float: -?digits.digits
    bool is_float = false, is_int = false;
    {
        const char* p = v.c_str();
        if (*p == '-') ++p;
        bool digits = false;
        while (std::isdigit(static_cast<unsigned char>(*p))) { digits = true; ++p; }
        if (digits) {
            if (*p == '.' && std::isdigit(static_cast<unsigned char>(*(p+1)))) {
                ++p;
                while (std::isdigit(static_cast<unsigned char>(*p))) ++p;
                is_float = (*p == '\0');
            } else {
                is_int = (*p == '\0');
            }
        }
    }
    if (is_float) return std::stod(v);
    if (is_int)   return static_cast<int64_t>(std::stoll(v));

    std::string low = v;
    for (char& c : low) c = static_cast<char>(std::tolower(static_cast<unsigned char>(c)));
    if (low == "true")  return true;
    if (low == "false") return false;
    if (low == "t")     return true;
    if (low == "f")     return false;

    return std::string(v);
}

ExtrasMap parse_extra(std::string_view raw_view) {
    ExtrasMap result;
    if (raw_view.empty()) return result;

    std::string raw(raw_view);
    // Split on ',' — but value may contain ':' so split key:value on first ':' only
    size_t start = 0;
    while (start <= raw.size()) {
        size_t comma = raw.find(',', start);
        if (comma == std::string::npos) comma = raw.size();
        std::string token = raw.substr(start, comma - start);
        start = comma + 1;
        if (token.empty()) continue;

        size_t colon = token.find(':');
        if (colon == std::string::npos) continue; // malformed token; skip

        std::string key   = token.substr(0, colon);
        std::string value = token.substr(colon + 1);

        key = remap_key(key);
        ExtraValue ev = parse_value(value);
        // Apply ASCII_MAP to string values
        if (std::holds_alternative<std::string>(ev)) {
            ev = strip_accents(std::get<std::string>(ev));
        }
        result.emplace_back(key, std::move(ev));
    }
    return result;
}

std::string cast_extra(const ExtrasMap& extras) {
    std::string out;
    bool first = true;
    for (auto& [key, val] : extras) {
        if (!first) out += ',';
        first = false;
        out += key + ':';
        std::visit([&](auto&& v) {
            using T = std::decay_t<decltype(v)>;
            if constexpr (std::is_same_v<T, bool>) {
                out += v ? "true" : "false";
            } else if constexpr (std::is_same_v<T, int64_t>) {
                out += std::to_string(v);
            } else if constexpr (std::is_same_v<T, double>) {
                out += python_repr_float(v);
            } else {
                out += v;
            }
        }, val);
    }
    return out;
}

} // namespace layrz::protocol
