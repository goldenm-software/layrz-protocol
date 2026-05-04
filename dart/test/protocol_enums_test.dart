import 'package:layrz_protocol/utils/protocol.dart';
import 'package:test/test.dart';

void main() {
  group('LayrzProtocolMode', () {
    test('tcp.value returns TCP', () {
      expect(LayrzProtocolMode.tcp.value, 'TCP');
    });

    test('http.value returns HTTP', () {
      expect(LayrzProtocolMode.http.value, 'HTTP');
    });

    test('fromString TCP', () {
      expect(LayrzProtocolMode.fromString('TCP'), LayrzProtocolMode.tcp);
    });

    test('fromString HTTP', () {
      expect(LayrzProtocolMode.fromString('HTTP'), LayrzProtocolMode.http);
    });

    test('fromString unknown defaults to http', () {
      expect(LayrzProtocolMode.fromString('UNKNOWN'), LayrzProtocolMode.http);
    });
  });

  group('LayrzProtocolVersion', () {
    test('v2.value returns v2', () {
      expect(LayrzProtocolVersion.v2.value, 'v2');
    });

    test('fromString v2', () {
      expect(LayrzProtocolVersion.fromString('v2'), LayrzProtocolVersion.v2);
    });

    test('fromString unknown defaults to v2', () {
      expect(LayrzProtocolVersion.fromString('unknown'), LayrzProtocolVersion.v2);
    });
  });
}
