part of '../../packets.dart';

class PdPacket extends ClientPacket {
  /// [timestamp] is the time of the packet.
  /// This is identified in the packet as `UNIX`
  final DateTime timestamp;

  /// [position] is the position of the packet.
  /// This is identified in the packet as `LAT`, `LON`, `ALT`, `SPD`, `DIR`, `SAT` and `HDOP`
  ///
  /// - `LAT` is the latitude of the packet.
  /// - `LON` is the longitude of the packet.
  /// - `ALT` is the altitude of the packet.
  /// - `SPD` is the speed of the packet.
  /// - `DIR` is the direction of the packet.
  /// - `SAT` is the number of satellites of the packet.
  /// - `HDOP` is the HDOP of the packet.
  ///
  /// All of the above values are in the packet separated by `;` and may be empty.
  final Position position;

  /// [extra] is the extra data of the packet.
  /// This is identified in the packet as `EXTRA+ARGS`.
  final Map<String, dynamic> extra;

  /// [PdPacket] is the data packet.
  ///
  /// This packet is part of the packet sent from the device to the server.
  ///
  ///  This packet should be sent passively by the device.
  PdPacket({
    required this.timestamp,
    required this.position,
    required this.extra,
  });

  /// [fromPacket] creates a [PdPacket] from a string packet in the format of `Layrz Protocol v3`.
  static PdPacket fromPacket(String raw) {
    if (!raw.startsWith('<Pd>') || !raw.endsWith('</Pd>')) {
      throw ParseException('Invalid identification packet, should be <Pd>...</Pd>');
    }

    final parts = raw.substring(4, raw.length - 5).split(';');
    if (parts.length != 10) {
      throw MalformedException('Invalid packet parts, should have 10 parts, received ${parts.length} parts');
    }

    int? receivedCrc = int.tryParse(parts[9], radix: 16);
    int? calculatedCrc = calculateCrc("${parts.sublist(0, 9).join(';')};".codeUnits);

    if (receivedCrc != calculatedCrc) {
      throw CrcException('Invalid CRC, received: $receivedCrc, calculated: $calculatedCrc');
    }

    DateTime timestamp;
    try {
      timestamp = DateTime.fromMillisecondsSinceEpoch(int.parse(parts[0]) * 1000);
    } catch (e) {
      throw MalformedException('Invalid timestamp');
    }

    final position = Position(
      latitude: double.tryParse(parts[1]),
      longitude: double.tryParse(parts[2]),
      altitude: double.tryParse(parts[3]),
      speed: double.tryParse(parts[4]),
      direction: double.tryParse(parts[5]),
      satellites: int.tryParse(parts[6]),
      hdop: double.tryParse(parts[7]),
    );

    return PdPacket(
      timestamp: timestamp,
      position: position,
      extra: Packet.parseExtraArgs(parts[8]),
    );
  }

  /// [toPacket] returns the packet in the format of `Layrz Protocol v3`.
  @override
  String toPacket() {
    String payload = '${timestamp.millisecondsSinceEpoch ~/ 1000};';
    payload += '${position.latitude ?? ''};';
    payload += '${position.longitude ?? ''};';
    payload += '${position.altitude ?? ''};';
    payload += '${position.speed ?? ''};';
    payload += '${position.direction ?? ''};';
    payload += '${position.satellites ?? ''};';
    payload += '${position.hdop ?? ''};';

    List<String> extraList = [];
    for (String key in extra.keys) {
      extraList.add('$key:${extra[key]}');
    }

    payload += '${extraList.join(',')};';

    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    return '<Pd>$payload$crc</Pd>';
  }

  @override
  String toString() => toPacket();

  @override
  PdPacket copyWith({
    DateTime? timestamp,
    Position? position,
    Map<String, dynamic>? extra,
  }) {
    return PdPacket(
      timestamp: timestamp ?? this.timestamp,
      position: position ?? this.position,
      extra: extra ?? this.extra,
    );
  }
}
