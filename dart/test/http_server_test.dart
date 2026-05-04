import 'dart:async';
import 'dart:convert';
import 'dart:io' as io;

import 'package:test/test.dart';
import 'package:layrz_protocol/servers/http.dart';
import 'package:layrz_protocol/packets/packets.dart';

Future<int> findFreePort() async {
  final server = await io.ServerSocket.bind(io.InternetAddress.loopbackIPv4, 0);
  final port = server.port;
  await server.close();
  return port;
}

// Helper to send HTTP requests
Future<(int, String)> sendHttpRequest(
  String method,
  String path,
  int port, {
  String? body,
  Map<String, String>? headers,
}) async {
  final completer = Completer<(int, String)>();

  try {
    final socket = await io.Socket.connect(
      io.InternetAddress.loopbackIPv4,
      port,
      timeout: const Duration(seconds: 2),
    );

    final req = StringBuffer();
    req.writeln('$method $path HTTP/1.1');
    req.writeln('Host: localhost:$port');
    req.writeln('Connection: close');

    if (headers != null) {
      headers.forEach((key, value) {
        req.writeln('$key: $value');
      });
    }

    if (body != null) {
      req.writeln('Content-Length: ${body.length}');
    }

    req.writeln();
    if (body != null) {
      req.write(body);
    }

    socket.write(req.toString());
    await socket.flush();

    final lines = <String>[];
    await socket.map((chunk) => utf8.decode(chunk)).join().then((String response) {
      lines.addAll(response.split('\n'));
    });

    socket.destroy();

    int statusCode = 0;
    int bodyStart = 0;

    for (int i = 0; i < lines.length; i++) {
      final line = lines[i];
      if (line.startsWith('HTTP/')) {
        statusCode = int.parse(line.split(' ')[1]);
      }
      if (line.trim().isEmpty) {
        bodyStart = i + 1;
        break;
      }
    }

    final responseBody = lines.sublist(bodyStart).join('\n').trim();
    completer.complete((statusCode, responseBody));
  } catch (e) {
    completer.completeError(e);
  }

  return completer.future;
}

