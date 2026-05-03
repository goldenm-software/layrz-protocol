#include "layrz_protocol/transport/tcp_comm.hpp"

#ifdef LAYRZ_PROTOCOL_CLIENTS

#include "layrz_protocol/parser.hpp"
#include "layrz_protocol/packets/pa.hpp"

#if defined(ESP_PLATFORM)
  #include <lwip/sockets.h>
  #include <lwip/netdb.h>
  #define LAYRZ_CLOSE(fd)       lwip_close(fd)
  #define LAYRZ_SHUTDOWN(fd, h) lwip_shutdown(fd, h)
#elif defined(ARDUINO)
  // Arduino TCP is handled via the WiFiClient/EthernetClient path; raw BSD
  // sockets are not available.  Stub everything out — Arduino users get a
  // compile-time "not implemented" at link time.
  // Full Arduino backend is a follow-up.
  #define LAYRZ_ARDUINO_STUB
#else
  #include <arpa/inet.h>
  #include <netdb.h>
  #include <sys/socket.h>
  #include <sys/types.h>
  #include <unistd.h>
  #define LAYRZ_CLOSE(fd)       ::close(fd)
  #define LAYRZ_SHUTDOWN(fd, h) ::shutdown(fd, h)
#endif

#include <chrono>
#include <cstring>
#include <thread>

namespace layrz::protocol::transport {

TcpComm::TcpComm(std::string host, int port,
                 std::string ident, std::string password,
                 int connect_timeout_secs)
    : host_(std::move(host))
    , port_(port)
    , ident_(std::move(ident))
    , password_(std::move(password))
    , connect_timeout_secs_(connect_timeout_secs)
{}

TcpComm::~TcpComm() {
    close();
}

void TcpComm::set_callback(std::function<void(AnyServerPacket)> cb) {
    callback_ = std::move(cb);
}

#ifdef LAYRZ_ARDUINO_STUB

Error TcpComm::connect()    { return Error::Unimplemented; }
Error TcpComm::send(const AnyClientPacket&) { return Error::Unimplemented; }
void  TcpComm::close()      {}
Error TcpComm::write_raw(const std::string&) { return Error::Unimplemented; }
void  TcpComm::listen_loop() {}

#else // POSIX or ESP-IDF (both expose BSD-ish socket API)

Error TcpComm::connect() {
    struct addrinfo hints{}, *res = nullptr;
    hints.ai_family   = AF_UNSPEC;
    hints.ai_socktype = SOCK_STREAM;
    std::string port_str = std::to_string(port_);
    if (::getaddrinfo(host_.c_str(), port_str.c_str(), &hints, &res) != 0)
        return Error::ServerError;

    sockfd_ = ::socket(res->ai_family, res->ai_socktype, res->ai_protocol);
    if (sockfd_ < 0) { ::freeaddrinfo(res); return Error::ServerError; }

    if (::connect(sockfd_, res->ai_addr, res->ai_addrlen) != 0) {
        ::freeaddrinfo(res);
        LAYRZ_CLOSE(sockfd_);
        sockfd_ = -1;
        return Error::ServerError;
    }
    ::freeaddrinfo(res);

    stop_.store(false);
    authenticated_.store(false);
    listener_ = std::thread(&TcpComm::listen_loop, this);

    packets::PaPacket pa;
    pa.ident    = ident_;
    pa.password = password_;
    if (auto e = write_raw(pa.to_packet() + "\r\n"); e != Error::Ok)
        return e;

    auto deadline = std::chrono::steady_clock::now()
                  + std::chrono::seconds(connect_timeout_secs_);
    while (!authenticated_.load()) {
        if (std::chrono::steady_clock::now() > deadline) return Error::ServerError;
        std::this_thread::sleep_for(std::chrono::milliseconds(100));
    }
    return Error::Ok;
}

Error TcpComm::send(const AnyClientPacket& packet) {
    auto r = parse_packet_to_string(packet);
    if (!r.ok()) return r.error;
    return write_raw(r.value + "\r\n");
}

void TcpComm::close() {
    stop_.store(true);
    if (sockfd_ >= 0) {
        LAYRZ_SHUTDOWN(sockfd_, 2 /* SHUT_RDWR */);
        LAYRZ_CLOSE(sockfd_);
        sockfd_ = -1;
    }
    if (listener_.joinable()) listener_.join();
}

Error TcpComm::write_raw(const std::string& frame) {
    if (sockfd_ < 0) return Error::ServerError;
    ssize_t sent = ::send(sockfd_, frame.data(), frame.size(), MSG_NOSIGNAL);
    return (sent == static_cast<ssize_t>(frame.size())) ? Error::Ok : Error::ServerError;
}

void TcpComm::listen_loop() {
    std::string buf;
    buf.reserve(4096);
    char tmp[4096];

    while (!stop_.load()) {
        ssize_t n = ::recv(sockfd_, tmp, sizeof(tmp), 0);
        if (n <= 0) break;
        buf.append(tmp, static_cast<size_t>(n));

        while (true) {
            size_t close_pos = std::string::npos;
            size_t close_len = 0;

            for (size_t i = 0; i + 5 <= buf.size(); ++i) {
                if (buf[i]   == '<' && buf[i+1] == '/' &&
                    std::isupper(static_cast<unsigned char>(buf[i+2])) &&
                    std::islower(static_cast<unsigned char>(buf[i+3])) &&
                    buf[i+4] == '>') {
                    close_pos = i;
                    close_len = 5;
                    break;
                }
            }
            if (close_pos == std::string::npos) break;

            std::string close_tag = buf.substr(close_pos, close_len);
            std::string open_tag  = "<" + close_tag.substr(2, 2) + ">";
            size_t open_pos = buf.rfind(open_tag, close_pos);
            if (open_pos == std::string::npos) {
                buf.erase(0, close_pos + close_len);
                continue;
            }

            std::string frame = buf.substr(open_pos, close_pos + close_len - open_pos);
            buf.erase(0, close_pos + close_len);

            auto parsed = handle_server_output(frame);
            if (!parsed.ok()) continue;

            if (std::holds_alternative<packets::AsPacket>(parsed.value)) {
                authenticated_.store(true);
                continue;
            }
            if (std::holds_alternative<packets::AuPacket>(parsed.value)) {
                authenticated_.store(true);
                continue;
            }

            if (callback_) callback_(std::move(parsed.value));
        }
    }
}

#endif // LAYRZ_ARDUINO_STUB

} // namespace layrz::protocol::transport

#endif // LAYRZ_PROTOCOL_CLIENTS
