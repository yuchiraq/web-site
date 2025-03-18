document.addEventListener('DOMContentLoaded', function () {
    // Элементы DOM
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
    if (phoneIcon) {
        phoneIcon.addEventListener('click', function (event) {
            event.preventDefault();
            modal.style.display = 'flex';
        });
    }

    // Открытие формы по всем кнопкам с классом open-modal-button
    const openModalButtons = document.querySelectorAll('.open-modal-button');
    openModalButtons.forEach(button => {
        button.addEventListener('click', function (event) {
            event.preventDefault();
            modal.style.display = 'flex';
        });
    });

    // Закрытие формы
    if (closeModal) {
        closeModal.addEventListener('click', function () {
            modal.style.display = 'none';
        });
    }

    // Закрытие формы при клике на размытый фон
    if (modalBackdrop) {
        modalBackdrop.addEventListener('click', function () {
            modal.style.display = 'none';
        });
    }

    // Отправка формы
    if (form) {
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
                    }, 2000);
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
                submitText.style.display = 'inline-block';
                loadingSpinner.style.display = 'none';
                submitButton.disabled = false;
            });
        });
    }

    // Плавная прокрутка к якорям
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        });
    });

    // Ленивая загрузка изображений
    document.querySelectorAll('img[data-src]').forEach(img => {
        const observer = new IntersectionObserver((entries, obs) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    img.setAttribute('src', img.getAttribute('data-src'));
                    img.onload = () => img.removeAttribute('data-src');
                    obs.unobserve(img);
                }
            });
        }, { threshold: 0.1 });
        observer.observe(img);
    });

    // Ripple-эффект при клике на кнопках
    document.querySelectorAll('button, .cta-button, .service-button, .order-button').forEach(button => {
        button.addEventListener('click', function(e) {
            const ripple = document.createElement('span');
            const rect = this.getBoundingClientRect();
            const size = Math.max(rect.width, rect.height);
            const x = e.clientX - rect.left - size / 2;
            const y = e.clientY - rect.top - size / 2;

            ripple.style.cssText = `
                position: absolute;
                background: rgba(255, 255, 255, 0.3);
                border-radius: 50%;
                width: ${size}px;
                height: ${size}px;
                top: ${y}px;
                left: ${x}px;
                transform: scale(0);
                animation: ripple 0.6s linear;
                pointer-events: none;
                z-index: 0;
            `;
            this.appendChild(ripple);

            ripple.addEventListener('animationend', () => ripple.remove());
        });
    });

    // Анимация при скролле (появление элементов)
    const observer = new IntersectionObserver((entries, obs) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('visible');
                obs.unobserve(entry.target);
            }
        });
    }, { threshold: 0.1 });

    document.querySelectorAll('.animate').forEach(element => {
        observer.observe(element);
    });
});

// Анимация для ripple-эффекта
const styleSheet = document.styleSheets[0];
styleSheet.insertRule(`
    @keyframes ripple {
        to {
            transform: scale(2);
            opacity: 0;
        }
    }
`, styleSheet.cssRules.length);

window.addEventListener('scroll', () => {
    const hero = document.querySelector('.hero');
    const scrollTop = window.pageYOffset;
    hero.style.backgroundPositionY = -scrollTop * 0.25 + 'px';
});