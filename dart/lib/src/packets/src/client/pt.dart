part of '../../packets.dart';

class PtPacket extends ClientPacket {
  /// [timestamp] is the time of the package.
  /// This is identified in the package as `UNIX`
  final DateTime timestamp;

  /// [tripId] is the trip identifier of the package.
  /// This is identified in the package as `TRIP_ID`
  final String tripId;

  /// [PtPacket] is the Trip start package.
  ///
  /// This package is part of the package sent from the device to the server.
  PtPacket({
    required this.timestamp,
    required this.tripId,
  });

  /// [fromPacket] creates a [PtPacket] from a string package in the format of `Layrz Protocol v3`.
  static PtPacket fromPacket(String raw) {
    if (!raw.startsWith('<Pt>') || !raw.endsWith('</Pt>')) {
      throw ParseException('Invalid identification package, should be <Pt>...</Pt>');
    }

    final parts = raw.substring(4, raw.length - 5).split(';');
    if (parts.length != 3) {
      throw MalformedException('Invalid package parts, should have 3 parts');
    }

    int? receivedCrc = int.tryParse(parts[2], radix: 16);
    int? calculatedCrc = calculateCrc("${parts.sublist(0, 2).join(';')};".codeUnits);

    if (receivedCrc != calculatedCrc) {
      throw CrcException('Invalid CRC, received: $receivedCrc, calculated: $calculatedCrc');
    }

    DateTime timestamp;
    try {
      timestamp = DateTime.fromMillisecondsSinceEpoch(int.parse(parts[0]) * 1000);
    } catch (e) {
      throw MalformedException('Invalid timestamp');
    }

    return PtPacket(
      timestamp: timestamp,
      tripId: parts[1],
    );
  }

  /// [toPacket] returns the package in the format of `Layrz Protocol v3`.
  @override
  String toPacket() {
    String payload = '${(timestamp.millisecondsSinceEpoch / 1000).round()};';
    payload += '$tripId;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    return '<Pt>$payload$crc</Pt>';
  }

  @override
  String toString() => toPacket();

  @override
  PtPacket copyWith({
    DateTime? timestamp,
    String? tripId,
  }) {
    return PtPacket(
      timestamp: timestamp ?? this.timestamp,
      tripId: tripId ?? this.tripId,
    );
  }
}
