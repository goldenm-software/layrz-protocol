part of '../../packets.dart';

class PaPacket extends ClientPacket {
  /// [ident] is the Unique identifier, sent as part of the package as `IMEI`
  final String ident;

  /// [password] is the password of the device.
  final String password;

  /// [PaPacket] is authentication package. Only used over TCP connections.
  PaPacket({
    required this.ident,
    required this.password,
  });

  /// [fromPacket] creates a [PaPacket] from a string package in the format of `Layrz Protocol v2`.
  static PaPacket fromPacket(String raw) {
    if (!raw.startsWith('<Pa>') || !raw.endsWith('</Pa>')) {
      throw ParseException('Invalid identification package, should be <Pa>...</Pa>');
    }

    final parts = raw.substring(4, raw.length - 5).split(';');
    int? receivedCrc = int.tryParse(parts[parts.length - 1], radix: 16);
    int? calculatedCrc = calculateCrc("${parts.sublist(0, parts.length - 1).join(';')};".codeUnits);

    if (receivedCrc != calculatedCrc) {
      throw CrcException('Invalid CRC, received: $receivedCrc, calculated: $calculatedCrc');
    }

    return PaPacket(ident: parts[0], password: parts[1]);
  }

  /// [toPacket] returns the package in the format of `Layrz Protocol v2`.
  ///
  /// Definition:
  /// `<Pa>IMEI;PASSWORD;CRC16</Pa>`
  @override
  String toPacket() {
    String payload = '$ident;$password;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    return '<Pa>$payload$crc</Pa>';
  }

  @override
  String toString() => toPacket();

  @override
  PaPacket copyWith({
    String? ident,
    String? password,
  }) {
    return PaPacket(
      ident: ident ?? this.ident,
      password: password ?? this.password,
    );
  }
}
