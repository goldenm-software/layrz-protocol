import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:test/test.dart';
import 'package:layrz_protocol/servers/tcp.dart';
import 'package:layrz_protocol/packets/packets.dart';

Future<int> findFreePort() async {
  final server = await ServerSocket.bind(InternetAddress.loopbackIPv4, 0);
  final port = server.port;
  await server.close();
  return port;
}

void main() {
  group('TcpServer', () {
    test('Construction validation: port 0 throws ArgumentError', () {
      expect(
        () => TcpServer(TcpConfig(port: 0, onNewPacket: (_, __) => null)),
        throwsA(isA<ArgumentError>()),
      );
    });

    test('Construction validation: port -1 throws ArgumentError', () {
      expect(
        () => TcpServer(TcpConfig(port: -1, onNewPacket: (_, __) => null)),
        throwsA(isA<ArgumentError>()),
      );
    });

    test('Construction validation: port 65535 throws ArgumentError', () {
      expect(
        () => TcpServer(TcpConfig(port: 65535, onNewPacket: (_, __) => null)),
        throwsA(isA<ArgumentError>()),
      );
    });

    test('Construction validation: port 70000 throws ArgumentError', () {
      expect(
        () => TcpServer(TcpConfig(port: 70000, onNewPacket: (_, __) => null)),
        throwsA(isA<ArgumentError>()),
      );
    });

    test('proxyProtocolV2 throws ArgumentError', () {
      expect(
        () => TcpServer(
          TcpConfig(
            port: 12345,
            proxyProtocolV2: true,
            onNewPacket: (_, __) => null,
          ),
        ),
        throwsA(isA<ArgumentError>()),
      );
    });

    test('Happy-path round trip', () async {
      final port = await findFreePort();
      bool handlerCalled = false;

      final server = TcpServer(
        TcpConfig(
          port: port,
          onNewPacket: (ClientPacket pkt, Socket conn) {
            handlerCalled = true;
            if (pkt is PaPacket) {
              return AsPacket();
            }
            return null;
          },
        ),
      );

      await server.start();
      addTearDown(() async => await server.close());

      final socket = await Socket.connect(InternetAddress.loopbackIPv4, port);
      addTearDown(() => socket.destroy());

      final paPacket = PaPacket(ident: 'test', password: 'pw');
      socket.write('${paPacket.toPacket()}\n');

      final responseLines = await socket.map((chunk) => utf8.decode(chunk)).first.timeout(const Duration(seconds: 2));

      expect(responseLines.contains('<As>'), true);
      expect(responseLines.contains('</As>'), true);
      expect(handlerCalled, true);
    });

    test('Multi-frame in one write', () async {
      final port = await findFreePort();
      final packetCalls = <ClientPacket>[];

      final server = TcpServer(
        TcpConfig(
          port: port,
          onNewPacket: (ClientPacket pkt, Socket conn) {
            packetCalls.add(pkt);
            if (pkt is PaPacket) {
              return AsPacket();
            }
            return null;
          },
        ),
      );

      await server.start();

      final socket = await Socket.connect(InternetAddress.loopbackIPv4, port);
      addTearDown(() => socket.destroy());

      final paPacket1 = PaPacket(ident: 'test1', password: 'pw1');
      final paPacket2 = PaPacket(ident: 'test2', password: 'pw2');

      socket.write('${paPacket1.toPacket()}${paPacket2.toPacket()}\n');

      await Future.delayed(const Duration(milliseconds: 200));

      expect(packetCalls.length, 2);
      expect(packetCalls[0], isA<PaPacket>());
      expect(packetCalls[1], isA<PaPacket>());
      expect((packetCalls[0] as PaPacket).ident, 'test1');
      expect((packetCalls[1] as PaPacket).ident, 'test2');

      addTearDown(() async => await server.close());
    });

    test('Decode error callback fires', () async {
      final port = await findFreePort();
      Object? capturedError;
      List<int>? capturedData;
      bool callbackCalled = false;

      final server = TcpServer(
        TcpConfig(
          port: port,
          onNewPacket: (_, __) => null,
          onDecodeError: (Object err, List<int> data, Socket conn) {
            callbackCalled = true;
            capturedError = err;
            capturedData = data;
          },
        ),
      );

      await server.start();
      addTearDown(() async => await server.close());

      final socket = await Socket.connect(InternetAddress.loopbackIPv4, port);
      addTearDown(() => socket.destroy());

      socket.write('garbage\n');

      await Future.delayed(const Duration(milliseconds: 150));

      expect(callbackCalled, true);
      expect(capturedError, isNotNull);
      expect(capturedData, isNotNull);

      socket.write('test\n');

      await Future.delayed(const Duration(milliseconds: 100));

      expect(socket.done, isA<Future>());
    });

    test('close() resolves serve()', () async {
      final port = await findFreePort();

      final server = TcpServer(
        TcpConfig(
          port: port,
          onNewPacket: (_, __) => null,
        ),
      );

      await server.start();
      addTearDown(() async => await server.close());

      final serveFuture = server.serve();

      await Future.delayed(const Duration(milliseconds: 50));

      expect(serveFuture, isA<Future<void>>());

      await server.close();

      await serveFuture.timeout(const Duration(seconds: 2));

      expect(true, true);
    });
  });
}
