// Check for saved theme preference, otherwise use system preference
const getPreferredTheme = () => {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme) {
        return savedTheme;
    }
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
};

// Apply theme to document
const applyTheme = (theme) => {
    console.log('Applying theme:', theme);
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
};

// Handle toggle click
const handleThemeToggle = () => {
    console.log('Toggle clicked');
    const currentTheme = document.documentElement.getAttribute('data-theme');
    console.log('Current theme:', currentTheme);
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    console.log('New theme:', newTheme);
    applyTheme(newTheme);
};

// Initialize theme
const initializeTheme = () => {
    const themeToggle = document.getElementById('theme-toggle');
    if (!themeToggle) {
        console.error('Theme toggle button not found');
        return;
    }
    
    // Apply initial theme
    const initialTheme = getPreferredTheme();
    console.log('Initial theme:', initialTheme);
    applyTheme(initialTheme);
    
    // Handle toggle click
    themeToggle.onclick = handleThemeToggle;
};

// Make sure we expose the initialization function globally
window.initializeTheme = initializeTheme;

// Run initialization when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializeTheme);
} else {
    initializeTheme();
} 