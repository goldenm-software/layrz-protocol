part of '../../packets.dart';

class ArPacket extends ServerPacket {
  final String reason;

  /// [ArPacket] is the error ACK packet.
  ///
  /// This packet is part of the packet sent from the server to the device.
  ArPacket({
    /// [reason] is the reason of the error.
    /// This is identified in the packet as `REASON`
    required this.reason,
  }) : super();

  /// [fromPacket] creates a [ArPacket] from a string packet in the format of `Layrz Protocol v3`.
  static ArPacket fromPacket(String raw) {
    if (!raw.startsWith('<Ar>') || !raw.endsWith('</Ar>')) {
      throw ParseException('Invalid identification packet, should be <Ar>...</Ar>');
    }

    final parts = raw.substring(4, raw.length - 5).split(';');
    if (parts.length != 2) {
      throw MalformedException('Invalid packet parts, should have 2 parts');
    }

    int? receivedCrc = int.tryParse(parts[1], radix: 16);
    int? calculatedCrc = calculateCrc("${parts.sublist(0, 1).join(';')};".codeUnits);

    if (receivedCrc != calculatedCrc) {
      throw CrcException('Invalid CRC, received: $receivedCrc, calculated: $calculatedCrc');
    }

    return ArPacket(reason: parts[0]);
  }

  /// [toPacket] returns the packet in the format of `Layrz Protocol v3`.
  @override
  String toPacket() {
    String payload = '$reason;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    return '<Ar>$payload$crc</Ar>';
  }

  @override
  String toString() => toPacket();

  @override
  ArPacket copyWith({
    String? reason,
  }) {
    return ArPacket(
      reason: reason ?? this.reason,
    );
  }
}
