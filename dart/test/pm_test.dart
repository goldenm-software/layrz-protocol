import 'dart:convert';
import 'dart:typed_data';
import 'package:layrz_protocol/utils/errors.dart';
import 'package:test/test.dart';
import 'package:layrz_protocol/packets/packets.dart';
import 'package:layrz_protocol/utils/crc.dart';

void main() {
  test('Packet.fromPacket() routes PmPacket', () {
    final original = PmPacket(
      filename: 'test.bin',
      contentType: 'application/octet-stream',
      data: Uint8List.fromList([0x01, 0x02]),
    );
    final parsed = Packet.fromPacket(original.toPacket());
    expect(parsed, isA<PmPacket>());
  });

  test('PmPacket.toPacket()', () {
    final data = Uint8List.fromList([0x48, 0x65, 0x6C, 0x6C, 0x6F]); // "Hello"
    final packet = PmPacket(filename: 'test.txt', contentType: 'text/plain', data: data);
    final raw = packet.toPacket();

    expect(raw.startsWith('<Pm>'), true);
    expect(raw.endsWith('</Pm>'), true);

    String b64 = base64Encode(data);
    String payload = 'test.txt;text/plain;$b64;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(raw, '<Pm>$payload$crc</Pm>');
  });

  test('PmPacket.fromPacket()', () {
    final data = Uint8List.fromList([0x48, 0x65, 0x6C, 0x6C, 0x6F]);
    String b64 = base64Encode(data);
    String payload = 'test.txt;text/plain;$b64;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    String raw = '<Pm>$payload$crc</Pm>';

    final packet = PmPacket.fromPacket(raw);
    expect(packet.filename, 'test.txt');
    expect(packet.contentType, 'text/plain');
    expect(packet.data, data);
    expect(packet.toPacket(), raw);
  });

  test('PmPacket.fromPacket() invalid format throws ParseException', () {
    expect(() => PmPacket.fromPacket('invalid'), throwsA(isA<ParseException>()));
  });

  test('PmPacket.fromPacket() wrong part count throws MalformedException', () {
    String payload = 'only;two;';
    String crc = calculateCrc(payload.codeUnits).toRadixString(16).padLeft(4, '0').toUpperCase();
    expect(() => PmPacket.fromPacket('<Pm>$payload$crc</Pm>'), throwsA(isA<MalformedException>()));
  });

  test('PmPacket.copyWith()', () {
    final data = Uint8List.fromList([0x01]);
    final original = PmPacket(filename: 'a.bin', contentType: 'application/octet-stream', data: data);
    final copy = original.copyWith(filename: 'b.bin');
    expect(copy.filename, 'b.bin');
    expect(copy.contentType, 'application/octet-stream');
    expect(copy.data, data);
  });
}
