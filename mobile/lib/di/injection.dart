import 'package:get_it/get_it.dart';

import '../core/ports/storage_port.dart';
import '../core/application/session_service.dart';
import '../adapters/storage/memory_storage_adapter.dart';

/// Global service locator instance
final getIt = GetIt.instance;

/// Configure all dependencies
///
/// Call this once at app startup before runApp().
/// GetIt provides simple service location without Flutter coupling.
Future<void> configureDependencies() async {
  // --- Adapters (implementations) ---

  // Storage adapter - swap this for Drift adapter when ready
  getIt.registerLazySingleton<StoragePort>(
    () => MemoryStorageAdapter(),
  );

  // --- Application Services ---

  getIt.registerLazySingleton<SessionService>(
    () => SessionService(getIt<StoragePort>()),
  );

  // --- Initialize storage ---
  await getIt<StoragePort>().init();
}

/// Reset all dependencies (useful for testing)
Future<void> resetDependencies() async {
  await getIt.reset();
}
