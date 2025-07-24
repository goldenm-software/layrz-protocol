import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('PdPacket.parse()', () {
    String payload = '0;'; // timestamp
    payload += '10.0;'; // LAT
    payload += '10.0;'; // LNG
    payload += '10.0;'; // ALT
    payload += '10.0;'; // SPD
    payload += '10.0;'; // DIR
    payload += '5;'; // SAT
    payload += '1.0;'; // HDOP

    Map<String, dynamic> extra = {
      'test.str': 'Hola mundo',
      'test.int': 1,
      'test.double': 1.0,
      'test.bool': true,
      'mac.address': '00:00:00:00:00:00',
    };

    List<String> extraList = [];
    for (String key in extra.keys) {
      extraList.add('$key:${extra[key]}');
    }

    payload += '${extraList.join(',')};'; // End of extras

    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    payload = '<Pd>$payload$crc</Pd>';

    PdPacket link = PdPacket.fromPacket(payload);

    expect(link.timestamp, DateTime.fromMillisecondsSinceEpoch(0));
    expect(link.position.latitude, 10.0);
    expect(link.position.longitude, 10.0);
    expect(link.position.altitude, 10.0);
    expect(link.position.speed, 10.0);
    expect(link.position.direction, 10.0);
    expect(link.position.satellites, 5);
    expect(link.position.hdop, 1.0);
    expect(link.extra, extra);

    String reversedPayload = link.toPacket();
    expect(reversedPayload, payload);
  });
}
