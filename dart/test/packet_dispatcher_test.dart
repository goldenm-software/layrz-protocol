import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/layrz_protocol.dart';

void main() {
  group('Packet.fromPacket() dispatcher', () {
    test('throws MalformedException for unknown packet', () {
      expect(() => Packet.fromPacket('<Xx>garbage</Xx>'), throwsA(isA<MalformedException>()));
    });
  });

  group('Packet.parseExtraArgs()', () {
    test('parses digital input', () {
      final result = Packet.parseExtraArgs('io1.di:true');
      expect(result['gpio.io1.digital.input'], true);
    });

    test('parses digital output', () {
      final result = Packet.parseExtraArgs('io2.do:false');
      expect(result['gpio.io2.digital.output'], false);
    });

    test('parses analog input', () {
      final result = Packet.parseExtraArgs('io1.ai:3.14');
      expect(result['gpio.io1.analog.input'], 3.14);
    });

    test('parses analog output', () {
      final result = Packet.parseExtraArgs('io3.ao:2.5');
      expect(result['gpio.io3.analog.output'], 2.5);
    });

    test('parses counter', () {
      final result = Packet.parseExtraArgs('io1.counter:42');
      expect(result['gpio.io1.event.count'], 42);
    });

    test('parses BLE id', () {
      final result = Packet.parseExtraArgs('ble.0.id:AA:BB:CC:DD:EE:FF');
      expect(result['ble.0.mac.address'], 'AA:BB:CC:DD:EE:FF');
    });

    test('parses BLE humidity', () {
      final result = Packet.parseExtraArgs('ble.0.hum:65.5');
      expect(result['ble.0.humidity'], 65.5);
    });

    test('parses BLE temperature celsius', () {
      final result = Packet.parseExtraArgs('ble.0.tempc:22.5');
      expect(result['ble.0.temperature.celsius'], 22.5);
    });

    test('parses BLE temperature fahrenheit', () {
      final result = Packet.parseExtraArgs('ble.0.tempf:72.5');
      expect(result['ble.0.temperature.fahrenheit'], 72.5);
    });

    test('parses BLE model id', () {
      final result = Packet.parseExtraArgs('ble.0.model_id:ELA_PUCK');
      expect(result['ble.0.model.id'], 'ELA_PUCK');
    });

    test('parses BLE battery level', () {
      final result = Packet.parseExtraArgs('ble.0.batt:80');
      expect(result['ble.0.battery.level'], 80);
    });

    test('parses BLE lux level', () {
      final result = Packet.parseExtraArgs('ble.0.lux:500');
      expect(result['ble.0.light.level.lux'], 500);
    });

    test('parses BLE voltage', () {
      final result = Packet.parseExtraArgs('ble.0.volt:3.7');
      expect(result['ble.0.voltage'], 3.7);
    });

    test('parses BLE rpm', () {
      final result = Packet.parseExtraArgs('ble.0.rpm:1200');
      expect(result['ble.0.rpm'], 1200);
    });

    test('parses BLE pressure', () {
      final result = Packet.parseExtraArgs('ble.0.press:101.3');
      expect(result['ble.0.pressure'], 101.3);
    });

    test('parses BLE event count', () {
      final result = Packet.parseExtraArgs('ble.0.counter:5');
      expect(result['ble.0.event.count'], 5);
    });

    test('parses BLE x acceleration', () {
      final result = Packet.parseExtraArgs('ble.0.x_acc:0.1');
      expect(result['ble.0.acceleration.x'], 0.1);
    });

    test('parses BLE y acceleration', () {
      final result = Packet.parseExtraArgs('ble.0.y_acc:0.2');
      expect(result['ble.0.acceleration.y'], 0.2);
    });

    test('parses BLE z acceleration', () {
      final result = Packet.parseExtraArgs('ble.0.z_acc:9.8');
      expect(result['ble.0.acceleration.z'], 9.8);
    });

    test('parses BLE message count', () {
      final result = Packet.parseExtraArgs('ble.0.msg_count:3');
      expect(result['ble.0.message.count'], 3);
    });

    test('parses BLE message', () {
      final result = Packet.parseExtraArgs('ble.0.msg:hello');
      expect(result['ble.0.message'], 'hello');
    });

    test('parses BLE magnetic event count', () {
      final result = Packet.parseExtraArgs('ble.0.mag_counter:2');
      expect(result['ble.0.magnetic.event.count'], 2);
    });

    test('parses BLE magnetic data', () {
      final result = Packet.parseExtraArgs('ble.0.mag_data:somedata');
      expect(result['ble.0.magnetic.data'], 'somedata');
    });

    test('parses BLE rssi', () {
      final result = Packet.parseExtraArgs('ble.0.rssi:-70');
      expect(result['ble.0.rssi.dbm'], -70);
    });

    test('parses report code', () {
      final result = Packet.parseExtraArgs('report:1');
      expect(result['report.code'], 1);
    });

    test('parses confiot_ble', () {
      final result = Packet.parseExtraArgs('confiot_ble:1');
      expect(result['ble.confiot.connection.status'], 1);
    });

    test('parses confiot_serial', () {
      final result = Packet.parseExtraArgs('confiot_serial:0');
      expect(result['serial.confiot.connection.status'], 0);
    });

    test('parses boolean true shorthand t', () {
      final result = Packet.parseExtraArgs('some.key:t');
      expect(result['some.key'], true);
    });

    test('parses boolean false shorthand f', () {
      final result = Packet.parseExtraArgs('some.key:f');
      expect(result['some.key'], false);
    });

    test('parses integer value', () {
      final result = Packet.parseExtraArgs('some.key:42');
      expect(result['some.key'], 42);
    });

    test('parses double value', () {
      final result = Packet.parseExtraArgs('some.key:3.14');
      expect(result['some.key'], 3.14);
    });

    test('parses string value', () {
      final result = Packet.parseExtraArgs('some.key:hello');
      expect(result['some.key'], 'hello');
    });

    test('parses multiple comma-separated args', () {
      final result = Packet.parseExtraArgs('io1.di:true,some.val:99');
      expect(result['gpio.io1.digital.input'], true);
      expect(result['some.val'], 99);
    });

    test('handles value with colon (joins remaining parts)', () {
      final result = Packet.parseExtraArgs('ble.0.id:AA:BB:CC');
      expect(result['ble.0.mac.address'], 'AA:BB:CC');
    });
  });
}
