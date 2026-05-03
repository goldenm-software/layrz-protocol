"""Test ImPacket with semicolon escaping"""

from datetime import datetime
from uuid import UUID

import pytest

from layrz_protocol.constants import UTC
from layrz_protocol.packets import ImPacket
from layrz_protocol.utils import calculate_crc
from layrz_protocol.utils.exceptions import CrcException, MalformedException


def test10_im_packet_basic() -> None:
  """Test ImPacket basic roundtrip"""
  chat_id = UUID('12345678-1234-5678-1234-567812345678')
  timestamp = datetime.fromtimestamp(1700000000, tz=UTC)
  message = 'Hello World'

  payload = f'{int(timestamp.timestamp())};{str(chat_id)};{message};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Im>{payload}{crc}</Im>'
  msg = ImPacket.from_packet(payload_with_crc)

  assert msg.chat_id == chat_id
  assert msg.timestamp == timestamp
  assert msg.message == message
  assert msg.to_packet() == payload_with_crc


def test10_im_packet_with_semicolons() -> None:
  """Test ImPacket with semicolons in message (should be escaped as |||)"""
  chat_id = UUID('87654321-4321-8765-4321-876543218765')
  timestamp = datetime.fromtimestamp(1700000000, tz=UTC)
  message = 'Hello; World; Test'

  # Message semicolons are escaped as ||| when building the packet
  escaped_msg = message.replace(';', '|||')
  payload = f'{int(timestamp.timestamp())};{str(chat_id)};{escaped_msg};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Im>{payload}{crc}</Im>'
  msg = ImPacket.from_packet(payload_with_crc)

  assert msg.message == message
  assert msg.to_packet() == payload_with_crc


def test10_im_packet_multiple_semicolons() -> None:
  """Test ImPacket with multiple consecutive semicolons"""
  chat_id = UUID('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
  timestamp = datetime.fromtimestamp(1700000000, tz=UTC)
  message = ';;; critical ;;;'

  escaped_msg = message.replace(';', '|||')
  payload = f'{int(timestamp.timestamp())};{str(chat_id)};{escaped_msg};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Im>{payload}{crc}</Im>'
  msg = ImPacket.from_packet(payload_with_crc)

  assert msg.message == message
  assert msg.to_packet() == payload_with_crc


def test10_im_packet_long_message() -> None:
  """Test ImPacket with long message"""
  chat_id = UUID('11111111-2222-3333-4444-555555555555')
  timestamp = datetime.fromtimestamp(1700000000, tz=UTC)
  message = 'This is a long message ' * 10

  escaped_msg = message.replace(';', '|||')
  payload = f'{int(timestamp.timestamp())};{str(chat_id)};{escaped_msg};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Im>{payload}{crc}</Im>'
  msg = ImPacket.from_packet(payload_with_crc)

  assert msg.message == message
  assert msg.to_packet() == payload_with_crc


def test10_im_packet_empty_message() -> None:
  """Test ImPacket with empty message"""
  chat_id = UUID('99999999-8888-7777-6666-555555555555')
  timestamp = datetime.fromtimestamp(1700000000, tz=UTC)
  message = ''

  payload = f'{int(timestamp.timestamp())};{str(chat_id)};{message};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Im>{payload}{crc}</Im>'
  msg = ImPacket.from_packet(payload_with_crc)

  assert msg.message == message
  assert msg.to_packet() == payload_with_crc


def test10_im_packet_special_chars() -> None:
  """Test ImPacket with special characters"""
  chat_id = UUID('deadbeef-dead-beef-dead-beefdeadbeef')
  timestamp = datetime.fromtimestamp(1700000000, tz=UTC)
  message = 'Test!@#$%^&*()[]{}=-+_~`'

  escaped_msg = message.replace(';', '|||')
  payload = f'{int(timestamp.timestamp())};{str(chat_id)};{escaped_msg};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Im>{payload}{crc}</Im>'
  msg = ImPacket.from_packet(payload_with_crc)

  assert msg.message == message
  assert msg.to_packet() == payload_with_crc


def test10_im_packet_unicode() -> None:
  """Test ImPacket with unicode characters"""
  chat_id = UUID('cafebabe-cafe-babe-cafe-babecafebabe')
  timestamp = datetime.fromtimestamp(1700000000, tz=UTC)
  message = 'Hello 世界 مرحبا мир 🚀'

  escaped_msg = message.replace(';', '|||')
  payload = f'{int(timestamp.timestamp())};{str(chat_id)};{escaped_msg};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Im>{payload}{crc}</Im>'
  msg = ImPacket.from_packet(payload_with_crc)

  assert msg.message == message
  assert msg.to_packet() == payload_with_crc


def test10_im_packet_invalid_format() -> None:
  """Test ImPacket with invalid format (missing closing tag)"""
  with pytest.raises(MalformedException):
    ImPacket.from_packet('<Im>1700000000;12345678-1234-5678-1234-567812345678;message;0000')


def test10_im_packet_invalid_opening_tag() -> None:
  """Test ImPacket with invalid opening tag"""
  with pytest.raises(MalformedException):
    ImPacket.from_packet('<Ix>1700000000;12345678-1234-5678-1234-567812345678;message;0000</Im>')


def test10_im_packet_wrong_parts_count() -> None:
  """Test ImPacket with wrong number of parts"""
  with pytest.raises(MalformedException):
    ImPacket.from_packet('<Im>1700000000;12345678-1234-5678-1234-567812345678;0000</Im>')


def test10_im_packet_invalid_timestamp() -> None:
  """Test ImPacket with invalid timestamp"""
  with pytest.raises(MalformedException):
    payload = 'notanumber;12345678-1234-5678-1234-567812345678;message;'
    crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)
    ImPacket.from_packet(f'<Im>{payload}{crc}</Im>')


def test10_im_packet_invalid_uuid() -> None:
  """Test ImPacket with invalid UUID"""
  with pytest.raises(MalformedException):
    payload = '1700000000;not-a-valid-uuid;message;'
    crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)
    ImPacket.from_packet(f'<Im>{payload}{crc}</Im>')


def test10_im_packet_bad_crc() -> None:
  """Test ImPacket with bad CRC"""
  with pytest.raises(CrcException):
    ImPacket.from_packet('<Im>1700000000;12345678-1234-5678-1234-567812345678;message;FFFF</Im>')


def test10_im_packet_escaping_roundtrip() -> None:
  """Test that semicolon escaping roundtrips correctly"""
  # Create packet with semicolons in message
  original_message = 'A;B;C;D'
  packet = ImPacket(
    chat_id=UUID('12345678-1234-5678-1234-567812345678'),
    timestamp=datetime.fromtimestamp(1700000000, tz=UTC),
    message=original_message
  )

  # Convert to packet string
  packet_str = packet.to_packet()

  # Parse back
  parsed = ImPacket.from_packet(packet_str)

  # Message should be identical
  assert parsed.message == original_message
  assert parsed.message == 'A;B;C;D'
