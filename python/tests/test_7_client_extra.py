"""Test Client Packets - Pa, Pm, Pr, Pb"""

import base64
from datetime import datetime

import pytest

from layrz_protocol.constants import UTC
from layrz_protocol.packets import PaPacket, PbPacket, PmPacket, PrPacket, BleAdvertisement
from layrz_protocol.utils import calculate_crc
from layrz_protocol.utils.exceptions import CrcException, MalformedException


def test7_pa_packet_basic() -> None:
  """Test PaPacket basic roundtrip"""
  payload = 'device123;secret_pass;'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Pa>{payload}{crc}</Pa>'
  msg: PaPacket = PaPacket.from_packet(payload_with_crc)

  assert msg.ident == 'device123'
  assert msg.password == 'secret_pass'
  assert msg.to_packet() == payload_with_crc


def test7_pa_packet_empty_ident() -> None:
  """Test PaPacket with empty ident"""
  payload = ';pass123;'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Pa>{payload}{crc}</Pa>'
  msg: PaPacket = PaPacket.from_packet(payload_with_crc)

  assert msg.ident == ''
  assert msg.password == 'pass123'
  assert msg.to_packet() == payload_with_crc


def test7_pa_packet_invalid_format() -> None:
  """Test PaPacket with invalid format"""
  with pytest.raises(MalformedException):
    PaPacket.from_packet('<Pa>ident;pass</Pa>')


def test7_pa_packet_bad_crc() -> None:
  """Test PaPacket with bad CRC"""
  with pytest.raises(CrcException):
    PaPacket.from_packet('<Pa>ident;pass;FFFF</Pa>')


def test7_pm_packet_basic() -> None:
  """Test PmPacket basic roundtrip"""
  data = b'Hello, World!'
  data_b64 = base64.b64encode(data).decode()

  payload = f'test.txt;text/plain;{data_b64};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Pm>{payload}{crc}</Pm>'
  msg: PmPacket = PmPacket.from_packet(payload_with_crc)

  assert msg.filename == 'test.txt'
  assert msg.content_type == 'text/plain'
  assert msg.data == data
  assert msg.to_packet() == payload_with_crc


def test7_pm_packet_binary_data() -> None:
  """Test PmPacket with binary data"""
  data = bytes([0x00, 0x01, 0x02, 0xFF, 0xFE])
  data_b64 = base64.b64encode(data).decode()

  payload = f'image.bin;application/octet-stream;{data_b64};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Pm>{payload}{crc}</Pm>'
  msg: PmPacket = PmPacket.from_packet(payload_with_crc)

  assert msg.filename == 'image.bin'
  assert msg.content_type == 'application/octet-stream'
  assert msg.data == data
  assert msg.to_packet() == payload_with_crc


def test7_pm_packet_empty_data() -> None:
  """Test PmPacket with empty data"""
  data = b''
  data_b64 = base64.b64encode(data).decode()

  payload = f'empty.bin;application/octet-stream;{data_b64};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Pm>{payload}{crc}</Pm>'
  msg: PmPacket = PmPacket.from_packet(payload_with_crc)

  assert msg.filename == 'empty.bin'
  assert msg.data == data
  assert msg.to_packet() == payload_with_crc


def test7_pm_packet_wrong_parts_count() -> None:
  """Test PmPacket with wrong number of parts"""
  with pytest.raises(MalformedException):
    PmPacket.from_packet('<Pm>file.txt;text/plain;0000</Pm>')


def test7_pr_packet_basic() -> None:
  """Test PrPacket to_packet roundtrip via to_packet"""
  msg = PrPacket()
  pkt = msg.to_packet()
  assert pkt.startswith('<Pr>')
  assert pkt.endswith('</Pr>')


def test7_pr_packet_bad_crc() -> None:
  """Test PrPacket with bad CRC raises CrcException"""
  inner = ';'
  crc = str(hex(calculate_crc(inner.encode())))[2:].upper().zfill(4)
  good_pkt = f'<Pr>{inner}{crc}</Pr>'
  bad_pkt = good_pkt[:-5] + 'FFFF' + good_pkt[-5:]
  with pytest.raises((CrcException, MalformedException)):
    PrPacket.from_packet(bad_pkt)


