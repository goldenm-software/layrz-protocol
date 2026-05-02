#include "layrz_protocol/transport/http_comm.hpp"
#include "layrz_protocol/parser.hpp"

// cpp-httplib: single-header HTTP library
// Included from the build-system-provided path; no HTTPS support (OPENSSL disabled)
#define CPPHTTPLIB_OPENSSL_SUPPORT 0
#include <httplib.h>

#include <stdexcept>

namespace layrz::protocol::transport {

HttpComm::HttpComm(HttpScheme scheme, std::string host,
                   std::string ident, std::string password)
    : scheme_(scheme)
    , host_(std::move(host))
    , ident_(std::move(ident))
    , password_(std::move(password))
{}

std::string HttpComm::make_url(const std::string& path) const {
    std::string scheme_str = (scheme_ == HttpScheme::Http) ? "http" : "https";
    return scheme_str + "://" + host_ + path;
}

std::string HttpComm::auth_header() const {
    return "LayrzAuth " + ident_ + ";" + password_;
}

Result<AnyServerPacket> HttpComm::send(const AnyClientPacket& packet) {
    auto frame_r = parse_packet_to_string(packet);
    if (!frame_r.ok()) return Result<AnyServerPacket>::fail(frame_r.error);
    const std::string& frame = frame_r.value;

    httplib::Client cli(host_);
    httplib::Headers headers = {{"Authorization", auth_header()}};
    auto res = cli.Post("/v2/message", headers, frame, "text/plain");
    if (!res) return Result<AnyServerPacket>::fail(Error::ServerError);
    if (res->status != 200) return Result<AnyServerPacket>::fail(Error::ServerError);

    return handle_server_output(res->body);
}

Result<AnyServerPacket> HttpComm::get_commands() {
    httplib::Client cli(host_);
    httplib::Headers headers = {{"Authorization", auth_header()}};
    auto res = cli.Get("/v2/commands", headers);
    if (!res) return Result<AnyServerPacket>::fail(Error::ServerError);
    if (res->status != 200) return Result<AnyServerPacket>::fail(Error::ServerError);

    return handle_server_output(res->body);
}

} // namespace layrz::protocol::transport
