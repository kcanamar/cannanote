// CannaNote Service Worker
// Intelligent caching with automatic updates

// Version is injected at runtime by the server
// Local dev: timestamp (changes every build)
// Production: git hash (changes only when code changes)
const SW_VERSION = '{{SW_VERSION}}';
const CACHE_NAME = `cannanote-v${SW_VERSION}`;

// Static assets - these use cache-first (they're fingerprinted/versioned)
const STATIC_ASSETS = [
  '/assets/css/output.css',
  '/assets/js/app.min.js',
  '/assets/js/theme.js',
  '/assets/js/htmx.min.js',
  '/assets/js/session-db.js',
  '/assets/js/session-notifications.js',
  '/assets/images/logos/logo-book.svg',
  '/assets/images/favicon/favicon.svg'
];

// Shell pages to pre-cache for offline access
const SHELL_PAGES = [
  '/',
  '/offline'
];

// Install: cache core assets
self.addEventListener('install', (event) => {
  console.log(`[SW] Installing version ${SW_VERSION}`);
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => {
        // Cache static assets and shell pages
        return cache.addAll([...STATIC_ASSETS, ...SHELL_PAGES]);
      })
      .then(() => {
        console.log(`[SW] Version ${SW_VERSION} installed`);
        // Force activation without waiting for existing tabs to close
        return self.skipWaiting();
      })
  );
});

// Activate: clean up old caches and take control
self.addEventListener('activate', (event) => {
  console.log(`[SW] Activating version ${SW_VERSION}`);
  event.waitUntil(
    caches.keys()
      .then((keys) => {
        return Promise.all(
          keys
            .filter((key) => key.startsWith('cannanote-') && key !== CACHE_NAME)
            .map((key) => {
              console.log(`[SW] Deleting old cache: ${key}`);
              return caches.delete(key);
            })
        );
      })
      .then(() => {
        console.log(`[SW] Version ${SW_VERSION} now active`);
        // Take control of all open tabs immediately
        return self.clients.claim();
      })
      .then(() => {
        // Notify all clients that SW has updated
        return self.clients.matchAll().then((clients) => {
          clients.forEach((client) => {
            client.postMessage({ type: 'SW_UPDATED', version: SW_VERSION });
          });
        });
      })
  );
});

// Fetch: different strategies for different content types
self.addEventListener('fetch', (event) => {
  const url = new URL(event.request.url);

  // Skip non-GET requests
  if (event.request.method !== 'GET') return;

  // Skip external requests (fonts, CDNs, etc.)
  if (url.origin !== self.location.origin) return;

  // Skip API calls - always need fresh data
  if (url.pathname.startsWith('/api/')) return;

  // Skip auth routes - always need fresh
  if (url.pathname.startsWith('/auth/')) return;
  if (url.pathname.startsWith('/login')) return;
  if (url.pathname.startsWith('/activate')) return;
  if (url.pathname.startsWith('/onboarding')) return;

  // Static assets: Cache-first (they're versioned)
  if (isStaticAsset(url.pathname)) {
    event.respondWith(cacheFirst(event.request));
    return;
  }

  // HTML pages (docs, etc.): Stale-while-revalidate
  // Serve cached version immediately, update cache in background
  if (event.request.mode === 'navigate' || url.pathname.startsWith('/docs')) {
    event.respondWith(staleWhileRevalidate(event.request));
    return;
  }

  // Default: Network-first with cache fallback
  event.respondWith(networkFirst(event.request));
});

// Check if URL is a static asset
function isStaticAsset(pathname) {
  return pathname.startsWith('/assets/') ||
         pathname.endsWith('.css') ||
         pathname.endsWith('.js') ||
         pathname.endsWith('.svg') ||
         pathname.endsWith('.png') ||
         pathname.endsWith('.jpg') ||
         pathname.endsWith('.woff2');
}

// Cache-first strategy: serve from cache, fallback to network
async function cacheFirst(request) {
  const cached = await caches.match(request);
  if (cached) {
    return cached;
  }

  try {
    const response = await fetch(request);
    if (response.ok) {
      const cache = await caches.open(CACHE_NAME);
      cache.put(request, response.clone());
    }
    return response;
  } catch (error) {
    // If both cache and network fail, return offline page for navigation
    if (request.mode === 'navigate') {
      return caches.match('/offline');
    }
    throw error;
  }
}

// Network-first strategy: try network, fallback to cache
async function networkFirst(request) {
  try {
    const response = await fetch(request);
    if (response.ok) {
      const cache = await caches.open(CACHE_NAME);
      cache.put(request, response.clone());
    }
    return response;
  } catch (error) {
    const cached = await caches.match(request);
    if (cached) {
      return cached;
    }
    if (request.mode === 'navigate') {
      return caches.match('/offline');
    }
    throw error;
  }
}

// Stale-while-revalidate: serve cached immediately, update in background
async function staleWhileRevalidate(request) {
  const cache = await caches.open(CACHE_NAME);
  const cached = await caches.match(request);

  // Start fetching fresh version in background
  const fetchPromise = fetch(request)
    .then((response) => {
      if (response.ok) {
        cache.put(request, response.clone());
      }
      return response;
    })
    .catch((error) => {
      console.log(`[SW] Background fetch failed for ${request.url}:`, error);
      return null;
    });

  // If we have cached content, return it immediately
  if (cached) {
    // Background update will refresh cache for next visit
    return cached;
  }

  // No cache - wait for network
  const response = await fetchPromise;
  if (response) {
    return response;
  }

  // Network failed and no cache - show offline page
  if (request.mode === 'navigate') {
    return caches.match('/offline');
  }

  throw new Error('No cached version available');
}

// Handle messages from the main thread
self.addEventListener('message', (event) => {
  if (event.data === 'skipWaiting') {
    self.skipWaiting();
  }

  if (event.data === 'clearCache') {
    caches.delete(CACHE_NAME).then(() => {
      console.log(`[SW] Cache cleared`);
      event.ports[0].postMessage({ success: true });
    });
  }

  if (event.data === 'getVersion') {
    event.ports[0].postMessage({ version: SW_VERSION });
  }

  // Handle account deletion - clear all local data
  if (event.data && event.data.type === 'ACCOUNT_DELETED') {
    console.log('[SW] Account deleted - clearing all local data');

    // Delete IndexedDB database
    if ('indexedDB' in self) {
      indexedDB.deleteDatabase('cannanote');
    }

    // Clear all caches
    caches.keys().then((names) => {
      return Promise.all(
        names.map((name) => {
          console.log(`[SW] Deleting cache: ${name}`);
          return caches.delete(name);
        })
      );
    }).then(() => {
      // Notify all clients that deletion is complete
      self.clients.matchAll().then((clients) => {
        clients.forEach((client) => {
          client.postMessage({ type: 'ACCOUNT_DELETED_COMPLETE' });
        });
      });
    });
  }
});

// Handle notification clicks (for session check-in notifications)
self.addEventListener('notificationclick', (event) => {
  console.log('[SW] Notification clicked:', event.notification.tag);

  event.notification.close();

  // Get the URL from notification data
  const url = event.notification.data?.url || '/app/sessions';

  event.waitUntil(
    clients.matchAll({ type: 'window', includeUncontrolled: true })
      .then((clientList) => {
        // Check if there's already a window open
        for (const client of clientList) {
          if (client.url.includes('/app') && 'focus' in client) {
            client.navigate(url);
            return client.focus();
          }
        }
        // No window open, open a new one
        if (clients.openWindow) {
          return clients.openWindow(url);
        }
      })
  );
});
