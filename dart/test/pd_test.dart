import 'package:test/test.dart';
import 'package:layrz_protocol/packets/packets.dart';
import 'package:layrz_protocol/utils/crc.dart';

void main() {
  test('Packet.fromPacket() routes PdPacket', () {
    final ts = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    final original = PdPacket(
      timestamp: ts,
      position: Position(latitude: 10.0, longitude: 20.0),
      extra: {'test.key': 1},
    );
    final parsed = Packet.fromPacket(original.toPacket());
    expect(parsed, isA<PdPacket>());
  });

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
      // Colons in key and value are escaped as `___` on the wire (see PdPacket.toPacket).
      final escapedKey = key.replaceAll(':', '___');
      final escapedValue = '${extra[key]}'.replaceAll(':', '___');
      extraList.add('$escapedKey:$escapedValue');
    }

    payload += '${extraList.join(',')};'; // End of extras

    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    payload = '<Pd>$payload$crc</Pd>';

    PdPacket link = PdPacket.fromPacket(payload);

    expect(link.timestamp, DateTime.fromMillisecondsSinceEpoch(0, isUtc: true));
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
