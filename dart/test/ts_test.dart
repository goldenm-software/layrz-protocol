import 'package:layrz_protocol/utils/errors.dart';
import 'package:test/test.dart';
import 'package:layrz_protocol/packets/packets.dart';
import 'package:layrz_protocol/utils/crc.dart';

void main() {
  test('Packet.fromPacket() routes TsPacket', () {
    final ts = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    final original = TsPacket(timestamp: ts, tripId: 'TRIP-001');
    final parsed = Packet.fromPacket(original.toPacket());
    expect(parsed, isA<TsPacket>());
  });

  test('TsPacket.toPacket()', () {
    final timestamp = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    final packet = TsPacket(timestamp: timestamp, tripId: 'TRIP-001');
    final raw = packet.toPacket();

    expect(raw.startsWith('<Ts>'), true);
    expect(raw.endsWith('</Ts>'), true);

    String payload = '1700000000;TRIP-001;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(raw, '<Ts>$payload$crc</Ts>');
  });

  test('TsPacket.fromPacket()', () {
    String payload = '1700000000;TRIP-001;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    String raw = '<Ts>$payload$crc</Ts>';

    final packet = TsPacket.fromPacket(raw);
    expect(packet.timestamp, DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true));
    expect(packet.tripId, 'TRIP-001');
    expect(packet.toPacket(), raw);
  });

  test('TsPacket.fromPacket() invalid format throws ParseException', () {
    expect(() => TsPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
  });

  test('TsPacket.fromPacket() wrong part count throws MalformedException', () {
    String payload = ';';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(() => TsPacket.fromPacket('<Ts>$payload$crc</Ts>'), throwsA(isA<MalformedException>()));
  });

  test('TsPacket.fromPacket() invalid CRC throws CrcException', () {
    expect(() => TsPacket.fromPacket('<Ts>1700000000;TRIP-001;FFFF</Ts>'), throwsA(isA<CrcException>()));
  });

  test('TsPacket.fromPacket() invalid timestamp throws MalformedException', () {
    String payload = 'NOTANUMBER;TRIP-001;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(() => TsPacket.fromPacket('<Ts>$payload$crc</Ts>'), throwsA(isA<MalformedException>()));
  });

  test('TsPacket.copyWith()', () {
    final ts = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    final original = TsPacket(timestamp: ts, tripId: 'A');
    final copy = original.copyWith(tripId: 'B');
    expect(copy.tripId, 'B');
    expect(copy.timestamp, ts);
  });
}
