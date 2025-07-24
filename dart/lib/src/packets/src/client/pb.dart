part of '../../packets.dart';

class PbPacket extends ClientPacket {
  /// [advertisements] is the list of advertisements detected by the device.
  final List<BleAdvertisement> advertisements;

  /// [PbPacket] is the Bluetooth Low Energy Detection package.
  /// This package is sent by the device to the server.
  ///
  /// This package is part of the `Layrz Global Network`, a new initiative to create
  /// a global network of devices, anywhere, anytime.
  PbPacket({required this.advertisements});

  /// [fromPacket] creates a [PbPacket] from a string package in the format of `Layrz Protocol v2`.
  static PbPacket fromPacket(String raw) {
    if (!raw.startsWith('<Pb>') || !raw.endsWith('</Pb>')) {
      throw ParseException('Invalid identification package, should be <Pb>...</Pb>');
    }

    final parts = raw.substring(4, raw.length - 5).split(';');
    int? receivedCrc = int.tryParse(parts[parts.length - 1], radix: 16);
    int? calculatedCrc = calculateCrc("${parts.sublist(0, parts.length - 1).join(';')};".codeUnits);

    if (receivedCrc != calculatedCrc) {
      throw CrcException('Invalid CRC, received: $receivedCrc, calculated: $calculatedCrc');
    }

    return PbPacket(advertisements: BleAdvertisement.fromPacket(parts.sublist(0, parts.length - 1).join(';')));
  }

  /// [toPacket] returns the package in the format of `Layrz Protocol v2`.
  ///
  /// Definition:
  /// `<Pb>BLE_ADVERSIEMENT;BLE_ADVERSIEMENT;BLE_ADVERSIEMENT;CRC16</Pb>`
  @override
  String toPacket() {
    String payload = advertisements.map((adv) => adv.toPacket()).join(';');
    payload += ';';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    return '<Pb>$payload$crc</Pb>';
  }

  @override
  String toString() => toPacket();

  @override
  PbPacket copyWith({List<BleAdvertisement>? advertisements}) {
    return PbPacket(advertisements: advertisements ?? this.advertisements);
  }
}
