import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('PcPacket.parse()', () {
    String payload = '0;1;Hello world;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    payload = '<Pc>$payload$crc</Pc>';

    PcPacket link = PcPacket.fromPacket(payload);

    expect(link.timestamp, DateTime.fromMillisecondsSinceEpoch(0));
    expect(link.commandId, '1');
    expect(link.message, 'Hello world');

    String reversedPayload = link.toPacket();
    expect(reversedPayload, payload);
  });
}
