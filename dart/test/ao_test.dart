import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('AoPacket.parse()', () {
    String payload = '1;'; // timestamp
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    payload = '<Ao>$payload$crc</Ao>';

    AoPacket link = AoPacket.fromPacket(payload);
    expect(link.timestamp, '1');

    String reversedPayload = link.toPacket();
    expect(reversedPayload, payload);
  });
}
