import 'package:layrz_protocol/utils/errors.dart';
import 'package:test/test.dart';
import 'package:layrz_protocol/packets/packets.dart';
import 'package:layrz_protocol/utils/crc.dart';

void main() {
  group('PcPacket error paths', () {
    test('fromPacket invalid format throws ParseException', () {
      expect(() => PcPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
    });

    test('fromPacket wrong part count throws MalformedException', () {
      String payload = 'only;one;';
      String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
      expect(() => PcPacket.fromPacket('<Pc>$payload$crc</Pc>'), throwsA(isA<MalformedException>()));
    });

    test('fromPacket invalid CRC throws CrcException', () {
      expect(() => PcPacket.fromPacket('<Pc>0;1;msg;FFFF</Pc>'), throwsA(isA<CrcException>()));
    });

    test('copyWith works', () {
      final ts = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
      final original = PcPacket(timestamp: ts, commandId: 1, message: 'orig');
      final copy = original.copyWith(message: 'new');
      expect(copy.message, 'new');
      expect(copy.commandId, 1);
      expect(copy.timestamp, ts);
    });
  });

  group('PiPacket error paths', () {
    test('fromPacket invalid format throws ParseException', () {
      expect(() => PiPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
    });

    test('fromPacket wrong part count throws MalformedException', () {
      String payload = 'only;two;';
      String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
      expect(() => PiPacket.fromPacket('<Pi>$payload$crc</Pi>'), throwsA(isA<MalformedException>()));
    });

    test('fromPacket invalid CRC throws CrcException', () {
      expect(() => PiPacket.fromPacket('<Pi>a;b;c;d;e;f;g;h;FFFF</Pi>'), throwsA(isA<CrcException>()));
    });

    test('copyWith works', () {
      final original = PiPacket(
        ident: 'IMEI',
        firmwareId: 'fw1',
        firmwareBuild: 1,
        deviceId: 'dev1',
        hardwareId: 'hw1',
        modelId: 'model1',
        firmwareBranch: FirmwareBranch.stable,
        fotaEnabled: false,
      );
      final copy = original.copyWith(fotaEnabled: true, firmwareBuild: 99);
      expect(copy.fotaEnabled, true);
      expect(copy.firmwareBuild, 99);
      expect(copy.ident, 'IMEI');
    });

    test('FirmwareBranch.development.toPacket() returns 1', () {
      expect(FirmwareBranch.development.toPacket(), '1');
    });

    test('FirmwareBranch.stable.toPacket() returns 0', () {
      expect(FirmwareBranch.stable.toPacket(), '0');
    });
  });

  group('PsPacket error paths', () {
    test('fromPacket invalid format throws ParseException', () {
      expect(() => PsPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
    });

    test('fromPacket wrong part count throws MalformedException', () {
      String payload = 'only;';
      String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
      expect(() => PsPacket.fromPacket('<Ps>$payload$crc</Ps>'), throwsA(isA<MalformedException>()));
    });

    test('fromPacket invalid CRC throws CrcException', () {
      expect(() => PsPacket.fromPacket('<Ps>0;extra:val;FFFF</Ps>'), throwsA(isA<CrcException>()));
    });

    test('copyWith works', () {
      final ts = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
      final original = PsPacket(timestamp: ts, params: {'key': 'val'});
      final newTs = DateTime.fromMillisecondsSinceEpoch(1700000001 * 1000, isUtc: true);
      final copy = original.copyWith(timestamp: newTs);
      expect(copy.timestamp, newTs);
      expect(copy.params, {'key': 'val'});
    });
  });

  group('ArPacket error paths', () {
    test('fromPacket invalid format throws ParseException', () {
      expect(() => ArPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
    });

    test('fromPacket invalid CRC throws CrcException', () {
      expect(() => ArPacket.fromPacket('<Ar>reason;FFFF</Ar>'), throwsA(isA<CrcException>()));
    });

    test('copyWith works', () {
      final original = ArPacket(reason: 'old reason');
      final copy = original.copyWith(reason: 'new reason');
      expect(copy.reason, 'new reason');
    });
  });

  group('AoPacket error paths', () {
    test('fromPacket invalid format throws ParseException', () {
      expect(() => AoPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
    });

    test('fromPacket wrong part count throws MalformedException', () {
      String payload = 'one;two;';
      String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
      expect(() => AoPacket.fromPacket('<Ao>$payload$crc</Ao>'), throwsA(isA<MalformedException>()));
    });

    test('fromPacket invalid CRC throws CrcException', () {
      expect(() => AoPacket.fromPacket('<Ao>0;FFFF</Ao>'), throwsA(isA<CrcException>()));
    });

    test('copyWith works', () {
      final ts = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
      final original = AoPacket(timestamp: ts);
      final newTs = DateTime.fromMillisecondsSinceEpoch(1700000001 * 1000, isUtc: true);
      final copy = original.copyWith(timestamp: newTs);
      expect(copy.timestamp, newTs);
    });
  });

  group('ImPacket error paths', () {
    test('fromPacket invalid format throws ParseException', () {
      expect(() => ImPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
    });

    test('fromPacket wrong part count throws MalformedException', () {
      String payload = 'one;two;';
      String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
      expect(() => ImPacket.fromPacket('<Im>$payload$crc</Im>'), throwsA(isA<MalformedException>()));
    });

    test('fromPacket invalid CRC throws CrcException', () {
      expect(() => ImPacket.fromPacket('<Im>0;chat;msg;FFFF</Im>'), throwsA(isA<CrcException>()));
    });

    test('copyWith works', () {
      final ts = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
      final original = ImPacket(timestamp: ts, chatId: 'c1', message: 'm1');
      final copy = original.copyWith(message: 'new');
      expect(copy.message, 'new');
      expect(copy.chatId, 'c1');
    });
  });

  group('AbPacket error paths', () {
    test('fromPacket invalid format throws ParseException', () {
      expect(() => AbPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
    });

    test('fromPacket invalid CRC throws CrcException', () {
      expect(() => AbPacket.fromPacket('<Ab>1234567890AB:GENERIC;FFFF</Ab>'), throwsA(isA<CrcException>()));
    });

    test('copyWith works', () {
      final original = AbPacket(devices: []);
      final copy = original.copyWith(devices: []);
      expect(copy.devices, isEmpty);
    });
  });

  group('AcPacket error paths', () {
    test('fromPacket invalid format throws ParseException', () {
      expect(() => AcPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
    });

    test('fromPacket invalid CRC throws CrcException', () {
      expect(() => AcPacket.fromPacket('<Ac>payload;FFFF</Ac>'), throwsA(isA<CrcException>()));
    });

    test('copyWith works', () {
      final original = AcPacket(
        commands: [Command(commandId: '1', commandName: 'ping', args: {})],
      );
      final copy = original.copyWith(commands: []);
      expect(copy.commands, isEmpty);
    });
  });

  group('PdPacket error paths', () {
    test('fromPacket invalid format throws ParseException', () {
      expect(() => PdPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
    });

    test('fromPacket wrong part count throws MalformedException', () {
      String payload = 'one;two;';
      String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
      expect(() => PdPacket.fromPacket('<Pd>$payload$crc</Pd>'), throwsA(isA<MalformedException>()));
    });

    test('copyWith works', () {
      final ts = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
      final original = PdPacket(timestamp: ts, position: Position(), extra: {'k': 1});
      final copy = original.copyWith(extra: {'k': 2});
      expect(copy.extra['k'], 2);
      expect(copy.timestamp, ts);
    });
  });

  group('PbPacket error paths', () {
    test('fromPacket invalid format throws ParseException', () {
      expect(() => PbPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
    });

    test('fromPacket invalid CRC throws CrcException', () {
      expect(() => PbPacket.fromPacket('<Pb>somedata;FFFF</Pb>'), throwsA(isA<CrcException>()));
    });

    test('copyWith works', () {
      final adv = BleAdvertisement(
        deviceName: 'Dev',
        macAddress: '12:34:56:78:90:AB',
        timestamp: DateTime.fromMillisecondsSinceEpoch(1738276597 * 1000, isUtc: true),
        rssi: -50,
        model: 'ELA_PUCK_RHT',
        txPower: -1,
        manufacturerData: [],
        serviceData: [],
      );
      final original = PbPacket(advertisements: [adv]);
      final copy = original.copyWith(advertisements: []);
      expect(copy.advertisements, isEmpty);
    });
  });
}
