document.addEventListener('DOMContentLoaded', function () {
    const header = document.querySelector('.header');
    const phoneIcon = document.querySelector('.phone-icon');

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

    // Обработчик для иконки телефона
    phoneIcon.addEventListener('click', function (event) {
        event.preventDefault(); // Предотвращаем переход по ссылке
        const phoneNumber = phoneIcon.getAttribute('href').replace('tel:', '');
        window.location.href = `tel:${phoneNumber}`; // Открываем системный выбор телефона
    });
});