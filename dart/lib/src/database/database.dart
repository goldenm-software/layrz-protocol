library;

import 'package:drift/drift.dart';
import 'package:drift_flutter/drift_flutter.dart';
import 'package:flutter/foundation.dart' show visibleForTesting;
import 'package:path_provider/path_provider.dart';

part 'database.g.dart';

part 'src/message.dart';

@DriftDatabase(tables: [Messages])
class LinkDatabase extends _$LinkDatabase {
  LinkDatabase() : super(_openConnection());

  @visibleForTesting
  LinkDatabase.fromExecutor(super.executor);

  @override
  int get schemaVersion => 1;

  static QueryExecutor _openConnection() {
    return driftDatabase(
      name: 'link_database',
      native: const DriftNativeOptions(databaseDirectory: getApplicationSupportDirectory),
    );
  }
}
