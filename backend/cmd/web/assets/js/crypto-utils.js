// CannaNote Crypto Utilities
// WebCrypto with graceful degradation for privacy browsers

/**
 * Check if WebCrypto is available
 * Some privacy browsers restrict or block crypto.subtle
 */
function hasWebCrypto() {
  try {
    return typeof crypto !== 'undefined'
      && typeof crypto.subtle !== 'undefined'
      && typeof crypto.getRandomValues === 'function';
  } catch (e) {
    return false;
  }
}

/**
 * Check if we're using fallback (insecure) crypto
 * Used to show warnings to users
 */
let usingFallbackCrypto = false;

function isUsingFallbackCrypto() {
  return usingFallbackCrypto;
}

/**
 * Generate cryptographically random bytes
 * Falls back to Math.random if crypto.getRandomValues unavailable
 *
 * @param {number} length - Number of bytes
 * @returns {Uint8Array}
 */
function getRandomBytes(length) {
  const array = new Uint8Array(length);

  if (typeof crypto !== 'undefined' && crypto.getRandomValues) {
    crypto.getRandomValues(array);
  } else {
    usingFallbackCrypto = true;
    console.warn('[Crypto] Using Math.random fallback - not cryptographically secure');
    for (let i = 0; i < length; i++) {
      array[i] = Math.floor(Math.random() * 256);
    }
  }

  return array;
}

/**
 * Generate a random salt for password hashing
 * @returns {string} Hex-encoded salt
 */
function generateSalt() {
  const bytes = getRandomBytes(16);
  return Array.from(bytes).map(b => b.toString(16).padStart(2, '0')).join('');
}

/**
 * Hash a password with salt using SHA-256
 * Falls back to simple hash if WebCrypto unavailable
 *
 * @param {string} password - Password to hash
 * @param {string} salt - Salt to use
 * @returns {Promise<string>} Hex-encoded hash
 */
async function hashPassword(password, salt) {
  // Try WebCrypto first
  if (hasWebCrypto()) {
    try {
      const encoder = new TextEncoder();
      const data = encoder.encode(password + salt);
      const hashBuffer = await crypto.subtle.digest('SHA-256', data);
      const hashArray = Array.from(new Uint8Array(hashBuffer));
      return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
    } catch (err) {
      console.warn('[Crypto] WebCrypto hash failed, using fallback:', err.message);
    }
  }

  // Fallback: Simple hash (NOT cryptographically secure)
  // This allows offline auth to work in restricted environments
  // but users should be warned about reduced security
  usingFallbackCrypto = true;
  console.warn('[Crypto] Using simple hash fallback - reduced security');
  return simpleHash(password + salt);
}

/**
 * Simple hash function for fallback
 * NOT cryptographically secure - use only when WebCrypto unavailable
 *
 * @param {string} str - String to hash
 * @returns {string} Hex hash
 */
function simpleHash(str) {
  let hash1 = 0;
  let hash2 = 0;

  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i);
    hash1 = ((hash1 << 5) - hash1) + char;
    hash1 = hash1 & hash1;
    hash2 = ((hash2 << 7) + hash2) ^ char;
    hash2 = hash2 & hash2;
  }

  // Combine two hashes for slightly better distribution
  const combined = Math.abs(hash1).toString(16).padStart(8, '0') +
                   Math.abs(hash2).toString(16).padStart(8, '0');
  return combined;
}

/**
 * Generate a UUID v4
 * Uses crypto.getRandomValues if available
 *
 * @returns {string} UUID string
 */
function generateUUID() {
  if (typeof crypto !== 'undefined' && crypto.randomUUID) {
    return crypto.randomUUID();
  }

  // Fallback UUID generation
  const bytes = getRandomBytes(16);

  // Set version (4) and variant (RFC 4122)
  bytes[6] = (bytes[6] & 0x0f) | 0x40;
  bytes[8] = (bytes[8] & 0x3f) | 0x80;

  const hex = Array.from(bytes).map(b => b.toString(16).padStart(2, '0')).join('');
  return [
    hex.slice(0, 8),
    hex.slice(8, 12),
    hex.slice(12, 16),
    hex.slice(16, 20),
    hex.slice(20, 32)
  ].join('-');
}

/**
 * Get crypto status for display
 * @returns {Object} Crypto status info
 */
function getCryptoStatus() {
  return {
    webCryptoAvailable: hasWebCrypto(),
    usingFallback: usingFallbackCrypto,
    randomAvailable: typeof crypto !== 'undefined' && typeof crypto.getRandomValues === 'function',
    uuidAvailable: typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function'
  };
}

// Export for module use
export {
  hasWebCrypto,
  isUsingFallbackCrypto,
  getRandomBytes,
  generateSalt,
  hashPassword,
  simpleHash,
  generateUUID,
  getCryptoStatus
};

// Also make available globally for non-module scripts
if (typeof window !== 'undefined') {
  window.CryptoUtils = {
    hasWebCrypto,
    isUsingFallbackCrypto,
    getRandomBytes,
    generateSalt,
    hashPassword,
    simpleHash,
    generateUUID,
    getCryptoStatus
  };
}
