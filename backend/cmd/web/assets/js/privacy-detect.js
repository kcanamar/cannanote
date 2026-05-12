// CannaNote Privacy Mode Detection
// Detects private browsing and storage limitations

/**
 * Privacy detection results
 */
const privacyState = {
  isPrivateMode: false,
  storageType: 'unknown',
  browserType: 'unknown',
  checked: false
};

/**
 * Detect browser type
 * @returns {string} Browser identifier
 */
function detectBrowser() {
  const ua = navigator.userAgent.toLowerCase();

  // Check for Brave first (has navigator.brave)
  if ('brave' in navigator) {
    return 'brave';
  }

  // DuckDuckGo browser
  if (ua.includes('duckduckgo')) {
    return 'duckduckgo';
  }

  // Firefox Focus / Klar
  if (ua.includes('focus') || ua.includes('klar')) {
    return 'firefox-focus';
  }

  // Safari private (detected via storage quota)
  if (ua.includes('safari') && !ua.includes('chrome') && !ua.includes('chromium')) {
    return 'safari';
  }

  // Firefox
  if (ua.includes('firefox')) {
    return 'firefox';
  }

  // Chrome/Chromium
  if (ua.includes('chrome') || ua.includes('chromium')) {
    return 'chrome';
  }

  return 'unknown';
}

/**
 * Detect private browsing mode using multiple heuristics
 * @returns {Promise<boolean>}
 */
async function detectPrivateMode() {
  if (privacyState.checked) {
    return privacyState.isPrivateMode;
  }

  privacyState.browserType = detectBrowser();
  let isPrivate = false;

  // Safari: Check storage quota (0 in private mode)
  if ('storage' in navigator && navigator.storage.estimate) {
    try {
      const estimate = await navigator.storage.estimate();
      if (estimate.quota === 0) {
        isPrivate = true;
      } else if (estimate.quota < 120 * 1024 * 1024) {
        // Chrome private has limited quota
        isPrivate = true;
      }
    } catch (e) {
      // Storage estimate not available
    }
  }

  // Firefox: Try IndexedDB (throws in strict private mode)
  if (!isPrivate && privacyState.browserType === 'firefox') {
    try {
      const testRequest = indexedDB.open('__private_test__');
      await new Promise((resolve, reject) => {
        testRequest.onsuccess = () => {
          testRequest.result.close();
          indexedDB.deleteDatabase('__private_test__');
          resolve();
        };
        testRequest.onerror = () => reject(new Error('IndexedDB blocked'));
      });
    } catch (e) {
      isPrivate = true;
    }
  }

  // localStorage test
  try {
    const testKey = '__private_test__';
    localStorage.setItem(testKey, 'test');
    localStorage.removeItem(testKey);
  } catch (e) {
    // localStorage blocked - likely private mode
    isPrivate = true;
  }

  privacyState.isPrivateMode = isPrivate;
  privacyState.checked = true;

  return isPrivate;
}

/**
 * Get storage type being used
 * @returns {Promise<string>}
 */
async function getStorageType() {
  // Check for window.sessionService from session-service.js
  if (typeof window !== 'undefined' && window.sessionService) {
    try {
      await window.sessionService.init();
      privacyState.storageType = window.sessionService.getStorageType();
      return privacyState.storageType;
    } catch (e) {
      // Service not initialized
    }
  }

  // Fallback detection
  if ('storage' in navigator && navigator.storage.getDirectory) {
    try {
      await navigator.storage.getDirectory();
      privacyState.storageType = 'opfs';
      return 'opfs';
    } catch (e) {
      // OPFS not available
    }
  }

  if ('indexedDB' in window) {
    try {
      const request = indexedDB.open('__type_test__', 1);
      await new Promise((resolve, reject) => {
        request.onsuccess = resolve;
        request.onerror = reject;
      });
      request.result.close();
      indexedDB.deleteDatabase('__type_test__');
      privacyState.storageType = 'indexeddb';
      return 'indexeddb';
    } catch (e) {
      // IndexedDB not available
    }
  }

  privacyState.storageType = 'memory';
  return 'memory';
}

