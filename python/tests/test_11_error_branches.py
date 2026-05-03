"""Test error-branch paths across all packet types."""

import pytest

from layrz_protocol.packets.client import PaPacket, PbPacket, PcPacket, PdPacket, PiPacket, PmPacket, PrPacket, PsPacket
from layrz_protocol.packets.server import AbPacket, AcPacket, AoPacket, ArPacket, AsPacket, AuPacket
from layrz_protocol.packets.trips.te import TePacket
from layrz_protocol.packets.trips.ts import TsPacket
from layrz_protocol.utils import MalformedException, CrcException, calculate_crc


def _crc(payload: str) -> str:
  return str(hex(calculate_crc(payload.encode())))[2:].upper().zfill(4)


# ---------------------------------------------------------------------------
# PrPacket
# ---------------------------------------------------------------------------

def test11_pr_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    PrPacket.from_packet('<Xx>;0000</Xx>')


def test11_pr_bad_parts() -> None:
  payload = 'extra;extra;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PrPacket.from_packet(f'<Pr>{payload}{crc}</Pr>')


def test11_pr_bad_crc_hex() -> None:
  with pytest.raises(CrcException):
    PrPacket.from_packet('<Pr>ZZZZ</Pr>')


def test11_pr_crc_mismatch() -> None:
  with pytest.raises(CrcException):
    PrPacket.from_packet('<Pr>DEAD</Pr>')


def test11_pr_str() -> None:
  p = PrPacket()
  assert str(p) == p.to_packet()


# ---------------------------------------------------------------------------
# PaPacket
# ---------------------------------------------------------------------------

def test11_pa_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    PaPacket.from_packet('<Xx>a;b;0000</Xx>')


def test11_pa_bad_parts() -> None:
  payload = 'a;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PaPacket.from_packet(f'<Pa>{payload}{crc}</Pa>')


def test11_pa_bad_crc_hex() -> None:
  with pytest.raises(CrcException):
    PaPacket.from_packet('<Pa>a;b;ZZZZ</Pa>')


def test11_pa_crc_mismatch() -> None:
  with pytest.raises(CrcException):
    PaPacket.from_packet('<Pa>a;b;DEAD</Pa>')


# ---------------------------------------------------------------------------
# PcPacket
# ---------------------------------------------------------------------------

def test11_pc_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    PcPacket.from_packet('<Xx>0;1;hi;0000</Xx>')


def test11_pc_bad_parts() -> None:
  payload = '0;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PcPacket.from_packet(f'<Pc>{payload}{crc}</Pc>')


def test11_pc_bad_crc_hex() -> None:
  with pytest.raises(CrcException):
    PcPacket.from_packet('<Pc>0;1;hi;ZZZZ</Pc>')


def test11_pc_crc_mismatch() -> None:
  with pytest.raises(CrcException):
    PcPacket.from_packet('<Pc>0;1;hi;DEAD</Pc>')


def test11_pc_bad_timestamp() -> None:
  payload = 'notts;1;hi;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PcPacket.from_packet(f'<Pc>{payload}{crc}</Pc>')


def test11_pc_bad_command_id() -> None:
  payload = '0;notanint;hi;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PcPacket.from_packet(f'<Pc>{payload}{crc}</Pc>')


def test11_pc_str() -> None:
  payload = '0;1;Hello;'
  crc = _crc(payload)
  p = PcPacket.from_packet(f'<Pc>{payload}{crc}</Pc>')
  assert str(p) == p.to_packet()


# ---------------------------------------------------------------------------
# PsPacket
# ---------------------------------------------------------------------------

def test11_ps_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    PsPacket.from_packet('<Xx>0;;0000</Xx>')


def test11_ps_bad_parts() -> None:
  payload = 'only_one_part;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PsPacket.from_packet(f'<Ps>{payload}{crc}</Ps>')


def test11_ps_bad_crc_hex() -> None:
  with pytest.raises(CrcException):
    PsPacket.from_packet('<Ps>0;;ZZZZ</Ps>')


