part of '../../layrz_protocol.dart';

class LayrzProtocolHttp {
  /// [ident] is the identifier of the device, this [ident] should exists in the Layrz ecosystem.
  final String ident;

  /// [password] is the password of the device, this [password] can be empty if the device does not have a password.
  final String password;

  /// [httpUrl] defines the URL of the endpoint to interact with the comm interface.
  final String server;

  /// [version] is the version of the protocol.
  final LayrzProtocolVersion version;

  /// [LayrzProtocolHttp] is the class that contains the methods to interact with the Layrz ecosystem.
  ///
  /// All of the methods are asynchronous and return a [Future<void>] that is a list of [Command].
  /// Also, all of them may throw an exception if something goes wrong:
  /// - [ServerException] if the server returns an error 500.
  /// - [MalformedException] if the message is malformed, or when the server returns an unexpected response.
  /// - [ParseException] if the message is not well formatted.
  /// - [CrcException] if the CRC is not valid.
  /// - [CommandException] if the command is not well formatted.
  LayrzProtocolHttp({
    required this.ident,
    this.password = '',
    required this.server,
    this.version = LayrzProtocolVersion.v2,
  }) : assert(ident.isNotEmpty) {
    _baseUrl = server;
    _dio = Dio(
      BaseOptions(
        baseUrl: baseUrl,
        headers: headers,
        responseDecoder: (bytes, options, responseBody) => utf8.decode(bytes, allowMalformed: true),
        requestEncoder: (data, options) => utf8.encode(data),
      ),
    );
  }

  late final Dio _dio;
  late String _baseUrl;

  /// [baseUrl] defines the URL of the endpoint to interact with the comm interface.
  /// This URL is built with the [server] and the [version].
  String get baseUrl => 'https://$_baseUrl/${version.value}';

  /// [headers] is the headers that will be sent in the request
  Map<String, String> get headers => {
    'Authorization': 'LayrzAuth $ident;$password',
    'Content-Type': 'text/plain',
  };

  /// [sendData] sends a plain message to the Layrz ecosystem.
  HttpPacketResponse sendData(ClientPacket message) async => _sendToLayrz(message);

  /// [sendSos] sends an SOS message to the Layrz ecosystem.
  HttpPacketResponse sendSos([PdPacket? message]) async {
    PdPacket msg = message ?? composeEmptyPd();
    Map<String, dynamic> extra = {...msg.extra};
    extra['alarm.event'] = true;
    return sendData(msg.copyWith(extra: extra));
  }

  /// [sendText] sends a text message to the Layrz ecosystem.
  HttpPacketResponse sendText(String text, {PdPacket? message}) async {
    PdPacket msg = message ?? composeEmptyPd();
    Map<String, dynamic> extra = {...msg.extra};
    extra['driver.message'] = text;
    return sendData(msg.copyWith(extra: extra));
  }

  /// [sendImage] sends an image to the Layrz ecosystem.
  HttpPacketResponse sendImage({
    /// [bytes] is the list of bytes of the image.
    required List<int> bytes,

    /// [filename] is the name of the file without the extension.
    required String filename,

    /// [contentType] is the content type of the image.
    /// By default is 'image/jpeg'.
    String contentType = 'image/jpeg',
  }) async {
    return sendData(PmPacket(filename: filename, contentType: contentType, data: Uint8List.fromList(bytes)));
  }

  /// [getCommands] ask to the server if there are commands to execute.
  HttpPacketResponse getCommands() async {
    final response = await _dio.get('/commands');
    return _processResponse(response);
  }

  /// [getBleDevices] ask to the server the BLE devices that should sniff.
  HttpPacketResponse getBleDevices() async {
    final response = await _dio.get('/ble');
    return _processResponse(response);
  }

  /// [_sendToLayrz] sends the final and formatted message to Layrz
  HttpPacketResponse _sendToLayrz(ClientPacket packet) async {
    final parsedMessage = packet.toPacket();
    final response = await _dio.post('/message', data: parsedMessage);
    return _processResponse(response);
  }

  HttpPacketResponse _processResponse(Response response) {
    if (response.statusCode == 500) {
      throw ServerException('Server returned an error 500');
    }

    try {
      return Future.value(Packet.fromPacket(response.data));
    } on FormatException {
      throw ParseException('The message is not well formatted');
    } on CrcException {
      throw CrcException('The CRC is not valid');
    } on CommandException {
      throw CommandException('The command is not well formatted');
    } catch (e) {
      throw MalformedException('The message is malformed');
    }
  }

  PdPacket composeEmptyPd() {
    return PdPacket(timestamp: DateTime.now(), position: Position(), extra: {});
  }
}
