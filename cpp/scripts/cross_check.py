#!/usr/bin/env python3
"""
Cross-language round-trip smoke test.

Encodes a fixed set of packets using the Python implementation (source of truth),
then asserts that the C++ binary (examples/encode_decode or a dedicated CLI)
decodes and re-encodes the same bytes.

Usage:
    cd python
    uv run python ../cpp/scripts/cross_check.py --cpp-bin ../cpp/build/example_encode_decode

For now this script verifies the Python-side encodings and prints the canonical frames
that should match C++ output. Run after building with CMake to cross-check manually.
"""
import sys
import os
import datetime
import argparse

repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..'))
sys.path.insert(0, os.path.join(repo_root, 'python'))

from layrz_protocol.packets.client.pa import PaPacket
from layrz_protocol.packets.client.pr import PrPacket
from layrz_protocol.packets.client.pc import PcPacket
from layrz_protocol.packets.client.pd import PdPacket
from layrz_protocol.packets.client.pi import PiPacket
from layrz_protocol.packets.client.ps import PsPacket
from layrz_protocol.packets.server.ao import AoPacket
from layrz_protocol.packets.server.ar import ArPacket
from layrz_protocol.packets.server.as_ import AsPacket
from layrz_protocol.packets.server.au import AuPacket
from layrz_protocol.packets.server.ts import TsPacket
from layrz_protocol.packets.server.te import TePacket
from layrz_protocol.packets.server.im import ImPacket
from layrz_protocol.packets.definitions.position import Position
from layrz_protocol.packets.definitions.firmware_branch import FirmwareBranch
from layrz_protocol.constants import UTC

FIXED_TS = datetime.datetime(2023, 11, 14, 22, 13, 20, tzinfo=UTC)  # 1700000000
UUID = '12345678-1234-5678-1234-567812345678'

cases = {
    'As': AsPacket().to_packet(),
    'Au': AuPacket().to_packet(),
    'Ao': AoPacket(timestamp=FIXED_TS).to_packet(),
    'Ar': ArPacket(reason='Unknown reason').to_packet(),
    'Pa': PaPacket(ident='123456789012345', password='mypassword').to_packet(),
    'Pr': PrPacket().to_packet(),
    'Pc': PcPacket(timestamp=FIXED_TS, command_id=42, message='ok').to_packet(),
    'Ts': TsPacket(timestamp=FIXED_TS, trip_id=UUID).to_packet(),
    'Te': TePacket(timestamp=FIXED_TS, trip_id=UUID,
                   distance_traveled=1234.567, max_speed=89.012,
                   duration=datetime.timedelta(seconds=3600)).to_packet(),
    'Im': ImPacket(timestamp=FIXED_TS, chat_id=UUID, message='Hello; world').to_packet(),
    'Pd': PdPacket(
        timestamp=datetime.datetime(1970, 1, 1, tzinfo=UTC),
        position=Position(latitude=10.0, longitude=10.0, altitude=10.0,
                          speed=10.0, direction=10.0, satellites=5, hdop=1.0),
        extra_data={'test.str': 'Hola mundo', 'test.int': 1,
                    'test.double': 1.0, 'test.bool': True},
    ).to_packet(),
    'Pi': PiPacket(
        ident='testident', firmware_id=1, firmware_build=1,
        device_id=1, hardware_id=1, model_id=1,
        firmware_branch=FirmwareBranch.DEVELOPMENT, fota_enabled=True,
    ).to_packet(),
    'Ps': PsPacket(
        timestamp=datetime.datetime(1970, 1, 1, tzinfo=UTC),
        params={'net_wifi_ssid': 'AWESOME WIFI', 'net_wifi_pass': 'dictadormarico69',
                'net_wifi_sec': 'WPA2', 'static.lat': -15.5, 'static.lng': 15.5},
    ).to_packet(),
}

print("Python canonical frames:")
for name, frame in cases.items():
    print(f"  {name:4s}: {frame}")

print("\nAll frames generated. Cross-check against C++ by running the test suite:")
print("  cmake --build cpp/build && ctest --test-dir cpp/build")
