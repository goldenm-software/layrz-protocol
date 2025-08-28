part of '../../layrz_protocol.dart';

class LayrzProtocolSocket {
  /// [ident] is the identifier of the device, this [ident] should exists in the Layrz ecosystem.
  final String ident;

  /// [password] is the password of the device, this [password] can be empty if the device does not have a password.
  final String password;

  /// [httpUrl] defines the URL of the endpoint to interact with the comm interface.
  final String server;

  /// [version] is the version of the protocol.
  final LayrzProtocolVersion version;

  /// [LayrzProtocolSocket] is the class that contains the methods to interact with the Layrz ecosystem.
  ///
  /// Also, all of them may throw an exception if something goes wrong:
  /// - [ServerException] if the server returns an error 500.
  /// - [MalformedException] if the message is malformed, or when the server returns an unexpected response.
  /// - [ParseException] if the message is not well formatted.
  /// - [CrcException] if the CRC is not valid.
  /// - [CommandException] if the command is not well formatted.
  LayrzProtocolSocket({
    required this.ident,
    this.password = '',
    required this.server,
    this.version = LayrzProtocolVersion.v2,
  }) : assert(ident.isNotEmpty) {
    if (server.contains(':')) {
      _host = server.split(':')[0];
      _port = int.tryParse(server.split(':')[1]) ?? 0;
      if (_port <= 0) {
        throw ArgumentError('The port should be a number and greater than 0');
      }
    } else {
      throw ArgumentError('The server should contain the port');
    }

    try {
      _db = LinkDatabase();
    } catch (_) {
      LayrzLogging.critical('Error initializing the database, continuing without Blackbox support');
      _db = null;
    }
  }

  /// [_db] is the database to store the messages (Blackbox)
  LinkDatabase? _db;

  /// [blackboxSending] refers to the state of Blackbox processing.
  static bool blackboxSending = false;

  /// [_eventController] is the controller of the events.
  final _eventController = StreamController<LayrzTcpEvent>.broadcast();

  /// [onEvent] is the stream of events.
  Stream<LayrzTcpEvent> get onEvent => _eventController.stream;

  /// [isActive] indicates if the service is enabled.
  /// This is used to avoid miss-savings on blackbox
  static bool isActive = false;

  /// [splitRegExp] is the regular expression to split the packets.
  /// Sometimes, the socket connection sent multiple packets at once.
  RegExp get splitRegExp => RegExp(r'(?=<(?:A\w{1})>)');

  /// [_port] is the port of the server.
  /// This should be a number greater than 0.
  late int _port;

  /// [_host] is the host of the server.
  /// This should be a valid hostname or IP address.
  late String _host;

  /// [_socket] is the socket connection to the server.
  /// It is used to send and receive messages.
  Socket? _socket;

  /// [connect] connects to the server.
  Future<bool> connect({Duration timeout = const Duration(seconds: 5)}) async {
    try {
      Completer<bool> completer = Completer<bool>();
      _socket = await Socket.connect(_host, _port, timeout: timeout);
      _eventController.add(TcpConnected());
      _socket!.listen(
        (List<int> event) {
          String raw = utf8.decode(event);

          final packets = raw
              .split(splitRegExp)
              .where((message) {
                return message.isNotEmpty;
              })
              .map((message) {
                return message.trim();
              })
              .toList();

          for (final packet in packets) {
            try {
              final parsedPacket = Packet.fromPacket(packet);
              if (parsedPacket is AuPacket) {
                LayrzLogging.info('AuPacket deprecated, skipping');
                return;
              }

              if (parsedPacket is AsPacket) {
                LayrzProtocolSocket.isActive = true;
                if (!completer.isCompleted) completer.complete(true);
              }
              _eventController.add(MessageReceived(message: parsedPacket));
            } catch (e) {
              LayrzLogging.critical('Error parsing packet: "$packet" - $e');
              disconnect();
            }
          }
        },
        onError: (err) async {
          LayrzLogging.debug('------------------> onError $err');
          await disconnect();
        },
        onDone: () async {
          LayrzLogging.debug('------------------> onDone');
          await disconnect();
        },
      );

      final auth = PaPacket(ident: ident, password: password).toPacket();
      _socket!.writeln(auth);

      return await completer.future.timeout(
        timeout,
        onTimeout: () {
          if (!completer.isCompleted) {
            completer.complete(false);
          }
          return false;
        },
      );
    } catch (e) {
      LayrzLogging.critical('Error connecting to the server: $e');
      _eventController.add(TcpDisconnected());
      LayrzProtocolSocket.isActive = false;
      return false;
    }
  }

