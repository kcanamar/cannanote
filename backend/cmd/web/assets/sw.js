// CannaNote Service Worker
// Simple cache-first strategy for offline support

const CACHE_NAME = 'cannanote-v1';

// Assets to cache on install
const ASSETS_TO_CACHE = [
  '/',
  '/app',
  '/offline',
  '/assets/css/output.css',
  '/assets/js/app.min.js',
  '/assets/js/theme.js',
  '/assets/js/htmx.min.js',
  '/assets/images/logos/logo-book.svg',
  '/assets/images/favicon/favicon.svg'
];

// Install: cache core assets
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => cache.addAll(ASSETS_TO_CACHE))
      .then(() => self.skipWaiting())
  );
});

// Activate: clean up old caches
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys()
      .then((keys) => Promise.all(
        keys
          .filter((key) => key !== CACHE_NAME)
          .map((key) => caches.delete(key))
      ))
      .then(() => self.clients.claim())
  );
});

// Fetch: cache-first for assets, network-first for pages
self.addEventListener('fetch', (event) => {
  const url = new URL(event.request.url);

  // Skip non-GET requests
  if (event.request.method !== 'GET') return;

  // Skip external requests (fonts, CDNs, etc.) - let browser handle them
  if (url.origin !== self.location.origin) return;

  // Skip API calls - let them go to network
  if (url.pathname.startsWith('/api/')) return;

  // Skip auth routes - always need fresh
  if (url.pathname.startsWith('/auth/')) return;

  event.respondWith(
    caches.match(event.request)
      .then((cached) => {
        if (cached) {
          return cached;
        }

        return fetch(event.request)
          .then((response) => {
            // Don't cache non-success responses
            if (!response || response.status !== 200) {
              return response;
            }

            // Cache successful responses
            const responseToCache = response.clone();
            caches.open(CACHE_NAME)
              .then((cache) => cache.put(event.request, responseToCache));

            return response;
          })
          .catch(() => {
            // Network failed, try to serve offline page for navigation requests
            if (event.request.mode === 'navigate') {
              return caches.match('/offline');
            }
          });
      })
  );
});
