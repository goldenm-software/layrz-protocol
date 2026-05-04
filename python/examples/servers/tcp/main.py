from __future__ import annotations

import asyncio
import logging
import signal
from datetime import UTC, datetime

from layrz_protocol.packets.client import ClientPacket, PaPacket
from layrz_protocol.packets.server import AoPacket, AsPacket, ServerPacket
from layrz_protocol.servers import TcpConfig, TcpServer

logging.basicConfig(level=logging.INFO)


async def on_packet(packet: ClientPacket, peer: tuple[str, int]) -> ServerPacket | None:
  if isinstance(packet, PaPacket):
    logging.info('Pa received from %s:%d', *peer)
    return AsPacket()
  logging.info('Packet %s received from %s:%d', type(packet).__name__, *peer)
  return AoPacket(timestamp=datetime.now(UTC))


async def main() -> None:
  server = TcpServer(TcpConfig(port=12345, on_new_packet=on_packet))
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
