import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'dart:typed_data';

import 'package:layrz_protocol/packets/packets.dart';

typedef TcpOnNewPacket = FutureOr<ServerPacket?> Function(ClientPacket packet, Socket conn);
typedef TcpOnDecodeError = void Function(Object err, List<int> data, Socket conn);

class TcpConfig {
  final int port;
  final TcpOnNewPacket onNewPacket;
  final TcpOnDecodeError? onDecodeError;
  final bool proxyProtocolV2;

  const TcpConfig({
    required this.port,
    required this.onNewPacket,
    this.onDecodeError,
    this.proxyProtocolV2 = false,
  });
}

class TcpServer {
  final TcpConfig _cfg;
  ServerSocket? _listener;
  final Set<Socket> _conns = {};
  final Completer<void> _done = Completer<void>();

  TcpServer._(this._cfg);

  factory TcpServer(TcpConfig cfg) {
    if (cfg.port <= 0 || cfg.port >= 65535) {
      throw ArgumentError('Port must be between 1 and 65534 inclusive');
    }
    if (cfg.proxyProtocolV2) {
      throw ArgumentError('Proxy Protocol v2 is not supported on Dart');
    }
    return TcpServer._(cfg);
  }

  int get port => _listener?.port ?? 0;

  Future<void> start() async {
    _listener = await ServerSocket.bind(InternetAddress.anyIPv4, _cfg.port);
    _acceptLoop();
  }

  Future<void> serve() async {
    await _done.future;
  }

  Future<void> close() async {
    final listener = _listener;
    if (listener != null) {
      await listener.close();
      _listener = null;
    }
    for (final socket in List<Socket>.from(_conns)) {
      socket.destroy();
    }
    _conns.clear();
    if (!_done.isCompleted) _done.complete();
  }

  void _acceptLoop() {
    _listener?.listen(
      _handleConnection,
      onError: (Object err) {
        if (!_done.isCompleted) _done.completeError(err);
      },
      onDone: () {
        if (!_done.isCompleted) _done.complete();
      },
    );
  }

  static final RegExp _splitter = RegExp(r'(?=<P[A-Za-z]>)');

  void _handleConnection(Socket socket) {
    _conns.add(socket);
    final buffer = BytesBuilder(copy: false);

    socket.listen(
      (List<int> chunk) async {
        buffer.add(chunk);
        if (!chunk.contains(0x0A)) return;

        final bytes = buffer.takeBytes();
        String text;
        try {
          text = utf8.decode(bytes, allowMalformed: true);
        } catch (e) {
          final cb = _cfg.onDecodeError;
          if (cb != null) {
            cb(e, bytes, socket);
          } else {
            stderr.writeln('LayrzProtocol TcpServer decode error: $e');
          }
          return;
        }

        final lastNl = text.lastIndexOf('\n');
        String complete = text;
        String tail = '';
        if (lastNl >= 0) {
          complete = text.substring(0, lastNl + 1);
          tail = text.substring(lastNl + 1);
        }
        if (tail.isNotEmpty) buffer.add(utf8.encode(tail));

        for (final raw in complete.split(_splitter)) {
          final trimmed = raw.replaceAll('\n', '').replaceAll('\r', '').trim();
          if (trimmed.isEmpty) continue;

          ClientPacket packet;
          try {
            packet = Packet.decodeClientPacket(trimmed);
          } catch (e) {
            final data = utf8.encode(trimmed);
            final cb = _cfg.onDecodeError;
            if (cb != null) {
              cb(e, data, socket);
            } else {
              stderr.writeln('LayrzProtocol TcpServer decode error: $e data: $trimmed');
            }
            continue;
          }

          try {
            final response = await _cfg.onNewPacket(packet, socket);
            if (response != null) {
              socket.add(utf8.encode(response.toPacket()));
            }
          } catch (e) {
            stderr.writeln('LayrzProtocol TcpServer handler error: $e');
          }
        }
      },
      onError: (Object err) {
        socket.destroy();
        _conns.remove(socket);
      },
      onDone: () {
        socket.destroy();
        _conns.remove(socket);
      },
      cancelOnError: true,
    );
  }
}
