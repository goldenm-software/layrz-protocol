part of '../database.dart';

class Messages extends Table {
  IntColumn get id => integer().autoIncrement()();
  TextColumn get message => text()();
  DateTimeColumn get createdAt => dateTime()();
}
