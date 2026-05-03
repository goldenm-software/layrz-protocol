import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('Packet.fromPacket() routes TePacket', () {
    final ts = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    final original = TePacket(
      timestamp: ts,
      tripId: 'TRIP-001',
      distanceTraveled: 100.0,
      maxSpeed: 60.0,
      duration: const Duration(seconds: 600),
    );
    final parsed = Packet.fromPacket(original.toPacket());
    expect(parsed, isA<TePacket>());
  });

  test('TePacket.toPacket()', () {
    final timestamp = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    final packet = TePacket(
      timestamp: timestamp,
      tripId: 'TRIP-001',
      distanceTraveled: 12345.678,
      maxSpeed: 120.5,
      duration: const Duration(seconds: 3600),
    );
    final raw = packet.toPacket();

    expect(raw.startsWith('<Te>'), true);
    expect(raw.endsWith('</Te>'), true);

    String payload = '1700000000;TRIP-001;12345.678;120.500;3600;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(raw, '<Te>$payload$crc</Te>');
  });

  test('TePacket.fromPacket()', () {
    String payload = '1700000000;TRIP-001;12345.678;120.500;3600;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    String raw = '<Te>$payload$crc</Te>';

    final packet = TePacket.fromPacket(raw);
    expect(packet.timestamp, DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true));
    expect(packet.tripId, 'TRIP-001');
    expect(packet.distanceTraveled, 12345.678);
    expect(packet.maxSpeed, 120.5);
    expect(packet.duration, const Duration(seconds: 3600));
    expect(packet.toPacket(), raw);
  });

  test('TePacket.fromPacket() invalid format throws ParseException', () {
    expect(() => TePacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
  });

  test('TePacket.fromPacket() wrong part count throws MalformedException', () {
    String payload = 'only;two;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(() => TePacket.fromPacket('<Te>$payload$crc</Te>'), throwsA(isA<MalformedException>()));
  });

  test('TePacket.fromPacket() invalid CRC throws CrcException', () {
    expect(() => TePacket.fromPacket('<Te>0;trip;100.0;50.0;600;FFFF</Te>'), throwsA(isA<CrcException>()));
  });

  test('TePacket.fromPacket() invalid timestamp throws MalformedException', () {
    String payload = 'NOTANUMBER;TRIP-001;100.000;50.000;600;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(() => TePacket.fromPacket('<Te>$payload$crc</Te>'), throwsA(isA<MalformedException>()));
  });

  test('TePacket.fromPacket() invalid distance throws MalformedException', () {
    String payload = '1700000000;TRIP-001;NOTADOUBLE;50.000;600;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(() => TePacket.fromPacket('<Te>$payload$crc</Te>'), throwsA(isA<MalformedException>()));
  });

  test('TePacket.fromPacket() invalid maxSpeed throws MalformedException', () {
    String payload = '1700000000;TRIP-001;100.000;NOTADOUBLE;600;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(() => TePacket.fromPacket('<Te>$payload$crc</Te>'), throwsA(isA<MalformedException>()));
  });

  test('TePacket.fromPacket() invalid duration throws MalformedException', () {
    String payload = '1700000000;TRIP-001;100.000;50.000;NOTANINT;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(() => TePacket.fromPacket('<Te>$payload$crc</Te>'), throwsA(isA<MalformedException>()));
  });

  test('TePacket.copyWith()', () {
    final ts = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    final original = TePacket(
      timestamp: ts,
      tripId: 'A',
      distanceTraveled: 100.0,
      maxSpeed: 50.0,
      duration: const Duration(seconds: 600),
    );
    final copy = original.copyWith(tripId: 'B', maxSpeed: 80.0);
    expect(copy.tripId, 'B');
    expect(copy.maxSpeed, 80.0);
    expect(copy.distanceTraveled, 100.0);
    expect(copy.timestamp, ts);
  });
}
