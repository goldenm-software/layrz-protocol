part of '../../packets.dart';

class PsPacket extends ClientPacket {
  /// [timestamp] is the time of the package.
  /// This is identified in the package as `UNIX`
  final DateTime timestamp;

  /// [params] is the configuration parameters of the package.
  /// This is identified in the package as `EXTRA+ARGS`
  final Map<String, dynamic> params;

  /// [PsPacket] is the configuration package.
  ///
  /// This package is part of the package sent from the device to the server.
  ///
  /// Also, this packet only will be sent when `get_config` or `set_config` command is received.
  PsPacket({
    required this.timestamp,
    required this.params,
  });

  /// [fromPacket] creates a [PsPacket] from a string package in the format of `Layrz Protocol v2`.
  static PsPacket fromPacket(String raw) {
    if (!raw.startsWith('<Ps>') || !raw.endsWith('</Ps>')) {
      throw ParseException('Invalid identification package, should be <Ps>...</Ps>');
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

    return PsPacket(
      timestamp: timestamp,
      params: Packet.parseExtraArgs(parts[1]),
    );
  }

  /// [toPacket] returns the package in the format of `Layrz Protocol v2`.
  @override
  String toPacket() {
    String payload = '${(timestamp.millisecondsSinceEpoch / 1000).round()};';
    List<String> extraList = [];
    params.forEach((key, value) {
      if (value is bool) {
        extraList.add('$key:${value ? 'true' : 'false'}');
      } else if (value is int) {
        extraList.add('$key:$value');
      } else if (value is double) {
        extraList.add('$key:$value');
      } else {
        extraList.add('$key:$value');
      }
    });

    payload += "${extraList.join(',')};";

    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    return '<Ps>$payload$crc</Ps>';
  }

  @override
  String toString() => toPacket();

  @override
  PsPacket copyWith({
    DateTime? timestamp,
    Map<String, dynamic>? params,
  }) {
    return PsPacket(
      timestamp: timestamp ?? this.timestamp,
      params: params ?? this.params,
    );
  }
}
