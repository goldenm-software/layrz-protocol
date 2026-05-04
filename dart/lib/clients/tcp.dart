// coverage:ignore-file
import 'dart:io';
import 'dart:async';
import 'dart:convert';
import 'dart:typed_data';
import 'package:collection/collection.dart';
import 'package:logging/logging.dart';

import '../packets/packets.dart';
import '../utils/protocol.dart';

export '../packets/packets.dart';
export '../utils/errors.dart';
export '../utils/protocol.dart';

final Logger _log = Logger('layrz_protocol.tcp');

/// Callback invoked when the socket is disconnected and a packet needs to be queued.
/// The consumer is responsible for persisting [packet] for later replay.
typedef BlackboxStoreCallback = Future<void> Function(String packet);

/// Callback invoked on reconnect to retrieve all queued packets.
/// The consumer should return the stored packets and clear their store.
typedef BlackboxFetchCallback = Future<List<String>> Function();

class LayrzProtocolSocket {
  final String ident;
  final String password;
  final String server;
  final LayrzProtocolVersion version;

  /// Called when a packet needs to be saved to the Blackbox (socket is down).
  final BlackboxStoreCallback? onBlackboxStore;

  /// Called on reconnect to retrieve and flush queued Blackbox packets.
  final BlackboxFetchCallback? onBlackboxFetch;

  LayrzProtocolSocket({
    required this.ident,
    this.password = '',
    required this.server,
    this.version = LayrzProtocolVersion.v2,
    this.onBlackboxStore,
    this.onBlackboxFetch,
  }) : assert(ident.isNotEmpty) {
    if (server.contains(':')) {
      _host = server.split(':')[0];
      _port = int.tryParse(server.split(':')[1]) ?? 0;
      if (_port <= 0) throw ArgumentError('The port should be a number and greater than 0');
    } else {
      throw ArgumentError('The server should contain the port');
    }
  }

  static bool blackboxSending = false;
  static bool isActive = false;

  final _eventController = StreamController<LayrzTcpEvent>.broadcast();
  Stream<LayrzTcpEvent> get onEvent => _eventController.stream;

  RegExp get splitRegExp => RegExp(r'(?=<(?:A\w{1})>)');

  late int _port;
  late String _host;
  Socket? _socket;

  Future<InternetAddress> _resolveHost() async {
    final addresses = await InternetAddress.lookup(_host);
    final ipv6 = addresses.firstWhereOrNull((addr) => addr.type == InternetAddressType.IPv6);
    if (ipv6 != null) {
      _log.info('Resolved host $_host to IPv6: ${ipv6.address}');
      return ipv6;
    }
    _log.info('Resolved host $_host to IPv4: ${addresses.first.address}');
    return addresses.first;
  }

  Future<bool> connect({Duration timeout = const Duration(seconds: 5)}) async {
    try {
      final completer = Completer<bool>();
      final address = await _resolveHost();
      _socket = await Socket.connect(address, _port, timeout: timeout);
      _eventController.add(TcpConnected());
      _socket!.listen(
        (List<int> event) {
          final raw = utf8.decode(event);
          final packets = raw.split(splitRegExp).where((m) => m.isNotEmpty).map((m) => m.trim()).toList();

          for (final packet in packets) {
            try {
              final parsed = Packet.fromPacket(packet);
              if (parsed is AuPacket) {
                _log.info('AuPacket deprecated, skipping');
                continue;
              }
              if (parsed is AsPacket) {
                LayrzProtocolSocket.isActive = true;
                if (!completer.isCompleted) completer.complete(true);
              }
              _eventController.add(MessageReceived(message: parsed));
            } catch (e) {
              _log.severe('Error parsing packet: "$packet" - $e');
              disconnect();
            }
          }
        },
        onError: (err) async {
          _log.fine('onError: $err');
          await disconnect();
        },
        onDone: () async {
          _log.fine('onDone');
          await disconnect();
        },
      );

      _socket!.writeln(PaPacket(ident: ident, password: password).toPacket());

      return await completer.future.timeout(
        timeout,
        onTimeout: () {
          if (!completer.isCompleted) completer.complete(false);
          return false;
        },
      );
    } catch (e) {
      _log.severe('Error connecting to the server: $e');
      _eventController.add(TcpDisconnected());
      LayrzProtocolSocket.isActive = false;
      return false;
    }
  }

  Future<bool> disconnect() async {
    LayrzProtocolSocket.isActive = false;
    _log.info('Disconnecting from the server');
    await _socket?.close();
    _socket?.destroy();
    _socket = null;
    _eventController.add(TcpDisconnected());
    _log.info('Disconnected from the server');
    return true;
  }

  Future<void> sendData(ClientPacket message) async {
    if (_socket == null) {
      _log.warning('Socket not connected, saving to Blackbox');
      final store = onBlackboxStore;
      if (store == null) {
        _log.warning('No Blackbox store callback, ignoring message');
        return;
      }
      await store(message.toPacket());
      _log.info('Saved to Blackbox');
      return;
    }

    final packet = message.toPacket();
    try {
      _socket?.writeln(packet);
    } catch (e) {
      _log.severe('Error sending packet: $packet - $e');
      await disconnect();
    }

    _flushBlackbox();
  }

  void _flushBlackbox() async {
    if (!LayrzProtocolSocket.isActive) {
      _log.warning('Service not active, skipping Blackbox flush');
      return;
    }

    final fetch = onBlackboxFetch;
    if (fetch == null) return;

    if (LayrzProtocolSocket.blackboxSending) return;
    LayrzProtocolSocket.blackboxSending = true;

    final messages = await fetch();
    if (messages.isEmpty) {
      LayrzProtocolSocket.blackboxSending = false;
      return;
    }

    _log.info('Flushing ${messages.length} messages from Blackbox');
    for (final raw in messages) {
      final packet = Packet.fromPacket(raw);
      if (packet is ClientPacket) await sendData(packet);
    }

    LayrzProtocolSocket.blackboxSending = false;
  }

  Future<void> sendSos([PdPacket? message]) {
    final msg = (message ?? composeEmptyPd()).copyWith(
      extra: {
        ...(message ?? composeEmptyPd()).extra,
        'alarm.event': true,
      },
    );
    return sendData(msg);
  }

  Future<void> sendImage({
    required List<int> bytes,
    required String filename,
    String contentType = 'image/jpeg',
  }) {
    return sendData(PmPacket(filename: filename, contentType: contentType, data: Uint8List.fromList(bytes)));
  }

  PdPacket composeEmptyPd() => PdPacket(timestamp: DateTime.now(), position: Position(), extra: {});
}

abstract class LayrzTcpEvent {
  @override
  String toString() => 'LayrzTcpEvent()';
}

class TcpConnected extends LayrzTcpEvent {
  @override
  String toString() => 'TcpConnected()';
}

class TcpDisconnected extends LayrzTcpEvent {
  @override
  String toString() => 'TcpDisconnected()';
}

class MessageReceived extends LayrzTcpEvent {
  final Packet message;
  MessageReceived({required this.message});
  @override
  String toString() => 'MessageReceived(message: $message)';
}
