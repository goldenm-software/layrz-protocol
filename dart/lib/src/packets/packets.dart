library;

import 'dart:convert';
import 'dart:typed_data';

import 'package:layrz_logging/layrz_logging.dart';
import 'package:layrz_protocol/layrz_protocol.dart';
import 'package:layrz_protocol/src/utils/constants.dart';

part 'src/client/pi.dart';
part 'src/client/pd.dart';
part 'src/client/pc.dart';
part 'src/client/ps.dart';
part 'src/client/pb.dart';
part 'src/client/pa.dart';
part 'src/client/pr.dart';
part 'src/client/pm.dart';
part 'src/client/pt.dart';
part 'src/client/pe.dart';

part 'src/server/ao.dart';
part 'src/server/ac.dart';
part 'src/server/ar.dart';
part 'src/server/ab.dart';
part 'src/server/as.dart';
part 'src/server/au.dart';

part 'src/utils/command.dart';
part 'src/utils/position.dart';
part 'src/utils/ble_advertisement.dart';

class Packet {
  static Packet fromPacket(String raw) {
    // Client packets
    if (raw.startsWith('<Pc>') && raw.endsWith('</Pc>')) return PcPacket.fromPacket(raw);
    if (raw.startsWith('<Pi>') && raw.endsWith('</Pi>')) return PiPacket.fromPacket(raw);
    if (raw.startsWith('<Pd>') && raw.endsWith('</Pd>')) return PdPacket.fromPacket(raw);
    if (raw.startsWith('<Ps>') && raw.endsWith('</Ps>')) return PsPacket.fromPacket(raw);
    if (raw.startsWith('<Pb>') && raw.endsWith('</Pb>')) return PbPacket.fromPacket(raw);
    if (raw.startsWith('<Pa>') && raw.endsWith('</Pa>')) return PaPacket.fromPacket(raw);
    if (raw.startsWith('<Pr>') && raw.endsWith('</Pr>')) return PrPacket.fromPacket(raw);
    if (raw.startsWith('<Pm>') && raw.endsWith('</Pm>')) return PmPacket.fromPacket(raw);

    // Server packets
    if (raw.startsWith('<Ao>') && raw.endsWith('</Ao>')) return AoPacket.fromPacket(raw);
    if (raw.startsWith('<Ac>') && raw.endsWith('</Ac>')) return AcPacket.fromPacket(raw);
    if (raw.startsWith('<Ar>') && raw.endsWith('</Ar>')) return ArPacket.fromPacket(raw);
    if (raw.startsWith('<Ab>') && raw.endsWith('</Ab>')) return AbPacket.fromPacket(raw);
    if (raw.startsWith('<As>') && raw.endsWith('</As>')) return AsPacket.fromPacket(raw);
    if (raw.startsWith('<Au>') && raw.endsWith('</Au>')) return AuPacket.fromPacket(raw);

    LayrzLogging.critical('Invalid packet: $raw');
    throw MalformedException('Invalid packet type');
  }

  String toPacket() => '';

