import '../domain/session.dart';
import '../ports/storage_port.dart';

/// Application service - orchestrates domain + ports
///
/// This is where business logic lives. It depends on ports (interfaces),
/// not concrete implementations, enabling easy testing and swapping.
class SessionService {
  final StoragePort _storage;

  SessionService(this._storage);

  /// Log a new session
  Future<Session> logSession({
    required String method,
    required String amount,
    required String unit,
    String notes = '',
  }) async {
    final session = Session(
      method: method,
      amount: amount,
      unit: unit,
      notes: notes,
    );

    // Domain validation
    session.validate();

    // Persist via port
    return _storage.addSession(session);
  }

  /// Get all sessions, newest first
  Future<List<Session>> getAllSessions() {
    return _storage.getAllSessions();
  }

  /// Get single session by ID
  Future<Session?> getSession(String id) {
    return _storage.getSession(id);
  }

  /// Complete a session
  Future<Session> completeSession(String id) async {
    final session = await _storage.getSession(id);

    if (session == null) {
      throw StorageException('Session not found: $id', code: 'NOT_FOUND');
    }

    session.complete();
    return _storage.updateSession(session);
  }

  /// Delete a session
  Future<void> deleteSession(String id) {
    return _storage.deleteSession(id);
  }

  /// Get sessions pending sync
  Future<List<Session>> getPendingSessions() {
    return _storage.getPendingSessions();
  }
}
