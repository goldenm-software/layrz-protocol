part of '../../packets.dart';

class PrPacket extends ClientPacket {
  /// [PrPacket] is the authentication success packet.
  ///
  /// This packet is part of the packet sent from the server to the device.
  PrPacket() : super();

  /// [fromPacket] creates a [PrPacket] from a string packet in the format of `Layrz Protocol v3`.
  static PrPacket fromPacket(String raw) {
    if (!raw.startsWith('<Pr>') || !raw.endsWith('</Pr>')) {
      throw ParseException('Invalid identification packet, should be <Pr>...</Pr>');
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

    return PrPacket();
  }

  /// [toPacket] returns the packet in the format of `Layrz Protocol v3`.
  @override
  String toPacket() {
    String payload = ';';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    return '<Pr>$payload$crc</Pr>';
  }

  @override
  String toString() => toPacket();

  @override
  PrPacket copyWith() {
    return PrPacket();
  }
}
