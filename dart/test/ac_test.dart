import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  test('AcPacket.parse()', () {
    Command cmd1 = Command(commandId: '1', commandName: 'get_msg', args: {});

    String payload1 = '${cmd1.commandId};${cmd1.commandName};;';
    String crc = calculateCrc(payload1.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    payload1 = '$payload1$crc';
    expect(payload1, cmd1.toPacket());

    Command cmd2 = Command(commandId: '2', commandName: 'set_config', args: {'wifi_ssid': 'test'});
    String payload2 = '${cmd2.commandId};${cmd2.commandName};wifi_ssid:test;';
    crc = calculateCrc(payload2.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    payload2 = '$payload2$crc';
    expect(payload2, cmd2.toPacket());

    String payload = '$payload1;$payload2;';
    crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    payload = '<Ac>$payload$crc</Ac>';

    AcPacket link = AcPacket.fromPacket(payload);

    expect(link.commands.length, 2);
    expect(link.commands[0].commandId, '1');
    expect(link.commands[0].commandName, 'get_msg');
    expect(link.commands[0].args, {});

    expect(link.commands[1].commandId, '2');
    expect(link.commands[1].commandName, 'set_config');
    expect(link.commands[1].args, {'wifi_ssid': 'test'});

    String reversedPayload = link.toPacket();
    expect(reversedPayload, payload);

    print(
      AcPacket(
        commands: [Command(commandId: '1', commandName: 'ping', args: {})],
      ).toPacket(),
    );
  });
}
