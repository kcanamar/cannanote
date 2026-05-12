// CannaNote Storage Port Interface
// Defines the contract all storage adapters must implement
// Part of hexagonal architecture for graceful degradation

/**
 * StoragePort defines the interface for all storage adapters.
 * Implementations: OPFS (primary), IndexedDB (fallback), Memory (last resort)
 *
 * @typedef {Object} Session
 * @property {string} id - UUID
 * @property {string|null} productId - Reference to product
 * @property {string} method - Consumption method
 * @property {string} amount - Amount consumed
 * @property {string} unit - Unit of measurement
 * @property {string|null} moodBefore - Mood before session
 * @property {string|null} mindBefore - Mind state before
 * @property {string|null} bodyBefore - Body state before
 * @property {string} checkInMode - 'none' | 'manual' | 'timed'
 * @property {number} checkInInterval - Minutes between check-ins
 * @property {number} timestamp - Unix timestamp
 * @property {string} status - 'active' | 'completed'
 * @property {string} syncStatus - 'local' | 'pending' | 'synced'
 * @property {string} notes - User notes
 * @property {number} createdAt - Unix timestamp
 * @property {number} updatedAt - Unix timestamp
 *
 * @typedef {Object} CheckIn
 * @property {string} id - UUID
 * @property {string} sessionId - Reference to session
 * @property {string} mood - Mood at check-in
 * @property {string} mind - Mind state
 * @property {string} body - Body state
 * @property {number} intensity - 1-10
 * @property {string|null} meetsExpectations - User feedback
 * @property {string} notes - Check-in notes
 * @property {number} timestamp - Unix timestamp
 *
 * @typedef {Object} Product
 * @property {string} id - UUID
 * @property {string} strainName - Name of strain
 * @property {string|null} producer - Producer name
 * @property {string|null} processor - Processor name
 * @property {string|null} distributor - Distributor name
 * @property {string|null} harvestDate - Harvest date
 * @property {string} cannabinoids - Compact format "THC:24.5,CBD:0.5"
 * @property {string} terpenes - Compact format "Myrcene:0.8,Limonene"
 * @property {number} timesUsed - Usage count
 * @property {number|null} lastUsed - Last used timestamp
 * @property {number} createdAt - Unix timestamp
 *
 * @typedef {Object} StorageCapabilities
 * @property {boolean} persistent - Data persists across sessions
 * @property {string} capacity - 'large' | 'medium' | 'session-only'
 * @property {boolean} quotaManaged - Adapter handles quota limits
 *
 * @typedef {'opfs'|'indexeddb'|'memory'} StorageType
 */

/**
 * Storage Port Interface
 * All storage adapters must implement these methods.
 * @interface
 */
