"""Test HTTP Server"""

from __future__ import annotations

import asyncio
import socket as _socket

import aiohttp
import pytest

from layrz_protocol.packets.client import PaPacket
from layrz_protocol.packets.server import AsPacket, ServerPacket
from layrz_protocol.servers import HttpConfig, HttpServer


def _free_port() -> int:
  with _socket.socket(_socket.AF_INET, _socket.SOCK_STREAM) as s:
    s.bind(('127.0.0.1', 0))
    return s.getsockname()[1]


async def test1_invalid_port_zero() -> None:
  async def handler(pkt: object, req: object) -> ServerPacket | None:
    return None

  with pytest.raises(ValueError):
    HttpServer(HttpConfig(port=0, on_new_packet=handler))


async def test2_invalid_port_negative() -> None:
  async def handler(pkt: object, req: object) -> ServerPacket | None:
    return None

  with pytest.raises(ValueError):
    HttpServer(HttpConfig(port=-1, on_new_packet=handler))


async def test3_invalid_port_65535() -> None:
  async def handler(pkt: object, req: object) -> ServerPacket | None:
    return None

  with pytest.raises(ValueError):
    HttpServer(HttpConfig(port=65535, on_new_packet=handler))


async def test4_missing_on_new_packet() -> None:
  port = _free_port()

  with pytest.raises(ValueError):
    HttpServer(HttpConfig(port=port, on_new_packet=None))  # type: ignore


async def test5_post_message_happy_path() -> None:
  port = _free_port()

  async def handler(pkt: object, req: object) -> ServerPacket | None:
    if isinstance(pkt, PaPacket):
      return AsPacket()
    return None

  server = HttpServer(HttpConfig(port=port, on_new_packet=handler))
  task = asyncio.create_task(server.start())
  await asyncio.sleep(0.15)
  try:
    async with aiohttp.ClientSession() as session:
      async with session.post(
        f'http://127.0.0.1:{port}/v2/message',
        headers={'Authorization': 'LayrzAuth dev;pw'},
        data=PaPacket(ident='dev', password='pw').to_packet(),
      ) as resp:
        assert resp.status == 200
        body = await resp.text()
        assert body == AsPacket().to_packet()
  finally:
    task.cancel()
    await server.close()
    await asyncio.gather(task, return_exceptions=True)


async def test6_missing_auth_401() -> None:
  port = _free_port()

  async def handler(pkt: object, req: object) -> ServerPacket | None:
    return None

  server = HttpServer(HttpConfig(port=port, on_new_packet=handler))
  task = asyncio.create_task(server.start())
  await asyncio.sleep(0.15)
  try:
    async with aiohttp.ClientSession() as session:
      async with session.post(
        f'http://127.0.0.1:{port}/v2/message',
        data=PaPacket(ident='dev', password='pw').to_packet(),
      ) as resp:
        assert resp.status == 401
  finally:
    task.cancel()
    await server.close()
    await asyncio.gather(task, return_exceptions=True)


async def test7_bad_auth_format_401() -> None:
  port = _free_port()

  async def handler(pkt: object, req: object) -> ServerPacket | None:
    return None

  server = HttpServer(HttpConfig(port=port, on_new_packet=handler))
  task = asyncio.create_task(server.start())
  await asyncio.sleep(0.15)
  try:
    async with aiohttp.ClientSession() as session:
      async with session.post(
        f'http://127.0.0.1:{port}/v2/message',
        headers={'Authorization': 'BadFormat'},
        data=PaPacket(ident='dev', password='pw').to_packet(),
      ) as resp:
        assert resp.status == 401
  finally:
    task.cancel()
    await server.close()
    await asyncio.gather(task, return_exceptions=True)


