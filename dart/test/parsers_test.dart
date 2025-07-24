import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('Packet.convertToDotCase()', () {
    Map<String, dynamic> data = {
      'dict': {
        'key1': 'value1',
        'key2': 'value2',
        'key3': 'value3',
      },
      'list': [
        'value1',
        'value2',
        'value3',
      ],
      'string': 'value',
      'int': 1,
      'float': 1.1,
      'bool': true,
    };

    Map<String, dynamic> converted = {};
    for (String key in data.keys) {
      converted.addAll(Packet.convertToDotCase(key, data[key]));
    }

    expect(converted['dict.key1'], data['dict']['key1']);
    expect(converted['dict.key2'], data['dict']['key2']);
    expect(converted['dict.key3'], data['dict']['key3']);
    expect(converted['list.0'], data['list'][0]);
    expect(converted['list.1'], data['list'][1]);
    expect(converted['list.2'], data['list'][2]);
    expect(converted['string'], data['string']);
    expect(converted['int'], data['int']);
    expect(converted['float'], data['float']);
    expect(converted['bool'], true);
  });
}
