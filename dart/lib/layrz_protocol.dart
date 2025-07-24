import 'dart:io';
import 'dart:async';
import 'dart:convert';
import 'dart:typed_data';
import 'package:dio/dio.dart';
import 'package:drift/drift.dart';
import 'package:layrz_logging/layrz_logging.dart';

import 'src/database/database.dart';
import 'src/utils/errors.dart';
import 'src/packets/packets.dart';

export 'package:layrz_models/layrz_models.dart' show BleManufacturerData, BleServiceData, FirmwareBranch;
export 'src/packets/packets.dart';
export 'src/utils/crc.dart';
export 'src/utils/errors.dart';

part 'src/clients/http.dart';
part 'src/clients/tcp.dart';

typedef HttpPacketResponse = Future<Packet>;

/// [LayrzProtocolMode] defines the protocol to use.
///
/// Please, read each mode to understand the differences between them.
enum LayrzProtocolMode {
  /// [tcp] defines the protocol to use TCP.
  /// Enhances realtime communication but requires a constant connection.
  tcp,

  /// [http] defines the protocol to use HTTP.
  /// Enhances compatibility with different devices but requires a polling mechanism.
  ///
  /// Also, this is the default mode.
  http;

  String get value {
    switch (this) {
      case LayrzProtocolMode.tcp:
        return 'TCP';
      case LayrzProtocolMode.http:
      default:
        return 'HTTP';
    }
  }

  static LayrzProtocolMode fromString(String value) {
    switch (value) {
      case 'TCP':
        return LayrzProtocolMode.tcp;
      case 'HTTP':
      default:
        return LayrzProtocolMode.http;
    }
  }
}

/// [LayrzProtocolVersion] defines the version of the protocol.
enum LayrzProtocolVersion {
  /// [v2] is the version 2 of the protocol. Enhances capabilities like commands, images, and BLE devices.
  ///
  /// Also, this is the default version.
  v2;

  String get value {
    switch (this) {
      case LayrzProtocolVersion.v2:
      default:
        return 'v2';
    }
  }

  static LayrzProtocolVersion fromString(String value) {
    switch (value) {
      case 'v2':
      default:
        return LayrzProtocolVersion.v2;
    }
  }
}