const StoragePortInterface = {
  /**
   * Initialize the storage adapter
   * @returns {Promise<void>}
   * @throws {Error} If storage is unavailable
   */
  init: async () => {},

  // === Session Operations ===

  /**
   * Add a new session
   * @param {Object} sessionData - Session data without id
   * @returns {Promise<Session>} Created session with id
   */
  addSession: async (sessionData) => {},

  /**
   * Get session by ID
   * @param {string} id - Session UUID
   * @returns {Promise<Session|null>}
   */
  getSession: async (id) => {},

  /**
   * Update an existing session
   * @param {string} id - Session UUID
   * @param {Object} updates - Fields to update
   * @returns {Promise<Session>}
   */
  updateSession: async (id, updates) => {},

  /**
   * Delete a session
   * @param {string} id - Session UUID
   * @returns {Promise<void>}
   */
  deleteSession: async (id) => {},

  /**
   * Get all sessions
   * @returns {Promise<Session[]>}
   */
  getAllSessions: async () => {},

  /**
   * Get sessions by status
   * @param {string} status - 'active' | 'completed'
   * @returns {Promise<Session[]>}
   */
  getSessionsByStatus: async (status) => {},

  /**
   * Get recent sessions with limit
   * @param {number} limit - Max sessions to return
   * @returns {Promise<Session[]>}
   */
  getRecentSessions: async (limit) => {},

  // === Check-In Operations ===

  /**
   * Add a check-in for a session
   * @param {Object} checkInData - Check-in data without id
   * @returns {Promise<CheckIn>}
   */
  addCheckIn: async (checkInData) => {},

  /**
   * Get all check-ins for a session
   * @param {string} sessionId - Session UUID
   * @returns {Promise<CheckIn[]>}
   */
  getSessionCheckIns: async (sessionId) => {},

  // === Product Operations ===

  /**
   * Add a new product
   * @param {Object} productData - Product data without id
   * @returns {Promise<Product>}
   */
  addProduct: async (productData) => {},

  /**
   * Get product by ID
   * @param {string} id - Product UUID
   * @returns {Promise<Product|null>}
   */
  getProduct: async (id) => {},

  /**
   * Update a product
   * @param {string} id - Product UUID
   * @param {Object} updates - Fields to update
   * @returns {Promise<Product>}
   */
  updateProduct: async (id, updates) => {},

  /**
   * Get all products
   * @returns {Promise<Product[]>}
   */
  getAllProducts: async () => {},

  /**
   * Search products by name
   * @param {string} query - Search query
   * @returns {Promise<Product[]>}
   */
  searchProducts: async (query) => {},

  // === Terpenes & Cannabinoids ===

  /**
   * Get all terpenes
   * @returns {Promise<Array<{name: string, timesUsed: number}>>}
   */
  getTerpenes: async () => {},

  /**
   * Get all cannabinoids
   * @returns {Promise<Array<{name: string, timesUsed: number}>>}
   */
  getCannabinoids: async () => {},

  /**
   * Increment terpene usage count
   * @param {string} name - Terpene name
   * @returns {Promise<void>}
   */
  incrementTerpeneUsage: async (name) => {},

  /**
   * Increment cannabinoid usage count
   * @param {string} name - Cannabinoid name
   * @returns {Promise<void>}
   */
  incrementCannabinoidUsage: async (name) => {},

  // === Autocomplete Stores ===

  /**
   * Get autocomplete suggestions
   * @param {string} storeName - 'producers' | 'processors' | 'distributors'
   * @returns {Promise<string[]>}
   */
  getAutocomplete: async (storeName) => {},

  /**
   * Add to autocomplete store
   * @param {string} storeName - Store name
   * @param {string} value - Value to add
   * @returns {Promise<void>}
   */
  addToAutocomplete: async (storeName, value) => {},

  // === Preferences ===

  /**
   * Get a preference value
   * @param {string} key - Preference key
   * @param {*} defaultValue - Default if not found
   * @returns {Promise<*>}
   */
  getPreference: async (key, defaultValue) => {},

  /**
   * Set a preference value
   * @param {string} key - Preference key
   * @param {*} value - Value to store
   * @returns {Promise<void>}
   */
  setPreference: async (key, value) => {},

  // === Notifications ===

  /**
   * Schedule a notification
   * @param {string} sessionId - Session UUID
   * @param {number} scheduledTime - Unix timestamp
   * @returns {Promise<Object>}
   */
  scheduleNotification: async (sessionId, scheduledTime) => {},

  /**
   * Get pending notifications
   * @returns {Promise<Object[]>}
   */
  getPendingNotifications: async () => {},

  /**
   * Mark notification as shown
   * @param {string} notificationId - Notification UUID
   * @returns {Promise<void>}
   */
  markNotificationShown: async (notificationId) => {},

  /**
   * Clear all notifications for a session
   * @param {string} sessionId - Session UUID
   * @returns {Promise<void>}
   */
  clearSessionNotifications: async (sessionId) => {},

  // === Offline Auth ===

  /**
   * Save credentials for offline auth
   * @param {string} email - User email
   * @param {string} password - Password (will be hashed)
   * @returns {Promise<boolean>}
   */
  saveOfflineCredentials: async (email, password) => {},

  /**
   * Verify offline credentials
   * @param {string} email - User email
   * @param {string} password - Password to verify
   * @returns {Promise<{valid: boolean, reason?: string, email?: string}>}
   */
  verifyOfflineCredentials: async (email, password) => {},

  /**
   * Check if offline credentials exist
   * @param {string} email - User email
   * @returns {Promise<boolean>}
   */
  hasOfflineCredentials: async (email) => {},

  /**
   * Clear all offline credentials
   * @returns {Promise<void>}
   */
  clearOfflineCredentials: async () => {},

  // === Sync Support ===

  /**
   * Get sessions pending sync
   * @returns {Promise<Session[]>}
   */
  getPendingSessions: async () => {},

  /**
   * Mark sessions as synced
   * @param {string[]} ids - Session UUIDs
   * @returns {Promise<void>}
   */
  markSessionsSynced: async (ids) => {},

  // === Export ===

  /**
   * Export all data as JSON
   * @returns {Promise<Object>}
   */
  exportData: async () => {},

  /**
   * Export sessions as CSV
   * @returns {Promise<string>}
   */
  exportCSV: async () => {},

  // === Metadata ===

  /**
   * Get storage type identifier
   * @returns {StorageType}
   */
  getStorageType: () => {},

  /**
   * Get storage capabilities
   * @returns {StorageCapabilities}
   */
  getCapabilities: () => {},

  // === Utility ===

  /**
   * Clear all data (use with caution)
   * @param {string} storeName - Store to clear, or 'all'
   * @returns {Promise<void>}
   */
  clear: async (storeName) => {},

  /**
   * Seed default data (terpenes, cannabinoids)
   * @returns {Promise<void>}
   */
  seedDefaults: async () => {},
};

// Default terpenes for seeding
const DEFAULT_TERPENES = [
  'Myrcene', 'Limonene', 'Caryophyllene', 'Pinene', 'Linalool',
  'Humulene', 'Terpinolene', 'Ocimene', 'Bisabolol', 'Eucalyptol'
];

// Default cannabinoids for seeding
const DEFAULT_CANNABINOIDS = [
  'THC', 'CBD', 'CBN', 'CBG', 'THCA', 'CBDA', 'THCV', 'CBC', 'Delta-8'
];

// UUID v4 generator (no crypto dependency)
function uuid() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0;
    const v = c === 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

// Compact format helpers for cannabinoids/terpenes
// Format: "NAME:PERCENTAGE,NAME:PERCENTAGE" or "NAME" if no percentage
function parseCompact(str) {
  if (!str) return [];
  return str.split(',').map(item => {
    const [name, pct] = item.split(':');
    return { name: name.trim(), percentage: pct ? parseFloat(pct) : null };
  });
}

function serializeCompact(arr) {
  if (!arr || arr.length === 0) return '';
  return arr.map(item => {
    if (typeof item === 'string') return item;
    return item.percentage != null ? `${item.name}:${item.percentage}` : item.name;
  }).join(',');
}

// Export for use by adapters
export {
  StoragePortInterface,
  DEFAULT_TERPENES,
  DEFAULT_CANNABINOIDS,
  uuid,
  parseCompact,
  serializeCompact
};
