/// [LayrzProtocolMode] defines the protocol to use.
///
/// Please, read each mode to understand the differences between them.
enum LayrzProtocolMode {
  /// [tcp] defines the protocol to use TCP.
  tcp,

  /// [http] defines the protocol to use HTTP.
  http;

  String get value {
    switch (this) {
      case LayrzProtocolMode.tcp:
        return 'TCP';
      case LayrzProtocolMode.http:
      default:
        return 'HTTP';
    }
  }

  static LayrzProtocolMode fromString(String value) {
    switch (value) {
      case 'TCP':
        return LayrzProtocolMode.tcp;
      case 'HTTP':
      default:
        return LayrzProtocolMode.http;
    }
  }
}

/// [LayrzProtocolVersion] defines the version of the protocol.
enum LayrzProtocolVersion {
  /// [v2] is the version 2 of the protocol. Enhances capabilities like commands, images, and BLE devices.
  v2;

  String get value {
    switch (this) {
      case LayrzProtocolVersion.v2:
      default:
        return 'v2';
    }
  }

  static LayrzProtocolVersion fromString(String value) {
    switch (value) {
      case 'v2':
      default:
        return LayrzProtocolVersion.v2;
    }
  }
}
