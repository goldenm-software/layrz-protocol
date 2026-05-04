part of '../../packets.dart';

class PiPacket extends ClientPacket {
  /// [ident] is the Unique identifier, sent as part of the packet as `IMEI`
  final String ident;

  /// [firmwareId] is the firmware internal ID, this is in newer versions is a Layrz ID, otherwise is a unique
  /// identifier set by the hardware department of Layrz LTD. Also, is idenfified in the packet as `FW_ID`
  final String firmwareId;

  /// [firmwareBuild] is the firmware version, is an incremental number that is increased in each release.
  /// This is identified in the packet as `FW_BUILD`
  final int firmwareBuild;

  /// [deviceId] is the device internal ID, this is in newer versions is a Layrz ID, otherwise is a unique
  /// identifier set by the hardware department of Layrz LTD. Also, is idenfified in the packet as `SYS_DEV_ID`
  final String deviceId;

  /// [hardwareId] is the hardware internal ID, this is in newer versions is a Layrz ID, otherwise is a unique
  /// identifier set by the hardware department of Layrz LTD. Also, is idenfified in the packet as `SYS_DEV_HW_ID`
  final String hardwareId;

  /// [modelId] is the model internal ID, this is in newer versions is a Layrz ID, otherwise is a unique
  /// identifier set by the hardware department of Layrz LTD. Also, is idenfified in the packet
  /// as `SYS_DEV_MODEL_ID`
  final String modelId;

  /// [firmwareBranch] is the branch of the firmware, this is identified in the packet as `SYS_DEV_FW_BRANCH`
  final FirmwareBranch firmwareBranch;

  /// [fotaEnabled] is a boolean that indicates if the device is capable of receiving FOTA updates.
  /// This is identified in the packet as `FOTA_ENABLED`
  final bool fotaEnabled;

  /// [PiPacket] is the identification packet.
  ///
  /// This packet is part of the packet sent from the device to the server.
  ///
  /// This packet should be sent when the device starts up or the `get_info` command is received.
  PiPacket({
    required this.ident,
    required this.firmwareId,
    required this.firmwareBuild,
    required this.deviceId,
    required this.hardwareId,
    required this.modelId,
    required this.firmwareBranch,
    required this.fotaEnabled,
  });

  /// [fromPacket] creates a [PiPacket] from a string packet in the format of `Layrz Protocol v3`.
  static PiPacket fromPacket(String raw) {
    if (!raw.startsWith('<Pi>') || !raw.endsWith('</Pi>')) {
      throw ParseException('Invalid identification packet, should be <Pi>...</Pi>');
    }

    final parts = raw.substring(4, raw.length - 5).split(';');
    if (parts.length != 9) {
      throw MalformedException('Invalid packet parts, should have 6 parts');
    }

    int? receivedCrc = int.tryParse(parts[8], radix: 16);
    int? calculatedCrc = calculateCrc("${parts.sublist(0, 8).join(';')};".codeUnits);

    if (receivedCrc != calculatedCrc) {
      throw CrcException('Invalid CRC, received: $receivedCrc, calculated: $calculatedCrc');
    }

    int? firmwareBuild = int.tryParse(parts[2]);
    if (firmwareBuild == null) {
      throw MalformedException('Invalid firmware build number');
    }

    bool fotaEnabled = parts[7].toLowerCase() == 'true' || parts[7] == '1';

    return PiPacket(
      ident: parts[0],
      firmwareId: parts[1],
      firmwareBuild: firmwareBuild,
      deviceId: parts[3],
      hardwareId: parts[4],
      modelId: parts[5],
      firmwareBranch: parts[6] == '1' ? FirmwareBranch.development : FirmwareBranch.stable,
      fotaEnabled: fotaEnabled,
    );
  }

  /// [toPacket] returns the packet in the format of `Layrz Protocol v3`.
  ///
  /// Definition:
  /// `<Pi>IMEI;FW_ID;FW_BUILD;SYS_DEV_ID;SYS_DEV_HW_ID;SYS_DEV_MD_ID;SYS_DEV_FW_BRANCH;FOTA_ENABLED;CRC16</Pi>`
  @override
  String toPacket() {
    String payload = '$ident;';
    payload += '$firmwareId;';
    payload += '$firmwareBuild;';
    payload += '$deviceId;';
    payload += '$hardwareId;';
    payload += '$modelId;';
    payload += '${firmwareBranch.toPacket()};';
    payload += '${fotaEnabled ? '1' : '0'};';

    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    return '<Pi>$payload$crc</Pi>';
  }

  @override
  String toString() => toPacket();

  @override
  PiPacket copyWith({
    String? ident,
    String? firmwareId,
    int? firmwareBuild,
    String? deviceId,
    String? hardwareId,
    String? modelId,
    FirmwareBranch? firmwareBranch,
    bool? fotaEnabled,
  }) {
    return PiPacket(
      ident: ident ?? this.ident,
      firmwareId: firmwareId ?? this.firmwareId,
      firmwareBuild: firmwareBuild ?? this.firmwareBuild,
      deviceId: deviceId ?? this.deviceId,
      hardwareId: hardwareId ?? this.hardwareId,
      modelId: modelId ?? this.modelId,
      firmwareBranch: firmwareBranch ?? this.firmwareBranch,
      fotaEnabled: fotaEnabled ?? this.fotaEnabled,
    );
  }
}

extension FirmwareBranchExtension on FirmwareBranch {
  /// [toPacket] converts a [FirmwareBranch] to a string.
  String toPacket() {
    switch (this) {
      case FirmwareBranch.development:
        return '1';
      case FirmwareBranch.stable:
      default:
        return '0';
    }
  }
}
