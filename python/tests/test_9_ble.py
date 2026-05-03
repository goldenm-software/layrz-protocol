"""Test BLE Advertisement, BleManufacturerData, BleServiceData"""

import pytest

from layrz_protocol.packets import BleAdvertisement, BleManufacturerData, BleServiceData
from layrz_protocol.utils import calculate_crc
from layrz_protocol.utils.exceptions import CrcException, MalformedException


def test9_ble_manufacturer_data_basic() -> None:
  """Test BleManufacturerData basic roundtrip"""
  mf_data = BleManufacturerData(company_id=0x004C, data=[0x02, 0x01, 0x06])
  packet = mf_data.to_packet()

  assert packet == '004C:020106'

  parsed = BleManufacturerData.from_packet(packet)
  assert parsed is not None
  assert parsed.company_id == 0x004C
  assert parsed.data == [0x02, 0x01, 0x06]


def test9_ble_manufacturer_data_empty() -> None:
  """Test BleManufacturerData with no data"""
  mf_data = BleManufacturerData(company_id=0x0001, data=[])
  packet = mf_data.to_packet()

  assert packet == '0001:'

  parsed = BleManufacturerData.from_packet(packet)
  assert parsed is not None
  assert parsed.company_id == 0x0001
  assert parsed.data == []


def test9_ble_manufacturer_data_from_empty_string() -> None:
  """Test BleManufacturerData.from_packet with empty string returns None"""
  result = BleManufacturerData.from_packet('')
  assert result is None


def test9_ble_manufacturer_data_invalid_format() -> None:
  """Test BleManufacturerData with invalid format"""
  with pytest.raises(MalformedException):
    BleManufacturerData.from_packet('004C:020106:extra')


def test9_ble_manufacturer_data_invalid_company_id() -> None:
  """Test BleManufacturerData with invalid company ID"""
  with pytest.raises(MalformedException):
    BleManufacturerData.from_packet('XXXX:020106')


def test9_ble_service_data_basic() -> None:
  """Test BleServiceData basic roundtrip"""
  svc_data = BleServiceData(uuid=0x180A, data=[0x00, 0x01, 0x02])
  packet = svc_data.to_packet()

  assert packet == '180A:000102'

  parsed = BleServiceData.from_packet(packet)
  assert parsed is not None
  assert parsed.uuid == 0x180A
  assert parsed.data == [0x00, 0x01, 0x02]


def test9_ble_service_data_empty() -> None:
  """Test BleServiceData with no data"""
  svc_data = BleServiceData(uuid=0x180D, data=[])
  packet = svc_data.to_packet()

  assert packet == '180D:'

  parsed = BleServiceData.from_packet(packet)
  assert parsed is not None
  assert parsed.uuid == 0x180D
  assert parsed.data == []


def test9_ble_service_data_from_empty_string() -> None:
  """Test BleServiceData.from_packet with empty string returns None"""
  result = BleServiceData.from_packet('')
  assert result is None


def test9_ble_service_data_invalid_format() -> None:
  """Test BleServiceData with invalid format"""
  with pytest.raises(MalformedException):
    BleServiceData.from_packet('180A:000102:extra')


def test9_ble_service_data_invalid_uuid() -> None:
  """Test BleServiceData with invalid UUID"""
  with pytest.raises(MalformedException):
    BleServiceData.from_packet('XXXX:000102')


def test9_ble_advertisement_no_position() -> None:
  """Test BleAdvertisement without position data"""
  raw = '001122334455;1762214400;;;;MODEL;DevName;-50;;004C:020106;180A:000102;'
  crc = str(hex(calculate_crc(raw.encode())))[2:].upper().zfill(4)

  raw_with_crc = f'{raw}{crc}'

  adv = BleAdvertisement.from_packet(raw_with_crc)

  assert adv.mac_address == '00:11:22:33:44:55'
  assert adv.timestamp.timestamp() == 1762214400
  assert adv.latitude is None
  assert adv.longitude is None
  assert adv.altitude is None
  assert adv.model == 'MODEL'
  assert adv.device_name == 'DevName'
  assert adv.rssi == -50
  assert adv.tx_power is None
  assert len(adv.manufacturer_data) == 1
  assert len(adv.service_data) == 1

  assert adv.to_packet() == raw_with_crc