  static Map<String, dynamic> parseExtraArgs(String raw) {
    final extra = <String, dynamic>{};

    final extraParts = raw.split(',');
    for (String extraPart in extraParts) {
      List<String> extraPartParts = extraPart.split(':');

      if (extraPartParts.length > 2) {
        // Join the extra args after the first colon
        extraPartParts = [
          extraPartParts[0],
          extraPartParts.sublist(1).join(':'),
        ];
      }

      String key;

      final RegExp digitalInput = RegExp(r'^io[0-9]+\.di$');
      final RegExp digitalOutput = RegExp(r'^io[0-9]+\.do$');
      final RegExp analogInput = RegExp(r'^io[0-9]+\.ai$');
      final RegExp analogOutput = RegExp(r'^io[0-9]+\.ao$');
      final RegExp counter = RegExp(r'^io[0-9]+\.counter$');
      final RegExp bleId = RegExp(r'^ble.[0-9]+\.id$');
      final RegExp bleHumidity = RegExp(r'^ble.[0-9]+\.hum$');
      final RegExp bleTempC = RegExp(r'^ble.[0-9]+\.tempc$');
      final RegExp bleTempF = RegExp(r'^ble.[0-9]+\.tempf$');
      final RegExp bleModelId = RegExp(r'^ble.[0-9]+\.model_id$');
      final RegExp bleBatteryLevel = RegExp(r'^ble.[0-9]+\.batt$');
      final RegExp bleLuxLevel = RegExp(r'^ble.[0-9]+\.lux$');
      final RegExp bleVoltageLevel = RegExp(r'^ble.[0-9]+\.volt$');
      final RegExp bleRpm = RegExp(r'^ble.[0-9]+\.rpm$');
      final RegExp blePressure = RegExp(r'^ble.[0-9]+\.press$');
      final RegExp bleEventCount = RegExp(r'^ble.[0-9]+\.counter$');
      final RegExp bleXAccel = RegExp(r'^ble.[0-9]+\.x_acc$');
      final RegExp bleYAccel = RegExp(r'^ble.[0-9]+\.y_acc$');
      final RegExp bleZAccel = RegExp(r'^ble.[0-9]+\.z_acc$');
      final RegExp bleMsgCount = RegExp(r'^ble.[0-9]+\.msg_count');
      final RegExp bleMsg = RegExp(r'^ble.[0-9]+\.msg');
      final RegExp bleMagCount = RegExp(r'^ble.[0-9]+\.mag_counter');
      final RegExp bleMagData = RegExp(r'^ble.[0-9]+\.mag_data');
      final RegExp bleRssi = RegExp(r'^ble.[0-9]+\.rssi');

      if (digitalInput.hasMatch(extraPartParts[0])) {
        final gpio = extraPartParts[0].replaceAll('io.', '').replaceAll('.di', '');
        key = 'gpio.$gpio.digital.input';
      } else if (digitalOutput.hasMatch(extraPartParts[0])) {
        final gpio = extraPartParts[0].replaceAll('io.', '').replaceAll('.do', '');
        key = 'gpio.$gpio.digital.output';
      } else if (analogInput.hasMatch(extraPartParts[0])) {
        final gpio = extraPartParts[0].replaceAll('io.', '').replaceAll('.ai', '');
        key = 'gpio.$gpio.analog.input';
      } else if (analogOutput.hasMatch(extraPartParts[0])) {
        final gpio = extraPartParts[0].replaceAll('io.', '').replaceAll('.ao', '');
        key = 'gpio.$gpio.analog.output';
      } else if (counter.hasMatch(extraPartParts[0])) {
        final gpio = extraPartParts[0].replaceAll('io.', '').replaceAll('.counter', '');
        key = 'gpio.$gpio.event.count';
      } else if (bleId.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.id', '');
        key = 'ble.$ble.mac.address';
      } else if (bleHumidity.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.hum', '');
        key = 'ble.$ble.humidity';
      } else if (bleTempC.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.tempc', '');
        key = 'ble.$ble.temperature.celsius';
      } else if (bleTempF.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.tempf', '');
        key = 'ble.$ble.temperature.fahrenheit';
      } else if (bleModelId.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.model_id', '');
        key = 'ble.$ble.model.id';
      } else if (bleBatteryLevel.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.batt', '');
        key = 'ble.$ble.battery.level';
      } else if (bleLuxLevel.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.lux', '');
        key = 'ble.$ble.light.level.lux';
      } else if (bleVoltageLevel.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.volt', '');
        key = 'ble.$ble.voltage';
      } else if (bleRpm.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.rpm', '');
        key = 'ble.$ble.rpm';
      } else if (blePressure.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.press', '');
        key = 'ble.$ble.pressure';
      } else if (bleEventCount.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.counter', '');
        key = 'ble.$ble.event.count';
      } else if (bleXAccel.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.x_acc', '');
        key = 'ble.$ble.acceleration.x';
      } else if (bleYAccel.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.y_acc', '');
        key = 'ble.$ble.acceleration.y';
      } else if (bleZAccel.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.z_acc', '');
        key = 'ble.$ble.acceleration.z';
      } else if (bleMsgCount.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.msg_count', '');
        key = 'ble.$ble.message.count';
      } else if (bleMsg.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.msg', '');
        key = 'ble.$ble.message';
      } else if (bleMagCount.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.mag_counter', '');
        key = 'ble.$ble.magnetic.event.count';
      } else if (bleMagData.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.mag_data', '');
        key = 'ble.$ble.magnetic.data';
      } else if (bleRssi.hasMatch(extraPartParts[0])) {
        final ble = extraPartParts[0].replaceAll('ble.', '').replaceAll('.rssi', '');
        key = 'ble.$ble.rssi.dbm';
      } else if (extraPartParts[0] == 'report') {
        key = 'report.code';
      } else if (extraPartParts[0] == 'confiot_ble') {
        key = 'ble.confiot.connection.status';
      } else if (extraPartParts[0] == 'confiot_serial') {
        key = 'serial.confiot.connection.status';
      } else {
        key = extraPartParts[0];
      }

      if (['true', 't'].contains(extraPartParts[1].toString().toLowerCase())) {
        extra[key] = true;
      } else if (['false', 'f'].contains(extraPartParts[1].toString().toLowerCase())) {
        extra[key] = false;
      } else if (RegExp(r'^-?\d+\.\d+$').hasMatch(extraPartParts[1])) {
        extra[key] = double.tryParse(extraPartParts[1]);
      } else if (RegExp(r'^-?\d+$').hasMatch(extraPartParts[1])) {
        extra[key] = int.tryParse(extraPartParts[1]);
      } else {
        extra[key] = extraPartParts[1];
      }
    }

    return extra;
  }

  static Map<String, dynamic> convertToDotCase(String key, dynamic value) {
    if (value == null) return {};

    try {
      final data = List.from(value);

      Map<String, dynamic> result = {};

      for (int i = 0; i < data.length; i++) {
        result.addAll(convertToDotCase('$key.$i', data[i]));
      }

      return result;
    } catch (e) {
      // Empty list
    }

    try {
      final data = Map<String, dynamic>.from(value);

      Map<String, dynamic> result = {};

      for (String k in data.keys) {
        result.addAll(convertToDotCase('$key.$k', data[k]));
      }

      return result;
    } catch (e) {
      // Empty map
    }

    if (value is String) {
      for (String ascii in asciiMap.keys) {
        value = value.replaceAll(ascii, asciiMap[ascii]);
      }
    }

    return {key: value};
  }

  @override
  String toString() {
    return toPacket();
  }

  Packet copyWith() {
    return Packet();
  }
}

class ServerPacket extends Packet {}

class ClientPacket extends Packet {}
