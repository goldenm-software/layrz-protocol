part of '../../packets.dart';

class AbPacket extends ServerPacket {
  final List<BleData> devices;

  /// [AbPacket] is the device list packet.
  ///
  /// This packet is part of the packet sent from the server to the device.
  AbPacket({
    /// [devices] is the list of devices that are being ACKed.
    /// This is identified in the packet as `BLE+DATA`
    required this.devices,
  }) : super();

  /// [fromPacket] creates a [AbPacket] from a string packet in the format of `Layrz Protocol v3`.
  static AbPacket fromPacket(String raw) {
    if (!raw.startsWith('<Ab>') || !raw.endsWith('</Ab>')) {
      throw ParseException('Invalid identification packet, should be <Ab>...</Ab>');
    }

    final parts = raw.substring(4, raw.length - 5).split(';');
    int? receivedCrc = int.tryParse(parts[parts.length - 1], radix: 16);
    int? calculatedCrc = calculateCrc("${parts.sublist(0, parts.length - 1).join(';')};".codeUnits);

    if (receivedCrc != calculatedCrc) {
      throw CrcException('Invalid CRC, received: $receivedCrc, calculated: $calculatedCrc');
    }

    return AbPacket(devices: BleData.fromPackets(parts.sublist(0, parts.length - 1).join(';')));
  }

  /// [toPacket] returns the packet in the format of `Layrz Link Protocol v2`.
  @override
  String toPacket() {
    String payload = devices.map((e) => e.toPacket()).join(';');
    payload += ';';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    return '<Ab>$payload$crc</Ab>';
  }

  @override
  String toString() => toPacket();

  @override
  AbPacket copyWith({
    List<BleData>? devices,
  }) {
    return AbPacket(devices: devices ?? this.devices);
  }
}
