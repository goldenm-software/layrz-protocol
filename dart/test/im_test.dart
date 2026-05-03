import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('Packet.fromPacket() routes ImPacket', () {
    final ts = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    final original = ImPacket(timestamp: ts, chatId: 'chat1', message: 'hello');
    final parsed = Packet.fromPacket(original.toPacket());
    expect(parsed, isA<ImPacket>());
  });

  test('ImPacket.parse()', () {
    final packet = ImPacket(
      timestamp: DateTime.utc(2025, 11, 4, 0, 0, 0),
      chatId: 'chat123',
      message: 'Hello; World',
    );

    final rawPacket = packet.toPacket();
    expect(rawPacket, '<Im>1762214400;chat123;Hello||| World;4C25</Im>');
    final parsedPacket = ImPacket.fromPacket(rawPacket);
    expect(parsedPacket.timestamp, packet.timestamp);
    expect(parsedPacket.chatId, packet.chatId);
    expect(parsedPacket.message, packet.message);
  });
}
