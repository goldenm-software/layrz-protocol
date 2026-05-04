"""Test TCP Server"""

from __future__ import annotations

import asyncio
import socket as _socket

import pytest

from layrz_protocol.packets.client import PaPacket
from layrz_protocol.packets.server import AsPacket, AoPacket, ServerPacket
from layrz_protocol.servers import TcpConfig, TcpServer


def _free_port() -> int:
  with _socket.socket(_socket.AF_INET, _socket.SOCK_STREAM) as s:
    s.bind(('127.0.0.1', 0))
    return s.getsockname()[1]


async def test1_invalid_port_zero() -> None:
  async def handler(pkt: object, peer: object) -> ServerPacket | None:
    return None

  with pytest.raises(ValueError):
    TcpServer(TcpConfig(port=0, on_new_packet=handler))


async def test2_invalid_port_negative() -> None:
  async def handler(pkt: object, peer: object) -> ServerPacket | None:
    return None

  with pytest.raises(ValueError):
    TcpServer(TcpConfig(port=-1, on_new_packet=handler))


async def test3_invalid_port_65535() -> None:
  async def handler(pkt: object, peer: object) -> ServerPacket | None:
    return None

  with pytest.raises(ValueError):
    TcpServer(TcpConfig(port=65535, on_new_packet=handler))


async def test4_invalid_port_too_high() -> None:
  async def handler(pkt: object, peer: object) -> ServerPacket | None:
    return None

  with pytest.raises(ValueError):
    TcpServer(TcpConfig(port=70000, on_new_packet=handler))


async def test5_missing_on_new_packet() -> None:
  port = _free_port()

  with pytest.raises(ValueError):
    TcpServer(TcpConfig(port=port, on_new_packet=None))  # type: ignore


async def test6_happy_path_round_trip() -> None:
  port = _free_port()

  async def handler(pkt: object, peer: object) -> ServerPacket | None:
    if isinstance(pkt, PaPacket):
      return AsPacket()
    return None

  server = TcpServer(TcpConfig(port=port, on_new_packet=handler))
  task = asyncio.create_task(server.start())
  await asyncio.sleep(0.1)
  try:
    reader, writer = await asyncio.open_connection('127.0.0.1', port)
    pkt = PaPacket(ident='dev', password='pw').to_packet() + '\n'
    writer.write(pkt.encode())
    await writer.drain()
    data = await asyncio.wait_for(reader.read(256), timeout=2.0)
    assert data.decode() == AsPacket().to_packet()
    writer.close()
    await writer.wait_closed()
  finally:
    task.cancel()
    await server.close()
    await asyncio.gather(task, return_exceptions=True)


async def test7_multi_frame_split() -> None:
  port = _free_port()
  call_count: list[int] = [0]

  async def handler(pkt: object, peer: object) -> ServerPacket | None:
    call_count[0] += 1
    if isinstance(pkt, PaPacket):
      return AsPacket()
    return None

  server = TcpServer(TcpConfig(port=port, on_new_packet=handler))
  task = asyncio.create_task(server.start())
  await asyncio.sleep(0.1)
  try:
    reader, writer = await asyncio.open_connection('127.0.0.1', port)
    pkt1 = PaPacket(ident='dev1', password='pw1').to_packet()
    pkt2 = PaPacket(ident='dev2', password='pw2').to_packet()
    writer.write((pkt1 + '\n' + pkt2 + '\n').encode())
    await writer.drain()
    data = await asyncio.wait_for(reader.read(512), timeout=2.0)
    assert call_count[0] == 2
    assert data.decode().count('<As>') == 2
    writer.close()
    await writer.wait_closed()
  finally:
    task.cancel()
    await server.close()
    await asyncio.gather(task, return_exceptions=True)


async def test8_decode_error_callback() -> None:
  port = _free_port()
  error_called: list[bool] = [False]

  async def on_decode_error(exc: Exception, data: bytes, peer: object) -> None:
    error_called[0] = True

  async def handler(pkt: object, peer: object) -> ServerPacket | None:
    if isinstance(pkt, PaPacket):
      return AsPacket()
    return None

  server = TcpServer(
    TcpConfig(port=port, on_new_packet=handler, on_decode_error=on_decode_error)
  )
  task = asyncio.create_task(server.start())
  await asyncio.sleep(0.1)
  try:
    reader, writer = await asyncio.open_connection('127.0.0.1', port)
    writer.write(b'garbage\n')
    await writer.drain()
    await asyncio.sleep(0.1)
    assert error_called[0]
    pkt = PaPacket(ident='dev', password='pw').to_packet() + '\n'
    writer.write(pkt.encode())
    await writer.drain()
    data = await asyncio.wait_for(reader.read(256), timeout=2.0)
    assert data.decode() == AsPacket().to_packet()
    writer.close()
    await writer.wait_closed()
  finally:
    task.cancel()
    await server.close()
    await asyncio.gather(task, return_exceptions=True)


async def test9_close_stops_server() -> None:
  port = _free_port()

  async def handler(pkt: object, peer: object) -> ServerPacket | None:
    return None

  server = TcpServer(TcpConfig(port=port, on_new_packet=handler))
  task = asyncio.create_task(server.start())
  await asyncio.sleep(0.1)
  await server.close()
  try:
    await asyncio.wait_for(task, timeout=2.0)
  except asyncio.CancelledError:
    pass
