import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('BleAdvertisement.toPacket()', () {
    final packet = BleAdvertisement(
      deviceName: 'P RHT',
      macAddress: '12:34:56:78:90:AB',
      timestamp: DateTime.fromMillisecondsSinceEpoch(1738276597 * 1000),
      rssi: -50,
      model: 'ELA_PUCK_RHT',
      txPower: -1,
      manufacturerData: [
        BleManufacturerData(companyId: 0x0757, data: [0xf2, 0xc8, 0x0b]),
      ],
      serviceData: [
        BleServiceData(uuid: 0x2a6e, data: [0x40, 0x0a]),
        BleServiceData(uuid: 0x2a6f, data: [0x2b]),
      ],
    );

    expect(packet.macAddress, "12:34:56:78:90:AB");
    expect(packet.timestamp, DateTime.fromMillisecondsSinceEpoch(1738276597 * 1000));
    expect(packet.latitude, null);
    expect(packet.longitude, null);
    expect(packet.altitude, null);
    expect(packet.deviceName, "P RHT");
    expect(packet.rssi, -50);
    expect(packet.txPower, -1);
    expect(packet.model, "ELA_PUCK_RHT");
    expect(packet.manufacturerData.length, 1);
    expect(packet.manufacturerData[0].companyId, 0x0757);
    expect(packet.manufacturerData[0].data, [0xf2, 0xc8, 0x0b]);
    expect(packet.serviceData.length, 2);
    expect(packet.serviceData[0].uuid, 0x2a6e);
    expect(packet.serviceData[0].data, [0x40, 0x0a]);
    expect(packet.serviceData[1].uuid, 0x2a6f);
    expect(packet.serviceData[1].data, [0x2b]);

    final rawPacket = packet.toPacket();
    String rawElaAdvPaket = "1234567890AB;1738276597;;;;ELA_PUCK_RHT;P RHT;-50;-1;0757:F2C80B;2A6E:400A,2A6F:2B;410C";
    expect(rawPacket, rawElaAdvPaket);
  });

  test('PbPacket.parse()', () {
    String rawElaAdvPaket = "1234567890AB;1738276597;;;;ELA_PUCK_RHT;P RHT;-50;-1;0757:F2C80B;2A6E:400A,2A6F:2B;410C";

    final packets = BleAdvertisement.fromPacket(rawElaAdvPaket);

    expect(packets.length, 1);

    final elaAdv = packets[0];
    expect(elaAdv.macAddress, "12:34:56:78:90:AB");
    expect(elaAdv.timestamp, DateTime.fromMillisecondsSinceEpoch(1738276597 * 1000));
    expect(elaAdv.latitude, null);
    expect(elaAdv.longitude, null);
    expect(elaAdv.altitude, null);
    expect(elaAdv.deviceName, "P RHT");
    expect(elaAdv.rssi, -50);
    expect(elaAdv.txPower, -1);
    expect(elaAdv.manufacturerData.length, 1);
    expect(elaAdv.manufacturerData[0].companyId, 0x0757);
    expect(elaAdv.manufacturerData[0].data, [0xf2, 0xc8, 0x0b]);
    expect(elaAdv.serviceData.length, 2);
    expect(elaAdv.serviceData[0].uuid, 0x2a6e);
    expect(elaAdv.serviceData[0].data, [0x40, 0x0a]);
    expect(elaAdv.serviceData[1].uuid, 0x2a6f);
    expect(elaAdv.serviceData[1].data, [0x2b]);

    final packet = PbPacket(advertisements: packets);

    final rawPacket = packet.toPacket();
    expect(rawPacket, "<Pb>$rawElaAdvPaket;7E33</Pb>");
  });
}
