use async_trait::async_trait;
use std::collections::HashMap;
use std::sync::RwLock;
use uuid::Uuid;

use crate::core::domain::{Session, SyncStatus};
use crate::core::ports::{StoragePort, StorageError};

/// Simple in-memory storage for testing
pub struct MemoryStorageAdapter {
    sessions: RwLock<HashMap<Uuid, Session>>,
}

impl MemoryStorageAdapter {
    pub fn new() -> Self {
        Self {
            sessions: RwLock::new(HashMap::new()),
        }
    }
}

impl Default for MemoryStorageAdapter {
    fn default() -> Self {
        Self::new()
    }
}

#[async_trait]
impl StoragePort for MemoryStorageAdapter {
    async fn init(&self) -> Result<(), StorageError> {
        Ok(()) // Nothing to initialize
    }

    async fn add_session(&self, session: &Session) -> Result<Session, StorageError> {
        let mut sessions = self.sessions.write().unwrap();
        sessions.insert(session.id, session.clone());
        Ok(session.clone())
    }

    async fn get_session(&self, id: Uuid) -> Result<Option<Session>, StorageError> {
        let sessions = self.sessions.read().unwrap();
        Ok(sessions.get(&id).cloned())
    }

    async fn get_all_sessions(&self) -> Result<Vec<Session>, StorageError> {
        let sessions = self.sessions.read().unwrap();
        let mut result: Vec<Session> = sessions.values().cloned().collect();
        // Sort by timestamp descending (newest first)
        result.sort_by(|a, b| b.timestamp.cmp(&a.timestamp));
        Ok(result)
    }

    async fn update_session(&self, session: &Session) -> Result<Session, StorageError> {
        let mut sessions = self.sessions.write().unwrap();
        sessions.insert(session.id, session.clone());
        Ok(session.clone())
    }

    async fn delete_session(&self, id: Uuid) -> Result<(), StorageError> {
        let mut sessions = self.sessions.write().unwrap();
        sessions.remove(&id);
        Ok(())
    }

    async fn get_pending_sessions(&self) -> Result<Vec<Session>, StorageError> {
        let sessions = self.sessions.read().unwrap();
        Ok(sessions
            .values()
            .filter(|s| s.sync_status == SyncStatus::Pending)
            .cloned()
            .collect())
    }

    fn storage_type(&self) -> &'static str {
        "memory"
    }
}
