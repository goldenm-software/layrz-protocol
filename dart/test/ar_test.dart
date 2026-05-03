import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('Packet.fromPacket() routes ArPacket', () {
    final original = ArPacket(reason: 'error');
    final parsed = Packet.fromPacket(original.toPacket());
    expect(parsed, isA<ArPacket>());
  });

  test('ArPacket.parse()', () {
    String payload = 'CRC mismatch;'; // timestamp
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    payload = '<Ar>$payload$crc</Ar>';

    ArPacket link = ArPacket.fromPacket(payload);
    expect(link.reason, 'CRC mismatch');

    String reversedPayload = link.toPacket();
    expect(reversedPayload, payload);
  });
}
