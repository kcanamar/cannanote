  use chrono::{DateTime, Utc};
  use uuid::Uuid;

  /// Session status in the logging lifecycle
  #[derive(Debug, Clone, PartialEq)]
  pub enum SessionStatus {
      Active,
      Completed,
  }

  /// Sync status for offline-first architecture
  #[derive(Debug, Clone, PartialEq)]
  pub enum SyncStatus {
      Local,
      Pending,
      Synced,
  }

  /// Core domain entity - no I/O, no frameworks
  #[derive(Debug, Clone)]
  pub struct Session {
      pub id: Uuid,
      pub method: String,
      pub amount: String,
      pub unit: String,
      pub notes: String,
      pub status: SessionStatus,
      pub sync_status: SyncStatus,
      pub timestamp: DateTime<Utc>,
      pub created_at: DateTime<Utc>,
      pub updated_at: DateTime<Utc>,
  }

  impl Session {
      /// Create new session with defaults
      pub fn new(method: String, amount: String, unit: String) -> Self {
          let now = Utc::now();
          Self {
              id: Uuid::new_v4(),
              method,
              amount,
              unit,
              notes: String::new(),
              status: SessionStatus::Active,
              sync_status: SyncStatus::Local,
              timestamp: now,
              created_at: now,
              updated_at: now,
          }
      }

      /// Domain validation - no external dependencies
      pub fn validate(&self) -> Result<(), String> {
          if self.method.is_empty() {
              return Err("Method is required".into());
          }
          if self.amount.is_empty() {
              return Err("Amount is required".into());
          }
          Ok(())
      }

      /// Mark session complete
      pub fn complete(&mut self) {
          self.status = SessionStatus::Completed;
          self.updated_at = Utc::now();
      }
  }
