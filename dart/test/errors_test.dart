import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('ParseException.toString()', () {
    final except = ParseException('test');
    expect(except.toString(), 'ParseException: test');
    expect(except, isA<Exception>());
  });
  test('CrcException.toString()', () {
    final except = CrcException('test');
    expect(except.toString(), 'CrcException: test');
    expect(except, isA<Exception>());
  });
  test('CommandException.toString()', () {
    final except = CommandException('test');
    expect(except.toString(), 'CommandException: test');
    expect(except, isA<Exception>());
  });
  test('ServerException.toString()', () {
    final except = ServerException('test');
    expect(except.toString(), 'ServerException: test');
    expect(except, isA<Exception>());
  });
  test('MalformedException.toString()', () {
    final except = MalformedException('test');
    expect(except.toString(), 'MalformedException: test');
    expect(except, isA<Exception>());
  });
}
