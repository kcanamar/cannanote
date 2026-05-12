// CannaNote Memory Storage Adapter
// Last resort fallback - data lives only in memory
// Used when OPFS and IndexedDB are both unavailable (private browsing)

import {
  DEFAULT_TERPENES,
  DEFAULT_CANNABINOIDS,
  uuid,
  parseCompact,
  serializeCompact
} from '../../core/ports/storage-port.js';

/**
 * In-Memory Storage Adapter
 * Used as last resort when persistent storage is unavailable
 *
 * Characteristics:
 * - Data lost on page refresh
 * - Always available (no browser restrictions)
 * - Good for demo/private browsing modes
 */
class MemoryAdapter {
  constructor() {
    this.data = this._getEmptyData();
    this._initialized = false;
  }

  /**
   * Initialize memory storage (no-op, always succeeds)
   */
  async init() {
    if (this._initialized) return;
    this._initialized = true;
    console.log('[Memory] Initialized - data will not persist across sessions');
  }

  /**
   * Get empty data structure
   */
  _getEmptyData() {
    return {
      version: 1,
      sessions: [],
      checkIns: [],
      products: [],
      terpenes: [],
      cannabinoids: [],
      producers: [],
      processors: [],
      distributors: [],
      preferences: {},
      pendingNotifications: [],
      offlineAuth: {}
    };
  }

  // === Session Operations ===

  async addSession(sessionData) {
    const session = {
      id: uuid(),
      productId: sessionData.productId || null,
      method: sessionData.method,
      amount: sessionData.amount,
      unit: sessionData.unit,
      moodBefore: sessionData.moodBefore || null,
      mindBefore: sessionData.mindBefore || null,
      bodyBefore: sessionData.bodyBefore || null,
      checkInMode: sessionData.checkInMode || 'none',
      checkInInterval: sessionData.checkInInterval || 30,
      timestamp: Date.now(),
      status: 'active',
      syncStatus: 'local',
      notes: sessionData.notes || '',
      createdAt: Date.now(),
      updatedAt: Date.now()
    };

    this.data.sessions.push(session);

    // Update product usage if exists
    if (session.productId) {
      const product = this.data.products.find(p => p.id === session.productId);
      if (product) {
        product.timesUsed = (product.timesUsed || 0) + 1;
        product.lastUsed = Date.now();
      }
    }

    return session;
  }

  async getSession(id) {
    return this.data.sessions.find(s => s.id === id) || null;
  }

  async updateSession(id, updates) {
    const index = this.data.sessions.findIndex(s => s.id === id);
    if (index === -1) throw new Error('Session not found');

    this.data.sessions[index] = {
      ...this.data.sessions[index],
      ...updates,
      updatedAt: Date.now()
    };
    return this.data.sessions[index];
  }

  async deleteSession(id) {
    const index = this.data.sessions.findIndex(s => s.id === id);
    if (index !== -1) {
      this.data.sessions.splice(index, 1);
    }
  }

  async getAllSessions() {
    return [...this.data.sessions].sort((a, b) => b.timestamp - a.timestamp);
  }

  async getSessionsByStatus(status) {
    return this.data.sessions
      .filter(s => s.status === status)
      .sort((a, b) => b.timestamp - a.timestamp);
  }

  async getRecentSessions(limit = 20) {
    return [...this.data.sessions]
      .sort((a, b) => b.timestamp - a.timestamp)
      .slice(0, limit);
  }

  // === Check-In Operations ===

  async addCheckIn(checkInData) {
    const checkIn = {
      id: uuid(),
      sessionId: checkInData.sessionId,
      mood: checkInData.mood,
      mind: checkInData.mind,
      body: checkInData.body,
      intensity: checkInData.intensity,
      meetsExpectations: checkInData.meetsExpectations || null,
      notes: checkInData.notes || '',
      timestamp: Date.now()
    };

    this.data.checkIns.push(checkIn);
    return checkIn;
  }

  async getSessionCheckIns(sessionId) {
    return this.data.checkIns
      .filter(c => c.sessionId === sessionId)
      .sort((a, b) => a.timestamp - b.timestamp);
  }

  // === Product Operations ===

