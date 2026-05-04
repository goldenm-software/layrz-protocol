import 'dart:io';

import 'package:layrz_protocol/servers/tcp.dart';
import 'package:layrz_protocol/packets/packets.dart';

Future<void> main() async {
  final server = TcpServer(
    TcpConfig(
      port: 12345,
      onNewPacket: (packet, conn) {
        if (packet is PaPacket) {
          print('Pa received from ${conn.remoteAddress.address}');
          return AsPacket();
        }
        if (packet is PbPacket) print('Pb received');
        if (packet is PcPacket) print('Pc received');
        if (packet is PdPacket) print('Pd received');
        if (packet is PiPacket) print('Pi received');
        if (packet is PmPacket) print('Pm received');
        if (packet is PrPacket) print('Pr received');
        if (packet is PsPacket) print('Ps received');
        return AoPacket(timestamp: DateTime.now());
      },
      onDecodeError: (err, data, conn) {
        print('Decode error: $err');
      },
    ),
  );

  await server.start();
  print('Listening on :12345');

  ProcessSignal.sigint.watch().listen((_) async {
    print('\nShutting down...');
    await server.close();
    exit(0);
  });

  await server.serve();
}
