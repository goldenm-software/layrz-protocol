part of '../../packets.dart';

class ImPacket extends AIPacket {
  /// [timestamp] is the time of the packet.
  /// This is identified in the packet as `UNIX`
  final DateTime timestamp;

  /// [chatId] is the unique identifier for the chat.
  /// This is identified in the packet as `CHAT_ID`
  final String chatId;

  /// [message] is the chat message content.
  /// This is identified in the packet as `MESSAGE`
  ///
  /// Note, this message replaces any semicolons with the string `|||` to avoid packet parsing issues.
  final String message;

  /// [ImPacket] is the chat message packet
  ///
  /// This packet is part of the packet sent between Layrz services to stream AI conversations
  ImPacket({
    required this.timestamp,
    required this.chatId,
    required this.message,
  });

  /// [fromPacket] creates a [ImPacket] from a string packet in the format of `Layrz Protocol v3`.
  static ImPacket fromPacket(String raw) {
    if (!raw.startsWith('<Im>') || !raw.endsWith('</Im>')) {
      throw ParseException('Invalid identification packet, should be <Im>...</Im>');
    }

    final parts = raw.substring(4, raw.length - 5).split(';');
    if (parts.length != 4) {
      throw MalformedException('Invalid packet parts, should have 4 parts');
    }

    int? receivedCrc = int.tryParse(parts.last, radix: 16);
    int? calculatedCrc = calculateCrc("${parts.sublist(0, 3).join(';')};".codeUnits);

    if (receivedCrc != calculatedCrc) {
      throw CrcException('Invalid CRC, received: $receivedCrc, calculated: $calculatedCrc');
    }

    DateTime timestamp;
    try {
      timestamp = DateTime.fromMillisecondsSinceEpoch(int.parse(parts[0]) * 1000);
    } catch (e) {
      throw MalformedException('Invalid timestamp');
    }

    return ImPacket(
      timestamp: timestamp,
      chatId: parts[1],
      message: parts[2].replaceAll('|||', ';'),
    );
  }

  /// [toPacket] returns the packet in the format of `Layrz Protocol v3`.
  @override
  String toPacket() {
    String payload = '${(timestamp.millisecondsSinceEpoch / 1000).round()};';
    payload += '$chatId;';
    payload += '${message.replaceAll(';', '|||')};';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();

    return '<Im>$payload$crc</Im>';
  }

  @override
  String toString() => toPacket();

  @override
  ImPacket copyWith({
    DateTime? timestamp,
    String? chatId,
    String? message,
  }) {
    return ImPacket(
      timestamp: timestamp ?? this.timestamp,
      chatId: chatId ?? this.chatId,
      message: message ?? this.message,
    );
  }
}
