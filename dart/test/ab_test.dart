import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('AbPacket.parse()', () {
    String payload = "<Ab>1234567890AB:GENERIC;BC0987654321:GENERIC;C1BE</Ab>";

    AbPacket link = AbPacket.fromPacket(payload);
    expect(link.devices.length, 2);
    expect(link.devices[0].macAddress, "12:34:56:78:90:AB");
    expect(link.devices[0].model, "GENERIC");

    expect(link.devices[1].macAddress, "BC:09:87:65:43:21");
    expect(link.devices[1].model, "GENERIC");

    expect(link.toPacket(), payload);
  });
}