def test11_ps_crc_mismatch() -> None:
  with pytest.raises(CrcException):
    PsPacket.from_packet('<Ps>0;;DEAD</Ps>')


def test11_ps_bad_timestamp() -> None:
  payload = 'notts;;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PsPacket.from_packet(f'<Ps>{payload}{crc}</Ps>')


def test11_ps_str() -> None:
  payload = '0;;'
  crc = _crc(payload)
  p = PsPacket.from_packet(f'<Ps>{payload}{crc}</Ps>')
  assert str(p) == p.to_packet()


# ---------------------------------------------------------------------------
# PmPacket
# ---------------------------------------------------------------------------

def test11_pm_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    PmPacket.from_packet('<Xx>f;ct;ZA==;0000</Xx>')


def test11_pm_bad_parts() -> None:
  payload = 'file;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PmPacket.from_packet(f'<Pm>{payload}{crc}</Pm>')


def test11_pm_bad_crc_hex() -> None:
  import base64
  data = base64.b64encode(b'hi').decode()
  with pytest.raises(CrcException):
    PmPacket.from_packet(f'<Pm>f;ct;{data};ZZZZ</Pm>')


def test11_pm_crc_mismatch() -> None:
  import base64
  data = base64.b64encode(b'hi').decode()
  with pytest.raises(CrcException):
    PmPacket.from_packet(f'<Pm>f;ct;{data};DEAD</Pm>')


def test11_pm_str() -> None:
  import base64
  data = base64.b64encode(b'hi').decode()
  payload = f'file;image/jpeg;{data};'
  crc = _crc(payload)
  p = PmPacket.from_packet(f'<Pm>{payload}{crc}</Pm>')
  assert str(p) == p.to_packet()


# ---------------------------------------------------------------------------
# PiPacket
# ---------------------------------------------------------------------------

def test11_pi_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    PiPacket.from_packet('<Xx>1;2;3;4;5;6;0;1;0000</Xx>')


def test11_pi_bad_parts() -> None:
  payload = 'a;b;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PiPacket.from_packet(f'<Pi>{payload}{crc}</Pi>')


def test11_pi_bad_crc_hex() -> None:
  with pytest.raises(CrcException):
    PiPacket.from_packet('<Pi>ident;1;1;1;1;1;0;1;ZZZZ</Pi>')


def test11_pi_crc_mismatch() -> None:
  with pytest.raises(CrcException):
    PiPacket.from_packet('<Pi>ident;1;1;1;1;1;0;1;DEAD</Pi>')


def test11_pi_bad_firmware_build() -> None:
  payload = 'ident;1;notint;1;1;1;0;1;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PiPacket.from_packet(f'<Pi>{payload}{crc}</Pi>')


def test11_pi_bad_device_id() -> None:
  payload = 'ident;1;1;notint;1;1;0;1;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PiPacket.from_packet(f'<Pi>{payload}{crc}</Pi>')


def test11_pi_bad_hardware_id() -> None:
  payload = 'ident;1;1;1;notint;1;0;1;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PiPacket.from_packet(f'<Pi>{payload}{crc}</Pi>')


def test11_pi_bad_model_id() -> None:
  payload = 'ident;1;1;1;1;notint;0;1;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PiPacket.from_packet(f'<Pi>{payload}{crc}</Pi>')


def test11_pi_str_firmware_id_fallback() -> None:
  """firmware_id falls back to str when not an int; firmware_branch falls back to STABLE."""
  payload = 'ident;notanint;1;1;1;1;999;1;'
  crc = _crc(payload)
  p = PiPacket.from_packet(f'<Pi>{payload}{crc}</Pi>')
  assert p.firmware_id == 'notanint'
  assert str(p) == p.to_packet()


# ---------------------------------------------------------------------------
# PdPacket
# ---------------------------------------------------------------------------

def test11_pd_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    PdPacket.from_packet('<Xx>0;1.0;2.0;3.0;4.0;5.0;6;7.0;;0000</Xx>')


