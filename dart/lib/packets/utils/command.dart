part of '../../packets.dart';

class Command {
  /// [commandId] is the Layrz ID to refer to this command.
  /// This ID is unique in the Layrz ecosystem and should be used to send the ACK.
  final String commandId;

  /// [commandName] is the name of the command.
  final String commandName;

  /// [args] is the arguments of the command.
  final Map<String, dynamic> args;

  Command({
    required this.commandId,
    required this.commandName,
    required this.args,
  });

  /// [formatAck] formats the ACK message to send to the Layrz ecosystem.
  String formatAck(String response) {
    String parsedMessage = '';
    DateTime now = DateTime.now();
    parsedMessage += '${now.millisecondsSinceEpoch ~/ 1000};';
    parsedMessage += '$commandId;';
    parsedMessage += '$response;';

    final crc16 = calculateCrc(parsedMessage.codeUnits).toRadixString(16).padLeft(4, '0');
    parsedMessage += crc16;

    return "<Pc>$parsedMessage</Pc>";
  }

  /// [fromPackets] creates a [List<Command>] from a raw message following this structure:
  /// commandId;commandName;args;crc16
  static List<Command> fromPackets(String raw) {
    final parts = raw.split(';');
    if (parts.isEmpty) return [];

    if (parts.length == 1) return [];
    // The parts should be divisible by 4, can be multiple groups

    if (parts.length % 4 != 0) {
      throw ParseException('Invalid command definition');
    }

    // Separate each group
    final List<Command> commands = [];

    for (int i = 0; i < parts.length; i += 4) {
      final commandId = parts[i];
      final commandName = parts[i + 1];
      final rawArgs = parts[i + 2];
      final receivedCrc = int.tryParse(parts[i + 3], radix: 16) ?? 0;
      final calculatedCrc = calculateCrc("$commandId;$commandName;$rawArgs;".codeUnits);
      if (receivedCrc != calculatedCrc) {
        throw CrcException(
          'Invalid CRC, expected ${receivedCrc.toRadixString(16)}, '
          'got ${calculatedCrc.toRadixString(16)}',
        );
      }

      Map<String, dynamic> args = {};

      if (rawArgs.isNotEmpty) {
        // Split each raw args by ','
        final argParts = rawArgs.split(',');
        for (var arg in argParts) {
          final parts = arg.split(':');
          final key = parts[0];
          final value = parts[1];

          final intRegex = RegExp(r'^\d+$');
          final doubleRegex = RegExp(r'^\d+\.\d+$');
          if (intRegex.hasMatch(value)) {
            args[key] = int.parse(value);
          } else if (doubleRegex.hasMatch(value)) {
            args[key] = double.parse(value);
          } else if (value == 'true' || value == 'false') {
            args[key] = value == 'true';
          } else {
            args[key] = value;
          }
        }
      }

      commands.add(Command(commandId: commandId, commandName: commandName, args: args));
    }

    return commands;
  }

  /// [toPacket] converts the [Command] to a raw message following this structure:
  /// commandId;commandName;args;crc16
  String toPacket() {
    String payload = '$commandId;$commandName;';

    payload += args.entries.map((e) {
      if (e.value is bool) {
        return '${e.key}:${e.value ? 'true' : 'false'}';
      } else if (e.value is int) {
        return '${e.key}:${e.value}';
      } else if (e.value is double) {
        return '${e.key}:${e.value}';
      } else {
        return '${e.key}:${e.value}';
      }
    }).join(',');

    payload += ';';

    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    return '$payload$crc';
  }
}
