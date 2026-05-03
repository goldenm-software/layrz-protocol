import 'dart:ffi';
import 'dart:io';
import 'package:drift/drift.dart';
import 'package:drift/native.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:layrz_protocol/src/database/database.dart';
import 'package:sqlite3/open.dart';

void _setupSqlite3() {
  if (Platform.isLinux) {
    open.overrideFor(OperatingSystem.linux, () {
      try {
        return DynamicLibrary.open('libsqlite3.so');
      } catch (_) {
        return DynamicLibrary.open('libsqlite3.so.0');
      }
    });
  }
}

LinkDatabase _openTestDatabase() {
  return LinkDatabase.fromExecutor(NativeDatabase.memory());
}

void main() {
  setUpAll(() {
    _setupSqlite3();
  });

  late LinkDatabase db;

  setUp(() {
    db = _openTestDatabase();
  });

  tearDown(() async {
    await db.close();
  });

  test('Messages table inserts and retrieves a message', () async {
    await db
        .into(db.messages)
        .insert(
          MessagesCompanion.insert(
            message: '<Ao>1;A1B2</Ao>',
            createdAt: DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true),
          ),
        );

    final rows = await db.select(db.messages).get();
    expect(rows.length, 1);
    expect(rows[0].message, '<Ao>1;A1B2</Ao>');
  });

  test('Messages table inserts multiple rows', () async {
    for (int i = 0; i < 3; i++) {
      await db
          .into(db.messages)
          .insert(
            MessagesCompanion.insert(
              message: 'msg$i',
              createdAt: DateTime.fromMillisecondsSinceEpoch((1700000000 + i) * 1000, isUtc: true),
            ),
          );
    }

    final rows = await db.select(db.messages).get();
    expect(rows.length, 3);
  });

  test('Messages table deletes a row', () async {
    final id = await db
        .into(db.messages)
        .insert(
          MessagesCompanion.insert(
            message: 'to delete',
            createdAt: DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true),
          ),
        );

    final rows = await db.select(db.messages).get();
    expect(rows.length, 1);

    await (db.delete(db.messages)..where((t) => t.id.equals(id))).go();
    final afterDelete = await db.select(db.messages).get();
    expect(afterDelete.length, 0);
  });

  test('Messages table row has correct fields', () async {
    final createdAt = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    await db
        .into(db.messages)
        .insert(
          MessagesCompanion.insert(
            message: 'test msg',
            createdAt: createdAt,
          ),
        );

    final rows = await db.select(db.messages).get();
    expect(rows[0].message, 'test msg');
    expect(rows[0].createdAt.millisecondsSinceEpoch, createdAt.millisecondsSinceEpoch);
    expect(rows[0].id, greaterThan(0));
  });

  test('Message.copyWith works', () async {
    final createdAt = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    await db.into(db.messages).insert(MessagesCompanion.insert(message: 'orig', createdAt: createdAt));
    final row = (await db.select(db.messages).get())[0];
    final copy = row.copyWith(message: 'modified');
    expect(copy.message, 'modified');
    expect(copy.id, row.id);
  });

  test('Message.toCompanion works', () async {
    final createdAt = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    await db.into(db.messages).insert(MessagesCompanion.insert(message: 'msg', createdAt: createdAt));
    final row = (await db.select(db.messages).get())[0];
    final companion = row.toCompanion(false);
    expect(companion.message.value, 'msg');
  });

  test('Message.toJson and fromJson roundtrip', () async {
    final createdAt = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    await db.into(db.messages).insert(MessagesCompanion.insert(message: 'json msg', createdAt: createdAt));
    final row = (await db.select(db.messages).get())[0];
    final json = row.toJson();
    expect(json['message'], 'json msg');
    final restored = Message.fromJson(json);
    expect(restored.message, 'json msg');
    expect(restored.id, row.id);
  });

  test('Message.toString includes fields', () async {
    final createdAt = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    await db.into(db.messages).insert(MessagesCompanion.insert(message: 'strtest', createdAt: createdAt));
    final row = (await db.select(db.messages).get())[0];
    expect(row.toString(), contains('strtest'));
  });

  test('Message equality and hashCode', () async {
    final createdAt = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    await db.into(db.messages).insert(MessagesCompanion.insert(message: 'eq test', createdAt: createdAt));
    final rows = await db.select(db.messages).get();
    expect(rows[0] == rows[0], true);
    expect(rows[0].hashCode, isA<int>());
  });

  test('MessagesCompanion.copyWith works', () {
    final createdAt = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    final companion = MessagesCompanion.insert(message: 'orig', createdAt: createdAt);
    final copy = companion.copyWith(message: const Value('modified'));
    expect(copy.message.value, 'modified');
  });

  test('MessagesCompanion.toString includes fields', () {
    final createdAt = DateTime.fromMillisecondsSinceEpoch(1700000000 * 1000, isUtc: true);
    final companion = MessagesCompanion.insert(message: 'tostr', createdAt: createdAt);
    expect(companion.toString(), contains('tostr'));
  });

  test('Database managers accessor works', () {
    expect(db.managers, isA<$LinkDatabaseManager>());
  });

  test('Database allTables is not empty', () {
    expect(db.allTables.toList(), isNotEmpty);
  });
}
