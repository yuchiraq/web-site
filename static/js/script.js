document.addEventListener('DOMContentLoaded', function() {
    const menuToggle = document.querySelector('.menu-toggle');
    const nav = document.querySelector('.nav');

    menuToggle.addEventListener('click', function() {
        nav.classList.toggle('active');
        menuToggle.classList.toggle('active');
    });

    // Автоматическое переключение темы
    const prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)');
    const theme = prefersDarkScheme.matches ? 'dark' : 'light';
    document.documentElement.setAttribute('data-theme', theme);
});