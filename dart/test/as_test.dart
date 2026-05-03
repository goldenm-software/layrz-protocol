import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('Packet.fromPacket() routes AsPacket', () {
    final original = AsPacket();
    final parsed = Packet.fromPacket(original.toPacket());
    expect(parsed, isA<AsPacket>());
  });

  test('AsPacket.toPacket()', () {
    final packet = AsPacket();
    final raw = packet.toPacket();

    expect(raw.startsWith('<As>'), true);
    expect(raw.endsWith('</As>'), true);

    String payload = ';';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(raw, '<As>$payload$crc</As>');
  });

  test('AsPacket.fromPacket()', () {
    String payload = ';';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    String raw = '<As>$payload$crc</As>';

    final packet = AsPacket.fromPacket(raw);
    expect(packet.toPacket(), raw);
  });

  test('AsPacket.fromPacket() invalid format throws ParseException', () {
    expect(() => AsPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
  });

  test('AsPacket.copyWith()', () {
    final original = AsPacket();
    final copy = original.copyWith();
    expect(copy.toPacket(), original.toPacket());
  });
}
