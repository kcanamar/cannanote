// CannaNote IndexedDB Storage Adapter
// Fallback storage when OPFS is unavailable
// Based on original session-db.js, refactored to implement StoragePort

import {
  DEFAULT_TERPENES,
  DEFAULT_CANNABINOIDS,
  uuid,
  parseCompact,
  serializeCompact
} from '../../core/ports/storage-port.js';

/**
 * IndexedDB Storage Adapter
 * Fallback when OPFS is unavailable
 *
 * Features:
 * - Persistent storage in most browsers
 * - Index-based queries
 * - Schema versioning
 */
class IndexedDBAdapter {
  constructor() {
    this.DB_NAME = 'cannanote';
    this.DB_VERSION = 3;
    this.db = null;
    this._initialized = false;
  }

  // Store definitions with indexes
  get STORES() {
    return {
      sessions: {
        keyPath: 'id',
        indexes: [
          { name: 'timestamp', keyPath: 'timestamp' },
          { name: 'productId', keyPath: 'productId' },
          { name: 'syncStatus', keyPath: 'syncStatus' },
          { name: 'status', keyPath: 'status' }
        ]
      },
      checkIns: {
        keyPath: 'id',
        indexes: [
          { name: 'sessionId', keyPath: 'sessionId' },
          { name: 'timestamp', keyPath: 'timestamp' }
        ]
      },
      products: {
        keyPath: 'id',
        indexes: [
          { name: 'strainName', keyPath: 'strainName' },
          { name: 'producer', keyPath: 'producer' },
          { name: 'lastUsed', keyPath: 'lastUsed' },
          { name: 'timesUsed', keyPath: 'timesUsed' }
        ]
      },
      producers: { keyPath: 'name' },
      processors: { keyPath: 'name' },
      distributors: { keyPath: 'name' },
      terpenes: {
        keyPath: 'name',
        indexes: [{ name: 'timesUsed', keyPath: 'timesUsed' }]
      },
      cannabinoids: {
        keyPath: 'name',
        indexes: [{ name: 'timesUsed', keyPath: 'timesUsed' }]
      },
      preferences: { keyPath: 'key' },
      pendingNotifications: {
        keyPath: 'id',
        indexes: [
          { name: 'sessionId', keyPath: 'sessionId' },
          { name: 'scheduledTime', keyPath: 'scheduledTime' }
        ]
      },
      offlineAuth: { keyPath: 'email' }
    };
  }

  /**
   * Initialize IndexedDB
   * @throws {Error} If IndexedDB unavailable or blocked
   */
  async init(isRetry = false) {
    if (this._initialized && this.db) return;

    // Check IndexedDB support
    if (!('indexedDB' in window)) {
      throw new Error('IndexedDB not supported');
    }

    return new Promise((resolve, reject) => {
      let request;
      try {
        request = indexedDB.open(this.DB_NAME, this.DB_VERSION);
      } catch (err) {
        reject(new Error('IndexedDB blocked: ' + err.message));
        return;
      }

      request.onerror = (event) => {
        const error = event.target.error;

        // Handle version mismatch
        if (error.name === 'VersionError' && !isRetry) {
          console.warn('[IndexedDB] Version mismatch, resetting database');
          const deleteRequest = indexedDB.deleteDatabase(this.DB_NAME);
          deleteRequest.onsuccess = () => {
            console.log('[IndexedDB] Database reset, reinitializing');
            this.init(true).then(resolve).catch(reject);
          };
          deleteRequest.onerror = () => reject(deleteRequest.error);
          return;
        }

        reject(error);
      };

      request.onsuccess = () => {
        this.db = request.result;
        this._initialized = true;
        console.log('[IndexedDB] Initialized successfully');
        resolve();
      };

      request.onupgradeneeded = (event) => {
        const database = event.target.result;

        for (const [storeName, config] of Object.entries(this.STORES)) {
          if (!database.objectStoreNames.contains(storeName)) {
            const store = database.createObjectStore(storeName, { keyPath: config.keyPath });

            if (config.indexes) {
              for (const index of config.indexes) {
                store.createIndex(index.name, index.keyPath, { unique: false });
              }
            }
          }
        }
      };
    });
  }

