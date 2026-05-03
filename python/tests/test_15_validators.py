"""Test pydantic field validators in Position and BleAdvertisement."""

from typing import Any

import pytest
from pydantic import ValidationError

from layrz_protocol.packets.definitions.ble_advertisement import BleAdvertisement
from layrz_protocol.packets.definitions.position import Position


# ---------------------------------------------------------------------------
# Position validators
# ---------------------------------------------------------------------------

def test15_position_defaults() -> None:
  p = Position()
  assert p.latitude is None
  assert p.longitude is None


def test15_position_latitude_not_float() -> None:
  bad: Any = 'not_a_float'
  with pytest.raises(ValidationError):
    Position(latitude=bad)


def test15_position_latitude_out_of_range_high() -> None:
  with pytest.raises(ValidationError):
    Position(latitude=91.0)


def test15_position_latitude_out_of_range_low() -> None:
  with pytest.raises(ValidationError):
    Position(latitude=-91.0)


def test15_position_longitude_not_float() -> None:
  bad: Any = 'not_a_float'
  with pytest.raises(ValidationError):
    Position(longitude=bad)


def test15_position_longitude_out_of_range_high() -> None:
  with pytest.raises(ValidationError):
    Position(longitude=181.0)


def test15_position_longitude_out_of_range_low() -> None:
  with pytest.raises(ValidationError):
    Position(longitude=-181.0)


def test15_position_direction_not_numeric() -> None:
  bad: Any = 'north'
  with pytest.raises(ValidationError):
    Position(direction=bad)


def test15_position_direction_out_of_range_high() -> None:
  with pytest.raises(ValidationError):
    Position(direction=361.0)


def test15_position_direction_out_of_range_low() -> None:
  with pytest.raises(ValidationError):
    Position(direction=-1.0)


def test15_position_direction_int_accepted() -> None:
  p = Position(direction=180)
  assert p.direction == 180.0


def test15_position_hdop_not_float() -> None:
  bad: Any = 'bad'
  with pytest.raises(ValidationError):
    Position(hdop=bad)


def test15_position_hdop_negative() -> None:
  with pytest.raises(ValidationError):
    Position(hdop=-1.0)


# ---------------------------------------------------------------------------
# BleAdvertisement validators
# ---------------------------------------------------------------------------

def test15_ble_advertisement_timestamp_fallback() -> None:
  # Non-int non-datetime triggers the fallback branch → defaults to now()
  bad_ts: Any = 'bad'
  adv = BleAdvertisement(mac_address='AA:BB:CC:DD:EE:FF', model='TEST', timestamp=bad_ts)
  assert adv.timestamp is not None


def test15_ble_advertisement_latitude_not_float() -> None:
  bad: Any = 'x'
  with pytest.raises(ValidationError):
    BleAdvertisement(mac_address='AA:BB:CC:DD:EE:FF', model='TEST', latitude=bad)


def test15_ble_advertisement_latitude_out_of_range() -> None:
  with pytest.raises(ValidationError):
    BleAdvertisement(mac_address='AA:BB:CC:DD:EE:FF', model='TEST', latitude=91.0)


def test15_ble_advertisement_longitude_not_float() -> None:
  bad: Any = 'x'
  with pytest.raises(ValidationError):
    BleAdvertisement(mac_address='AA:BB:CC:DD:EE:FF', model='TEST', longitude=bad)


def test15_ble_advertisement_longitude_out_of_range() -> None:
  with pytest.raises(ValidationError):
    BleAdvertisement(mac_address='AA:BB:CC:DD:EE:FF', model='TEST', longitude=181.0)