def test11_pd_bad_parts() -> None:
  payload = '0;1.0;2.0;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PdPacket.from_packet(f'<Pd>{payload}{crc}</Pd>')


def test11_pd_bad_crc_hex() -> None:
  with pytest.raises(CrcException):
    PdPacket.from_packet('<Pd>0;1.0;2.0;3.0;4.0;5.0;6;7.0;;ZZZZ</Pd>')


def test11_pd_crc_mismatch() -> None:
  with pytest.raises(CrcException):
    PdPacket.from_packet('<Pd>0;1.0;2.0;3.0;4.0;5.0;6;7.0;;DEAD</Pd>')


def test11_pd_bad_timestamp() -> None:
  payload = 'notts;1.0;2.0;3.0;4.0;5.0;6;7.0;;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PdPacket.from_packet(f'<Pd>{payload}{crc}</Pd>')


def test11_pd_bad_latitude() -> None:
  payload = '0;notfloat;2.0;3.0;4.0;5.0;6;7.0;;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PdPacket.from_packet(f'<Pd>{payload}{crc}</Pd>')


def test11_pd_bad_longitude() -> None:
  payload = '0;1.0;notfloat;3.0;4.0;5.0;6;7.0;;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PdPacket.from_packet(f'<Pd>{payload}{crc}</Pd>')


def test11_pd_bad_altitude() -> None:
  payload = '0;1.0;2.0;notfloat;4.0;5.0;6;7.0;;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PdPacket.from_packet(f'<Pd>{payload}{crc}</Pd>')


def test11_pd_bad_speed() -> None:
  payload = '0;1.0;2.0;3.0;notfloat;5.0;6;7.0;;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PdPacket.from_packet(f'<Pd>{payload}{crc}</Pd>')


def test11_pd_bad_direction() -> None:
  payload = '0;1.0;2.0;3.0;4.0;notfloat;6;7.0;;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PdPacket.from_packet(f'<Pd>{payload}{crc}</Pd>')


def test11_pd_bad_satellites() -> None:
  payload = '0;1.0;2.0;3.0;4.0;5.0;notint;7.0;;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PdPacket.from_packet(f'<Pd>{payload}{crc}</Pd>')


def test11_pd_bad_hdop() -> None:
  payload = '0;1.0;2.0;3.0;4.0;5.0;6;notfloat;;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PdPacket.from_packet(f'<Pd>{payload}{crc}</Pd>')


def test11_pd_str() -> None:
  payload = '0;1.0;2.0;3.0;4.0;5.0;6;7.0;;'
  crc = _crc(payload)
  p = PdPacket.from_packet(f'<Pd>{payload}{crc}</Pd>')
  assert str(p) == p.to_packet()


# ---------------------------------------------------------------------------
# PbPacket
# ---------------------------------------------------------------------------

def test11_pb_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    PbPacket.from_packet('<Xx>;0000</Xx>')


def test11_pb_bad_parts() -> None:
  payload = 'a;b;c;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    PbPacket.from_packet(f'<Pb>{payload}{crc}</Pb>')


def test11_pb_bad_crc_hex() -> None:
  # Need (len(parts[:-1]) % 12) == 0. With 12 fields + CRC in the body:
  # body = f1;f2;...;f12;CRC → split → 13 parts, parts[:-1] = 12.
  body = 'AABBCCDDEEFF;0;;;3.0;MODEL;;-70;;;+DATA;+SVC;'
  with pytest.raises(CrcException):
    PbPacket.from_packet(f'<Pb>{body}ZZZZ</Pb>')


# ---------------------------------------------------------------------------
# TsPacket
# ---------------------------------------------------------------------------

def test11_ts_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    TsPacket.from_packet('<Xx>0;00000000-0000-0000-0000-000000000000;0000</Xx>')


def test11_ts_bad_parts() -> None:
  payload = '0;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    TsPacket.from_packet(f'<Ts>{payload}{crc}</Ts>')