  // === Generic CRUD Operations ===

  async _add(storeName, data) {
    return new Promise((resolve, reject) => {
      const tx = this.db.transaction(storeName, 'readwrite');
      const store = tx.objectStore(storeName);
      const request = store.add(data);

      request.onsuccess = () => resolve(data);
      request.onerror = () => reject(request.error);
    });
  }

  async _get(storeName, key) {
    return new Promise((resolve, reject) => {
      const tx = this.db.transaction(storeName, 'readonly');
      const store = tx.objectStore(storeName);
      const request = store.get(key);

      request.onsuccess = () => resolve(request.result);
      request.onerror = () => reject(request.error);
    });
  }

  async _update(storeName, data) {
    return new Promise((resolve, reject) => {
      const tx = this.db.transaction(storeName, 'readwrite');
      const store = tx.objectStore(storeName);
      const request = store.put(data);

      request.onsuccess = () => resolve(data);
      request.onerror = () => reject(request.error);
    });
  }

  async _remove(storeName, key) {
    return new Promise((resolve, reject) => {
      const tx = this.db.transaction(storeName, 'readwrite');
      const store = tx.objectStore(storeName);
      const request = store.delete(key);

      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
  }

  async _getAll(storeName) {
    return new Promise((resolve, reject) => {
      const tx = this.db.transaction(storeName, 'readonly');
      const store = tx.objectStore(storeName);
      const request = store.getAll();

      request.onsuccess = () => resolve(request.result || []);
      request.onerror = () => reject(request.error);
    });
  }

  async _getByIndex(storeName, indexName, value) {
    return new Promise((resolve, reject) => {
      const tx = this.db.transaction(storeName, 'readonly');
      const store = tx.objectStore(storeName);
      const index = store.index(indexName);
      const request = index.getAll(value);

      request.onsuccess = () => resolve(request.result || []);
      request.onerror = () => reject(request.error);
    });
  }

  async _clear(storeName) {
    return new Promise((resolve, reject) => {
      const tx = this.db.transaction(storeName, 'readwrite');
      const store = tx.objectStore(storeName);
      const request = store.clear();

      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
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

    await this._add('sessions', session);

    // Update product usage
    if (session.productId) {
      const product = await this._get('products', session.productId);
      if (product) {
        product.timesUsed = (product.timesUsed || 0) + 1;
        product.lastUsed = Date.now();
        await this._update('products', product);
      }
    }

    return session;
  }

  async getSession(id) {
    return await this._get('sessions', id) || null;
  }

  async updateSession(id, updates) {
    const session = await this._get('sessions', id);
    if (!session) throw new Error('Session not found');

    const updated = { ...session, ...updates, updatedAt: Date.now() };
    await this._update('sessions', updated);
    return updated;
  }

  async deleteSession(id) {
    await this._remove('sessions', id);
  }

  async getAllSessions() {
    const sessions = await this._getAll('sessions');
    return sessions.sort((a, b) => b.timestamp - a.timestamp);
  }

  async getSessionsByStatus(status) {
    const sessions = await this._getByIndex('sessions', 'status', status);
    return sessions.sort((a, b) => b.timestamp - a.timestamp);
  }

  async getRecentSessions(limit = 20) {
    const sessions = await this._getAll('sessions');
    return sessions
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

    await this._add('checkIns', checkIn);
    return checkIn;
  }

  async getSessionCheckIns(sessionId) {
    const checkIns = await this._getByIndex('checkIns', 'sessionId', sessionId);
    return checkIns.sort((a, b) => a.timestamp - b.timestamp);
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

    await this._add('products', product);

    // Update autocomplete stores
    if (productData.producer) await this.addToAutocomplete('producers', productData.producer);
    if (productData.processor) await this.addToAutocomplete('processors', productData.processor);
    if (productData.distributor) await this.addToAutocomplete('distributors', productData.distributor);

    // Update usage counts
    for (const terp of parseCompact(terpenesStr)) {
      await this.incrementTerpeneUsage(terp.name);
    }
    for (const canna of parseCompact(cannabinoidsStr)) {
      await this.incrementCannabinoidUsage(canna.name);
    }

    return product;
  }

  async getProduct(id) {
    return await this._get('products', id) || null;
  }

  async updateProduct(id, updates) {
    const product = await this._get('products', id);
    if (!product) throw new Error('Product not found');

    const updated = { ...product, ...updates };
    await this._update('products', updated);
    return updated;
  }

  async getAllProducts() {
    const products = await this._getAll('products');
    return products.sort((a, b) => (b.lastUsed || 0) - (a.lastUsed || 0));
  }

  async searchProducts(query) {
    const products = await this._getAll('products');
    const lowerQuery = query.toLowerCase();
    return products.filter(p =>
      p.strainName.toLowerCase().includes(lowerQuery) ||
      (p.producer && p.producer.toLowerCase().includes(lowerQuery))
    );
  }

  // === Terpenes & Cannabinoids ===

  async getTerpenes() {
    const terpenes = await this._getAll('terpenes');
    return terpenes.sort((a, b) => (b.timesUsed || 0) - (a.timesUsed || 0));
  }

  async getCannabinoids() {
    const cannabinoids = await this._getAll('cannabinoids');
    return cannabinoids.sort((a, b) => (b.timesUsed || 0) - (a.timesUsed || 0));
  }

  async incrementTerpeneUsage(name) {
    let terp = await this._get('terpenes', name);
    if (terp) {
      terp.timesUsed = (terp.timesUsed || 0) + 1;
      await this._update('terpenes', terp);
    } else {
      await this._add('terpenes', { name, timesUsed: 1, createdAt: Date.now() });
    }
  }

  async incrementCannabinoidUsage(name) {
    let canna = await this._get('cannabinoids', name);
    if (canna) {
      canna.timesUsed = (canna.timesUsed || 0) + 1;
      await this._update('cannabinoids', canna);
    } else {
      await this._add('cannabinoids', { name, timesUsed: 1, createdAt: Date.now() });
    }
  }

  // === Autocomplete Stores ===

  async getAutocomplete(storeName) {
    const items = await this._getAll(storeName);
    return items.map(i => i.name).sort();
  }

  async addToAutocomplete(storeName, value) {
    const existing = await this._get(storeName, value);
    if (!existing) {
      await this._add(storeName, { name: value, createdAt: Date.now() });
    }
  }

  // === Preferences ===

  async getPreference(key, defaultValue = null) {
    const pref = await this._get('preferences', key);
    return pref ? pref.value : defaultValue;
  }

  async setPreference(key, value) {
    await this._update('preferences', { key, value, updatedAt: Date.now() });
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
    await this._add('pendingNotifications', notification);
    return notification;
  }

  async getPendingNotifications() {
    const notifications = await this._getAll('pendingNotifications');
    const now = Date.now();
    return notifications.filter(n => n.status === 'pending' && n.scheduledTime <= now);
  }

  async markNotificationShown(notificationId) {
    const notification = await this._get('pendingNotifications', notificationId);
    if (notification) {
      notification.status = 'shown';
      notification.shownAt = Date.now();
      await this._update('pendingNotifications', notification);
    }
  }

  async clearSessionNotifications(sessionId) {
    const notifications = await this._getByIndex('pendingNotifications', 'sessionId', sessionId);
    for (const n of notifications) {
      await this._remove('pendingNotifications', n.id);
    }
  }

  // === Offline Auth ===

  async saveOfflineCredentials(email, password) {
    const salt = this._generateSalt();
    const passwordHash = await this._hashPassword(password, salt);

    const credentials = {
      email: email.toLowerCase(),
      passwordHash,
      salt,
      savedAt: Date.now()
    };

    await this._update('offlineAuth', credentials);
    console.log('[IndexedDB] Offline credentials saved for:', email);
    return true;
  }

  async verifyOfflineCredentials(email, password) {
    const credentials = await this._get('offlineAuth', email.toLowerCase());

    if (!credentials) {
      console.log('[IndexedDB] No offline credentials found for:', email);
      return { valid: false, reason: 'no_credentials' };
    }

    const passwordHash = await this._hashPassword(password, credentials.salt);

    if (passwordHash === credentials.passwordHash) {
      console.log('[IndexedDB] Offline login successful for:', email);
      return { valid: true, email: credentials.email };
    }

    console.log('[IndexedDB] Offline login failed - wrong password');
    return { valid: false, reason: 'wrong_password' };
  }

  async hasOfflineCredentials(email) {
    const credentials = await this._get('offlineAuth', email.toLowerCase());
    return !!credentials;
  }

  async clearOfflineCredentials() {
    await this._clear('offlineAuth');
    console.log('[IndexedDB] Offline credentials cleared');
  }

  // Hash helpers with WebCrypto fallback
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

  async _hashPassword(password, salt) {
    // Try WebCrypto first
    if (typeof crypto !== 'undefined' && crypto.subtle) {
      try {
        const encoder = new TextEncoder();
        const data = encoder.encode(password + salt);
        const hashBuffer = await crypto.subtle.digest('SHA-256', data);
        const hashArray = Array.from(new Uint8Array(hashBuffer));
        return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
      } catch (err) {
        console.warn('[IndexedDB] WebCrypto hash failed, using fallback');
      }
    }

    // Simple fallback hash
    return this._simpleHash(password + salt);
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
    return await this._getByIndex('sessions', 'syncStatus', 'pending');
  }

  async markSessionsSynced(ids) {
    for (const id of ids) {
      const session = await this._get('sessions', id);
      if (session) {
        session.syncStatus = 'synced';
        await this._update('sessions', session);
      }
    }
  }

  // === Export ===

  async exportData() {
    return {
      version: this.DB_VERSION,
      exportedAt: new Date().toISOString(),
      storageType: 'indexeddb',
      sessions: await this._getAll('sessions'),
      checkIns: await this._getAll('checkIns'),
      products: await this._getAll('products'),
      preferences: await this._getAll('preferences')
    };
  }

  async exportCSV() {
    const sessions = await this._getAll('sessions');
    const products = await this._getAll('products');

    const productMap = {};
    for (const p of products) {
      productMap[p.id] = p;
    }

    const rows = [
      ['Date', 'Time', 'Strain', 'Producer', 'Processor', 'Distributor', 'Cannabinoids', 'Terpenes', 'Method', 'Amount', 'Unit', 'Mood Before', 'Mind Before', 'Body Before', 'Status', 'Notes']
    ];

    for (const s of sessions.sort((a, b) => a.timestamp - b.timestamp)) {
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
    return 'indexeddb';
  }

  getCapabilities() {
    return {
      persistent: true,
      capacity: 'large',
      quotaManaged: false
    };
  }

  // === Utility ===

  async clear(storeName) {
    if (storeName === 'all') {
      for (const store of Object.keys(this.STORES)) {
        await this._clear(store);
      }
    } else {
      await this._clear(storeName);
    }
  }

  async seedDefaults() {
    // Seed terpenes
    const existingTerpenes = await this._getAll('terpenes');
    if (existingTerpenes.length === 0) {
      for (const name of DEFAULT_TERPENES) {
        await this._add('terpenes', { name, timesUsed: 0, createdAt: Date.now() });
      }
    }

    // Seed cannabinoids
    const existingCannabinoids = await this._getAll('cannabinoids');
    if (existingCannabinoids.length === 0) {
      for (const name of DEFAULT_CANNABINOIDS) {
        await this._add('cannabinoids', { name, timesUsed: 0, createdAt: Date.now() });
      }
    }
  }
}

export { IndexedDBAdapter };