  async addProduct(productData) {
    const cannabinoidsStr = typeof productData.cannabinoids === 'string'
      ? productData.cannabinoids
      : serializeCompact(productData.cannabinoids);
    const terpenesStr = typeof productData.terpenes === 'string'
      ? productData.terpenes
      : serializeCompact(productData.terpenes);

    const product = {
      id: uuid(),
      strainName: productData.strainName,
      producer: productData.producer || null,
      processor: productData.processor || null,
      distributor: productData.distributor || null,
      harvestDate: productData.harvestDate || null,
      cannabinoids: cannabinoidsStr,
      terpenes: terpenesStr,
      timesUsed: 0,
      lastUsed: null,
      createdAt: Date.now()
    };

    this.data.products.push(product);

    // Update autocomplete stores
    if (productData.producer) await this.addToAutocomplete('producers', productData.producer);
    if (productData.processor) await this.addToAutocomplete('processors', productData.processor);
    if (productData.distributor) await this.addToAutocomplete('distributors', productData.distributor);

    // Update terpene/cannabinoid usage
    for (const terp of parseCompact(terpenesStr)) {
      await this.incrementTerpeneUsage(terp.name);
    }
    for (const canna of parseCompact(cannabinoidsStr)) {
      await this.incrementCannabinoidUsage(canna.name);
    }

    return product;
  }

  async getProduct(id) {
    return this.data.products.find(p => p.id === id) || null;
  }

  async updateProduct(id, updates) {
    const index = this.data.products.findIndex(p => p.id === id);
    if (index === -1) throw new Error('Product not found');

    this.data.products[index] = { ...this.data.products[index], ...updates };
    return this.data.products[index];
  }

  async getAllProducts() {
    return [...this.data.products].sort((a, b) => (b.lastUsed || 0) - (a.lastUsed || 0));
  }

  async searchProducts(query) {
    const lowerQuery = query.toLowerCase();
    return this.data.products.filter(p =>
      p.strainName.toLowerCase().includes(lowerQuery) ||
      (p.producer && p.producer.toLowerCase().includes(lowerQuery))
    );
  }

  // === Terpenes & Cannabinoids ===

  async getTerpenes() {
    return [...this.data.terpenes].sort((a, b) => (b.timesUsed || 0) - (a.timesUsed || 0));
  }

  async getCannabinoids() {
    return [...this.data.cannabinoids].sort((a, b) => (b.timesUsed || 0) - (a.timesUsed || 0));
  }

  async incrementTerpeneUsage(name) {
    let terp = this.data.terpenes.find(t => t.name === name);
    if (terp) {
      terp.timesUsed = (terp.timesUsed || 0) + 1;
    } else {
      this.data.terpenes.push({ name, timesUsed: 1, createdAt: Date.now() });
    }
  }

  async incrementCannabinoidUsage(name) {
    let canna = this.data.cannabinoids.find(c => c.name === name);
    if (canna) {
      canna.timesUsed = (canna.timesUsed || 0) + 1;
    } else {
      this.data.cannabinoids.push({ name, timesUsed: 1, createdAt: Date.now() });
    }
  }

  // === Autocomplete Stores ===

  async getAutocomplete(storeName) {
    const store = this.data[storeName] || [];
    return store.map(i => typeof i === 'string' ? i : i.name).sort();
  }

  async addToAutocomplete(storeName, value) {
    if (!this.data[storeName]) {
      this.data[storeName] = [];
    }
    const existing = this.data[storeName].find(i =>
      (typeof i === 'string' ? i : i.name) === value
    );
    if (!existing) {
      this.data[storeName].push({ name: value, createdAt: Date.now() });
    }
  }

  // === Preferences ===

  async getPreference(key, defaultValue = null) {
    return this.data.preferences[key] !== undefined
      ? this.data.preferences[key]
      : defaultValue;
  }

  async setPreference(key, value) {
    this.data.preferences[key] = value;
  }

  // === Notifications ===

  async scheduleNotification(sessionId, scheduledTime) {
    const notification = {
      id: uuid(),
      sessionId,
      scheduledTime,
      status: 'pending',
      createdAt: Date.now()
    };
    this.data.pendingNotifications.push(notification);
    return notification;
  }

  async getPendingNotifications() {
    const now = Date.now();
    return this.data.pendingNotifications.filter(n =>
      n.status === 'pending' && n.scheduledTime <= now
    );
  }

  async markNotificationShown(notificationId) {
    const notification = this.data.pendingNotifications.find(n => n.id === notificationId);
    if (notification) {
      notification.status = 'shown';
      notification.shownAt = Date.now();
    }
  }

  async clearSessionNotifications(sessionId) {
    this.data.pendingNotifications = this.data.pendingNotifications.filter(
      n => n.sessionId !== sessionId
    );
  }

  // === Offline Auth ===
  // Note: In memory mode, credentials won't persist across sessions
  // This still allows testing the auth flow within a single session

  async saveOfflineCredentials(email, password) {
    const salt = this._generateSalt();
    const passwordHash = this._simpleHash(password + salt);

    this.data.offlineAuth[email.toLowerCase()] = {
      email: email.toLowerCase(),
      passwordHash,
      salt,
      savedAt: Date.now()
    };
    console.log('[Memory] Offline credentials saved (session only) for:', email);
    return true;
  }

