"""Test Server Packets - As, Au"""

import pytest

from layrz_protocol.packets import AsPacket, AuPacket, ArPacket
from layrz_protocol.utils import calculate_crc
from layrz_protocol.utils.exceptions import CrcException, MalformedException


def test8_as_packet_basic() -> None:
  """Test AsPacket basic roundtrip"""
  payload = ';'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<As>{payload}{crc}</As>'
  msg: AsPacket = AsPacket.from_packet(payload_with_crc)

  assert msg.to_packet() == payload_with_crc


def test8_as_packet_bad_crc() -> None:
  """Test AsPacket with bad CRC"""
  with pytest.raises(CrcException):
    AsPacket.from_packet('<As>;FFFF</As>')


def test8_as_packet_invalid_format() -> None:
  """Test AsPacket with invalid format"""
  with pytest.raises(MalformedException):
    AsPacket.from_packet('<Ax>;0000</Ax>')


def test8_au_packet_basic() -> None:
  """Test AuPacket basic roundtrip"""
  payload = ';'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Au>{payload}{crc}</Au>'
  msg: AuPacket = AuPacket.from_packet(payload_with_crc)

  assert msg.to_packet() == payload_with_crc


def test8_au_packet_bad_crc() -> None:
  """Test AuPacket with bad CRC"""
  with pytest.raises(CrcException):
    AuPacket.from_packet('<Au>;FFFF</Au>')


def test8_au_packet_invalid_format() -> None:
  """Test AuPacket with invalid format"""
  with pytest.raises(MalformedException):
    AuPacket.from_packet('notapacket')


def test8_ar_packet_custom_reason() -> None:
  """Test ArPacket with custom reason text"""
  reason = 'Device offline for 5 minutes'
  payload = f'{reason};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Ar>{payload}{crc}</Ar>'
  msg: ArPacket = ArPacket.from_packet(payload_with_crc)

  assert msg.reason == reason
  assert msg.to_packet() == payload_with_crc


def test8_ar_packet_long_reason() -> None:
  """Test ArPacket with long reason text"""
  reason = 'Communication lost. Last known position: 10.5°N, 20.3°W. Attempting reconnection...'
  payload = f'{reason};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Ar>{payload}{crc}</Ar>'
  msg: ArPacket = ArPacket.from_packet(payload_with_crc)

  assert msg.reason == reason
  assert msg.to_packet() == payload_with_crc


def test8_ar_packet_empty_reason() -> None:
  """Test ArPacket with empty reason (should use default)"""
  payload = ';'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Ar>{payload}{crc}</Ar>'
  msg: ArPacket = ArPacket.from_packet(payload_with_crc)

  assert msg.reason == ''
  assert msg.to_packet() == payload_with_crc
