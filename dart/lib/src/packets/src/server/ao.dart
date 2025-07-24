part of '../../packets.dart';

class AoPacket extends ServerPacket {
  final DateTime timestamp;

  /// [AoPacket] is the ACK package.
  ///
  /// This package is part of the package sent from the server to the device.
  ///
  /// Specifically, this package is sent when the server doesn't have commands waiting for the device.
  AoPacket({
    /// [messageId] is the message ID that is being ACKed.
    /// This is identified in the package as `MSG_ID`
    required this.timestamp,
  }) : super();

  /// [fromPacket] creates a [AoPacket] from a string package in the format of `Layrz Protocol v2`.
  static AoPacket fromPacket(String raw) {
    if (!raw.startsWith('<Ao>') || !raw.endsWith('</Ao>')) {
      throw ParseException('Invalid identification package, should be <Ao>...</Ao>');
    }

    final parts = raw.substring(4, raw.length - 5).split(';');
    if (parts.length != 2) {
      throw MalformedException('Invalid package parts, should have 2 parts');
    }

    int? receivedCrc = int.tryParse(parts[1], radix: 16);
    int? calculatedCrc = calculateCrc("${parts.sublist(0, 1).join(';')};".codeUnits);

    if (receivedCrc != calculatedCrc) {
      throw CrcException('Invalid CRC, received: $receivedCrc, calculated: $calculatedCrc');
    }

    try {
      return AoPacket(timestamp: DateTime.parse(parts[0]));
    } catch (e) {
      throw ParseException('Invalid timestamp');
    }
  }

  /// [toPacket] returns the package in the format of `Layrz Protocol v2`.
  @override
  String toPacket() {
    String payload = '${timestamp.millisecondsSinceEpoch ~/ 1000};';

    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    return '<Ao>$payload$crc</Ao>';
  }

  @override
  String toString() => toPacket();

  @override
  AoPacket copyWith({
    DateTime? timestamp,
  }) {
    return AoPacket(
      timestamp: timestamp ?? this.timestamp,
    );
  }
}
