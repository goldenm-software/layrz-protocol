"""Test parse_extra key-translation branches and convert_to_dotcase."""

import pytest

from layrz_protocol.utils.parsers import cast_extra, convert_to_dotcase, parse_extra


@pytest.mark.parametrize('raw,expected_key,expected_val', [
  # Digital / Analog GPIO branches
  ('io1.di:1', 'gpio.1.digital.input', True),
  ('io2.do:0', 'gpio.2.digital.output', False),
  ('io3.ai:3.5', 'gpio.3.analog.input', 3.5),
  ('io4.ao:7', 'gpio.4.analog.output', 7),
  ('io5.counter:42', 'gpio.5.event.count', 42),
  # BLE branches
  ('ble.0.id:AABBCCDDEEFF', 'ble.0.mac.address', 'AABBCCDDEEFF'),
  ('ble.0.hum:55', 'ble.0.humidity', 55),
  ('ble.0.tempc:23.5', 'ble.0.temperature.celsius', 23.5),
  ('ble.0.tempf:74.3', 'ble.0.temperature.fahrenheit', 74.3),
  ('ble.0.model_id:9', 'ble.0.model.id', 9),
  ('ble.0.batt:80', 'ble.0.battery.level', 80),
  ('ble.0.lux:100', 'ble.0.light.level.lux', 100),
  ('ble.0.volt:3.7', 'ble.0.voltage', 3.7),
  ('ble.0.rpm:1200', 'ble.0.rpm', 1200),
  ('ble.0.press:1013', 'ble.0.pressure', 1013),
  ('ble.0.counter:5', 'ble.0.event.count', 5),
  ('ble.0.x_acc:0.1', 'ble.0.acceleration.x', 0.1),
  ('ble.0.y_acc:0.2', 'ble.0.acceleration.y', 0.2),
  ('ble.0.z_acc:0.3', 'ble.0.acceleration.z', 0.3),
  ('ble.0.msg_count:3', 'ble.0.message.count', 3),
  ('ble.0.msg:hello', 'ble.0.message', 'hello'),
  ('ble.0.mag_counter:2', 'ble.0.magnetic.event.count', 2),
  ('ble.0.mag_data:ff00', 'ble.0.magnetic.data', 'ff00'),
  ('ble.0.rssi:-70', 'ble.0.rssi.dbm', -70),
  # Well-known string keys
  ('report:5', 'report.code', 5),
  ('confiot_ble:1', 'ble.confiot.connection.status', 1),
  ('confiot_serial:0', 'serial.confiot.connection.status', 0),
])
def test13_parse_extra_branch(raw: str, expected_key: str, expected_val: object) -> None:
  result = parse_extra(raw)
  assert expected_key in result
  assert result[expected_key] == expected_val


def test13_parse_extra_empty() -> None:
  assert parse_extra('') == {}


def test13_parse_extra_unknown_key_str() -> None:
  result = parse_extra('foo:bar')
  assert result['foo'] == 'bar'


def test13_parse_extra_float_value() -> None:
  result = parse_extra('custom:3.14')
  assert result['custom'] == pytest.approx(3.14)


def test13_parse_extra_bool_true() -> None:
  result = parse_extra('flag:true')
  assert result['flag'] is True


def test13_parse_extra_bool_false() -> None:
  result = parse_extra('flag:false')
  assert result['flag'] is False


def test13_parse_extra_bool_t_f() -> None:
  # 't' converts to False (value.lower() == 'true' is False); 'f' also False
  assert parse_extra('a:t')['a'] is False
  assert parse_extra('a:f')['a'] is False


def test13_convert_to_dotcase_none() -> None:
  assert convert_to_dotcase('key', None) == {}


def test13_convert_to_dotcase_ascii_replacement() -> None:
  result = convert_to_dotcase('name', 'café')
  assert result == {'name': 'cafe'}


def test13_convert_to_dotcase_list() -> None:
  result = convert_to_dotcase('x', [1, 2])
  assert result == {'x.0': 1, 'x.1': 2}


def test13_convert_to_dotcase_dict() -> None:
  result = convert_to_dotcase('x', {'a': 1, 'b': 2})
  assert result == {'x.a': 1, 'x.b': 2}


def test13_cast_extra_bool_true() -> None:
  raw = cast_extra({'flag': True})
  assert 'flag:true' in raw


def test13_cast_extra_bool_false() -> None:
  raw = cast_extra({'flag': False})
  assert 'flag:false' in raw