  async verifyOfflineCredentials(email, password) {
    const credentials = this.data.offlineAuth[email.toLowerCase()];

    if (!credentials) {
      console.log('[Memory] No offline credentials found for:', email);
      return { valid: false, reason: 'no_credentials' };
    }

    const passwordHash = this._simpleHash(password + credentials.salt);

    if (passwordHash === credentials.passwordHash) {
      console.log('[Memory] Offline login successful for:', email);
      return { valid: true, email: credentials.email };
    }

    console.log('[Memory] Offline login failed - wrong password');
    return { valid: false, reason: 'wrong_password' };
  }

  async hasOfflineCredentials(email) {
    return !!this.data.offlineAuth[email.toLowerCase()];
  }

  async clearOfflineCredentials() {
    this.data.offlineAuth = {};
    console.log('[Memory] Offline credentials cleared');
  }

  _generateSalt() {
    const array = new Uint8Array(16);
    if (typeof crypto !== 'undefined' && crypto.getRandomValues) {
      crypto.getRandomValues(array);
    } else {
      for (let i = 0; i < 16; i++) {
        array[i] = Math.floor(Math.random() * 256);
      }
    }
    return Array.from(array).map(b => b.toString(16).padStart(2, '0')).join('');
  }

  _simpleHash(str) {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      const char = str.charCodeAt(i);
      hash = ((hash << 5) - hash) + char;
      hash = hash & hash;
    }
    return Math.abs(hash).toString(16).padStart(8, '0');
  }

  // === Sync Support ===

  async getPendingSessions() {
    return this.data.sessions.filter(s => s.syncStatus === 'pending');
  }

  async markSessionsSynced(ids) {
    for (const id of ids) {
      const session = this.data.sessions.find(s => s.id === id);
      if (session) {
        session.syncStatus = 'synced';
      }
    }
  }

  // === Export ===

  async exportData() {
    return {
      version: this.data.version,
      exportedAt: new Date().toISOString(),
      storageType: 'memory',
      sessions: this.data.sessions,
      checkIns: this.data.checkIns,
      products: this.data.products,
      preferences: this.data.preferences
    };
  }

  async exportCSV() {
    const sessions = this.data.sessions;
    const products = this.data.products;

    const productMap = {};
    for (const p of products) {
      productMap[p.id] = p;
    }

    const rows = [
      ['Date', 'Time', 'Strain', 'Producer', 'Processor', 'Distributor', 'Cannabinoids', 'Terpenes', 'Method', 'Amount', 'Unit', 'Mood Before', 'Mind Before', 'Body Before', 'Status', 'Notes']
    ];

    for (const s of [...sessions].sort((a, b) => a.timestamp - b.timestamp)) {
      const product = productMap[s.productId] || {};
      const date = new Date(s.timestamp);
      rows.push([
        date.toLocaleDateString(),
        date.toLocaleTimeString(),
        product.strainName || '',
        product.producer || '',
        product.processor || '',
        product.distributor || '',
        product.cannabinoids || '',
        product.terpenes || '',
        s.method || '',
        s.amount || '',
        s.unit || '',
        s.moodBefore || '',
        s.mindBefore || '',
        s.bodyBefore || '',
        s.status || '',
        (s.notes || '').replace(/"/g, '""')
      ]);
    }

    return rows.map(row => row.map(cell => `"${cell}"`).join(',')).join('\n');
  }

  // === Metadata ===

  getStorageType() {
    return 'memory';
  }

  getCapabilities() {
    return {
      persistent: false,
      capacity: 'session-only',
      quotaManaged: false
    };
  }

  // === Utility ===

  async clear(storeName) {
    if (storeName === 'all') {
      this.data = this._getEmptyData();
    } else if (this.data[storeName]) {
      if (Array.isArray(this.data[storeName])) {
        this.data[storeName] = [];
      } else {
        this.data[storeName] = {};
      }
    }
  }

  async seedDefaults() {
    // Seed terpenes
    if (this.data.terpenes.length === 0) {
      for (const name of DEFAULT_TERPENES) {
        this.data.terpenes.push({ name, timesUsed: 0, createdAt: Date.now() });
      }
    }

    // Seed cannabinoids
    if (this.data.cannabinoids.length === 0) {
      for (const name of DEFAULT_CANNABINOIDS) {
        this.data.cannabinoids.push({ name, timesUsed: 0, createdAt: Date.now() });
      }
    }
  }
}

export { MemoryAdapter };
