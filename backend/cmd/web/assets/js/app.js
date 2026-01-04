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
  const menuOpenIcon = document.getElementById('menu-open-icon');
  const menuCloseIcon = document.getElementById('menu-close-icon');
  
  if (mobileMenu) {
    mobileMenu.classList.toggle('hidden');
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
  const mobileMenuToggle = document.querySelector('[onclick*="toggleMobileMenu"]');
  
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
    }
  }
}

// Initialize everything when DOM is ready
document.addEventListener('DOMContentLoaded', function() {
  // Initialize theme
  initTheme();
  
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
  
  // Close mobile menu on escape key
  document.addEventListener('keydown', function(event) {
    if (event.key === 'Escape') {
      const mobileMenu = document.getElementById('mobile-menu');
      if (mobileMenu && !mobileMenu.classList.contains('hidden')) {
        toggleMobileMenu();
      }
    }
  });
});