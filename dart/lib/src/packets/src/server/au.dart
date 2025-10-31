part of '../../packets.dart';

@Deprecated('This packet is deprecated and will be removed in v4.0')
class AuPacket extends ServerPacket {
  /// [AuPacket] is the authentication request packet.
  ///
  /// This packet is part of the packet sent from the server to the device.
  @Deprecated('This packet is deprecated and will be removed in v4.0')
  AuPacket() : super();

  /// [fromPacket] creates a [AuPacket] from a string packet in the format of `Layrz Protocol v3`.
  static AuPacket fromPacket(String raw) {
    if (!raw.startsWith('<Au>') || !raw.endsWith('</Au>')) {
      throw ParseException('Invalid identification packet, should be <Au>...</Au>');
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

    return AuPacket();
  }

  /// [toPacket] returns the packet in the format of `Layrz Protocol v3`.
  @override
  String toPacket() {
    String payload = ';';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    return '<Au>$payload$crc</Au>';
  }

  @override
  String toString() => toPacket();

  @override
  AuPacket copyWith() {
    return AuPacket();
  }
}
