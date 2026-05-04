import 'dart:async';
import 'dart:convert';
import 'dart:io' as io;
import 'dart:typed_data';

import 'package:layrz_protocol/packets/packets.dart';

typedef HttpOnNewPacket = FutureOr<ServerPacket?> Function(ClientPacket packet, io.HttpRequest req);
typedef HttpOnPullCommands = FutureOr<ServerPacket?> Function(String ident, String passwd, io.HttpRequest req);
typedef HttpOnAuthenticate = FutureOr<bool> Function(String ident, String passwd, io.HttpRequest req);
typedef HttpOnDecodeError = void Function(Object err, List<int> data, io.HttpRequest req);

class HttpConfig {
  final int port;
  final HttpOnNewPacket onNewPacket;
  final HttpOnPullCommands? onPullCommands;
  final HttpOnAuthenticate? onAuthenticate;
  final HttpOnDecodeError? onDecodeError;
  final int maxBodyBytes;
  final Duration shutdownTimeout;

  const HttpConfig({
    required this.port,
    required this.onNewPacket,
    this.onPullCommands,
    this.onAuthenticate,
    this.onDecodeError,
    this.maxBodyBytes = 1 << 20,
    this.shutdownTimeout = const Duration(seconds: 5),
  });
}

class HttpServer {
  final HttpConfig _cfg;
  io.HttpServer? _ioServer;
  final Completer<void> _done = Completer<void>();

  HttpServer._(this._cfg);

  factory HttpServer(HttpConfig cfg) {
    if (cfg.port <= 0 || cfg.port >= 65535) {
      throw ArgumentError('Port must be between 1 and 65534 inclusive');
    }
    return HttpServer._(cfg);
  }

  int get port => _ioServer?.port ?? 0;

  Future<void> start() async {
    _ioServer = await io.HttpServer.bind(io.InternetAddress.anyIPv4, _cfg.port);
    _serveRequests(_ioServer!);
  }

  Future<void> serve() async {
    await _done.future;
  }

  Future<void> close() async {
    try {
      await _ioServer?.close(force: false).timeout(_cfg.shutdownTimeout);
    } catch (_) {
      await _ioServer?.close(force: true);
    }
    if (!_done.isCompleted) _done.complete();
  }

  void _serveRequests(io.HttpServer server) {
    server.listen(
      (io.HttpRequest req) async {
        try {
          if (req.uri.path == '/v2/message') {
            await _handleMessage(req);
          } else if (req.uri.path == '/v2/commands') {
            await _handleCommands(req);
          } else {
            req.response.statusCode = io.HttpStatus.notFound;
            await req.response.close();
          }
        } catch (e) {
          try {
            req.response.statusCode = io.HttpStatus.internalServerError;
            await req.response.close();
          } catch (_) {}
        }
      },
      onError: (Object err) {
        if (!_done.isCompleted) _done.completeError(err);
      },
      onDone: () {
        if (!_done.isCompleted) _done.complete();
      },
    );
  }

  Future<void> _handleMessage(io.HttpRequest req) async {
    if (req.method != 'POST') {
      req.response.statusCode = io.HttpStatus.methodNotAllowed;
      await req.response.close();
      return;
    }

    final creds = _parseLayrzAuth(req.headers.value('Authorization'));
    if (creds == null) {
      req.response.statusCode = io.HttpStatus.unauthorized;
      req.response.write('unauthorized');
      await req.response.close();
      return;
    }
    final (ident, passwd) = creds;

    final onAuth = _cfg.onAuthenticate;
    if (onAuth != null && !await onAuth(ident, passwd, req)) {
      req.response.statusCode = io.HttpStatus.unauthorized;
      req.response.write('unauthorized');
      await req.response.close();
      return;
    }

    final bodyBuilder = BytesBuilder();
    int total = 0;
    await for (final chunk in req) {
      total += chunk.length;
      if (total > _cfg.maxBodyBytes) {
        req.response.statusCode = io.HttpStatus.requestEntityTooLarge;
        req.response.write('request entity too large');
        await req.response.close();
        return;
      }
      bodyBuilder.add(chunk);
    }
    final bodyBytes = bodyBuilder.takeBytes();

    String bodyStr;
    try {
      bodyStr = utf8.decode(bodyBytes);
    } catch (e) {
      final cb = _cfg.onDecodeError;
      if (cb != null) cb(e, bodyBytes, req);
      req.response.statusCode = io.HttpStatus.badRequest;
      req.response.write('invalid encoding');
      await req.response.close();
      return;
    }

    ClientPacket packet;
    try {
      packet = Packet.decodeClientPacket(bodyStr.trim());
    } catch (e) {
      final cb = _cfg.onDecodeError;
      if (cb != null) cb(e, bodyBytes, req);
      req.response.statusCode = io.HttpStatus.badRequest;
      req.response.write('invalid packet');
      await req.response.close();
      return;
    }

    final response = await _cfg.onNewPacket(packet, req);
    if (response == null) {
      req.response.statusCode = io.HttpStatus.noContent;
      await req.response.close();
      return;
    }

    req.response.statusCode = io.HttpStatus.ok;
    req.response.headers.contentType = io.ContentType('text', 'plain', charset: 'utf-8');
    req.response.write(response.toPacket());
    await req.response.close();
  }

  Future<void> _handleCommands(io.HttpRequest req) async {
    if (req.method != 'GET') {
      req.response.statusCode = io.HttpStatus.methodNotAllowed;
      await req.response.close();
      return;
    }

    final creds = _parseLayrzAuth(req.headers.value('Authorization'));
    if (creds == null) {
      req.response.statusCode = io.HttpStatus.unauthorized;
      req.response.write('unauthorized');
      await req.response.close();
      return;
    }
    final (ident, passwd) = creds;

    final onAuth = _cfg.onAuthenticate;
    if (onAuth != null && !await onAuth(ident, passwd, req)) {
      req.response.statusCode = io.HttpStatus.unauthorized;
      req.response.write('unauthorized');
      await req.response.close();
      return;
    }

    final onPull = _cfg.onPullCommands;
    if (onPull == null) {
      req.response.statusCode = io.HttpStatus.noContent;
      await req.response.close();
      return;
    }

    final response = await onPull(ident, passwd, req);
    if (response == null) {
      req.response.statusCode = io.HttpStatus.noContent;
      await req.response.close();
      return;
    }

    req.response.statusCode = io.HttpStatus.ok;
    req.response.headers.contentType = io.ContentType('text', 'plain', charset: 'utf-8');
    req.response.write(response.toPacket());
    await req.response.close();
  }
}

(String, String)? _parseLayrzAuth(String? header) {
  if (header == null || !header.startsWith('LayrzAuth ')) return null;
  final rest = header.substring('LayrzAuth '.length);
  final sep = rest.indexOf(';');
  if (sep < 0) return null;
  return (rest.substring(0, sep), rest.substring(sep + 1));
}
