part of '../../packets.dart';

class BleData {
  /// [macAddress] is the detected Mac Address. This Mac Adress came from the detected device, not the
  /// device that is detecting.
  final String macAddress;

  /// [model] is the model of the detected device. This model should be equals to the model of the device
  /// and the model defined by Layrz.
  final String model;

  /// [BleData] is the BLE data packet.
  ///
  /// This packet is part of the packet sent from the server to the device.
  const BleData({
    required this.macAddress,
    required this.model,
  });

  /// [fromPackets] creates a [List<BleData>] from a raw message following this structure:
  /// MAC_ADDRESS:MODEL
  static List<BleData> fromPackets(String raw) {
    final parts = raw.split(';');
    if (parts.isEmpty) return [];

    if (parts.length == 1) return [];
    // The parts should be divisible by 2, can be multiple groups

    // Separate each group
    final List<BleData> devices = [];

    for (int i = 0; i < parts.length; i++) {
      final part = parts[i].split(':');
      if (part.length != 2) {
        throw MalformedException('Invalid BLE data definition');
      }

      final rawMacAddress = part[0];
      String macAddress = '';
      for (int i = 0; i < rawMacAddress.length; i += 2) {
        macAddress += '${rawMacAddress[i]}${rawMacAddress[i + 1]}';
        if (i + 2 < rawMacAddress.length) {
          macAddress += ':';
        }
      }

      final model = part[1];

      devices.add(BleData(macAddress: macAddress, model: model));
    }

    return devices;
  }

  /// [toPacket] converts the [BleData] to a raw message following this structure:
  /// MAC_ADDRESS:MODEL
  String toPacket() {
    return '${macAddress.replaceAll(':', '').toUpperCase()}:$model';
  }
}

class BleAdvertisement {
  /// [macAddress] is the detected Mac Address. This Mac Adress came from the detected device, not the
  /// device that is detecting.
  final String macAddress;

  /// [timestamp] is when the device was detected.
  final DateTime timestamp;

  /// [latitude] is the closest latitude of the device. Defined by the device that is detecting.
  /// This value is optional
  final double? latitude;

  /// [longitude] is the closest longitude of the device. Defined by the device that is detecting.
  /// This value is optional
  final double? longitude;

  /// [altitude] is the closest altitude of the device. Defined by the device that is detecting.
  /// This value is optional
  final double? altitude;

  /// [rssi] is the signal strength of the detected device.
  final int rssi;

  /// [txPower] is the transmission power of the detected device.
  /// This value is optional
  final int? txPower;

  /// [model] is the model of the detected device. This model should be equals to the model of the device
  /// and the model defined by Layrz.
  final String model;

  /// [manufacturerData] is the list of manufacturer data advertised by the device.
  List<BleManufacturerData> manufacturerData;

  /// [serviceData] is the list of service data advertised by the device.
  List<BleServiceData> serviceData;

  /// [deviceName] is the name of the device. This name is optional.
  final String? deviceName;

  BleAdvertisement({
    required this.macAddress,
    required this.timestamp,
    this.latitude,
    this.longitude,
    this.altitude,
    required this.rssi,
    this.txPower,
    required this.model,
    this.manufacturerData = const [],
    this.serviceData = const [],
    this.deviceName,
  });

