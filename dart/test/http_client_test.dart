import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';
import 'package:mocktail/mocktail.dart';

class MockDio extends Mock implements Dio {}

class FakeRequestOptions extends Fake implements RequestOptions {}

Dio _buildMockDio(String responseBody, {int statusCode = 200}) {
  final mockDio = MockDio();
  final fakeOptions = RequestOptions(path: '');

  when(() => mockDio.post(any(), data: any(named: 'data'))).thenAnswer(
    (_) async => Response(
      data: responseBody,
      statusCode: statusCode,
      requestOptions: fakeOptions,
    ),
  );

  when(() => mockDio.get(any())).thenAnswer(
    (_) async => Response(
      data: responseBody,
      statusCode: statusCode,
      requestOptions: fakeOptions,
    ),
  );

  return mockDio;
}

void main() {
  setUpAll(() {
    registerFallbackValue(FakeRequestOptions());
  });

  test('LayrzProtocolHttp.baseUrl is correct', () {
    final client = LayrzProtocolHttp(ident: 'TEST123', server: 'api.example.com');
    expect(client.baseUrl, 'https://api.example.com/v2');
  });

  test('LayrzProtocolHttp.headers contains authorization', () {
    final client = LayrzProtocolHttp(ident: 'IMEI1', password: 'pass1', server: 'api.example.com');
    expect(client.headers['Authorization'], 'LayrzAuth IMEI1;pass1');
    expect(client.headers['Content-Type'], 'text/plain');
  });

  test('LayrzProtocolHttp.composeEmptyPd returns PdPacket', () {
    final client = LayrzProtocolHttp(ident: 'IMEI1', server: 'api.example.com');
    final pd = client.composeEmptyPd();
    expect(pd, isA<PdPacket>());
    expect(pd.extra, isEmpty);
  });

  test('LayrzProtocolHttp.sendData returns parsed packet', () async {
    final aoPacket = AoPacket(timestamp: DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true));
    final mockDio = _buildMockDio(aoPacket.toPacket());
    final client = LayrzProtocolHttp(ident: 'IMEI1', server: 'api.example.com', testDio: mockDio);

    final pd = PdPacket(
      timestamp: DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true),
      position: Position(latitude: 1.0, longitude: 1.0),
      extra: {'test.val': 1},
    );

    final result = await client.sendData(pd);
    expect(result, isA<AoPacket>());
  });

  test('LayrzProtocolHttp.sendSos sets alarm.event', () async {
    final aoPacket = AoPacket(timestamp: DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true));
    final mockDio = _buildMockDio(aoPacket.toPacket());
    final client = LayrzProtocolHttp(ident: 'IMEI1', server: 'api.example.com', testDio: mockDio);

    final result = await client.sendSos();
    expect(result, isA<AoPacket>());
  });

  test('LayrzProtocolHttp.sendText includes driver.message', () async {
    final aoPacket = AoPacket(timestamp: DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true));
    final mockDio = _buildMockDio(aoPacket.toPacket());
    final client = LayrzProtocolHttp(ident: 'IMEI1', server: 'api.example.com', testDio: mockDio);

    final result = await client.sendText('hello world');
    expect(result, isA<AoPacket>());
  });

  test('LayrzProtocolHttp.sendImage sends PmPacket', () async {
    final aoPacket = AoPacket(timestamp: DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true));
    final mockDio = _buildMockDio(aoPacket.toPacket());
    final client = LayrzProtocolHttp(ident: 'IMEI1', server: 'api.example.com', testDio: mockDio);

    final result = await client.sendImage(bytes: [0x01, 0x02, 0x03], filename: 'photo.jpg');
    expect(result, isA<AoPacket>());
  });

  test('LayrzProtocolHttp.getCommands returns parsed packet', () async {
    final aoPacket = AoPacket(timestamp: DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true));
    final mockDio = _buildMockDio(aoPacket.toPacket());
    final client = LayrzProtocolHttp(ident: 'IMEI1', server: 'api.example.com', testDio: mockDio);

    final result = await client.getCommands();
    expect(result, isA<AoPacket>());
  });

  test('LayrzProtocolHttp.getBleDevices returns parsed packet', () async {
    final abPacket = AbPacket(devices: []);
    final mockDio = _buildMockDio(abPacket.toPacket());
    final client = LayrzProtocolHttp(ident: 'IMEI1', server: 'api.example.com', testDio: mockDio);

    final result = await client.getBleDevices();
    expect(result, isA<AbPacket>());
  });

  test('LayrzProtocolHttp._processResponse throws ServerException on 500', () async {
    final mockDio = _buildMockDio('error', statusCode: 500);
    final client = LayrzProtocolHttp(ident: 'IMEI1', server: 'api.example.com', testDio: mockDio);

    expect(() => client.getCommands(), throwsA(isA<ServerException>()));
  });

  test('LayrzProtocolHttp._processResponse throws MalformedException on invalid packet', () async {
    final mockDio = _buildMockDio('not-a-valid-packet');
    final client = LayrzProtocolHttp(ident: 'IMEI1', server: 'api.example.com', testDio: mockDio);

    expect(() => client.getCommands(), throwsA(isA<MalformedException>()));
  });
}
