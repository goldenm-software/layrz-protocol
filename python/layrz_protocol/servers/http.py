from __future__ import annotations

import asyncio
import logging
import sys
from collections.abc import Awaitable, Callable
from dataclasses import dataclass, field

if sys.version_info >= (3, 11):
  from typing import Self
else:
  from typing_extensions import Self

from aiohttp import web

from layrz_protocol.packets.client import ClientPacket, decode_client_packet
from layrz_protocol.packets.server import ServerPacket

_log = logging.getLogger(__name__)

OnNewPacketHttp = Callable[[ClientPacket, web.Request], Awaitable[ServerPacket | None]]
OnPullCommands = Callable[[str, str, web.Request], Awaitable[ServerPacket | None]]
OnAuthenticate = Callable[[str, str, web.Request], Awaitable[bool]]
OnDecodeErrorHttp = Callable[[Exception, bytes, web.Request], Awaitable[None]]


def _parse_layrz_auth(header: str | None) -> tuple[str, str] | None:
  if not header or not header.startswith('LayrzAuth '):
    return None
  rest = header[len('LayrzAuth '):]
  sep = rest.find(';')
  if sep < 0:
    return None
  return rest[:sep], rest[sep + 1:]


@dataclass
class HttpConfig:
  """Configuration for the Layrz Protocol HTTP server."""

  port: int
  on_new_packet: OnNewPacketHttp
  on_pull_commands: OnPullCommands | None = field(default=None)
  on_authenticate: OnAuthenticate | None = field(default=None)
  on_decode_error: OnDecodeErrorHttp | None = field(default=None)
  max_body_bytes: int = field(default=1 << 20)
  shutdown_timeout: float = field(default=5.0)


class HttpServer:
  """Aiohttp-based HTTP server for the Layrz Protocol."""

  def __init__(self, cfg: HttpConfig) -> None:
    if cfg.on_new_packet is None:
      raise ValueError('on_new_packet handler is not set')
    if cfg.port <= 0 or cfg.port >= 65535:
      raise ValueError('port is not valid')
    self._cfg = cfg
    self._runner: web.AppRunner | None = None
    self._site: web.TCPSite | None = None
    self._stopped: asyncio.Event | None = None

  async def start(self) -> None:
    """Start the HTTP server. Blocks until close() is called."""
    app = web.Application(client_max_size=self._cfg.max_body_bytes)
    app.router.add_post('/v2/message', self._handle_message)
    app.router.add_get('/v2/commands', self._handle_commands)

    self._runner = web.AppRunner(app)
    await self._runner.setup()
    self._site = web.TCPSite(
      self._runner,
      host='0.0.0.0',
      port=self._cfg.port,
      shutdown_timeout=self._cfg.shutdown_timeout,
    )
    await self._site.start()

    self._stopped = asyncio.Event()
    await self._stopped.wait()

  async def close(self) -> None:
    """Gracefully stop the HTTP server."""
    if self._site is not None:
      await self._site.stop()
    if self._runner is not None:
      await self._runner.cleanup()
    if self._stopped is not None:
      self._stopped.set()

  async def __aenter__(self) -> Self:
    return self

  async def __aexit__(self, *exc: object) -> None:
    await self.close()

  async def _handle_message(self, req: web.Request) -> web.Response:
    creds = _parse_layrz_auth(req.headers.get('Authorization'))
    if creds is None:
      return web.Response(status=401, text='unauthorized')

    ident, passwd = creds

    if self._cfg.on_authenticate is not None:
      if not await self._cfg.on_authenticate(ident, passwd, req):
        return web.Response(status=401, text='unauthorized')

    try:
      body = await req.read()
    except web.HTTPRequestEntityTooLarge:
      return web.Response(status=413, text='request entity too large')

    try:
      packet = decode_client_packet(body.decode('utf-8').strip())
    except Exception as exc:
      if self._cfg.on_decode_error is not None:
        await self._cfg.on_decode_error(exc, body, req)
      else:
        _log.warning('Decode error: %s', exc)
      return web.Response(status=400, text='invalid packet')

    try:
      response = await self._cfg.on_new_packet(packet, req)
    except Exception:
      _log.exception('Handler error in on_new_packet')
      return web.Response(status=500, text='internal server error')

    if response is None:
      return web.Response(status=204)

    return web.Response(status=200, content_type='text/plain', charset='utf-8', text=response.to_packet())

  async def _handle_commands(self, req: web.Request) -> web.Response:
    creds = _parse_layrz_auth(req.headers.get('Authorization'))
    if creds is None:
      return web.Response(status=401, text='unauthorized')

    ident, passwd = creds

    if self._cfg.on_authenticate is not None:
      if not await self._cfg.on_authenticate(ident, passwd, req):
        return web.Response(status=401, text='unauthorized')

    if self._cfg.on_pull_commands is None:
      return web.Response(status=204)

    try:
      response = await self._cfg.on_pull_commands(ident, passwd, req)
    except Exception:
      _log.exception('Handler error in on_pull_commands')
      return web.Response(status=500, text='internal server error')

    if response is None:
      return web.Response(status=204)

    return web.Response(status=200, content_type='text/plain', charset='utf-8', text=response.to_packet())
