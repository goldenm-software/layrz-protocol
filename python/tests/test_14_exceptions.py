"""Test exception class __init__ and __str__ paths."""

from layrz_protocol.utils.exceptions import (
  CommandException,
  CrcException,
  MalformedException,
  ParseException,
  ServerException,
  UnimplementedException,
)


def test14_parse_exception() -> None:
  exc = ParseException('bad packet')
  assert exc.message == 'bad packet'
  assert str(exc) == 'ParseException: bad packet'


def test14_crc_exception() -> None:
  exc = CrcException('CRC mismatch', received=0xDEAD, calculated=0xBEEF)
  assert exc.message == 'CRC mismatch'
  assert exc.received == 0xDEAD
  assert exc.calculated == 0xBEEF
  text = str(exc)
  assert 'DEAD' in text
  assert 'BEEF' in text


def test14_command_exception() -> None:
  exc = CommandException('bad cmd')
  assert exc.message == 'bad cmd'
  assert str(exc) == 'bad cmd'


def test14_server_exception() -> None:
  exc = ServerException('server error')
  assert exc.message == 'server error'
  assert str(exc) == 'server error'


def test14_malformed_exception() -> None:
  exc = MalformedException('malformed')
  assert exc.message == 'malformed'
  assert str(exc) == 'malformed'


def test14_unimplemented_exception() -> None:
  exc = UnimplementedException('not implemented')
  assert exc.message == 'not implemented'
  assert str(exc) == 'not implemented'
