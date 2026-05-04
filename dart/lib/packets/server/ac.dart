part of '../../packets.dart';

class AcPacket extends ServerPacket {
  final List<Command> commands;

  /// [AcPacket] is the error ACK packet.
  ///
  /// This packet is part of the packet sent from the server to the device.
  AcPacket({
    /// [commands] is the list of commands that are being ACKed.
    /// This is identified in the packet as `CMD+DEFINITION`
    required this.commands,
  }) : super();

  /// [fromPacket] creates a [AcPacket] from a string packet in the format of `Layrz Protocol v3`.
  static AcPacket fromPacket(String raw) {
    if (!raw.startsWith('<Ac>') || !raw.endsWith('</Ac>')) {
      throw ParseException('Invalid identification packet, should be <Ac>...</Ac>');
    }

    final parts = raw.substring(4, raw.length - 5).split(';');
    int? receivedCrc = int.tryParse(parts[parts.length - 1], radix: 16);
    int? calculatedCrc = calculateCrc("${parts.sublist(0, parts.length - 1).join(';')};".codeUnits);

    if (receivedCrc != calculatedCrc) {
      throw CrcException('Invalid CRC, received: $receivedCrc, calculated: $calculatedCrc');
    }

    return AcPacket(commands: Command.fromPackets(parts.sublist(0, parts.length - 1).join(';')));
  }

  /// [toPacket] returns the packet in the format of `Layrz Protocol v3`.
  @override
  String toPacket() {
    String payload = commands.map((e) => e.toPacket()).join(';');
    payload += ';';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    return '<Ac>$payload$crc</Ac>';
  }

  @override
  String toString() => toPacket();

  @override
  AcPacket copyWith({
    List<Command>? commands,
  }) {
    return AcPacket(
      commands: commands ?? this.commands,
    );
  }
}
