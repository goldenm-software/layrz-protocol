part of '../../packets.dart';

class TePacket extends ClientPacket {
  /// [timestamp] is the time of the packet.
  /// This is identified in the packet as `UNIX`
  final DateTime timestamp;

  /// [tripId] is the trip identifier of the packet.
  /// This is identified in the packet as `TRIP_ID`
  final String tripId;

  /// [TePacket] is the Trip end packet.
  ///
  /// This packet is part of the packet sent between Layrz services to identify trips.
  TePacket({
    required this.timestamp,
    required this.tripId,
  });

  /// [fromPacket] creates a [TePacket] from a string packet in the format of `Layrz Protocol v3`.
  static TePacket fromPacket(String raw) {
    if (!raw.startsWith('<Te>') || !raw.endsWith('</Te>')) {
      throw ParseException('Invalid identification packet, should be <Te>...</Te>');
    }

    final parts = raw.substring(4, raw.length - 5).split(';');
    if (parts.length != 3) {
      throw MalformedException('Invalid packet parts, should have 3 parts');
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

    return TePacket(
      timestamp: timestamp,
      tripId: parts[1],
    );
  }

  /// [toPacket] returns the packet in the format of `Layrz Protocol v3`.
  @override
  String toPacket() {
    String payload = '${(timestamp.millisecondsSinceEpoch / 1000).round()};';
    payload += '$tripId;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    return '<Te>$payload$crc</Te>';
  }

  @override
  String toString() => toPacket();

  @override
  TePacket copyWith({
    DateTime? timestamp,
    String? tripId,
  }) {
    return TePacket(
      timestamp: timestamp ?? this.timestamp,
      tripId: tripId ?? this.tripId,
    );
  }
}
