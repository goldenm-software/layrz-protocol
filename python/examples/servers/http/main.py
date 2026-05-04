from __future__ import annotations

import asyncio
import logging
import signal
from datetime import UTC, datetime

from aiohttp import web

from layrz_protocol.packets.client import ClientPacket, PaPacket
from layrz_protocol.packets.server import AoPacket, AsPacket, ServerPacket
from layrz_protocol.servers import HttpConfig, HttpServer

logging.basicConfig(level=logging.INFO)


async def on_authenticate(ident: str, passwd: str, req: web.Request) -> bool:
  return ident == 'device001' and passwd == 'secret'


async def on_packet(packet: ClientPacket, req: web.Request) -> ServerPacket | None:
  if isinstance(packet, PaPacket):
    logging.info('Pa received')
    return AsPacket()
  logging.info('Packet %s received', type(packet).__name__)
  return AoPacket(timestamp=datetime.now(UTC))


async def on_pull_commands(ident: str, passwd: str, req: web.Request) -> ServerPacket | None:
  logging.info('Commands requested by %s', ident)
  return None


async def main() -> None:
  server = HttpServer(HttpConfig(
    port=12345,
    on_authenticate=on_authenticate,
    on_new_packet=on_packet,
    on_pull_commands=on_pull_commands,
  ))
  loop = asyncio.get_running_loop()

  def _stop() -> None:
    logging.info('Shutting down...')
    loop.create_task(server.close())

  loop.add_signal_handler(signal.SIGINT, _stop)
  loop.add_signal_handler(signal.SIGTERM, _stop)

  logging.info('Listening on :12345')
  try:
    await server.start()
  except asyncio.CancelledError:
    pass
  finally:
    await server.close()


if __name__ == '__main__':
  asyncio.run(main())
