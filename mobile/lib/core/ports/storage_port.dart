import '../domain/session.dart';
import '../domain/product.dart';

/// Error type for storage operations
class StorageException implements Exception {
  final String message;
  final String? code;

  StorageException(this.message, {this.code});

  @override
  String toString() => 'StorageException: $message';
}

/// Port defining what storage must do - not HOW
///
/// This is the hexagonal architecture "port" - an interface that
/// defines capabilities without implementation details.
abstract class StoragePort {
  /// Initialize storage (create tables, etc)
  Future<void> init();

  // --- Sessions ---

  /// Add new session
  Future<Session> addSession(Session session);

  /// Get session by ID
  Future<Session?> getSession(String id);

  /// Get all sessions, newest first
  Future<List<Session>> getAllSessions();

  /// Update existing session
  Future<Session> updateSession(Session session);

  /// Delete session
  Future<void> deleteSession(String id);

  /// Get sessions pending sync
  Future<List<Session>> getPendingSessions();

  /// Mark sessions as synced
  Future<void> markSessionsSynced(List<String> ids);

  // --- Products ---

  /// Add new product
  Future<Product> addProduct(Product product);

  /// Get product by ID
  Future<Product?> getProduct(String id);

  /// Get all products
  Future<List<Product>> getAllProducts();

  /// Update existing product
  Future<Product> updateProduct(Product product);

  /// Delete product
  Future<void> deleteProduct(String id);

  // --- Metadata ---

  /// Storage type identifier
  String get storageType;
}
