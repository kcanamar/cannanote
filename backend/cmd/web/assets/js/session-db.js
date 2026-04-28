// CannaNote Session Database
// IndexedDB storage adapter for offline-first session logging
// Zero dependencies, ~200 lines, promise-based API

const SessionDB = (function() {
  const DB_NAME = 'cannanote';
  const DB_VERSION = 3;

  let db = null;

  // Compact format helpers for cannabinoids/terpenes
  // Format: "NAME:PERCENTAGE,NAME:PERCENTAGE" or "NAME" if no percentage
  // Example: "THC:24.5,CBD:0.5,CBG" or "Myrcene:0.8,Limonene,Caryophyllene:0.3"
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

  // Store definitions with indexes
  const STORES = {
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
    producers: {
      keyPath: 'name'
    },
    processors: {
      keyPath: 'name'
    },
    distributors: {
      keyPath: 'name'
    },
    terpenes: {
      keyPath: 'name',
      indexes: [
        { name: 'timesUsed', keyPath: 'timesUsed' }
      ]
    },
    cannabinoids: {
      keyPath: 'name',
      indexes: [
        { name: 'timesUsed', keyPath: 'timesUsed' }
      ]
    },
    preferences: {
      keyPath: 'key'
    },
    pendingNotifications: {
      keyPath: 'id',
      indexes: [
        { name: 'sessionId', keyPath: 'sessionId' },
        { name: 'scheduledTime', keyPath: 'scheduledTime' }
      ]
    },
    offlineAuth: {
      keyPath: 'email'
    }
  };

  // Default terpenes list
  const DEFAULT_TERPENES = [
    'Myrcene', 'Limonene', 'Caryophyllene', 'Pinene', 'Linalool',
    'Humulene', 'Terpinolene', 'Ocimene', 'Bisabolol', 'Eucalyptol'
  ];

  // Default cannabinoids list
  const DEFAULT_CANNABINOIDS = [
    'THC', 'CBD', 'CBN', 'CBG', 'THCA', 'CBDA', 'THCV', 'CBC', 'Delta-8'
  ];

  // Generate UUID v4
  function uuid() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
      const r = Math.random() * 16 | 0;
      const v = c === 'x' ? r : (r & 0x3 | 0x8);
      return v.toString(16);
    });
  }

  // Initialize database
  // Handles version mismatch gracefully by resetting if needed
  async function init(isRetry = false) {
    if (db) return db;

    return new Promise((resolve, reject) => {
      const request = indexedDB.open(DB_NAME, DB_VERSION);

      request.onerror = (event) => {
        const error = event.target.error;

        // Handle version mismatch (e.g., browser has higher version than code expects)
        // This can happen when site data is partially cleared or code is rolled back
        if (error.name === 'VersionError' && !isRetry) {
          console.warn('[SessionDB] Version mismatch detected, resetting database');
          const deleteRequest = indexedDB.deleteDatabase(DB_NAME);
          deleteRequest.onsuccess = () => {
            console.log('[SessionDB] Database reset, reinitializing');
            init(true).then(resolve).catch(reject);
          };
          deleteRequest.onerror = () => reject(deleteRequest.error);
          return;
        }

        reject(error);
      };

      request.onsuccess = () => {
        db = request.result;
        resolve(db);
      };

      request.onupgradeneeded = (event) => {
        const database = event.target.result;

        // Create object stores
        for (const [storeName, config] of Object.entries(STORES)) {
          if (!database.objectStoreNames.contains(storeName)) {
            const store = database.createObjectStore(storeName, { keyPath: config.keyPath });

            // Create indexes
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

  // Seed default data
  async function seedDefaults() {
    await init();

    // Seed terpenes
    const existingTerpenes = await getAll('terpenes');
    if (existingTerpenes.length === 0) {
      for (const name of DEFAULT_TERPENES) {
        await add('terpenes', { name, timesUsed: 0, createdAt: Date.now() });
      }
    }

    // Seed cannabinoids
    const existingCannabinoids = await getAll('cannabinoids');
    if (existingCannabinoids.length === 0) {
      for (const name of DEFAULT_CANNABINOIDS) {
        await add('cannabinoids', { name, timesUsed: 0, createdAt: Date.now() });
      }
    }
  }

  // Generic CRUD operations
  async function add(storeName, data) {
    await init();
    return new Promise((resolve, reject) => {
      const tx = db.transaction(storeName, 'readwrite');
      const store = tx.objectStore(storeName);
      const request = store.add(data);

      request.onsuccess = () => resolve(data);
      request.onerror = () => reject(request.error);
    });
  }

  async function get(storeName, key) {
    await init();
    return new Promise((resolve, reject) => {
      const tx = db.transaction(storeName, 'readonly');
      const store = tx.objectStore(storeName);
      const request = store.get(key);

      request.onsuccess = () => resolve(request.result);
      request.onerror = () => reject(request.error);
    });
  }

  async function update(storeName, data) {
    await init();
    return new Promise((resolve, reject) => {
      const tx = db.transaction(storeName, 'readwrite');
      const store = tx.objectStore(storeName);
      const request = store.put(data);

      request.onsuccess = () => resolve(data);
      request.onerror = () => reject(request.error);
    });
  }

  async function remove(storeName, key) {
    await init();
    return new Promise((resolve, reject) => {
      const tx = db.transaction(storeName, 'readwrite');
      const store = tx.objectStore(storeName);
      const request = store.delete(key);

      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
  }

  async function getAll(storeName) {
    await init();
    return new Promise((resolve, reject) => {
      const tx = db.transaction(storeName, 'readonly');
      const store = tx.objectStore(storeName);
      const request = store.getAll();

      request.onsuccess = () => resolve(request.result || []);
      request.onerror = () => reject(request.error);
    });
  }

  async function getByIndex(storeName, indexName, value) {
    await init();
    return new Promise((resolve, reject) => {
      const tx = db.transaction(storeName, 'readonly');
      const store = tx.objectStore(storeName);
      const index = store.index(indexName);
      const request = index.getAll(value);

      request.onsuccess = () => resolve(request.result || []);
      request.onerror = () => reject(request.error);
    });
  }

  async function clear(storeName) {
    await init();
    return new Promise((resolve, reject) => {
      const tx = db.transaction(storeName, 'readwrite');
      const store = tx.objectStore(storeName);
      const request = store.clear();

      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
  }

  // Session-specific operations
  async function createSession(data) {
    const session = {
      id: uuid(),
      productId: data.productId || null,
      method: data.method,
      amount: data.amount,
      unit: data.unit,
      moodBefore: data.moodBefore || null,
      mindBefore: data.mindBefore || null,
      bodyBefore: data.bodyBefore || null,
      checkInMode: data.checkInMode || 'none',
      checkInInterval: data.checkInInterval || 30,
      timestamp: Date.now(),
      status: 'active',
      syncStatus: 'local',
      notes: data.notes || '',
      createdAt: Date.now(),
      updatedAt: Date.now()
    };

    await add('sessions', session);

    // Update product usage if exists
    if (session.productId) {
      const product = await get('products', session.productId);
      if (product) {
        product.timesUsed = (product.timesUsed || 0) + 1;
        product.lastUsed = Date.now();
        await update('products', product);
      }
    }

    return session;
  }

  async function getActiveSessions() {
    const sessions = await getByIndex('sessions', 'status', 'active');
    return sessions.sort((a, b) => b.timestamp - a.timestamp);
  }

  async function getRecentSessions(limit = 20) {
    const sessions = await getAll('sessions');
    return sessions
      .sort((a, b) => b.timestamp - a.timestamp)
      .slice(0, limit);
  }

  async function completeSession(sessionId) {
    const session = await get('sessions', sessionId);
    if (session) {
      session.status = 'completed';
      session.completedAt = Date.now();
      session.updatedAt = Date.now();
      await update('sessions', session);
    }
    return session;
  }

  // Check-in operations
  async function createCheckIn(data) {
    const checkIn = {
      id: uuid(),
      sessionId: data.sessionId,
      mood: data.mood,
      mind: data.mind,
      body: data.body,
      intensity: data.intensity,
      meetsExpectations: data.meetsExpectations || null,
      notes: data.notes || '',
      timestamp: Date.now()
    };

    await add('checkIns', checkIn);
    return checkIn;
  }

  async function getSessionCheckIns(sessionId) {
    const checkIns = await getByIndex('checkIns', 'sessionId', sessionId);
    return checkIns.sort((a, b) => a.timestamp - b.timestamp);
  }

  // Product (personal collection) operations
  // Product fields: strainName, producer, processor, distributor, harvestDate
  // Cannabinoids/terpenes stored as compact strings: "THC:24.5,CBD:0.5,CBG"
  async function createProduct(data) {
    // Normalize cannabinoids/terpenes to compact string format
    const cannabinoidsStr = typeof data.cannabinoids === 'string'
      ? data.cannabinoids
      : serializeCompact(data.cannabinoids);
    const terpenesStr = typeof data.terpenes === 'string'
      ? data.terpenes
      : serializeCompact(data.terpenes);

    const product = {
      id: uuid(),
      strainName: data.strainName,
      producer: data.producer || null,
      processor: data.processor || null,
      distributor: data.distributor || null,
      harvestDate: data.harvestDate || null,
      cannabinoids: cannabinoidsStr,
      terpenes: terpenesStr,
      timesUsed: 0,
      lastUsed: null,
      createdAt: Date.now()
    };

    await add('products', product);

    // Update autocomplete stores
    if (data.producer) await addToAutocomplete('producers', data.producer);
    if (data.processor) await addToAutocomplete('processors', data.processor);
    if (data.distributor) await addToAutocomplete('distributors', data.distributor);

    // Update terpene usage counts (parse compact format)
    for (const terp of parseCompact(terpenesStr)) {
      await incrementTerpeneUsage(terp.name);
    }

    // Update cannabinoid usage counts (parse compact format)
    for (const canna of parseCompact(cannabinoidsStr)) {
      await incrementCannabinoidUsage(canna.name);
    }

    return product;
  }

  async function getProducts() {
    const products = await getAll('products');
    return products.sort((a, b) => (b.lastUsed || 0) - (a.lastUsed || 0));
  }

  // Get product with parsed cannabinoids/terpenes for display
  function getProductDisplay(product) {
    return {
      ...product,
      cannabinoidsParsed: parseCompact(product.cannabinoids),
      terpenesParsed: parseCompact(product.terpenes)
    };
  }

  async function searchProducts(query) {
    const products = await getAll('products');
    const lowerQuery = query.toLowerCase();
    return products.filter(p =>
      p.strainName.toLowerCase().includes(lowerQuery) ||
      (p.producer && p.producer.toLowerCase().includes(lowerQuery))
    );
  }

  // Autocomplete helpers
  async function addToAutocomplete(storeName, name) {
    const existing = await get(storeName, name);
    if (!existing) {
      await add(storeName, { name, createdAt: Date.now() });
    }
  }

  async function getAutocomplete(storeName) {
    const items = await getAll(storeName);
    return items.map(i => i.name).sort();
  }

  async function incrementTerpeneUsage(name) {
    let terp = await get('terpenes', name);
    if (terp) {
      terp.timesUsed = (terp.timesUsed || 0) + 1;
      await update('terpenes', terp);
    } else {
      await add('terpenes', { name, timesUsed: 1, createdAt: Date.now() });
    }
  }

  async function incrementCannabinoidUsage(name) {
    let canna = await get('cannabinoids', name);
    if (canna) {
      canna.timesUsed = (canna.timesUsed || 0) + 1;
      await update('cannabinoids', canna);
    } else {
      await add('cannabinoids', { name, timesUsed: 1, createdAt: Date.now() });
    }
  }

  async function getTerpenes() {
    const terpenes = await getAll('terpenes');
    return terpenes.sort((a, b) => (b.timesUsed || 0) - (a.timesUsed || 0));
  }

  async function getCannabinoids() {
    const cannabinoids = await getAll('cannabinoids');
    return cannabinoids.sort((a, b) => (b.timesUsed || 0) - (a.timesUsed || 0));
  }

  // Preferences
  async function getPreference(key, defaultValue = null) {
    const pref = await get('preferences', key);
    return pref ? pref.value : defaultValue;
  }

  async function setPreference(key, value) {
    await update('preferences', { key, value, updatedAt: Date.now() });
  }

  // Notification scheduling
  async function scheduleNotification(sessionId, scheduledTime) {
    const notification = {
      id: uuid(),
      sessionId,
      scheduledTime,
      status: 'pending',
      createdAt: Date.now()
    };
    await add('pendingNotifications', notification);
    return notification;
  }

  async function getPendingNotifications() {
    const notifications = await getByIndex('pendingNotifications', 'status', 'pending');
    return notifications.filter(n => n.scheduledTime <= Date.now());
  }

  async function markNotificationShown(notificationId) {
    const notification = await get('pendingNotifications', notificationId);
    if (notification) {
      notification.status = 'shown';
      notification.shownAt = Date.now();
      await update('pendingNotifications', notification);
    }
  }

  async function clearSessionNotifications(sessionId) {
    const notifications = await getByIndex('pendingNotifications', 'sessionId', sessionId);
    for (const n of notifications) {
      await remove('pendingNotifications', n.id);
    }
  }

  // Offline Authentication
  // Uses Web Crypto API to hash passwords for secure local storage

  async function hashPassword(password, salt) {
    const encoder = new TextEncoder();
    const data = encoder.encode(password + salt);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
  }

  function generateSalt() {
    const array = new Uint8Array(16);
    crypto.getRandomValues(array);
    return Array.from(array).map(b => b.toString(16).padStart(2, '0')).join('');
  }

  async function saveOfflineCredentials(email, password) {
    await init();
    const salt = generateSalt();
    const passwordHash = await hashPassword(password, salt);

    const credentials = {
      email: email.toLowerCase(),
      passwordHash,
      salt,
      savedAt: Date.now()
    };

    await update('offlineAuth', credentials);
    console.log('[SessionDB] Offline credentials saved for:', email);
    return true;
  }

  async function verifyOfflineCredentials(email, password) {
    await init();
    const credentials = await get('offlineAuth', email.toLowerCase());

    if (!credentials) {
      console.log('[SessionDB] No offline credentials found for:', email);
      return { valid: false, reason: 'no_credentials' };
    }

    const passwordHash = await hashPassword(password, credentials.salt);

    if (passwordHash === credentials.passwordHash) {
      console.log('[SessionDB] Offline login successful for:', email);
      return { valid: true, email: credentials.email };
    }

    console.log('[SessionDB] Offline login failed - wrong password');
    return { valid: false, reason: 'wrong_password' };
  }

  async function hasOfflineCredentials(email) {
    await init();
    const credentials = await get('offlineAuth', email.toLowerCase());
    return !!credentials;
  }

  async function clearOfflineCredentials() {
    await init();
    await clear('offlineAuth');
    console.log('[SessionDB] Offline credentials cleared');
  }

  // Data export
  async function exportData() {
    const data = {
      version: DB_VERSION,
      exportedAt: new Date().toISOString(),
      sessions: await getAll('sessions'),
      checkIns: await getAll('checkIns'),
      products: await getAll('products'),
      preferences: await getAll('preferences')
    };
    return data;
  }

  async function exportCSV() {
    const sessions = await getAll('sessions');
    const products = await getAll('products');
    const checkIns = await getAll('checkIns');

    // Create product lookup
    const productMap = {};
    for (const p of products) {
      productMap[p.id] = p;
    }

    // Build CSV rows
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

  // Public API
  return {
    init,
    seedDefaults,
    uuid,

    // Generic operations
    add,
    get,
    update,
    remove,
    getAll,
    getByIndex,
    clear,

    // Sessions
    createSession,
    getActiveSessions,
    getRecentSessions,
    completeSession,

    // Check-ins
    createCheckIn,
    getSessionCheckIns,

    // Products
    createProduct,
    getProducts,
    getProductDisplay,
    searchProducts,

    // Autocomplete
    getAutocomplete,

    // Terpenes & Cannabinoids
    getTerpenes,
    getCannabinoids,

    // Compact format helpers (for UI)
    parseCompact,
    serializeCompact,

    // Preferences
    getPreference,
    setPreference,

    // Notifications
    scheduleNotification,
    getPendingNotifications,
    markNotificationShown,
    clearSessionNotifications,

    // Export
    exportData,
    exportCSV,

    // Offline Auth
    saveOfflineCredentials,
    verifyOfflineCredentials,
    hasOfflineCredentials,
    clearOfflineCredentials
  };
})();

// Initialize on load
if (typeof window !== 'undefined') {
  window.SessionDB = SessionDB;

  document.addEventListener('DOMContentLoaded', async () => {
    try {
      await SessionDB.init();
      await SessionDB.seedDefaults();
      console.log('[SessionDB] Initialized');
    } catch (err) {
      console.error('[SessionDB] Init failed:', err);
    }
  });
}
