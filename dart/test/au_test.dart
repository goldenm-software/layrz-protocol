import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('Packet.fromPacket() routes AuPacket', () {
    // ignore: deprecated_member_use
    final original = AuPacket();
    final parsed = Packet.fromPacket(original.toPacket());
    // ignore: deprecated_member_use
    expect(parsed, isA<AuPacket>());
  });

  test('AuPacket.toPacket()', () {
    // ignore: deprecated_member_use
    final packet = AuPacket();
    final raw = packet.toPacket();

    expect(raw.startsWith('<Au>'), true);
    expect(raw.endsWith('</Au>'), true);

    String payload = ';';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(raw, '<Au>$payload$crc</Au>');
  });

  test('AuPacket.fromPacket()', () {
    String payload = ';';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    String raw = '<Au>$payload$crc</Au>';

    // ignore: deprecated_member_use
    final packet = AuPacket.fromPacket(raw);
    expect(packet.toPacket(), raw);
  });

  test('AuPacket.fromPacket() invalid format throws ParseException', () {
    // ignore: deprecated_member_use
    expect(() => AuPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
  });

  test('AuPacket.copyWith()', () {
    // ignore: deprecated_member_use
    final original = AuPacket();
    final copy = original.copyWith();
    expect(copy.toPacket(), original.toPacket());
  });
}