  /// [fromPacket] creates a [List<BleAdvertisement>] from a raw message following this structure:
  /// MAC_ADDRESS;UNIX;LAT;LNG;ALT;MODEL;RSSI;MANUFACTURER+DATA;SERVICE+DATA;CRC16
  static List<BleAdvertisement> fromPacket(String raw) {
    final parts = raw.split(';');
    if (parts.isEmpty) return [];

    if (parts.length == 1) return [];
    // The parts should be divisible by 4, can be multiple groups

    if (parts.length % 12 != 0) {
      throw ParseException('Invalid advertisement definition');
    }

    // Separate each group
    final List<BleAdvertisement> advertisements = [];

    for (int i = 0; i < parts.length; i += 12) {
      final rawMacAddress = parts[i];
      final rawUnix = parts[i + 1];
      final rawLatitude = parts[i + 2];
      final rawLongitude = parts[i + 3];
      final rawAltitude = parts[i + 4];
      final model = parts[i + 5];
      final deviceName = parts[i + 6];
      final rawRssi = parts[i + 7];
      final rawTxPower = parts[i + 8];
      final rawManufacturerData = parts[i + 9];
      final rawServiceData = parts[i + 10];

      final receivedCrc = int.tryParse(parts[i + 11], radix: 16) ?? 0;

      final calculatedCrc = calculateCrc('${parts.sublist(i, i + 11).join(';')};'.codeUnits);
      if (receivedCrc != calculatedCrc) {
        throw CrcException(
          'Invalid CRC, expected ${receivedCrc.toRadixString(16)}, '
          'got ${calculatedCrc.toRadixString(16)}',
        );
      }

      final macParts = rawMacAddress.split('');
      if (macParts.length != 12) {
        throw MalformedException('Invalid MAC Address');
      }

      String macAddress = '';
      for (int i = 0; i < macParts.length; i += 2) {
        macAddress += '${macParts[i]}${macParts[i + 1]}';
        if (i + 2 < macParts.length) {
          macAddress += ':';
        }
      }

      DateTime timestamp;
      try {
        timestamp = DateTime.fromMillisecondsSinceEpoch(int.parse(rawUnix) * 1000);
      } catch (e) {
        throw MalformedException('Invalid timestamp');
      }

      double? latitude;
      if (rawLatitude.isNotEmpty) {
        latitude = double.tryParse(rawLatitude);
      }

      double? longitude;
      if (rawLongitude.isNotEmpty) {
        longitude = double.tryParse(rawLongitude);
      }

      double? altitude;
      if (rawAltitude.isNotEmpty) {
        altitude = double.tryParse(rawAltitude);
      }

      int rssi = int.tryParse(rawRssi) ?? 0;
      int? txPower = int.tryParse(rawTxPower);

      List<BleManufacturerData> manufacturerData = [];
      for (String mfd in rawManufacturerData.split(',')) {
        if (mfd.isEmpty) continue;
        final subparts = mfd.split(':');
        if (subparts.isEmpty) continue;

        if (subparts.length != 2) {
          throw MalformedException('Invalid manufacturer data "$mfd"');
        }

        final companyId = int.tryParse(subparts[0], radix: 16) ?? 0;
        List<int> data = [];
        for (int i = 0; i < subparts[1].length; i += 2) {
          data.add(int.tryParse(subparts[1].substring(i, i + 2), radix: 16) ?? 0);
        }
        manufacturerData.add(BleManufacturerData(companyId: companyId, data: data));
      }

      List<BleServiceData> serviceData = [];

      for (String sfd in rawServiceData.split(',')) {
        if (sfd.isEmpty) continue;
        final subparts = sfd.split(':');
        if (subparts.isEmpty) continue;

        if (subparts.length != 2) {
          throw MalformedException('Invalid service data "$sfd"');
        }

        final uuid = int.tryParse(subparts[0], radix: 16) ?? 0;
        List<int> data = [];
        for (int i = 0; i < subparts[1].length; i += 2) {
          data.add(int.tryParse(subparts[1].substring(i, i + 2), radix: 16) ?? 0);
        }

        serviceData.add(BleServiceData(uuid: uuid, data: data));
      }

      advertisements.add(
        BleAdvertisement(
          deviceName: deviceName,
          macAddress: macAddress,
          timestamp: timestamp,
          rssi: rssi,
          model: model,
          latitude: latitude,
          longitude: longitude,
          altitude: altitude,
          txPower: txPower,
          manufacturerData: manufacturerData,
          serviceData: serviceData,
        ),
      );
    }

    return advertisements;
  }

  /// [toPacket] converts the [BleAdvertisement] to a raw message following this structure:
  /// MAC_ADDRESS;UNIX;LAT;LNG;ALT;MODEL;RSSI;MANUFACTURER+DATA;SERVICE+DATA;CRC16
  String toPacket() {
    String payload = '';

    payload += '${macAddress.replaceAll(':', '').toUpperCase()};';
    payload += '${timestamp.millisecondsSinceEpoch ~/ 1000};';
    if (latitude != null) {
      payload += '$latitude;';
    } else {
      payload += ';';
    }
    if (longitude != null) {
      payload += '$longitude;';
    } else {
      payload += ';';
    }
    if (altitude != null) {
      payload += '$altitude;';
    } else {
      payload += ';';
    }
    payload += '$model;';
    if (deviceName != null) {
      payload += '$deviceName;';
    } else {
      payload += ';';
    }
    payload += '$rssi;';
    if (txPower != null) {
      payload += '$txPower;';
    } else {
      payload += ';';
    }
    payload += '${manufacturerData.map((e) => e.toPacket()).join(',')};';
    payload += '${serviceData.map((e) => e.toPacket()).join(',')};';

    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    return '$payload$crc';
  }
}

extension LayrzProtocolManufacturerDataSpec on BleManufacturerData {
  /// [toPacket] converts the [BleManufacturerData] to a raw message following this structure:
  /// MANUFACTURER_ID:DATA
  String toPacket() {
    String message = '${companyId.toRadixString(16).padLeft(4, '0').toUpperCase()}:';
    message += (data ?? []).map((e) => e.toRadixString(16).padLeft(2, '0')).join('').toUpperCase();
    return message;
  }
}

extension LayrzProtocolServiceDataSpec on BleServiceData {
  /// [toPacket] converts the [BleServiceData] to a raw message following this structure:
  /// SERVICE_UUID:DATA
  String toPacket() {
    String message = '${uuid.toRadixString(16).padLeft(4, '0').toUpperCase()}:';
    message += (data ?? []).map((e) => e.toRadixString(16).padLeft(2, '0')).join('').toUpperCase();
    return message;
  }
}
