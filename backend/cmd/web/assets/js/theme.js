// Dark mode initialization script - must run before any content renders
(function() {
	const theme = localStorage.getItem('theme');
	const prefersDark = !theme && window.matchMedia('(prefers-color-scheme: dark)').matches;
	const isDark = theme === 'dark' || prefersDark;
	
	if (isDark) {
		document.documentElement.classList.add('dark');
	}
	
	// Set initial icon visibility after DOM loads
	document.addEventListener('DOMContentLoaded', function() {
		// Desktop icons
		const sunIcon = document.getElementById('sun-icon');
		const moonIcon = document.getElementById('moon-icon');
		// Mobile icons  
		const mobileSunIcon = document.getElementById('mobile-sun-icon');
		const mobileMoonIcon = document.getElementById('mobile-moon-icon');
		const mobileDarkText = document.getElementById('mobile-dark-text');
		const mobileLightText = document.getElementById('mobile-light-text');
		
		if (isDark) {
			// Show moon icon, hide sun icon
			if (sunIcon) sunIcon.classList.add('hidden');
			if (moonIcon) moonIcon.classList.remove('hidden');
			if (mobileSunIcon) mobileSunIcon.classList.add('hidden');
			if (mobileMoonIcon) mobileMoonIcon.classList.remove('hidden');
			if (mobileDarkText) mobileDarkText.classList.add('hidden');
			if (mobileLightText) mobileLightText.classList.remove('hidden');
		} else {
			// Show sun icon, hide moon icon (default state)
			if (sunIcon) sunIcon.classList.remove('hidden');
			if (moonIcon) moonIcon.classList.add('hidden');
			if (mobileSunIcon) mobileSunIcon.classList.remove('hidden');
			if (mobileMoonIcon) mobileMoonIcon.classList.add('hidden');
			if (mobileDarkText) mobileDarkText.classList.remove('hidden');
			if (mobileLightText) mobileLightText.classList.add('hidden');
		}
	});
})();

// Mobile navigation toggle (for docs)
function toggleMobileNav() {
	const overlay = document.getElementById('mobile-nav-overlay');
	const toggle = document.getElementById('mobile-nav-toggle');
	if (overlay) {
		const isHidden = overlay.classList.toggle('hidden');
		// Update aria-expanded state (WCAG 4.1.2)
		if (toggle) {
			toggle.setAttribute('aria-expanded', !isHidden);
		}
	}
}

// Close mobile nav on escape key (WCAG 2.1.2)
document.addEventListener('keydown', function(e) {
	if (e.key === 'Escape') {
		const overlay = document.getElementById('mobile-nav-overlay');
		const toggle = document.getElementById('mobile-nav-toggle');
		if (overlay && !overlay.classList.contains('hidden')) {
			overlay.classList.add('hidden');
			// Reset aria-expanded state
			if (toggle) {
				toggle.setAttribute('aria-expanded', 'false');
				toggle.focus(); // Return focus to trigger
			}
		}
	}
});