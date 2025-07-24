part of '../../packets.dart';

class Position {
  /// [latitude] is the latitude of the position.
  /// This value should be between -90 and 90.
  final double? latitude;

  /// [longitude] is the longitude of the position.
  /// This value should be between -180 and 180.
  final double? longitude;

  /// [altitude]  is the altitude of the position.
  final double? altitude;

  /// [speed] is the speed of the position.
  /// This value should be greater than or equals to 0.
  final double? speed;

  /// [direction] is the direction of the position.
  /// This value should be between 0 and 360.
  final double? direction;

  /// [satellites] is the number of satellites.
  /// This value should be greater than or equals to 0.
  final int? satellites;

  /// [hdop] is the horizontal dilution of precision.
  /// This value should be greater than or equals to 0.
  final double? hdop;

  /// [Position] is the definition of the position.
  /// This class is used to send the position to the Layrz ecosystem.
  Position({
    this.latitude,
    this.longitude,
    this.altitude,
    this.speed,
    this.direction,
    this.satellites,
    this.hdop,
  })  : assert(
          latitude == null || (latitude >= -90 && latitude <= 90),
          'latitude should be between -90 and 90',
        ),
        assert(
          longitude == null || (longitude >= -180 && longitude <= 180),
          'longitude should be between -180 and 180',
        ),
        assert(
          speed == null || speed >= 0,
          'speed should be greater than or equals to 0',
        ),
        assert(
          direction == null || (direction >= 0 && direction <= 360),
          'direction should be between 0 and 360',
        ),
        assert(
          satellites == null || satellites >= 0,
          'satellites should be greater than or equals to 0',
        ),
        assert(
          hdop == null || hdop >= 0,
          'hdop should be greater than or equals to 0',
        );
}
