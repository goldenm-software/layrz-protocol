"""Test Client Packets"""

from datetime import datetime
from uuid import UUID

import pytest

from layrz_protocol.constants import UTC
from layrz_protocol.packets import TePacket, TsPacket


def test5_pt() -> None:
  packet = TsPacket(
    trip_id=UUID('123e4567-e89b-12d3-a456-426614174000'),
    timestamp=datetime(2025, 1, 1, 0, 0, 0, tzinfo=UTC),
  )

  str_packet = packet.to_packet()
  assert isinstance(str_packet, str)
  assert str_packet == '<Ts>1735689600;123e4567-e89b-12d3-a456-426614174000;696E</Ts>'

  parsed_packet = TsPacket.from_packet(str_packet)
  assert isinstance(parsed_packet, TsPacket)
  assert parsed_packet.trip_id == packet.trip_id
  assert parsed_packet.timestamp == packet.timestamp


def test5_pe() -> None:
  packet = TePacket(
    trip_id=UUID('123e4567-e89b-12d3-a456-426614174000'),
    timestamp=datetime(2025, 1, 1, 0, 0, 0, tzinfo=UTC),
  )

  str_packet = packet.to_packet()
  assert isinstance(str_packet, str)
  assert str_packet == '<Te>1735689600;123e4567-e89b-12d3-a456-426614174000;696E</Te>'

  parsed_packet = TePacket.from_packet(str_packet)
  assert isinstance(parsed_packet, TePacket)
  assert parsed_packet.trip_id == packet.trip_id
  assert parsed_packet.timestamp == packet.timestamp


if __name__ == '__main__':
  pytest.main([__file__])