def test9_ble_advertisement_with_position() -> None:
  """Test BleAdvertisement with full position data"""
  raw = 'AABBCCDDEEFF;1762214400;10.5;-20.3;100.0;MODEL;Device;-60;5;;;'
  crc = str(hex(calculate_crc(raw.encode())))[2:].upper().zfill(4)

  raw_with_crc = f'{raw}{crc}'

  adv = BleAdvertisement.from_packet(raw_with_crc)

  assert adv.mac_address == 'AA:BB:CC:DD:EE:FF'
  assert adv.latitude == 10.5
  assert adv.longitude == -20.3
  assert adv.altitude == 100.0
  assert adv.rssi == -60
  assert adv.tx_power == 5

  assert adv.to_packet() == raw_with_crc


def test9_ble_advertisement_multiple_mfr_data() -> None:
  """Test BleAdvertisement with multiple manufacturer data entries"""
  raw = '112233445566;1762214400;;;;MODEL;Name;-70;;004C:0201,0005:AABBCC;180A:00;'
  crc = str(hex(calculate_crc(raw.encode())))[2:].upper().zfill(4)

  raw_with_crc = f'{raw}{crc}'

  adv = BleAdvertisement.from_packet(raw_with_crc)

  assert len(adv.manufacturer_data) == 2
  assert adv.manufacturer_data[0].company_id == 0x004C
  assert adv.manufacturer_data[1].company_id == 0x0005

  assert adv.to_packet() == raw_with_crc


def test9_ble_advertisement_multiple_svc_data() -> None:
  """Test BleAdvertisement with multiple service data entries"""
  raw = '223344556677;1762214400;;;;MODEL;Name;-55;;004C:01;180A:00,180D:FF;'
  crc = str(hex(calculate_crc(raw.encode())))[2:].upper().zfill(4)

  raw_with_crc = f'{raw}{crc}'

  adv = BleAdvertisement.from_packet(raw_with_crc)

  assert len(adv.service_data) == 2
  assert adv.service_data[0].uuid == 0x180A
  assert adv.service_data[1].uuid == 0x180D

  assert adv.to_packet() == raw_with_crc


def test9_ble_advertisement_invalid_parts() -> None:
  """Test BleAdvertisement with wrong number of parts"""
  with pytest.raises(MalformedException):
    BleAdvertisement.from_packet('001122334455;1762214400;;')


def test9_ble_advertisement_invalid_mac() -> None:
  """Test BleAdvertisement with invalid MAC address length"""
  with pytest.raises(MalformedException):
    raw = '0011223344;1762214400;;;;MODEL;Name;-50;;;'
    crc = str(hex(calculate_crc(raw.encode())))[2:].upper().zfill(4)
    BleAdvertisement.from_packet(f'{raw}{crc}')


def test9_ble_advertisement_invalid_timestamp() -> None:
  """Test BleAdvertisement with invalid timestamp"""
  with pytest.raises(MalformedException):
    raw = '001122334455;notanumber;;;;MODEL;Name;-50;;;'
    crc = str(hex(calculate_crc(raw.encode())))[2:].upper().zfill(4)
    BleAdvertisement.from_packet(f'{raw}{crc}')


def test9_ble_advertisement_invalid_rssi() -> None:
  """Test BleAdvertisement with invalid RSSI"""
  with pytest.raises(MalformedException):
    raw = '001122334455;1762214400;;;;MODEL;Name;badrssi;;;'
    crc = str(hex(calculate_crc(raw.encode())))[2:].upper().zfill(4)
    BleAdvertisement.from_packet(f'{raw}{crc}')


def test9_ble_advertisement_invalid_tx_power() -> None:
  """Test BleAdvertisement with invalid TX power"""
  with pytest.raises(MalformedException):
    raw = '001122334455;1762214400;;;;MODEL;Name;-50;badpower;;'
    crc = str(hex(calculate_crc(raw.encode())))[2:].upper().zfill(4)
    BleAdvertisement.from_packet(f'{raw}{crc}')


def test9_ble_advertisement_bad_crc() -> None:
  """Test BleAdvertisement with bad CRC"""
  with pytest.raises(CrcException):
    BleAdvertisement.from_packet('001122334455;1762214400;;;;MODEL;Name;-50;;;;FFFF')


def test9_ble_advertisement_partial_position() -> None:
  """Test BleAdvertisement with partial position (lat/lng but no alt)"""
  raw = '334455667788;1762214400;5.5;10.5;;MODEL;Device;-60;;;;'
  crc = str(hex(calculate_crc(raw.encode())))[2:].upper().zfill(4)

  raw_with_crc = f'{raw}{crc}'

  adv = BleAdvertisement.from_packet(raw_with_crc)

  assert adv.latitude == 5.5
  assert adv.longitude == 10.5
  assert adv.altitude is None

  assert adv.to_packet() == raw_with_crc
