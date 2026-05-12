// CannaNote Storage Factory
// Handles adapter selection with graceful degradation
// Priority: OPFS → IndexedDB → Memory

import { OPFSAdapter } from './opfs-adapter.js';
import { IndexedDBAdapter } from './indexeddb-adapter.js';
import { MemoryAdapter } from './memory-adapter.js';

/**
 * Storage initialization state
 */
let storageInstance = null;
let initializationPromise = null;

/**
 * Detect private browsing mode
 * @returns {Promise<boolean>}
 */
async function detectPrivateMode() {
  // Safari: quota is 0 in private mode
  if ('storage' in navigator && navigator.storage.estimate) {
    try {
      const estimate = await navigator.storage.estimate();
      // Safari private mode reports 0 quota
      if (estimate.quota === 0) {
        return true;
      }
      // Chrome private mode has limited quota (<120MB)
      if (estimate.quota < 120 * 1024 * 1024) {
        return true;
      }
    } catch (e) {
      // Estimate not available
    }
  }

  // Firefox: IndexedDB.open() throws in private mode (with partitioned storage)
  // This is checked during adapter initialization

  // Check for navigator.brave (Brave browser)
  // Brave in normal mode is fine, but shields may affect storage
  if ('brave' in navigator) {
    console.log('[Storage] Brave browser detected');
  }

  return false;
}

/**
 * Create and initialize the best available storage adapter
 * Priority: OPFS → IndexedDB → Memory
 *
 * @returns {Promise<OPFSAdapter|IndexedDBAdapter|MemoryAdapter>}
 */
async function createStorageAdapter() {
  const isPrivate = await detectPrivateMode();

  if (isPrivate) {
    console.log('[Storage] Private browsing detected, may have limited storage');
  }

  // Try OPFS first (best for privacy browsers)
  try {
    const opfs = new OPFSAdapter();
    await opfs.init();
    console.log('[Storage] Using OPFS (Origin Private File System)');
    return opfs;
  } catch (err) {
    console.warn('[Storage] OPFS unavailable:', err.message);
  }

  // Try IndexedDB second
  try {
    const idb = new IndexedDBAdapter();
    await idb.init();
    console.log('[Storage] Using IndexedDB');
    return idb;
  } catch (err) {
    console.warn('[Storage] IndexedDB unavailable:', err.message);
  }

  // Fall back to memory (always works)
  console.warn('[Storage] Using memory-only storage (data will not persist)');
  const memory = new MemoryAdapter();
  await memory.init();
  return memory;
}

/**
 * Get the singleton storage instance
 * Creates and initializes on first call
 *
 * @returns {Promise<OPFSAdapter|IndexedDBAdapter|MemoryAdapter>}
 */
async function getStorage() {
  if (storageInstance) {
    return storageInstance;
  }

  // Prevent multiple simultaneous initializations
  if (initializationPromise) {
    return initializationPromise;
  }

  initializationPromise = createStorageAdapter().then(adapter => {
    storageInstance = adapter;
    initializationPromise = null;
    return adapter;
  }).catch(err => {
    initializationPromise = null;
    throw err;
  });

  return initializationPromise;
}

/**
 * Get storage info for display
 * @returns {Promise<{type: string, capabilities: Object, isPrivate: boolean}>}
 */
async function getStorageInfo() {
  const storage = await getStorage();
  const isPrivate = await detectPrivateMode();

  return {
    type: storage.getStorageType(),
    capabilities: storage.getCapabilities(),
    isPrivate
  };
}

/**
 * Check if data will persist
 * @returns {Promise<boolean>}
 */
async function willDataPersist() {
  const storage = await getStorage();
  return storage.getCapabilities().persistent;
}

/**
 * Migrate data from IndexedDB to OPFS (upgrade path)
 * Called when OPFS becomes available for existing users
 *
 * @returns {Promise<{migrated: boolean, sessionCount: number}>}
 */
async function migrateToOPFS() {
  // Check if OPFS is available
  if (!('storage' in navigator) || !navigator.storage.getDirectory) {
    return { migrated: false, sessionCount: 0 };
  }

  // Check if we have IndexedDB data
  if (!('indexedDB' in window)) {
    return { migrated: false, sessionCount: 0 };
  }

  try {
    // Open existing IndexedDB
    const idb = new IndexedDBAdapter();
    await idb.init();

    const sessions = await idb.getAllSessions();
    const checkIns = await idb._getAll('checkIns');
    const products = await idb.getAllProducts();

    // If no data, nothing to migrate
    if (sessions.length === 0 && products.length === 0) {
      return { migrated: false, sessionCount: 0 };
    }

    // Create OPFS adapter and migrate data
    const opfs = new OPFSAdapter();
    await opfs.init();

    // Check if OPFS already has data (already migrated)
    const existingSessions = await opfs.getAllSessions();
    if (existingSessions.length > 0) {
      console.log('[Migration] OPFS already has data, skipping migration');
      return { migrated: false, sessionCount: 0 };
    }

    // Migrate sessions
    for (const session of sessions) {
      opfs.data.sessions.push(session);
    }

    // Migrate check-ins
    for (const checkIn of checkIns) {
      opfs.data.checkIns.push(checkIn);
    }

    // Migrate products
    for (const product of products) {
      opfs.data.products.push(product);
    }

    // Migrate preferences
    const prefs = await idb._getAll('preferences');
    for (const pref of prefs) {
      opfs.data.preferences[pref.key] = pref.value;
    }

    // Migrate terpenes and cannabinoids
    opfs.data.terpenes = await idb._getAll('terpenes');
    opfs.data.cannabinoids = await idb._getAll('cannabinoids');

    // Migrate autocomplete stores
    opfs.data.producers = await idb._getAll('producers');
    opfs.data.processors = await idb._getAll('processors');
    opfs.data.distributors = await idb._getAll('distributors');

    // Save migrated data
    await opfs._saveData();

    console.log(`[Migration] Migrated ${sessions.length} sessions to OPFS`);

    // Optionally clear IndexedDB after successful migration
    // await idb.clear('all');

    return { migrated: true, sessionCount: sessions.length };
  } catch (err) {
    console.error('[Migration] Failed to migrate to OPFS:', err);
    return { migrated: false, sessionCount: 0 };
  }
}

/**
 * Reset storage (for testing/debugging)
 * @param {boolean} confirm - Must be true to proceed
 */
async function resetStorage(confirm = false) {
  if (!confirm) {
    throw new Error('Must pass confirm=true to reset storage');
  }

  const storage = await getStorage();
  await storage.clear('all');
  storageInstance = null;
  console.log('[Storage] Storage reset');
}

export {
  getStorage,
  getStorageInfo,
  willDataPersist,
  migrateToOPFS,
  resetStorage,
  detectPrivateMode
};