/**
 * Get full privacy status
 * @returns {Promise<Object>}
 */
async function getPrivacyStatus() {
  const isPrivate = await detectPrivateMode();
  const storageType = await getStorageType();

  return {
    isPrivateMode: isPrivate,
    storageType: storageType,
    browserType: privacyState.browserType,
    dataPersists: storageType !== 'memory',
    recommendations: getRecommendations(isPrivate, storageType)
  };
}

/**
 * Get recommendations based on privacy status
 * @param {boolean} isPrivate
 * @param {string} storageType
 * @returns {string[]}
 */
function getRecommendations(isPrivate, storageType) {
  const recs = [];

  if (storageType === 'memory') {
    recs.push('Your data will be lost when you close this tab. Consider using regular browsing mode for persistent storage.');
    recs.push('Export your data before closing if you want to save it.');
  }

  if (isPrivate && storageType !== 'memory') {
    recs.push('Private browsing detected. Storage may be limited or cleared when you close the browser.');
  }

  return recs;
}

/**
 * Show privacy banner if needed
 * @param {HTMLElement} container - Container to insert banner
 */
async function showPrivacyBannerIfNeeded(container) {
  const status = await getPrivacyStatus();

  // Only show if using memory storage (data won't persist)
  if (status.storageType !== 'memory') {
    return;
  }

  // Check if banner already shown this session
  const bannerShown = sessionStorage.getItem('privacy_banner_shown');
  if (bannerShown) {
    return;
  }

  const banner = document.createElement('div');
  banner.id = 'privacy-banner';
  banner.className = 'fixed bottom-0 left-0 right-0 bg-amber-50 dark:bg-amber-900/20 border-t border-amber-200 dark:border-amber-800 p-4 z-40';
  banner.innerHTML = `
    <div class="max-w-4xl mx-auto flex items-center justify-between gap-4">
      <div class="flex items-center gap-3">
        <svg class="w-5 h-5 text-amber-600 dark:text-amber-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
        </svg>
        <div>
          <p class="text-sm font-medium text-amber-800 dark:text-amber-200">Private browsing detected</p>
          <p class="text-xs text-amber-600 dark:text-amber-400">Sessions stored temporarily. Export before closing.</p>
        </div>
      </div>
      <button onclick="dismissPrivacyBanner()" class="text-amber-600 dark:text-amber-400 hover:text-amber-800 dark:hover:text-amber-200 p-1" aria-label="Dismiss">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
        </svg>
      </button>
    </div>
  `;

  container.appendChild(banner);
}

/**
 * Dismiss privacy banner
 */
function dismissPrivacyBanner() {
  const banner = document.getElementById('privacy-banner');
  if (banner) {
    banner.remove();
    sessionStorage.setItem('privacy_banner_shown', 'true');
  }
}

/**
 * Show export prompt before leaving in private mode
 * @param {number} sessionCount - Number of sessions to warn about
 */
function showExportPromptIfNeeded(sessionCount) {
  if (privacyState.storageType !== 'memory' || sessionCount === 0) {
    return;
  }

  const message = `You have ${sessionCount} session${sessionCount > 1 ? 's' : ''} that will be lost. Export before closing?`;

  // Using beforeunload for browsers that support it
  window.addEventListener('beforeunload', (e) => {
    if (privacyState.storageType === 'memory' && sessionCount > 0) {
      e.preventDefault();
      e.returnValue = message;
      return message;
    }
  });
}

// Export for module use
export {
  detectPrivateMode,
  detectBrowser,
  getStorageType,
  getPrivacyStatus,
  showPrivacyBannerIfNeeded,
  dismissPrivacyBanner,
  showExportPromptIfNeeded,
  privacyState
};

// Make available globally
if (typeof window !== 'undefined') {
  window.privacyDetect = {
    detectPrivateMode,
    detectBrowser,
    getStorageType,
    getPrivacyStatus,
    showPrivacyBannerIfNeeded,
    dismissPrivacyBanner,
    showExportPromptIfNeeded
  };
  window.dismissPrivacyBanner = dismissPrivacyBanner;
}
