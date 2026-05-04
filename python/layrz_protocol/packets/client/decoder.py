from __future__ import annotations

from typing import cast

from layrz_protocol.utils.exceptions import MalformedException

from .base import ClientPacket
from .pa import PaPacket
from .pb import PbPacket
from .pc import PcPacket
from .pd import PdPacket
from .pi import PiPacket
from .pm import PmPacket
from .pr import PrPacket
from .ps import PsPacket

_DISPATCH: tuple[tuple[str, type[ClientPacket]], ...] = (
  ('<Pa>', PaPacket),
  ('<Pb>', PbPacket),
  ('<Pc>', PcPacket),
  ('<Pd>', PdPacket),
  ('<Pi>', PiPacket),
  ('<Pm>', PmPacket),
  ('<Pr>', PrPacket),
  ('<Ps>', PsPacket),
)


def decode_client_packet(raw: str) -> ClientPacket:
  """Decode a raw client packet string into the appropriate ClientPacket subclass."""
  raw = raw.strip()
  for prefix, cls in _DISPATCH:
    closing = f'</{prefix[1:]}'
    if raw.startswith(prefix) and raw.endswith(closing):
      return cast(ClientPacket, cls.from_packet(raw))
  raise MalformedException(f'Invalid client packet: {raw!r}')
