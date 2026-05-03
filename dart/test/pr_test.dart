import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('Packet.fromPacket() routes PrPacket', () {
    final original = PrPacket();
    final parsed = Packet.fromPacket(original.toPacket());
    expect(parsed, isA<PrPacket>());
  });

  test('PrPacket.toPacket()', () {
    final packet = PrPacket();
    final raw = packet.toPacket();

    expect(raw.startsWith('<Pr>'), true);
    expect(raw.endsWith('</Pr>'), true);

    String payload = ';';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(raw, '<Pr>$payload$crc</Pr>');
  });

  test('PrPacket.fromPacket()', () {
    String payload = ';';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    String raw = '<Pr>$payload$crc</Pr>';

    final packet = PrPacket.fromPacket(raw);
    expect(packet.toPacket(), raw);
  });

  test('PrPacket.fromPacket() invalid format throws ParseException', () {
    expect(() => PrPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
  });

  test('PrPacket.copyWith()', () {
    final original = PrPacket();
    final copy = original.copyWith();
    expect(copy.toPacket(), original.toPacket());
  });
}