  /// [disconnect] disconnects from the server.
  Future<bool> disconnect() async {
    LayrzProtocolSocket.isActive = false;
    LayrzLogging.info('Disconnecting from the server');
    await _socket?.close();
    _socket?.destroy();
    _socket = null;
    _eventController.add(TcpDisconnected());
    LayrzLogging.info('Disconnected from the server');
    return true;
  }

  /// [sendData] sends a plain message to the Layrz ecosystem.
  Future<void> sendData(ClientPacket message) async {
    if (_socket == null) {
      LayrzLogging.warning('The socket is not connected, saving on Blackbox');
      if (_db == null) {
        LayrzLogging.warning('Blackbox is not initialized, ignoring message');
        return;
      }
      await _db!
          .into(_db!.messages)
          .insert(
            MessagesCompanion.insert(
              message: message.toPacket(),
              createdAt: DateTime.now(),
            ),
          );

      LayrzLogging.info('Saved on Blackbox');
      return;
    }

    final packet = message.toPacket();
    try {
      _socket?.writeln(packet);
    } catch (e) {
      LayrzLogging.critical('Error sending packet: $packet - $e');
      await disconnect();
    }

    _validateAndSendBlackbox();
  }

  /// [_validateAndSendBlackbox] validates if the service is enabled and sends the messages on the Blackbox.
  void _validateAndSendBlackbox() async {
    if (!LayrzProtocolSocket.isActive) {
      LayrzLogging.warning('Service is not enabled, ignoring Blackbox');
      return;
    }

    if (_db == null) {
      LayrzLogging.warning('Blackbox is not initialized, ignoring message');
      return;
    }

    if (LayrzProtocolSocket.blackboxSending) return;

    LayrzProtocolSocket.blackboxSending = true;
    final messages = await _db!.select(_db!.messages).get();
    if (messages.isEmpty) return;

    LayrzLogging.info('Sending ${messages.length} messages from Blackbox');
    for (final message in messages) {
      final packet = Packet.fromPacket(message.message);
      if (packet is ClientPacket) {
        await sendData(packet);
      }
      await _db!.delete(_db!.messages).delete(message);
    }

    LayrzProtocolSocket.blackboxSending = false;
  }

  /// [sendSos] sends an SOS message to the Layrz ecosystem.
  Future<void> sendSos([PdPacket? message]) {
    PdPacket msg = message ?? composeEmptyPd();
    Map<String, dynamic> extra = {...msg.extra};
    extra['alarm.event'] = true;
    msg = msg.copyWith(extra: extra);
    return sendData(msg);
  }

  /// [sendImage] sends an image to the Layrz ecosystem.
  Future<void> sendImage({
    /// [bytes] is the list of bytes of the image.
    required List<int> bytes,

    /// [filename] is the name of the file without the extension.
    required String filename,

    /// [contentType] is the content type of the image.
    /// By default is 'image/jpeg'.
    String contentType = 'image/jpeg',
  }) async {
    final packet = PmPacket(
      filename: filename,
      contentType: contentType,
      data: Uint8List.fromList(bytes),
    );

    return sendData(packet);
  }

  /// [composeEmptyPd] composes an empty PdPacket with the current timestamp and a default position.
  PdPacket composeEmptyPd() {
    return PdPacket(timestamp: DateTime.now(), position: Position(), extra: {});
  }
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
