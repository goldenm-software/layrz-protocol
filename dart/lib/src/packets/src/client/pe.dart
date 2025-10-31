part of '../../packets.dart';

class PePacket extends ClientPacket {
  /// [timestamp] is the time of the package.
  /// This is identified in the package as `UNIX`
  final DateTime timestamp;

  /// [tripId] is the trip identifier of the package.
  /// This is identified in the package as `TRIP_ID`
  final String tripId;

  /// [PePacket] is the Trip end package.
  ///
  /// This package is part of the package sent from the device to the server.
  PePacket({
    required this.timestamp,
    required this.tripId,
  });

  /// [fromPacket] creates a [PePacket] from a string package in the format of `Layrz Protocol v3`.
  static PePacket fromPacket(String raw) {
    if (!raw.startsWith('<Pe>') || !raw.endsWith('</Pe>')) {
      throw ParseException('Invalid identification package, should be <Pe>...</Pe>');
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

    return PePacket(
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

    return '<Pe>$payload$crc</Pe>';
  }

  @override
  String toString() => toPacket();

  @override
  PePacket copyWith({
    DateTime? timestamp,
    String? tripId,
  }) {
    return PePacket(
      timestamp: timestamp ?? this.timestamp,
      tripId: tripId ?? this.tripId,
    );
  }
}
