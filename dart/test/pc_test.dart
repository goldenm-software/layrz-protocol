import 'package:test/test.dart';
import 'package:layrz_protocol/packets/packets.dart';
import 'package:layrz_protocol/utils/crc.dart';

void main() {
  test('Packet.fromPacket() routes PcPacket', () {
    final ts = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    final original = PcPacket(timestamp: ts, commandId: 1, message: 'ok');
    final parsed = Packet.fromPacket(original.toPacket());
    expect(parsed, isA<PcPacket>());
  });

  test('PcPacket.parse()', () {
    String payload = '0;1;Hello world;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    payload = '<Pc>$payload$crc</Pc>';

    PcPacket link = PcPacket.fromPacket(payload);

    expect(link.timestamp, DateTime.fromMillisecondsSinceEpoch(0, isUtc: true));
    expect(link.commandId, 1);
    expect(link.message, 'Hello world');

    String reversedPayload = link.toPacket();
    expect(reversedPayload, payload);
  });
}
