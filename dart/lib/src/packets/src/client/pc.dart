part of '../../packets.dart';

class PcPacket extends ClientPacket {
  /// [timestamp] is the time of the package.
  /// This is identified in the package as `UNIX`
  final DateTime timestamp;

  /// [commandId] is the command ID that is being ACKed.
  /// This is identified in the package as `CMD_ID`
  final String commandId;

  /// [message] is the message of the ACK.
  /// This is identified in the package as `MSG`
  final String message;

  /// [PcPacket] is the command ACK package.
  ///
  /// This package is part of the package sent from the device to the server.
  ///
  /// This package should be send after the execution of a command, received in [Ac] package.
  PcPacket({
    required this.timestamp,
    required this.commandId,
    required this.message,
  });

  /// [fromPacket] creates a [PcPacket] from a string package in the format of `Layrz Protocol v2`.
  static PcPacket fromPacket(String raw) {
    if (!raw.startsWith('<Pc>') || !raw.endsWith('</Pc>')) {
      throw ParseException('Invalid identification package, should be <Pc>...</Pc>');
    }

    final parts = raw.substring(4, raw.length - 5).split(';');
    if (parts.length != 4) {
      throw MalformedException('Invalid package parts, should have 4 parts');
    }

    int? receivedCrc = int.tryParse(parts[3], radix: 16);
    int? calculatedCrc = calculateCrc("${parts.sublist(0, 3).join(';')};".codeUnits);

    if (receivedCrc != calculatedCrc) {
      throw CrcException('Invalid CRC, received: $receivedCrc, calculated: $calculatedCrc');
    }

    DateTime timestamp;
    try {
      timestamp = DateTime.fromMillisecondsSinceEpoch(int.parse(parts[0]) * 1000);
    } catch (e) {
      throw MalformedException('Invalid timestamp');
    }

    return PcPacket(
      timestamp: timestamp,
      commandId: parts[1],
      message: parts[2],
    );
  }

  /// [toPacket] returns the package in the format of `Layrz Protocol v2`.
  @override
  String toPacket() {
    String payload = '${(timestamp.millisecondsSinceEpoch / 1000).round()};$commandId;$message;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    return '<Pc>$payload$crc</Pc>';
  }

  @override
  String toString() => toPacket();

  @override
  PcPacket copyWith({
    DateTime? timestamp,
    String? commandId,
    String? message,
  }) {
    return PcPacket(
      timestamp: timestamp ?? this.timestamp,
      commandId: commandId ?? this.commandId,
      message: message ?? this.message,
    );
  }
}
