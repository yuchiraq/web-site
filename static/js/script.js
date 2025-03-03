document.addEventListener('DOMContentLoaded', function () {
    const header = document.querySelector('.header');

    let lastScroll = 0;

    // Анимация навигации при скролле
    window.addEventListener('scroll', function () {
        const currentScroll = window.scrollY;

        if (currentScroll > lastScroll && currentScroll > 100) {
            // Скрываем навигацию при скролле вниз
            header.classList.add('scrolled');
        } else if (currentScroll < lastScroll) {
            // Показываем навигацию при скролле вверх
            header.classList.remove('scrolled');
        }

        lastScroll = currentScroll;
    });

    // Переключение темы
    const prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)');
    const theme = prefersDarkScheme.matches ? 'dark' : 'light';
    document.documentElement.setAttribute('data-theme', theme);
});