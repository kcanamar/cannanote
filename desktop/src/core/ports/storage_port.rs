  use async_trait::async_trait;
  use uuid::Uuid;

  use crate::core::domain::Session;

  /// Error type for storage operations
  #[derive(Debug, thiserror::Error)]
  pub enum StorageError {
      #[error("Session not found: {0}")]
      NotFound(Uuid),

      #[error("Database error: {0}")]
      Database(String),

      #[error("Initialization failed: {0}")]
      InitFailed(String),
  }

  /// Port defining what storage must do - not HOW
  #[async_trait]
  pub trait StoragePort: Send + Sync {
      /// Initialize storage (create tables, etc)
      async fn init(&self) -> Result<(), StorageError>;

      /// Add new session
      async fn add_session(&self, session: &Session) -> Result<Session, StorageError>;

      /// Get session by ID
      async fn get_session(&self, id: Uuid) -> Result<Option<Session>, StorageError>;

      /// Get all sessions
      async fn get_all_sessions(&self) -> Result<Vec<Session>, StorageError>;

      /// Update existing session
      async fn update_session(&self, session: &Session) -> Result<Session, StorageError>;

      /// Delete session
      async fn delete_session(&self, id: Uuid) -> Result<(), StorageError>;

      /// Get sessions pending sync
      async fn get_pending_sessions(&self) -> Result<Vec<Session>, StorageError>;

      /// Storage type identifier
      fn storage_type(&self) -> &'static str;
  }
