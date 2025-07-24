import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('PiPacket.parse()', () {
    String ident = 'testident';
    String firmwareId = '1';
    int firmwareBuild = 1;
    String deviceId = '1';
    String hardwareId = '1';
    String modelId = '1';
    FirmwareBranch firmwareBranch = FirmwareBranch.development;
    bool fotaEnabled = true;

    String payload = '$ident;';
    payload += '$firmwareId;';
    payload += '$firmwareBuild;';
    payload += '$deviceId;';
    payload += '$hardwareId;';
    payload += '$modelId;';
    payload += '${firmwareBranch.toJson()};';
    // payload += '${fotaEnabled ? '1' : '0'};';
    payload += '1;';

    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    payload = '<Pi>$payload$crc</Pi>';

    PiPacket link = PiPacket.fromPacket(payload);

    expect(link.ident, ident);
    expect(link.firmwareId, firmwareId);
    expect(link.firmwareBuild, firmwareBuild);
    expect(link.deviceId, deviceId);
    expect(link.hardwareId, hardwareId);
    expect(link.modelId, modelId);
    expect(link.firmwareBranch, firmwareBranch);
    expect(link.fotaEnabled, fotaEnabled);

    String reversedPayload = link.toPacket();
    expect(reversedPayload, payload);
  });
}
