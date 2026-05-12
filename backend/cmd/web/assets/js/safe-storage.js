// CannaNote Safe Storage
// localStorage wrapper with graceful degradation for private browsing

/**
 * Safe localStorage wrapper
 * Falls back to in-memory storage if localStorage is unavailable
 */
const safeStorage = (function() {
  // In-memory fallback
  const memoryStorage = {};
  let useMemory = false;
  let checkedAvailability = false;

  /**
   * Check if localStorage is available
   * @returns {boolean}
   */
  function isAvailable() {
    if (checkedAvailability) {
      return !useMemory;
    }

    checkedAvailability = true;

    try {
      const testKey = '__cannanote_storage_test__';
      localStorage.setItem(testKey, 'test');
      localStorage.removeItem(testKey);
      return true;
    } catch (e) {
      useMemory = true;
      console.warn('[SafeStorage] localStorage unavailable, using memory fallback');
      return false;
    }
  }

  /**
   * Get item from storage
   * @param {string} key - Storage key
   * @returns {string|null}
   */
  function get(key) {
    if (!isAvailable()) {
      return memoryStorage[key] || null;
    }

    try {
      return localStorage.getItem(key);
    } catch (e) {
      console.warn('[SafeStorage] Read failed:', e.message);
      return memoryStorage[key] || null;
    }
  }

  /**
   * Set item in storage
   * @param {string} key - Storage key
   * @param {string} value - Value to store
   * @returns {boolean} Success
   */
  function set(key, value) {
    // Always store in memory as backup
    memoryStorage[key] = value;

    if (!isAvailable()) {
      return true; // Memory storage succeeded
    }

    try {
      localStorage.setItem(key, value);
      return true;
    } catch (e) {
      // QuotaExceededError or other storage errors
      console.warn('[SafeStorage] Write failed:', e.message);
      return false; // Memory storage still has the value
    }
  }

  /**
   * Remove item from storage
   * @param {string} key - Storage key
   * @returns {boolean} Success
   */
  function remove(key) {
    delete memoryStorage[key];

    if (!isAvailable()) {
      return true;
    }

    try {
      localStorage.removeItem(key);
      return true;
    } catch (e) {
      console.warn('[SafeStorage] Remove failed:', e.message);
      return false;
    }
  }

  /**
   * Clear all items from storage
   * @returns {boolean} Success
   */
  function clear() {
    Object.keys(memoryStorage).forEach(key => delete memoryStorage[key]);

    if (!isAvailable()) {
      return true;
    }

    try {
      localStorage.clear();
      return true;
    } catch (e) {
      console.warn('[SafeStorage] Clear failed:', e.message);
      return false;
    }
  }

  /**
   * Get storage status for debugging
   * @returns {Object}
   */
  function getStatus() {
    return {
      localStorageAvailable: isAvailable(),
      usingMemoryFallback: useMemory,
      memoryItemCount: Object.keys(memoryStorage).length
    };
  }

  return {
    get,
    set,
    remove,
    clear,
    isAvailable,
    getStatus
  };
})();

// Export for module use
if (typeof module !== 'undefined' && module.exports) {
  module.exports = { safeStorage };
}

// Make available globally
if (typeof window !== 'undefined') {
  window.safeStorage = safeStorage;
}
