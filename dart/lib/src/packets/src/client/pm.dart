part of '../../packets.dart';

class PmPacket extends ClientPacket {
  /// [filename] is the name of the file.
  String filename;

  /// [contentType] is the content type of the file.
  String contentType;

  /// [data] is the content of the file.
  Uint8List data;

  /// [PmPacket] is the media package.
  ///
  /// This package is part of the package sent from the server to the device.
  PmPacket({
    required this.filename,
    required this.contentType,
    required this.data,
  }) : super();

  /// [fromPacket] creates a [PmPacket] from a string package in the format of `Layrz Protocol v3`.
  static PmPacket fromPacket(String raw) {
    if (!raw.startsWith('<Pm>') || !raw.endsWith('</Pm>')) {
      throw ParseException('Invalid identification package, should be <Pm>...</Pm>');
    }

    final parts = raw.substring(4, raw.length - 5).split(';');
    if (parts.length != 4) {
      throw MalformedException('Invalid package parts, should have 4 parts');
    }

    int? receivedCrc = int.tryParse(parts[3], radix: 16);
    int? calculatedCrc = calculateCrc("${parts.sublist(0, 1).join(';')};".codeUnits);

    if (receivedCrc != calculatedCrc) {
      throw CrcException('Invalid CRC, received: $receivedCrc, calculated: $calculatedCrc');
    }

    String filename = parts[0];
    String contentType = parts[1];
    Uint8List data;
    try {
      data = base64Decode(parts[2]);
    } catch (e) {
      throw ParseException('Invalid base64 data');
    }

    return PmPacket(
      filename: filename,
      contentType: contentType,
      data: data,
    );
  }

  /// [toPacket] returns the package in the format of `Layrz Protocol v3`.
  @override
  String toPacket() {
    String payload = '$filename;$contentType;${base64Encode(data)};';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    return '<Pm>$payload$crc</Pm>';
  }

  @override
  String toString() => toPacket();

  @override
  PmPacket copyWith({
    String? filename,
    String? contentType,
    Uint8List? data,
  }) {
    return PmPacket(
      filename: filename ?? this.filename,
      contentType: contentType ?? this.contentType,
      data: data ?? this.data,
    );
  }
}
