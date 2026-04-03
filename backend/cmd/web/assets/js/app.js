// CannaNote - Minimal Vanilla JavaScript
// Handles dark mode, mobile menu, and form validation without frameworks

// Theme Management
function initTheme() {
  const theme = localStorage.getItem('theme');
  const prefersDark = !theme && window.matchMedia('(prefers-color-scheme: dark)').matches;
  const isDark = theme === 'dark' || prefersDark;
  
  if (isDark) {
    document.documentElement.classList.add('dark');
  }
  
  updateThemeIcons(isDark);
}

function toggleTheme() {
  document.documentElement.classList.toggle('dark');
  const isDark = document.documentElement.classList.contains('dark');
  localStorage.theme = isDark ? 'dark' : 'light';
  updateThemeIcons(isDark);
}

function updateThemeIcons(isDark) {
  // Desktop icons
  const sunIcon = document.getElementById('sun-icon');
  const moonIcon = document.getElementById('moon-icon');
  
  // Mobile icons and text
  const mobileSunIcon = document.getElementById('mobile-sun-icon');
  const mobileMoonIcon = document.getElementById('mobile-moon-icon');
  const mobileDarkText = document.getElementById('mobile-dark-text');
  const mobileLightText = document.getElementById('mobile-light-text');
  
  if (isDark) {
    // Show moon (dark mode is ON)
    if (sunIcon) sunIcon.classList.add('hidden');
    if (moonIcon) moonIcon.classList.remove('hidden');
    if (mobileSunIcon) mobileSunIcon.classList.add('hidden');
    if (mobileMoonIcon) mobileMoonIcon.classList.remove('hidden');
    if (mobileDarkText) mobileDarkText.classList.add('hidden');
    if (mobileLightText) mobileLightText.classList.remove('hidden');
  } else {
    // Show sun (light mode is ON)
    if (sunIcon) sunIcon.classList.remove('hidden');
    if (moonIcon) moonIcon.classList.add('hidden');
    if (mobileSunIcon) mobileSunIcon.classList.remove('hidden');
    if (mobileMoonIcon) mobileMoonIcon.classList.add('hidden');
    if (mobileDarkText) mobileDarkText.classList.remove('hidden');
    if (mobileLightText) mobileLightText.classList.add('hidden');
  }
}

// Mobile Menu Management
function toggleMobileMenu() {
  const mobileMenu = document.getElementById('mobile-menu');
  const menuToggle = document.getElementById('mobile-menu-toggle');
  const menuOpenIcon = document.getElementById('menu-open-icon');
  const menuCloseIcon = document.getElementById('menu-close-icon');

  if (mobileMenu) {
    const isHidden = mobileMenu.classList.toggle('hidden');

    // Update aria-expanded state (WCAG 4.1.2)
    if (menuToggle) {
      menuToggle.setAttribute('aria-expanded', !isHidden);
    }
  }

  if (menuOpenIcon && menuCloseIcon) {
    menuOpenIcon.classList.toggle('hidden');
    menuCloseIcon.classList.toggle('hidden');
  }
}

// Form Validation
function validateBetaForm() {
  const emailInput = document.getElementById('email-input');
  const consentCheckbox = document.getElementById('consent-checkbox');
  const submitButton = document.getElementById('submit-button');
  
  if (emailInput && consentCheckbox && submitButton) {
    const isValid = emailInput.value.trim() !== '' && consentCheckbox.checked;
    submitButton.disabled = !isValid;
  }
}

// Click-away handler for mobile menu
function handleDocumentClick(event) {
  const mobileMenu = document.getElementById('mobile-menu');
  const mobileMenuToggle = document.getElementById('mobile-menu-toggle');

  if (mobileMenu && !mobileMenu.classList.contains('hidden')) {
    // Check if click was outside menu and toggle button
    if (!mobileMenu.contains(event.target) && !mobileMenuToggle?.contains(event.target)) {
      mobileMenu.classList.add('hidden');

      // Reset menu icons
      const menuOpenIcon = document.getElementById('menu-open-icon');
      const menuCloseIcon = document.getElementById('menu-close-icon');
      if (menuOpenIcon && menuCloseIcon) {
        menuOpenIcon.classList.remove('hidden');
        menuCloseIcon.classList.add('hidden');
      }

      // Reset aria-expanded state (WCAG 4.1.2)
      if (mobileMenuToggle) {
        mobileMenuToggle.setAttribute('aria-expanded', 'false');
      }
    }
  }
}