async def test8_on_authenticate_false_401() -> None:
  port = _free_port()

  async def authenticate(ident: str, passwd: str, req: object) -> bool:
    return False

  async def handler(pkt: object, req: object) -> ServerPacket | None:
    return None

  server = HttpServer(
    HttpConfig(
      port=port,
      on_new_packet=handler,
      on_authenticate=authenticate,
    )
  )
  task = asyncio.create_task(server.start())
  await asyncio.sleep(0.15)
  try:
    async with aiohttp.ClientSession() as session:
      async with session.post(
        f'http://127.0.0.1:{port}/v2/message',
        headers={'Authorization': 'LayrzAuth dev;pw'},
        data=PaPacket(ident='dev', password='pw').to_packet(),
      ) as resp:
        assert resp.status == 401
  finally:
    task.cancel()
    await server.close()
    await asyncio.gather(task, return_exceptions=True)


async def test9_wrong_method_405() -> None:
  port = _free_port()

  async def handler(pkt: object, req: object) -> ServerPacket | None:
    return None

  server = HttpServer(HttpConfig(port=port, on_new_packet=handler))
  task = asyncio.create_task(server.start())
  await asyncio.sleep(0.15)
  try:
    async with aiohttp.ClientSession() as session:
      async with session.get(
        f'http://127.0.0.1:{port}/v2/message',
        headers={'Authorization': 'LayrzAuth dev;pw'},
      ) as resp:
        assert resp.status == 405
  finally:
    task.cancel()
    await server.close()
    await asyncio.gather(task, return_exceptions=True)


async def test10_invalid_body_400() -> None:
  port = _free_port()

  async def handler(pkt: object, req: object) -> ServerPacket | None:
    return None

  server = HttpServer(HttpConfig(port=port, on_new_packet=handler))
  task = asyncio.create_task(server.start())
  await asyncio.sleep(0.15)
  try:
    async with aiohttp.ClientSession() as session:
      async with session.post(
        f'http://127.0.0.1:{port}/v2/message',
        headers={'Authorization': 'LayrzAuth dev;pw'},
        data='garbage',
      ) as resp:
        assert resp.status == 400
  finally:
    task.cancel()
    await server.close()
    await asyncio.gather(task, return_exceptions=True)


async def test11_commands_204_no_handler() -> None:
  port = _free_port()

  async def handler(pkt: object, req: object) -> ServerPacket | None:
    return None

  server = HttpServer(HttpConfig(port=port, on_new_packet=handler))
  task = asyncio.create_task(server.start())
  await asyncio.sleep(0.15)
  try:
    async with aiohttp.ClientSession() as session:
      async with session.get(
        f'http://127.0.0.1:{port}/v2/commands',
        headers={'Authorization': 'LayrzAuth dev;pw'},
      ) as resp:
        assert resp.status == 204
  finally:
    task.cancel()
    await server.close()
    await asyncio.gather(task, return_exceptions=True)


async def test12_commands_200_with_packet() -> None:
  port = _free_port()

  async def pull_commands(ident: str, passwd: str, req: object) -> ServerPacket | None:
    return AsPacket()

  async def handler(pkt: object, req: object) -> ServerPacket | None:
    return None

  server = HttpServer(
    HttpConfig(
      port=port,
      on_new_packet=handler,
      on_pull_commands=pull_commands,
    )
  )
  task = asyncio.create_task(server.start())
  await asyncio.sleep(0.15)
  try:
    async with aiohttp.ClientSession() as session:
      async with session.get(
        f'http://127.0.0.1:{port}/v2/commands',
        headers={'Authorization': 'LayrzAuth dev;pw'},
      ) as resp:
        assert resp.status == 200
        body = await resp.text()
        assert body == AsPacket().to_packet()
  finally:
    task.cancel()
    await server.close()
    await asyncio.gather(task, return_exceptions=True)


async def test13_close_stops_server() -> None:
  port = _free_port()

  async def handler(pkt: object, req: object) -> ServerPacket | None:
    return None

  server = HttpServer(HttpConfig(port=port, on_new_packet=handler))
  task = asyncio.create_task(server.start())
  await asyncio.sleep(0.15)
  await server.close()
  await asyncio.wait_for(task, timeout=2.0)