def test7_pr_packet_invalid_format() -> None:
  """Test PrPacket with invalid format"""
  with pytest.raises(MalformedException):
    PrPacket.from_packet('<Pr>extra;data;0000</Pr>')


def test7_pb_packet_single_advertisement() -> None:
  """Test PbPacket with single BLE advertisement"""
  ble = BleAdvertisement(
    mac_address='00:11:22:33:44:55',
    timestamp=datetime.fromtimestamp(1762214400, tz=UTC),
    model='MODEL',
    device_name='DevName',
    rssi=-50
  )

  ble_packet = ble.to_packet()
  payload = f'{ble_packet};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Pb>{payload}{crc}</Pb>'
  msg: PbPacket = PbPacket.from_packet(payload_with_crc)

  assert len(msg.advertisements) == 1
  assert msg.advertisements[0].mac_address == '00:11:22:33:44:55'
  assert msg.advertisements[0].model == 'MODEL'
  assert msg.advertisements[0].device_name == 'DevName'
  assert msg.advertisements[0].rssi == -50
  assert msg.to_packet() == payload_with_crc


def test7_pb_packet_with_position() -> None:
  """Test PbPacket with position data"""
  ble = BleAdvertisement(
    mac_address='AA:BB:CC:DD:EE:FF',
    timestamp=datetime.fromtimestamp(1762214400, tz=UTC),
    latitude=10.5,
    longitude=-20.3,
    altitude=100.0,
    model='MODEL',
    device_name='Device',
    rssi=-60,
    tx_power=5
  )

  ble_packet = ble.to_packet()
  payload = f'{ble_packet};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Pb>{payload}{crc}</Pb>'
  msg: PbPacket = PbPacket.from_packet(payload_with_crc)

  assert len(msg.advertisements) == 1
  assert msg.advertisements[0].latitude == 10.5
  assert msg.advertisements[0].longitude == -20.3
  assert msg.advertisements[0].altitude == 100.0
  assert msg.advertisements[0].tx_power == 5
  assert msg.to_packet() == payload_with_crc


def test7_pb_packet_multiple_advertisements() -> None:
  """Test PbPacket with multiple BLE advertisements"""
  ble1 = BleAdvertisement(
    mac_address='00:11:22:33:44:55',
    timestamp=datetime.fromtimestamp(1762214400, tz=UTC),
    model='MODEL1',
    device_name='Dev1',
    rssi=-50
  )

  ble2 = BleAdvertisement(
    mac_address='AA:BB:CC:DD:EE:FF',
    timestamp=datetime.fromtimestamp(1762214400, tz=UTC),
    model='MODEL2',
    device_name='Dev2',
    rssi=-60
  )

  ble_packet1 = ble1.to_packet()
  ble_packet2 = ble2.to_packet()

  payload = f'{ble_packet1};{ble_packet2};'
  crc = str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)

  payload_with_crc = f'<Pb>{payload}{crc}</Pb>'
  msg: PbPacket = PbPacket.from_packet(payload_with_crc)

  assert len(msg.advertisements) == 2
  assert msg.advertisements[0].mac_address == '00:11:22:33:44:55'
  assert msg.advertisements[1].mac_address == 'AA:BB:CC:DD:EE:FF'
  assert msg.to_packet() == payload_with_crc


def test7_pb_packet_invalid_parts_count() -> None:
  """Test PbPacket with invalid number of parts"""
  with pytest.raises(MalformedException):
    PbPacket.from_packet('<Pb>001122334455;1762214400;;;</Pb>')


def test7_pb_packet_bad_crc() -> None:
  """Test PbPacket with bad CRC"""
  ble = BleAdvertisement(
    mac_address='00:11:22:33:44:55',
    timestamp=datetime.fromtimestamp(1762214400, tz=UTC),
    model='MODEL',
    device_name='Device',
    rssi=-50
  )

  ble_packet = ble.to_packet()
  payload = f'{ble_packet};'

  with pytest.raises(CrcException):
    PbPacket.from_packet(f'<Pb>{payload}FFFF</Pb>')
