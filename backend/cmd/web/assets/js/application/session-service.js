// CannaNote Session Service
// Application layer providing unified API for session management
// Abstracts storage adapter details from UI code

import {
  getStorage,
  getStorageInfo,
  willDataPersist,
  migrateToOPFS
} from '../adapters/storage/storage-factory.js';
import { parseCompact, serializeCompact } from '../core/ports/storage-port.js';

/**
 * Session Service
 * Main entry point for all storage operations
 * Provides business logic layer between UI and storage adapters
 */
class SessionService {
  constructor() {
    this.storage = null;
    this._initialized = false;
    this._initPromise = null;
  }

  /**
   * Initialize the service
   * @returns {Promise<string>} Storage type being used
   */
  async init() {
    if (this._initialized) {
      return this.storage.getStorageType();
    }

    if (this._initPromise) {
      return this._initPromise;
    }

    this._initPromise = (async () => {
      // Try migration first (for existing users)
      try {
        const migration = await migrateToOPFS();
        if (migration.migrated) {
          console.log(`[SessionService] Migrated ${migration.sessionCount} sessions to OPFS`);
        }
      } catch (e) {
        console.warn('[SessionService] Migration check failed:', e);
      }

      // Get storage adapter
      this.storage = await getStorage();
      await this.storage.seedDefaults();

      this._initialized = true;
      this._initPromise = null;

      const storageType = this.storage.getStorageType();
      console.log(`[SessionService] Initialized with ${storageType} storage`);

      return storageType;
    })();

    return this._initPromise;
  }

  /**
   * Ensure service is initialized before operations
   */
  async _ensureInit() {
    if (!this._initialized) {
      await this.init();
    }
  }

  // === Session Operations ===

  /**
   * Log a new session
   * @param {Object} sessionData - Session data
   * @returns {Promise<Object>} Created session
   */
  async logSession(sessionData) {
    await this._ensureInit();
    return this.storage.addSession(sessionData);
  }

  /**
   * Get session by ID
   * @param {string} id - Session UUID
   * @returns {Promise<Object|null>}
   */
  async getSession(id) {
    await this._ensureInit();
    return this.storage.getSession(id);
  }

  /**
   * Get all active sessions
   * @returns {Promise<Object[]>}
   */
  async getActiveSessions() {
    await this._ensureInit();
    return this.storage.getSessionsByStatus('active');
  }

  /**
   * Get session history
   * @param {number} limit - Max sessions to return
   * @returns {Promise<Object[]>}
   */
  async getHistory(limit = 50) {
    await this._ensureInit();
    return this.storage.getRecentSessions(limit);
  }

  /**
   * Complete an active session
   * @param {string} sessionId - Session UUID
   * @returns {Promise<Object>}
   */
  async completeSession(sessionId) {
    await this._ensureInit();
    return this.storage.updateSession(sessionId, {
      status: 'completed',
      completedAt: Date.now()
    });
  }

  /**
   * Delete a session
   * @param {string} sessionId - Session UUID
   * @returns {Promise<void>}
   */
  async deleteSession(sessionId) {
    await this._ensureInit();
    // Clear notifications first
    await this.storage.clearSessionNotifications(sessionId);
    return this.storage.deleteSession(sessionId);
  }

  // === Check-In Operations ===

  /**
   * Add a check-in to a session
   * @param {Object} checkInData - Check-in data
   * @returns {Promise<Object>}
   */
  async addCheckIn(checkInData) {
    await this._ensureInit();
    return this.storage.addCheckIn(checkInData);
  }

  /**
   * Get all check-ins for a session
   * @param {string} sessionId - Session UUID
   * @returns {Promise<Object[]>}
   */
  async getSessionCheckIns(sessionId) {
    await this._ensureInit();
    return this.storage.getSessionCheckIns(sessionId);
  }

  // === Product Operations ===

  /**
   * Add a new product to the collection
   * @param {Object} productData - Product data
   * @returns {Promise<Object>}
   */
  async addProduct(productData) {
    await this._ensureInit();
    return this.storage.addProduct(productData);
  }

  /**
   * Get product by ID
   * @param {string} id - Product UUID
   * @returns {Promise<Object|null>}
   */
  async getProduct(id) {
    await this._ensureInit();
    return this.storage.getProduct(id);
  }

  /**
   * Get product with parsed cannabinoids/terpenes for display
   * @param {Object} product - Raw product
   * @returns {Object} Product with parsed fields
   */
  getProductDisplay(product) {
    return {
      ...product,
      cannabinoidsParsed: parseCompact(product.cannabinoids),
      terpenesParsed: parseCompact(product.terpenes)
    };
  }

  /**
   * Get all products
   * @returns {Promise<Object[]>}
   */
  async getProducts() {
    await this._ensureInit();
    return this.storage.getAllProducts();
  }

  /**
   * Search products
   * @param {string} query - Search query
   * @returns {Promise<Object[]>}
   */
  async searchProducts(query) {
    await this._ensureInit();
    return this.storage.searchProducts(query);
  }

  // === Terpenes & Cannabinoids ===

