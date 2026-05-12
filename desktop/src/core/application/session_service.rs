  use std::sync::Arc;
  use uuid::Uuid;

  use crate::core::domain::Session;
  use crate::core::ports::{StoragePort, StorageError};

  /// Application service - orchestrates domain + ports
  pub struct SessionService {
      storage: Arc<dyn StoragePort>,
  }

  impl SessionService {
      pub fn new(storage: Arc<dyn StoragePort>) -> Self {
          Self { storage }
      }

      /// Log a new session
      pub async fn log_session(
          &self,
          method: String,
          amount: String,
          unit: String,
      ) -> Result<Session, StorageError> {
          let session = Session::new(method, amount, unit);

          // Domain validation
          session.validate().map_err(|e| StorageError::Database(e))?;

          // Persist via port
          self.storage.add_session(&session).await
      }

      /// Get all sessions
      pub async fn get_all_sessions(&self) -> Result<Vec<Session>, StorageError> {
          self.storage.get_all_sessions().await
      }

      /// Complete a session
      pub async fn complete_session(&self, id: Uuid) -> Result<Session, StorageError> {
          let session = self.storage.get_session(id).await?;

          match session {
              Some(mut s) => {
                  s.complete();
                  self.storage.update_session(&s).await
              }
              None => Err(StorageError::NotFound(id)),
          }
      }
  }
