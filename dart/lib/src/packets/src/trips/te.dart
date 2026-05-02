part of '../../packets.dart';

class TePacket extends TripsPacket {
  /// [timestamp] is the time of the packet.
  /// This is identified in the packet as `UNIX`
  final DateTime timestamp;

  /// [tripId] is the trip identifier of the packet.
  /// This is identified in the packet as `TRIP_ID`
  final String tripId;

  /// [distanceTraveled] is the distance traveled during the trip in meters.
  /// This is identified in the packet as `DISTANCE_TRAVELED`
  final double distanceTraveled;

  /// [maxSpeed] is the maximum speed during the trip in km/h.
  /// This is identified in the packet as `MAX_SPEED`
  final double maxSpeed;

  /// [duration] is the duration of the trip in seconds.
  /// This is identified in the packet as `DURATION`
  final Duration duration;

  /// [TePacket] is the Trip end packet.
  ///
  /// This packet is part of the packet sent between Layrz services to identify trips.
  TePacket({
    required this.timestamp,
    required this.tripId,
    required this.distanceTraveled,
    required this.maxSpeed,
    required this.duration,
  });

  /// [fromPacket] creates a [TePacket] from a string packet in the format of `Layrz Protocol v3`.
  static TePacket fromPacket(String raw) {
    if (!raw.startsWith('<Te>') || !raw.endsWith('</Te>')) {
      throw ParseException('Invalid identification packet, should be <Te>...</Te>');
    }

    final parts = raw.substring(4, raw.length - 5).split(';');
    if (parts.length != 6) {
      throw MalformedException('Invalid packet parts, should have 6 parts');
    }

    int? receivedCrc = int.tryParse(parts.last, radix: 16);
    int? calculatedCrc = calculateCrc("${parts.sublist(0, 5).join(';')};".codeUnits);

    if (receivedCrc != calculatedCrc) {
      throw CrcException('Invalid CRC, received: $receivedCrc, calculated: $calculatedCrc');
    }

    DateTime timestamp;
    try {
      timestamp = DateTime.fromMillisecondsSinceEpoch(int.parse(parts[0]) * 1000);
    } catch (e) {
      throw MalformedException('Invalid timestamp');
    }

    double distanceTraveled;
    try {
      distanceTraveled = double.parse(parts[2]);
    } catch (e) {
      throw MalformedException('Invalid distance traveled ${parts[2]}');
    }

    double maxSpeed;
    try {
      maxSpeed = double.parse(parts[3]);
    } catch (e) {
      throw MalformedException('Invalid max speed');
    }

    Duration duration;
    try {
      duration = Duration(seconds: int.parse(parts[4]));
    } catch (e) {
      throw MalformedException('Invalid duration');
    }

    return TePacket(
      timestamp: timestamp,
      distanceTraveled: distanceTraveled,
      maxSpeed: maxSpeed,
      duration: duration,
      tripId: parts[1],
    );
  }

  /// [toPacket] returns the packet in the format of `Layrz Protocol v3`.
  @override
  String toPacket() {
    String payload = '${(timestamp.millisecondsSinceEpoch / 1000).round()};';
    payload += '$tripId;';
    payload += '${distanceTraveled.toStringAsFixed(3)};';
    payload += '${maxSpeed.toStringAsFixed(3)};';
    payload += '${duration.inSeconds};';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    return '<Te>$payload$crc</Te>';
  }

  @override
  String toString() => toPacket();

  @override
  TePacket copyWith({
    DateTime? timestamp,
    String? tripId,
    double? distanceTraveled,
    double? maxSpeed,
    Duration? duration,
  }) {
    return TePacket(
      timestamp: timestamp ?? this.timestamp,
      tripId: tripId ?? this.tripId,
      distanceTraveled: distanceTraveled ?? this.distanceTraveled,
      maxSpeed: maxSpeed ?? this.maxSpeed,
      duration: duration ?? this.duration,
    );
  }
}
