document.addEventListener('DOMContentLoaded', function () {
    const header = document.querySelector('.header');
    const phoneIcon = document.getElementById('openPhoneForm');
    const modal = document.getElementById('phoneFormModal');
    const modalBackdrop = document.querySelector('.modal-backdrop');
    const closeModal = document.querySelector('.close');
    const form = document.getElementById('contactForm');
    const responseMessage = document.getElementById('responseMessage');
    const submitButton = form.querySelector('button[type="submit"]');
    const submitText = submitButton.querySelector('.submit-text');
    const loadingSpinner = submitButton.querySelector('.loading-spinner');

    let lastScroll = 0;

    // Переключение темы
    const prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)');
    const theme = prefersDarkScheme.matches ? 'dark' : 'light';
    document.documentElement.setAttribute('data-theme', theme);

    // Анимация навигации при скролле
    window.addEventListener('scroll', function () {
        const currentScroll = window.scrollY;

        if (currentScroll > lastScroll && currentScroll > 100) {
            header.classList.add('scrolled');
        } else if (currentScroll < lastScroll) {
            header.classList.remove('scrolled');
        }

        lastScroll = currentScroll;
    });

    // Открытие формы по кнопке телефона
    phoneIcon.addEventListener('click', function (event) {
        event.preventDefault();
        modal.style.display = 'flex';
    });

    // Открытие формы по всем кнопкам с классом open-modal-button
    const openModalButtons = document.querySelectorAll('.open-modal-button');
    openModalButtons.forEach(button => {
        button.addEventListener('click', function (event) {
            event.preventDefault(); // Предотвращаем переход по ссылке
            modal.style.display = 'flex';
        });
    });

    // Закрытие формы
    closeModal.addEventListener('click', function () {
        modal.style.display = 'none';
    });

    // Закрытие формы при клике на размытый фон
    modalBackdrop.addEventListener('click', function () {
        modal.style.display = 'none';
    });

    // Отправка формы
    form.addEventListener('submit', function (event) {
        event.preventDefault();

        // Показываем индикатор загрузки
        submitText.style.display = 'none';
        loadingSpinner.style.display = 'inline-block';
        submitButton.disabled = true;

        const formData = new FormData(form);
        const name = formData.get('name');
        const phone = formData.get('phone');

        fetch('/submit', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name, phone }),
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка сети');
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                responseMessage.textContent = 'Спасибо! Мы свяжемся с вами в ближайшее время.';
                responseMessage.style.color = 'green';
                form.reset();
                setTimeout(() => {
                    modal.style.display = 'none';
                }, 2000); // Закрываем форму через 2 секунды
            } else {
                responseMessage.textContent = 'Ошибка при отправке. Попробуйте ещё раз.';
                responseMessage.style.color = 'red';
            }
        })
        .catch(error => {
            console.error('Ошибка:', error);
            responseMessage.textContent = 'Ошибка при отправке. Попробуйте ещё раз.';
            responseMessage.style.color = 'red';
        })
        .finally(() => {
            // Скрываем индикатор загрузки
            submitText.style.display = 'inline-block';
            loadingSpinner.style.display = 'none';
            submitButton.disabled = false;
        });
    });
});