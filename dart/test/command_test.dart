import 'package:layrz_protocol/utils/errors.dart';
import 'package:test/test.dart';
import 'package:layrz_protocol/packets/packets.dart';

void main() {
  test('Command.toPacket() with no args', () {
    final cmd = Command(commandId: '1', commandName: 'ping', args: {});
    final raw = cmd.toPacket();
    expect(raw.startsWith('1;ping;;'), true);
  });

  test('Command.toPacket() with int arg', () {
    final cmd = Command(commandId: '2', commandName: 'set', args: {'port': 8080});
    final raw = cmd.toPacket();
    expect(raw.contains('port:8080'), true);
  });

  test('Command.toPacket() with double arg', () {
    final cmd = Command(commandId: '3', commandName: 'set', args: {'threshold': 3.14});
    final raw = cmd.toPacket();
    expect(raw.contains('threshold:3.14'), true);
  });

  test('Command.toPacket() with bool arg', () {
    final cmd = Command(commandId: '4', commandName: 'set', args: {'enabled': true});
    final raw = cmd.toPacket();
    expect(raw.contains('enabled:true'), true);
  });

  test('Command.toPacket() with string arg', () {
    final cmd = Command(commandId: '5', commandName: 'set', args: {'ssid': 'mynet'});
    final raw = cmd.toPacket();
    expect(raw.contains('ssid:mynet'), true);
  });

  test('Command.formatAck() returns valid PcPacket string', () {
    final cmd = Command(commandId: '1', commandName: 'ping', args: {});
    final ack = cmd.formatAck('ok');
    expect(ack.startsWith('<Pc>'), true);
    expect(ack.endsWith('</Pc>'), true);
    expect(ack.contains(';1;ok;'), true);
  });

  test('Command.fromPackets() throws ParseException on invalid part count', () {
    expect(() => Command.fromPackets('one;two;three'), throwsA(isA<ParseException>()));
  });

  test('Command.fromPackets() with int, double, bool, string args', () {
    final cmd = Command(
      commandId: '1',
      commandName: 'config',
      args: {'port': 8080, 'ratio': 0.5, 'enabled': false, 'ssid': 'net'},
    );
    final raw = cmd.toPacket();
    final cmds = Command.fromPackets(raw);
    expect(cmds.length, 1);
    expect(cmds[0].args['port'], 8080);
    expect(cmds[0].args['ratio'], 0.5);
    expect(cmds[0].args['enabled'], false);
    expect(cmds[0].args['ssid'], 'net');
  });

  test('Command.fromPackets() throws CrcException on CRC mismatch', () {
    final cmd = Command(commandId: '1', commandName: 'ping', args: {});
    String raw = cmd.toPacket();
    raw = raw.replaceRange(raw.length - 4, raw.length, 'FFFF');
    expect(() => Command.fromPackets(raw), throwsA(isA<CrcException>()));
  });

  test('Command.fromPackets() returns empty list for empty string', () {
    expect(Command.fromPackets(''), isEmpty);
  });

  test('Command.fromPackets() returns empty list for single part', () {
    expect(Command.fromPackets('single'), isEmpty);
  });
}