void main() {
  group('HttpServer', () {
    test('Construction validation: port 0 throws ArgumentError', () {
      expect(
        () => HttpServer(HttpConfig(port: 0, onNewPacket: (_, __) => null)),
        throwsA(isA<ArgumentError>()),
      );
    });

    test('Construction validation: port -1 throws ArgumentError', () {
      expect(
        () => HttpServer(HttpConfig(port: -1, onNewPacket: (_, __) => null)),
        throwsA(isA<ArgumentError>()),
      );
    });

    test('Construction validation: port 65535 throws ArgumentError', () {
      expect(
        () => HttpServer(HttpConfig(port: 65535, onNewPacket: (_, __) => null)),
        throwsA(isA<ArgumentError>()),
      );
    });

    test('Construction validation: port 70000 throws ArgumentError', () {
      expect(
        () => HttpServer(HttpConfig(port: 70000, onNewPacket: (_, __) => null)),
        throwsA(isA<ArgumentError>()),
      );
    });

    test('POST /v2/message happy path', () async {
      final port = await findFreePort();
      bool handlerCalled = false;

      final server = HttpServer(
        HttpConfig(
          port: port,
          onNewPacket: (ClientPacket pkt, io.HttpRequest req) {
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

      final paPacket = PaPacket(ident: 'test', password: 'pw');
      final (statusCode, responseBody) = await sendHttpRequest(
        'POST',
        '/v2/message',
        port,
        body: paPacket.toPacket(),
        headers: {'Authorization': 'LayrzAuth test;pw'},
      );

      expect(statusCode, 200);
      expect(responseBody.contains('<As>'), true);
      expect(responseBody.contains('</As>'), true);
      expect(handlerCalled, true);
    });

    test('401 on missing auth', () async {
      final port = await findFreePort();

      final server = HttpServer(
        HttpConfig(
          port: port,
          onNewPacket: (_, __) => null,
        ),
      );

      await server.start();
      addTearDown(() async => await server.close());

      final (statusCode, _) = await sendHttpRequest(
        'POST',
        '/v2/message',
        port,
        body: 'garbage',
      );

      expect(statusCode, 401);
    });

    test('401 on bad auth format', () async {
      final port = await findFreePort();

      final server = HttpServer(
        HttpConfig(
          port: port,
          onNewPacket: (_, __) => null,
        ),
      );

      await server.start();
      addTearDown(() async => await server.close());

      final (statusCode, _) = await sendHttpRequest(
        'POST',
        '/v2/message',
        port,
        body: 'garbage',
        headers: {'Authorization': 'Bearer foo'},
      );

      expect(statusCode, 401);
    });

    test('401 on onAuthenticate returning false', () async {
      final port = await findFreePort();

      final server = HttpServer(
        HttpConfig(
          port: port,
          onNewPacket: (_, __) => null,
          onAuthenticate: (_, __, ___) async => false,
        ),
      );

      await server.start();
      addTearDown(() async => await server.close());

      final (statusCode, _) = await sendHttpRequest(
        'POST',
        '/v2/message',
        port,
        body: 'garbage',
        headers: {'Authorization': 'LayrzAuth test;pw'},
      );

      expect(statusCode, 401);
    });

    test('405 on GET /v2/message', () async {
      final port = await findFreePort();

      final server = HttpServer(
        HttpConfig(
          port: port,
          onNewPacket: (_, __) => null,
        ),
      );

      await server.start();
      addTearDown(() async => await server.close());

      final (statusCode, _) = await sendHttpRequest(
        'GET',
        '/v2/message',
        port,
        headers: {'Authorization': 'LayrzAuth test;pw'},
      );

      expect(statusCode, 405);
    });

    test('400 on invalid packet body', () async {
      final port = await findFreePort();
      bool decodeErrorCalled = false;

      final server = HttpServer(
        HttpConfig(
          port: port,
          onNewPacket: (_, __) => null,
          onDecodeError: (_, __, ___) {
            decodeErrorCalled = true;
          },
        ),
      );

      await server.start();
      addTearDown(() async => await server.close());

      final (statusCode, _) = await sendHttpRequest(
        'POST',
        '/v2/message',
        port,
        body: 'garbage',
        headers: {'Authorization': 'LayrzAuth test;pw'},
      );

      expect(statusCode, 400);
      expect(decodeErrorCalled, true);
    });

    test('GET /v2/commands 204 when onPullCommands is null', () async {
      final port = await findFreePort();

      final server = HttpServer(
        HttpConfig(
          port: port,
          onNewPacket: (_, __) => null,
        ),
      );

      await server.start();
      addTearDown(() async => await server.close());

      final (statusCode, _) = await sendHttpRequest(
        'GET',
        '/v2/commands',
        port,
        headers: {'Authorization': 'LayrzAuth test;pw'},
      );

      expect(statusCode, 204);
    });

    test('GET /v2/commands 200 with packet when onPullCommands returns one', () async {
      final port = await findFreePort();

      final server = HttpServer(
        HttpConfig(
          port: port,
          onNewPacket: (_, __) => null,
          onPullCommands: (_, __, ___) async => AsPacket(),
        ),
      );

      await server.start();
      addTearDown(() async => await server.close());

      final (statusCode, responseBody) = await sendHttpRequest(
        'GET',
        '/v2/commands',
        port,
        headers: {'Authorization': 'LayrzAuth test;pw'},
      );

      expect(statusCode, 200);
      expect(responseBody.contains('<As>'), true);
      expect(responseBody.contains('</As>'), true);
    });

    test('close() resolves serve()', () async {
      final port = await findFreePort();

      final server = HttpServer(
        HttpConfig(
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
