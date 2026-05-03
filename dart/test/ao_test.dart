import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('Packet.fromPacket() routes AoPacket', () {
    final original = AoPacket(timestamp: DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true));
    final parsed = Packet.fromPacket(original.toPacket());
    expect(parsed, isA<AoPacket>());
  });

  test('AoPacket.parse()', () {
    String payload = '1;'; // timestamp
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    payload = '<Ao>$payload$crc</Ao>';

    AoPacket link = AoPacket.fromPacket(payload);
    expect(link.timestamp, DateTime.fromMillisecondsSinceEpoch(1 * 1000, isUtc: true));

    String reversedPayload = link.toPacket();
    expect(reversedPayload, payload);
  });
}
