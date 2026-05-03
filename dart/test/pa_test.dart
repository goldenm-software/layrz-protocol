import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('Packet.fromPacket() routes PaPacket', () {
    final original = PaPacket(ident: '123456789012345', password: 'pass');
    final parsed = Packet.fromPacket(original.toPacket());
    expect(parsed, isA<PaPacket>());
  });

  test('PaPacket.toPacket()', () {
    final packet = PaPacket(ident: '123456789012345', password: 'secret');
    final raw = packet.toPacket();

    expect(raw.startsWith('<Pa>'), true);
    expect(raw.endsWith('</Pa>'), true);

    String payload = '123456789012345;secret;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(raw, '<Pa>$payload$crc</Pa>');
  });

  test('PaPacket.fromPacket()', () {
    String payload = '123456789012345;secret;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    String raw = '<Pa>$payload$crc</Pa>';

    final packet = PaPacket.fromPacket(raw);
    expect(packet.ident, '123456789012345');
    expect(packet.password, 'secret');

    expect(packet.toPacket(), raw);
  });

  test('PaPacket.fromPacket() invalid format throws ParseException', () {
    expect(() => PaPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
    expect(() => PaPacket.fromPacket('<Pa>missing</Pa>'), throwsA(isA<CrcException>()));
  });

  test('PaPacket.copyWith()', () {
    final original = PaPacket(ident: '111', password: 'pass');
    final copy = original.copyWith(ident: '222');
    expect(copy.ident, '222');
    expect(copy.password, 'pass');
  });
}
