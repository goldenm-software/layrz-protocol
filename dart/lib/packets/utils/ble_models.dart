part of '../packets.dart';

class BleManufacturerData {
  final int companyId;
  final List<int> data;
  const BleManufacturerData({required this.companyId, required this.data});
}

class BleServiceData {
  final int uuid;
  final List<int> data;
  const BleServiceData({required this.uuid, required this.data});
}
