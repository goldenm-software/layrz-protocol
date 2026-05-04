from __future__ import annotations

import asyncio
import logging
import re
import sys
from collections.abc import Awaitable, Callable
from dataclasses import dataclass, field

if sys.version_info >= (3, 11):
  from typing import Self
else:
  from typing_extensions import Self

from proxyprotocol import ProxyProtocol
from proxyprotocol.noop import ProxyProtocolNoop
from proxyprotocol.reader import ProxyProtocolReader
from proxyprotocol.sock import SocketInfo
from proxyprotocol.v2 import ProxyProtocolV2

from layrz_protocol.packets.client import ClientPacket, decode_client_packet
from layrz_protocol.packets.server import ServerPacket

_log = logging.getLogger(__name__)

_PACKET_TAG: re.Pattern[str] = re.compile(r'<P[A-Za-z]>')

PeerInfo = tuple[str, int]

OnNewPacketTcp = Callable[[ClientPacket, PeerInfo], Awaitable[ServerPacket | None]]
OnDecodeErrorTcp = Callable[[Exception, bytes, PeerInfo], Awaitable[None]]


def _split(data: str) -> list[str]:
  data = data.rstrip('\n\r')
  locs = [m.start() for m in _PACKET_TAG.finditer(data)]
  if not locs:
    s = data.strip()
    return [s] if s else []
  out: list[str] = []
  for i, start in enumerate(locs):
    end = locs[i + 1] if i + 1 < len(locs) else len(data)
    s = data[start:end].strip()
    if s:
      out.append(s)
  return out


@dataclass
class TcpConfig:
  """Configuration for the Layrz Protocol TCP server."""

  port: int
  on_new_packet: OnNewPacketTcp
  on_decode_error: OnDecodeErrorTcp | None = field(default=None)
  proxy_protocol_v2: bool = field(default=False)
  read_chunk: int = field(default=1024)


class TcpServer:
  """Asyncio-based TCP server for the Layrz Protocol."""

  def __init__(self, cfg: TcpConfig) -> None:
    if cfg.on_new_packet is None:
      raise ValueError('on_new_packet handler is not set')
    if cfg.port <= 0 or cfg.port >= 65535:
      raise ValueError('port is not valid')
    self._cfg = cfg
    self._server: asyncio.Server | None = None
    self._tasks: set[asyncio.Task[None]] = set()

  async def start(self) -> None:
    """Start the server. Blocks until close() is called."""
    protocol: ProxyProtocol
    if self._cfg.proxy_protocol_v2:
      protocol = ProxyProtocolV2()
      _log.info('Proxy Protocol v2 is enabled on the TCP server')
    else:
      protocol = ProxyProtocolNoop()
      _log.info('Proxy Protocol is disabled on the TCP server')

    callback = ProxyProtocolReader(pp=protocol).get_callback(self._handle)
    self._server = await asyncio.start_server(
      client_connected_cb=callback,
      host='0.0.0.0',
      port=self._cfg.port,
    )
    async with self._server:
      await self._server.serve_forever()

  async def close(self) -> None:
    """Stop the server and cancel all live connection tasks."""
    if self._server is not None:
      self._server.close()
      await self._server.wait_closed()
    tasks = list(self._tasks)
    for t in tasks:
      t.cancel()
    if tasks:
      await asyncio.gather(*tasks, return_exceptions=True)

  async def __aenter__(self) -> Self:
    return self

  async def __aexit__(self, *exc: object) -> None:
    await self.close()

  async def _handle(
    self,
    reader: asyncio.StreamReader,
    writer: asyncio.StreamWriter,
    sock_info: SocketInfo,
  ) -> None:
    task = asyncio.current_task()
    if task is not None:
      self._tasks.add(task)
      task.add_done_callback(self._tasks.discard)

    fallback: PeerInfo = writer.get_extra_info('peername') or ('', 0)
    peer: PeerInfo = (str(sock_info.peername_ip), sock_info.peername_port or 0) if sock_info.peername_ip else fallback

    buf = bytearray()
    try:
      while not reader.at_eof():
        chunk = await reader.read(self._cfg.read_chunk)
        if not chunk:
          break
        buf.extend(chunk)
        if b'\n' not in buf:
          continue
        text = buf.decode('utf-8', errors='replace')
        buf.clear()
        for msg in _split(text):
          await self._dispatch(msg, peer, writer)
    except asyncio.CancelledError:
      raise
    except Exception:
      _log.exception('Connection error from %s:%d', *peer)
    finally:
      writer.close()
      try:
        await writer.wait_closed()
      except Exception:
        pass

  async def _dispatch(self, msg: str, peer: PeerInfo, writer: asyncio.StreamWriter) -> None:
    try:
      packet = decode_client_packet(msg)
    except Exception as exc:
      if self._cfg.on_decode_error is not None:
        await self._cfg.on_decode_error(exc, msg.encode(), peer)
      else:
        _log.warning('Decode error from %s:%d: %s', *peer, exc)
      return
    try:
      response = await self._cfg.on_new_packet(packet, peer)
    except Exception:
      _log.exception('Handler error from %s:%d', *peer)
      return
    if response is not None:
      writer.write(response.to_packet().encode())
      await writer.drain()
