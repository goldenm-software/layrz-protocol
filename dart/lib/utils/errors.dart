// Exceptions

class ParseException implements Exception {
  final String message;

  ParseException(this.message);

  @override
  String toString() {
    return 'ParseException: $message';
  }
}

class CrcException implements Exception {
  final String message;

  CrcException(this.message);

  @override
  String toString() {
    return 'CrcException: $message';
  }
}

class CommandException implements Exception {
  final String message;

  CommandException(this.message);

  @override
  String toString() {
    return 'CommandException: $message';
  }
}

class ServerException implements Exception {
  final String message;

  ServerException(this.message);

  @override
  String toString() {
    return 'ServerException: $message';
  }
}

class MalformedException implements Exception {
  final String message;

  MalformedException(this.message);

  @override
  String toString() {
    return 'MalformedException: $message';
  }
}
