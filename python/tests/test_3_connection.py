"""Test Client Packets"""

from datetime import datetime

import pytest

from layrz_protocol import LayrzProtocol
from layrz_protocol.constants import UTC
from layrz_protocol.packets import PdPacket, Position


@pytest.mark.skip(reason='integration test requires a live Layrz server')
def test3_base() -> None:
  client = LayrzProtocol(ident='link_test', base_url='http://localhost')
  msg = PdPacket(
    timestamp=datetime.now(UTC),
    position=Position(latitude=10.0, longitude=10.0, speed=10.0),
    extra={'is.unit.test.python': True},
  )
  client.send_data(msg)