def test11_ts_bad_crc_hex() -> None:
  with pytest.raises(CrcException):
    TsPacket.from_packet('<Ts>0;00000000-0000-0000-0000-000000000000;ZZZZ</Ts>')


def test11_ts_crc_mismatch() -> None:
  with pytest.raises(CrcException):
    TsPacket.from_packet('<Ts>0;00000000-0000-0000-0000-000000000000;DEAD</Ts>')


def test11_ts_bad_timestamp() -> None:
  payload = 'notts;00000000-0000-0000-0000-000000000000;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    TsPacket.from_packet(f'<Ts>{payload}{crc}</Ts>')


def test11_ts_bad_uuid() -> None:
  payload = '0;not-a-uuid;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    TsPacket.from_packet(f'<Ts>{payload}{crc}</Ts>')


def test11_ts_str() -> None:
  payload = '0;00000000-0000-0000-0000-000000000000;'
  crc = _crc(payload)
  p = TsPacket.from_packet(f'<Ts>{payload}{crc}</Ts>')
  assert str(p) == p.to_packet()


# ---------------------------------------------------------------------------
# TePacket
# ---------------------------------------------------------------------------

def test11_te_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    TePacket.from_packet('<Xx>0;00000000-0000-0000-0000-000000000000;0.0;0.0;0;0000</Xx>')


def test11_te_bad_parts() -> None:
  payload = '0;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    TePacket.from_packet(f'<Te>{payload}{crc}</Te>')


def test11_te_bad_crc_hex() -> None:
  with pytest.raises(CrcException):
    TePacket.from_packet('<Te>0;00000000-0000-0000-0000-000000000000;0.0;0.0;0;ZZZZ</Te>')


def test11_te_crc_mismatch() -> None:
  with pytest.raises(CrcException):
    TePacket.from_packet('<Te>0;00000000-0000-0000-0000-000000000000;0.0;0.0;0;DEAD</Te>')


def test11_te_bad_timestamp() -> None:
  payload = 'notts;00000000-0000-0000-0000-000000000000;0.0;0.0;0;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    TePacket.from_packet(f'<Te>{payload}{crc}</Te>')


def test11_te_bad_uuid() -> None:
  payload = '0;not-a-uuid;0.0;0.0;0;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    TePacket.from_packet(f'<Te>{payload}{crc}</Te>')


def test11_te_bad_distance() -> None:
  payload = '0;00000000-0000-0000-0000-000000000000;notfloat;0.0;0;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    TePacket.from_packet(f'<Te>{payload}{crc}</Te>')


def test11_te_bad_max_speed() -> None:
  payload = '0;00000000-0000-0000-0000-000000000000;0.0;notfloat;0;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    TePacket.from_packet(f'<Te>{payload}{crc}</Te>')


def test11_te_bad_duration() -> None:
  payload = '0;00000000-0000-0000-0000-000000000000;0.0;0.0;notint;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    TePacket.from_packet(f'<Te>{payload}{crc}</Te>')


def test11_te_str() -> None:
  payload = '0;00000000-0000-0000-0000-000000000000;0.000;0.000;0;'
  crc = _crc(payload)
  p = TePacket.from_packet(f'<Te>{payload}{crc}</Te>')
  assert str(p) == p.to_packet()


# ---------------------------------------------------------------------------
# AbPacket
# ---------------------------------------------------------------------------

def test11_ab_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    AbPacket.from_packet('<Xx>;0000</Xx>')


def test11_ab_bad_crc_hex() -> None:
  with pytest.raises(CrcException):
    AbPacket.from_packet('<Ab>;ZZZZ</Ab>')


def test11_ab_crc_mismatch() -> None:
  with pytest.raises(CrcException):
    AbPacket.from_packet('<Ab>;DEAD</Ab>')


def test11_ab_odd_parts() -> None:
  payload = 'a;b;c;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    AbPacket.from_packet(f'<Ab>{payload}{crc}</Ab>')


def test11_ab_str() -> None:
  p = AbPacket(devices=[])
  assert str(p) == p.to_packet()


