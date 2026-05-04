import 'package:test/test.dart';
import 'package:layrz_protocol/clients/tcp.dart';

void main() {
  group('LayrzProtocolSocket constructor', () {
    test('parses host and port correctly', () {
      expect(
        () => LayrzProtocolSocket(ident: 'IMEI1', server: 'tcp.example.com:5000'),
        returnsNormally,
      );
    });

    test('throws ArgumentError when server has no port', () {
      expect(
        () => LayrzProtocolSocket(ident: 'IMEI1', server: 'tcp.example.com'),
        throwsA(isA<ArgumentError>()),
      );
    });

    test('throws ArgumentError when port is 0', () {
      expect(
        () => LayrzProtocolSocket(ident: 'IMEI1', server: 'tcp.example.com:0'),
        throwsA(isA<ArgumentError>()),
      );
    });
  });

  group('LayrzProtocolSocket properties', () {
    test('composeEmptyPd returns PdPacket', () {
      final client = LayrzProtocolSocket(ident: 'IMEI1', server: 'tcp.example.com:5000');
      final pd = client.composeEmptyPd();
      expect(pd, isA<PdPacket>());
      expect(pd.extra, isEmpty);
    });

    test('splitRegExp splits multiple packets', () {
      final client = LayrzProtocolSocket(ident: 'IMEI1', server: 'tcp.example.com:5000');
      final ao = AoPacket(timestamp: DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true));
      final ar = ArPacket(reason: 'error');
      final combined = '${ao.toPacket()}${ar.toPacket()}';
      final parts = combined.split(client.splitRegExp).where((s) => s.isNotEmpty).toList();
      expect(parts.length, 2);
    });

    test('onEvent returns a Stream', () {
      final client = LayrzProtocolSocket(ident: 'IMEI1', server: 'tcp.example.com:5000');
      expect(client.onEvent, isA<Stream<LayrzTcpEvent>>());
    });

    test('disconnect returns true when not connected', () async {
      final client = LayrzProtocolSocket(ident: 'IMEI1', server: 'tcp.example.com:5000');
      final result = await client.disconnect();
      expect(result, isTrue);
    });

    test('sendData with no socket and no store callback is a no-op', () async {
      final client = LayrzProtocolSocket(ident: 'IMEI1', server: 'tcp.example.com:5000');
      final pd = PdPacket(
        timestamp: DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true),
        position: Position(latitude: 1.0, longitude: 1.0),
        extra: {'x': 1},
      );
      await client.sendData(pd);
    });

    test('sendData with no socket calls onBlackboxStore callback', () async {
      final stored = <String>[];
      final client = LayrzProtocolSocket(
        ident: 'IMEI1',
        server: 'tcp.example.com:5000',
        onBlackboxStore: (packet) async => stored.add(packet),
        onBlackboxFetch: () async => [],
      );
      final pd = PdPacket(
        timestamp: DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true),
        position: Position(),
        extra: {},
      );
      await client.sendData(pd);
      expect(stored.length, 1);
      expect(stored.first, pd.toPacket());
    });

    test('sendSos with no active socket does not throw', () async {
      final client = LayrzProtocolSocket(ident: 'IMEI1', server: 'tcp.example.com:5000');
      await client.sendSos();
    });

    test('sendImage with no active socket does not throw', () async {
      final client = LayrzProtocolSocket(ident: 'IMEI1', server: 'tcp.example.com:5000');
      await client.sendImage(bytes: [0x01, 0x02], filename: 'img.jpg');
    });
  });

  group('LayrzTcpEvent types', () {
    test('TcpConnected toString', () {
      final event = TcpConnected();
      expect(event.toString(), 'TcpConnected()');
    });

    test('TcpDisconnected toString', () {
      final event = TcpDisconnected();
      expect(event.toString(), 'TcpDisconnected()');
    });

    test('MessageReceived toString', () {
      final ao = AoPacket(timestamp: DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true));
      final event = MessageReceived(message: ao);
      expect(event.toString(), startsWith('MessageReceived(message:'));
    });
  });
}
