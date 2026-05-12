import '../../core/domain/session.dart';
import '../../core/domain/product.dart';
import '../../core/ports/storage_port.dart';

/// Simple in-memory storage for testing and development
///
/// This adapter implements the StoragePort interface using
/// in-memory maps. Data persists only during app lifetime.
class MemoryStorageAdapter implements StoragePort {
  final Map<String, Session> _sessions = {};
  final Map<String, Product> _products = {};

  @override
  Future<void> init() async {
    // Nothing to initialize for memory storage
  }

  @override
  String get storageType => 'memory';

  // --- Sessions ---

  @override
  Future<Session> addSession(Session session) async {
    _sessions[session.id] = session;
    return session;
  }

  @override
  Future<Session?> getSession(String id) async {
    return _sessions[id];
  }

  @override
  Future<List<Session>> getAllSessions() async {
    final sessions = _sessions.values.toList();
    // Sort by timestamp descending (newest first)
    sessions.sort((a, b) => b.timestamp.compareTo(a.timestamp));
    return sessions;
  }

  @override
  Future<Session> updateSession(Session session) async {
    _sessions[session.id] = session;
    return session;
  }

  @override
  Future<void> deleteSession(String id) async {
    _sessions.remove(id);
  }

  @override
  Future<List<Session>> getPendingSessions() async {
    return _sessions.values
        .where((s) => s.syncStatus == SyncStatus.pending)
        .toList();
  }

  @override
  Future<void> markSessionsSynced(List<String> ids) async {
    for (final id in ids) {
      final session = _sessions[id];
      if (session != null) {
        _sessions[id] = session.copyWith(syncStatus: SyncStatus.synced);
      }
    }
  }

  // --- Products ---

  @override
  Future<Product> addProduct(Product product) async {
    _products[product.id] = product;
    return product;
  }

  @override
  Future<Product?> getProduct(String id) async {
    return _products[id];
  }

  @override
  Future<List<Product>> getAllProducts() async {
    return _products.values.toList();
  }

  @override
  Future<Product> updateProduct(Product product) async {
    _products[product.id] = product;
    return product;
  }

  @override
  Future<void> deleteProduct(String id) async {
    _products.remove(id);
  }
}
