#include "layrz_protocol/layrz_protocol.hpp"
#include <iostream>

int main() {
    using namespace layrz::protocol;
    using namespace layrz::protocol::packets;

    // Encode a Pa (auth) packet
    PaPacket pa;
    pa.ident    = "mydevice123";
    pa.password = "secret";
    std::cout << "Pa: " << pa.to_packet() << "\n";

    // Encode a Pr (keepalive) packet
    PrPacket pr;
    std::cout << "Pr: " << pr.to_packet() << "\n";

    // Encode a Pi (identification) packet
    PiPacket pi;
    pi.ident           = "mydevice123";
    pi.firmware_id     = 42;
    pi.firmware_build  = 1;
    pi.device_id       = 100;
    pi.hardware_id     = 200;
    pi.model_id        = 300;
    pi.firmware_branch = FirmwareBranch::Stable;
    pi.fota_enabled    = true;
    std::cout << "Pi: " << pi.to_packet() << "\n";

    // Decode a server As packet
    std::string as_frame = "<As>;7F28</As>";
    auto r = handle_server_output(as_frame);
    if (r.ok() && std::holds_alternative<AsPacket>(r.value)) {
        std::cout << "Decoded As packet (auth success)\n";
    }

    // Decode an Ar (error) packet
    std::string ar_frame = "<Ar>Unknown reason;1DFD</Ar>";
    auto r2 = handle_server_output(ar_frame);
    if (r2.ok() && std::holds_alternative<ArPacket>(r2.value)) {
        std::cout << "Ar reason: " << std::get<ArPacket>(r2.value).reason << "\n";
    }

    return 0;
}
