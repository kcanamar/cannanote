import 'package:uuid/uuid.dart';

/// Session status in the logging lifecycle
enum SessionStatus { active, completed }

/// Sync status for offline-first architecture
enum SyncStatus { local, pending, synced }

/// Core domain entity - no I/O, no frameworks
class Session {
  final String id;
  final String method;
  final String amount;
  final String unit;
  String notes;
  SessionStatus status;
  SyncStatus syncStatus;
  final DateTime timestamp;
  final DateTime createdAt;
  DateTime updatedAt;

  Session({
    String? id,
    required this.method,
    required this.amount,
    required this.unit,
    this.notes = '',
    this.status = SessionStatus.active,
    this.syncStatus = SyncStatus.local,
    DateTime? timestamp,
    DateTime? createdAt,
    DateTime? updatedAt,
  })  : id = id ?? const Uuid().v4(),
        timestamp = timestamp ?? DateTime.now(),
        createdAt = createdAt ?? DateTime.now(),
        updatedAt = updatedAt ?? DateTime.now();

  /// Domain validation - no external dependencies
  void validate() {
    if (method.isEmpty) {
      throw DomainException('Method is required');
    }
    if (amount.isEmpty) {
      throw DomainException('Amount is required');
    }
  }

  /// Mark session complete
  void complete() {
    status = SessionStatus.completed;
    updatedAt = DateTime.now();
  }

  /// Create copy with modifications
  Session copyWith({
    String? id,
    String? method,
    String? amount,
    String? unit,
    String? notes,
    SessionStatus? status,
    SyncStatus? syncStatus,
    DateTime? timestamp,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Session(
      id: id ?? this.id,
      method: method ?? this.method,
      amount: amount ?? this.amount,
      unit: unit ?? this.unit,
      notes: notes ?? this.notes,
      status: status ?? this.status,
      syncStatus: syncStatus ?? this.syncStatus,
      timestamp: timestamp ?? this.timestamp,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}

/// Domain exception for validation errors
class DomainException implements Exception {
  final String message;
  DomainException(this.message);

  @override
  String toString() => message;
}
