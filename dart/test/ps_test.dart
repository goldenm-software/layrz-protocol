import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('PsPacket.parse()', () {
    String payload = '0;'; // timestamp
    Map<String, dynamic> extra = {
      'net_wifi_ssid': 'AWESOME WIFI',
      'net_wifi_pass': 'dictadormarico69', // https://www.youtube.com/watch?v=kq0VUZXiUQs
      'net_wifi_sec': 'WPA2',
      'mac.address': '00:00:00:00:00:00',
      'static.lat': -15.0,
      'static.lon': 15.0,
      'static.alt': 0.0,
    };

    List<String> extraList = [];
    for (String key in extra.keys) {
      extraList.add('$key:${extra[key]}');
    }

    payload += '${extraList.join(',')};'; // End of extras
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    payload = '<Ps>$payload$crc</Ps>';

    PsPacket link = PsPacket.fromPacket(payload);

    expect(link.timestamp, DateTime.fromMillisecondsSinceEpoch(0));
    expect(link.params, extra);

    String reversedPayload = link.toPacket();
    expect(reversedPayload, payload);

    expect(link.params['net_wifi_ssid'], extra['net_wifi_ssid']);
    expect(link.params['net_wifi_pass'], extra['net_wifi_pass']);
    expect(link.params['net_wifi_sec'], extra['net_wifi_sec']);
    expect(link.params['mac.address'], extra['mac.address']);
    expect(link.params['static.lat'], extra['static.lat']);
    expect(link.params['static.lon'], extra['static.lon']);
  });
}
