/**
 * CSRF Protection Helper
 * Automatically injects CSRF tokens into forms and AJAX requests
 */
(function() {
	'use strict';

	const CSRF_COOKIE_NAME = 'cannanote_csrf';
	const CSRF_FIELD_NAME = 'csrf_token';
	const CSRF_HEADER_NAME = 'X-CSRF-Token';

	/**
	 * Get CSRF token from cookie
	 */
	function getCSRFToken() {
		const cookies = document.cookie.split(';');
		for (const cookie of cookies) {
			const [name, value] = cookie.trim().split('=');
			if (name === CSRF_COOKIE_NAME) {
				return decodeURIComponent(value);
			}
		}
		return null;
	}

	/**
	 * Inject CSRF token into a form
	 */
	function injectCSRFToken(form) {
		// Skip if already has CSRF field
		if (form.querySelector(`input[name="${CSRF_FIELD_NAME}"]`)) {
			return;
		}

		// Only inject for POST/PUT/DELETE forms
		const method = (form.method || 'GET').toUpperCase();
		if (method === 'GET' || method === 'HEAD') {
			return;
		}

		const token = getCSRFToken();
		if (!token) {
			console.warn('[CSRF] No CSRF token found in cookies');
			return;
		}

		const input = document.createElement('input');
		input.type = 'hidden';
		input.name = CSRF_FIELD_NAME;
		input.value = token;
		form.appendChild(input);
	}

	/**
	 * Initialize CSRF protection on page load
	 */
	function init() {
		// Inject CSRF tokens into all existing forms
		document.querySelectorAll('form').forEach(injectCSRFToken);

		// Watch for dynamically added forms
		const observer = new MutationObserver((mutations) => {
			mutations.forEach((mutation) => {
				mutation.addedNodes.forEach((node) => {
					if (node.nodeType === Node.ELEMENT_NODE) {
						if (node.tagName === 'FORM') {
							injectCSRFToken(node);
						}
						// Also check children of added nodes
						node.querySelectorAll?.('form').forEach(injectCSRFToken);
					}
				});
			});
		});

		observer.observe(document.body, {
			childList: true,
			subtree: true
		});

		// Add CSRF header to fetch requests
		const originalFetch = window.fetch;
		window.fetch = function(url, options = {}) {
			// Only add CSRF header for same-origin requests with mutating methods
			if (typeof url === 'string' && url.startsWith('/')) {
				const method = (options.method || 'GET').toUpperCase();
				if (method !== 'GET' && method !== 'HEAD') {
					const token = getCSRFToken();
					if (token) {
						options.headers = options.headers || {};
						if (options.headers instanceof Headers) {
							options.headers.set(CSRF_HEADER_NAME, token);
						} else {
							options.headers[CSRF_HEADER_NAME] = token;
						}
					}
				}
			}
			return originalFetch.call(this, url, options);
		};

		// Add CSRF header to XMLHttpRequest
		const originalOpen = XMLHttpRequest.prototype.open;
		const originalSend = XMLHttpRequest.prototype.send;

		XMLHttpRequest.prototype.open = function(method, url, ...args) {
			this._csrfMethod = method;
			this._csrfUrl = url;
			return originalOpen.apply(this, [method, url, ...args]);
		};

		XMLHttpRequest.prototype.send = function(body) {
			const method = (this._csrfMethod || 'GET').toUpperCase();
			const url = this._csrfUrl || '';

			// Only add CSRF header for same-origin mutating requests
			if (typeof url === 'string' && url.startsWith('/')) {
				if (method !== 'GET' && method !== 'HEAD') {
					const token = getCSRFToken();
					if (token) {
						this.setRequestHeader(CSRF_HEADER_NAME, token);
					}
				}
			}
			return originalSend.call(this, body);
		};
	}

	// Expose for manual use
	window.CSRF = {
		getToken: getCSRFToken,
		injectToken: injectCSRFToken,
		HEADER_NAME: CSRF_HEADER_NAME,
		FIELD_NAME: CSRF_FIELD_NAME
	};

	/**
	 * Configure HTMX to include CSRF token in headers
	 */
	function configureHTMX() {
		// Wait for HTMX to load
		if (typeof htmx === 'undefined') {
			return;
		}

		// Add CSRF token to all HTMX requests
		document.body.addEventListener('htmx:configRequest', function(event) {
			const token = getCSRFToken();
			if (token) {
				event.detail.headers[CSRF_HEADER_NAME] = token;
			}
		});
	}

	// Initialize on DOM ready
	if (document.readyState === 'loading') {
		document.addEventListener('DOMContentLoaded', () => {
			init();
			// Configure HTMX after it loads (it has defer attribute)
			setTimeout(configureHTMX, 100);
		});
	} else {
		init();
		setTimeout(configureHTMX, 100);
	}
})();