# ---------------------------------------------------------------------------
# AcPacket
# ---------------------------------------------------------------------------

def test11_ac_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    AcPacket.from_packet('<Xx>;0000</Xx>')


def test11_ac_bad_crc_hex() -> None:
  with pytest.raises(CrcException):
    AcPacket.from_packet('<Ac>;ZZZZ</Ac>')


def test11_ac_crc_mismatch() -> None:
  with pytest.raises(CrcException):
    AcPacket.from_packet('<Ac>;DEAD</Ac>')


def test11_ac_bad_parts() -> None:
  payload = 'a;b;c;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    AcPacket.from_packet(f'<Ac>{payload}{crc}</Ac>')


def test11_ac_str() -> None:
  p = AcPacket(commands=[])
  assert str(p) == p.to_packet()


# ---------------------------------------------------------------------------
# AoPacket
# ---------------------------------------------------------------------------

def test11_ao_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    AoPacket.from_packet('<Xx>0;0000</Xx>')


def test11_ao_bad_parts() -> None:
  payload = 'a;b;c;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    AoPacket.from_packet(f'<Ao>{payload}{crc}</Ao>')


def test11_ao_bad_crc_hex() -> None:
  with pytest.raises(CrcException):
    AoPacket.from_packet('<Ao>0;ZZZZ</Ao>')


def test11_ao_crc_mismatch() -> None:
  with pytest.raises(CrcException):
    AoPacket.from_packet('<Ao>0;DEAD</Ao>')


def test11_ao_bad_timestamp() -> None:
  payload = 'notts;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    AoPacket.from_packet(f'<Ao>{payload}{crc}</Ao>')


def test11_ao_str() -> None:
  payload = '0;'
  crc = _crc(payload)
  p = AoPacket.from_packet(f'<Ao>{payload}{crc}</Ao>')
  assert str(p) == p.to_packet()


# ---------------------------------------------------------------------------
# ArPacket
# ---------------------------------------------------------------------------

def test11_ar_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    ArPacket.from_packet('<Xx>reason;0000</Xx>')


def test11_ar_bad_parts() -> None:
  payload = 'a;b;'
  crc = _crc(payload)
  with pytest.raises(MalformedException):
    ArPacket.from_packet(f'<Ar>{payload}{crc}</Ar>')


def test11_ar_bad_crc_hex() -> None:
  with pytest.raises(CrcException):
    ArPacket.from_packet('<Ar>reason;ZZZZ</Ar>')


def test11_ar_crc_mismatch() -> None:
  with pytest.raises(CrcException):
    ArPacket.from_packet('<Ar>reason;DEAD</Ar>')


def test11_ar_str() -> None:
  payload = 'reason;'
  crc = _crc(payload)
  p = ArPacket.from_packet(f'<Ar>{payload}{crc}</Ar>')
  assert str(p) == p.to_packet()


# ---------------------------------------------------------------------------
# AsPacket
# ---------------------------------------------------------------------------

def test11_as_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    AsPacket.from_packet('<Xx>;0000</Xx>')


def test11_as_bad_crc_hex() -> None:
  with pytest.raises(CrcException):
    AsPacket.from_packet('<As>;ZZZZ</As>')


def test11_as_crc_mismatch() -> None:
  with pytest.raises(CrcException):
    AsPacket.from_packet('<As>;DEAD</As>')


def test11_as_str() -> None:
  p = AsPacket()
  assert str(p) == p.to_packet()


# ---------------------------------------------------------------------------
# AuPacket
# ---------------------------------------------------------------------------

def test11_au_wrong_tag() -> None:
  with pytest.raises(MalformedException):
    AuPacket.from_packet('<Xx>;0000</Xx>')


def test11_au_bad_crc_hex() -> None:
  with pytest.raises(CrcException):
    AuPacket.from_packet('<Au>;ZZZZ</Au>')


def test11_au_crc_mismatch() -> None:
  with pytest.raises(CrcException):
    AuPacket.from_packet('<Au>;DEAD</Au>')