  /**
   * Get all terpenes sorted by usage
   * @returns {Promise<Object[]>}
   */
  async getTerpenes() {
    await this._ensureInit();
    return this.storage.getTerpenes();
  }

  /**
   * Get all cannabinoids sorted by usage
   * @returns {Promise<Object[]>}
   */
  async getCannabinoids() {
    await this._ensureInit();
    return this.storage.getCannabinoids();
  }

  // === Autocomplete ===

  /**
   * Get autocomplete suggestions
   * @param {string} field - 'producers' | 'processors' | 'distributors'
   * @returns {Promise<string[]>}
   */
  async getAutocomplete(field) {
    await this._ensureInit();
    return this.storage.getAutocomplete(field);
  }

  // === Preferences ===

  /**
   * Get a preference value
   * @param {string} key - Preference key
   * @param {*} defaultValue - Default if not set
   * @returns {Promise<*>}
   */
  async getPreference(key, defaultValue = null) {
    await this._ensureInit();
    return this.storage.getPreference(key, defaultValue);
  }

  /**
   * Set a preference value
   * @param {string} key - Preference key
   * @param {*} value - Value to store
   * @returns {Promise<void>}
   */
  async setPreference(key, value) {
    await this._ensureInit();
    return this.storage.setPreference(key, value);
  }

  // === Notifications ===

  /**
   * Schedule a check-in notification
   * @param {string} sessionId - Session UUID
   * @param {number} minutesFromNow - Minutes until notification
   * @returns {Promise<Object>}
   */
  async scheduleCheckInNotification(sessionId, minutesFromNow) {
    await this._ensureInit();
    const scheduledTime = Date.now() + (minutesFromNow * 60 * 1000);
    return this.storage.scheduleNotification(sessionId, scheduledTime);
  }

  /**
   * Get pending notifications that should fire now
   * @returns {Promise<Object[]>}
   */
  async getPendingNotifications() {
    await this._ensureInit();
    return this.storage.getPendingNotifications();
  }

  /**
   * Mark a notification as shown
   * @param {string} notificationId - Notification UUID
   * @returns {Promise<void>}
   */
  async markNotificationShown(notificationId) {
    await this._ensureInit();
    return this.storage.markNotificationShown(notificationId);
  }

  // === Offline Auth ===

  /**
   * Save credentials for offline authentication
   * @param {string} email - User email
   * @param {string} password - User password
   * @returns {Promise<boolean>}
   */
  async saveOfflineCredentials(email, password) {
    await this._ensureInit();
    return this.storage.saveOfflineCredentials(email, password);
  }

  /**
   * Verify offline credentials
   * @param {string} email - User email
   * @param {string} password - User password
   * @returns {Promise<{valid: boolean, reason?: string, email?: string}>}
   */
  async verifyOfflineCredentials(email, password) {
    await this._ensureInit();
    return this.storage.verifyOfflineCredentials(email, password);
  }

  /**
   * Check if offline credentials exist for email
   * @param {string} email - User email
   * @returns {Promise<boolean>}
   */
  async hasOfflineCredentials(email) {
    await this._ensureInit();
    return this.storage.hasOfflineCredentials(email);
  }

  /**
   * Clear offline credentials
   * @returns {Promise<void>}
   */
  async clearOfflineCredentials() {
    await this._ensureInit();
    return this.storage.clearOfflineCredentials();
  }

  // === Export ===

  /**
   * Export all data as JSON
   * @returns {Promise<Object>}
   */
  async exportData() {
    await this._ensureInit();
    return this.storage.exportData();
  }

  /**
   * Export sessions as CSV
   * @returns {Promise<string>}
   */
  async exportCSV() {
    await this._ensureInit();
    return this.storage.exportCSV();
  }

  /**
   * Download data as file
   * @param {string} format - 'json' | 'csv'
   */
  async downloadExport(format = 'json') {
    await this._ensureInit();

    let content, filename, mimeType;

    if (format === 'csv') {
      content = await this.exportCSV();
      filename = `cannanote-export-${new Date().toISOString().split('T')[0]}.csv`;
      mimeType = 'text/csv';
    } else {
      const data = await this.exportData();
      content = JSON.stringify(data, null, 2);
      filename = `cannanote-export-${new Date().toISOString().split('T')[0]}.json`;
      mimeType = 'application/json';
    }

    const blob = new Blob([content], { type: mimeType });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }

  // === Storage Info ===

  /**
   * Get storage information
   * @returns {Promise<{type: string, capabilities: Object, isPrivate: boolean}>}
   */
  async getStorageInfo() {
    await this._ensureInit();
    return getStorageInfo();
  }

  /**
   * Check if data will persist
   * @returns {Promise<boolean>}
   */
  async willDataPersist() {
    await this._ensureInit();
    return willDataPersist();
  }

  /**
   * Get storage type
   * @returns {string}
   */
  getStorageType() {
    if (!this.storage) return 'unknown';
    return this.storage.getStorageType();
  }
}

// Export singleton instance
const sessionService = new SessionService();

// Also export class for testing
export { SessionService, sessionService };

// Make available globally for non-module scripts
if (typeof window !== 'undefined') {
  window.sessionService = sessionService;
}