// Service Worker Update Handling
function initServiceWorker() {
  if (!('serviceWorker' in navigator)) return;

  navigator.serviceWorker.register('/sw.js')
    .then((registration) => {
      console.log('[App] Service worker registered');

      // Check for updates periodically (every 5 minutes)
      setInterval(() => {
        registration.update();
      }, 5 * 60 * 1000);

      // Handle updates
      registration.addEventListener('updatefound', () => {
        const newWorker = registration.installing;
        console.log('[App] New service worker installing...');

        newWorker.addEventListener('statechange', () => {
          if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
            // New SW installed but waiting - show update notification
            showUpdateNotification();
          }
        });
      });
    })
    .catch((error) => {
      console.log('[App] Service worker registration failed:', error);
    });

  // Listen for messages from SW
  navigator.serviceWorker.addEventListener('message', (event) => {
    if (event.data.type === 'SW_UPDATED') {
      console.log(`[App] Service worker updated to version ${event.data.version}`);
      // Content will be fresh on next navigation
    }
  });

  // Handle controller change (new SW took over)
  navigator.serviceWorker.addEventListener('controllerchange', () => {
    console.log('[App] New service worker activated');
  });
}

// Show update notification to user
function showUpdateNotification() {
  // Check if notification already exists
  if (document.getElementById('sw-update-toast')) return;

  // Add keyframes animation if not already present
  if (!document.getElementById('sw-toast-styles')) {
    const style = document.createElement('style');
    style.id = 'sw-toast-styles';
    style.textContent = `
      @keyframes slideUp {
        from { transform: translateY(100%); opacity: 0; }
        to { transform: translateY(0); opacity: 1; }
      }
      .sw-toast-animate { animation: slideUp 0.3s ease-out; }
    `;
    document.head.appendChild(style);
  }

  const toast = document.createElement('div');
  toast.id = 'sw-update-toast';
  toast.className = 'fixed bottom-4 right-4 bg-sage-600 text-white px-4 py-3 rounded-lg shadow-lg z-50 flex items-center gap-3 sw-toast-animate';
  toast.innerHTML = `
    <svg class="w-5 h-5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
    </svg>
    <span>New content available!</span>
    <button onclick="window.location.reload()" class="bg-white text-sage-700 px-3 py-1 rounded text-sm font-medium hover:bg-sage-50 transition">
      Refresh
    </button>
    <button onclick="dismissUpdateToast()" class="text-sage-200 hover:text-white transition ml-1" aria-label="Dismiss">
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
      </svg>
    </button>
  `;

  document.body.appendChild(toast);
}

function dismissUpdateToast() {
  const toast = document.getElementById('sw-update-toast');
  if (toast) {
    toast.remove();
  }
}

// Initialize everything when DOM is ready
document.addEventListener('DOMContentLoaded', function() {
  // Initialize theme
  initTheme();

  // Register and manage service worker
  initServiceWorker();

  // Add form validation listeners
  const emailInput = document.getElementById('email-input');
  const consentCheckbox = document.getElementById('consent-checkbox');

  if (emailInput) {
    emailInput.addEventListener('input', validateBetaForm);
  }

  if (consentCheckbox) {
    consentCheckbox.addEventListener('change', validateBetaForm);
  }

  // Initial form validation
  validateBetaForm();

  // Add click-away listener for mobile menu
  document.addEventListener('click', handleDocumentClick);

  // Close mobile menu on escape key and return focus to trigger (WCAG 2.1.2)
  document.addEventListener('keydown', function(event) {
    if (event.key === 'Escape') {
      const mobileMenu = document.getElementById('mobile-menu');
      const mobileMenuToggle = document.getElementById('mobile-menu-toggle');
      if (mobileMenu && !mobileMenu.classList.contains('hidden')) {
        toggleMobileMenu();
        // Return focus to trigger button (WCAG focus management)
        if (mobileMenuToggle) {
          mobileMenuToggle.focus();
        }
      }
    }
  });
});