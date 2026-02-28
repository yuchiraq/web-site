document.addEventListener('DOMContentLoaded', function () {
    const header = document.querySelector('.header');
    const phoneIcon = document.getElementById('openPhoneForm');
    const modal = document.getElementById('phoneFormModal');
    const modalBackdrop = document.querySelector('.modal-backdrop');
    const closeModal = document.querySelector('.close');
    const form = document.getElementById('contactForm');
    const responseMessage = document.getElementById('responseMessage');

    const submitButton = form ? form.querySelector('button[type="submit"]') : null;
    const submitText = submitButton ? submitButton.querySelector('.submit-text') : null;
    const loadingSpinner = submitButton ? submitButton.querySelector('.loading-spinner') : null;

    let lastScroll = 0;

    function pushTrackingEvent(eventName, params = {}) {
        window.dataLayer = window.dataLayer || [];
        const payload = Object.assign({
            event: eventName,
            event_category: 'lead',
            page_path: window.location.pathname
        }, params);

        window.dataLayer.push(payload);

        if (typeof window.gtag === 'function') {
            window.gtag('event', eventName, params);
        }

        if (typeof window.ym === 'function') {
            window.ym(100205864, 'reachGoal', eventName, params);
        }
    }

    const prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)');
    document.documentElement.setAttribute('data-theme', prefersDarkScheme.matches ? 'dark' : 'light');

    const hero = document.querySelector('.hero');

    let isScrollScheduled = false;
    function onScroll() {
        const currentScroll = window.scrollY;

        if (header) {
            if (currentScroll > lastScroll && currentScroll > 100) {
                header.classList.add('scrolled');
            } else if (currentScroll < lastScroll) {
                header.classList.remove('scrolled');
            }
        }

        if (hero) {
            if (window.innerWidth <= 768) {
                hero.style.backgroundPosition = 'center';
            } else {
                hero.style.backgroundPosition = `center ${-window.pageYOffset * 0.25}px`;
            }
        }

        lastScroll = currentScroll;
        isScrollScheduled = false;
    }

    window.addEventListener('scroll', () => {
        if (isScrollScheduled) return;
        isScrollScheduled = true;
        window.requestAnimationFrame(onScroll);
    }, { passive: true });

    function openModal() {
        if (modal) {
            modal.style.display = 'flex';
            pushTrackingEvent('order_click', { source: 'modal_open' });
        }
    }

    function closeModalAction() {
        if (modal) {
            modal.style.display = 'none';
        }
    }

    if (phoneIcon) {
        phoneIcon.addEventListener('click', function (event) {
            event.preventDefault();
            openModal();
            pushTrackingEvent('phone_click', { source: 'phone_floating_button' });
        });
    }

    document.querySelectorAll('.open-modal-button').forEach(button => {
        button.addEventListener('click', function (event) {
            event.preventDefault();
            openModal();
            pushTrackingEvent('order_click', {
                source: this.dataset.service || 'unknown_service',
                button_text: this.textContent.trim()
            });
        });
    });

    if (closeModal) closeModal.addEventListener('click', closeModalAction);
    if (modalBackdrop) modalBackdrop.addEventListener('click', closeModalAction);

    if (form) {
        form.addEventListener('submit', function (event) {
            event.preventDefault();

            if (submitText) submitText.style.display = 'none';
            if (loadingSpinner) loadingSpinner.style.display = 'inline-block';
            if (submitButton) submitButton.disabled = true;

            const formData = new FormData(form);
            const name = formData.get('name');
            const phone = formData.get('phone');

            fetch('/submit', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name, phone })
            })
                .then(response => {
                    if (!response.ok) throw new Error('Ошибка сети');
                    return response.json();
                })
                .then(data => {
                    if (data.success) {
                        if (responseMessage) {
                            responseMessage.textContent = 'Спасибо! Мы свяжемся с вами в ближайшее время.';
                            responseMessage.style.color = 'green';
                        }
                        pushTrackingEvent('form_submit', {
                            form_id: 'contactForm',
                            lead_type: 'callback'
                        });
                        form.reset();
                        setTimeout(closeModalAction, 2000);
                    } else if (responseMessage) {
                        responseMessage.textContent = 'Ошибка при отправке. Попробуйте ещё раз.';
                        responseMessage.style.color = 'red';
                    }
                })
                .catch(error => {
                    console.error('Ошибка:', error);
                    if (responseMessage) {
                        responseMessage.textContent = 'Ошибка при отправке. Попробуйте ещё раз.';
                        responseMessage.style.color = 'red';
                    }
                })
                .finally(() => {
                    if (submitText) submitText.style.display = 'inline-block';
                    if (loadingSpinner) loadingSpinner.style.display = 'none';
                    if (submitButton) submitButton.disabled = false;
                });
        });
    }

    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            const target = document.querySelector(this.getAttribute('href'));
            if (!target) return;
            e.preventDefault();
            target.scrollIntoView({ behavior: 'smooth', block: 'start' });
        });
    });

    const lazyImages = document.querySelectorAll('img[data-src]');
    if (lazyImages.length) {
        const imageObserver = new IntersectionObserver((entries, obs) => {
            entries.forEach(entry => {
                if (!entry.isIntersecting) return;
                const img = entry.target;
                img.setAttribute('src', img.getAttribute('data-src'));
                img.onload = () => img.removeAttribute('data-src');
                obs.unobserve(img);
            });
        }, { threshold: 0.1 });

        lazyImages.forEach(img => imageObserver.observe(img));
    }

    document.querySelectorAll('button, .cta-button, .service-button, .order-button').forEach(button => {
        button.addEventListener('click', function (e) {
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

    const animatedElements = document.querySelectorAll('.animate');
    if (animatedElements.length) {
        const animationObserver = new IntersectionObserver((entries, obs) => {
            entries.forEach(entry => {
                if (!entry.isIntersecting) return;
                entry.target.classList.add('visible');
                obs.unobserve(entry.target);
            });
        }, { threshold: 0.1 });

        animatedElements.forEach(element => animationObserver.observe(element));
    }


    const projectImageModal = document.getElementById('projectImageModal');
    const projectImageModalContent = document.getElementById('projectImageModalContent');
    const projectImageCloseButton = document.querySelector('.project-image-close');

    function openProjectImageModal(sourceImage) {
        if (!projectImageModal || !projectImageModalContent || !sourceImage) return;
        projectImageModalContent.src = sourceImage.currentSrc || sourceImage.src;
        projectImageModalContent.alt = sourceImage.alt || 'Увеличенное изображение проекта';
        projectImageModal.classList.add('is-open');
        projectImageModal.setAttribute('aria-hidden', 'false');
        document.body.style.overflow = 'hidden';
    }

    function closeProjectImageModal() {
        if (!projectImageModal || !projectImageModalContent) return;
        projectImageModal.classList.remove('is-open');
        projectImageModal.setAttribute('aria-hidden', 'true');
        projectImageModalContent.src = '';
        document.body.style.overflow = '';
    }

    document.querySelectorAll('.project-carousel-image').forEach(image => {
        if (image.getAttribute('aria-hidden') === 'true') return;
        image.addEventListener('click', () => openProjectImageModal(image));
        image.addEventListener('keydown', event => {
            if (event.key === 'Enter' || event.key === ' ') {
                event.preventDefault();
                openProjectImageModal(image);
            }
        });
    });

    if (projectImageCloseButton) {
        projectImageCloseButton.addEventListener('click', closeProjectImageModal);
    }

    if (projectImageModal) {
        projectImageModal.addEventListener('click', event => {
            if (event.target.closest('[data-close-project-modal="true"]')) {
                closeProjectImageModal();
            }
        });
    }

    document.addEventListener('keydown', event => {
        if (event.key === 'Escape' && projectImageModal && projectImageModal.classList.contains('is-open')) {
            closeProjectImageModal();
        }
    });


    const homepageCarouselTracks = document.querySelectorAll('.projects .carousel-track');

    function setHomepageCarouselSpeed() {
        homepageCarouselTracks.forEach(track => {
            const visibleItemsCount = track.querySelectorAll('.project-carousel-image:not([aria-hidden="true"])').length;
            if (!visibleItemsCount) return;
            const isMobile = window.innerWidth <= 768;
            const baseDuration = isMobile ? 180 : 120;
            const itemFactor = isMobile ? 12 : 8;
            let duration = Math.max(baseDuration, Math.round(visibleItemsCount * itemFactor));
            if (track.closest('.carousel-row-secondary')) {
                duration += isMobile ? 24 : 20;
            }
            if (isMobile) {
                duration *= 3;
            }
            track.style.setProperty('--carousel-duration', `${duration}s`);
        });
    }

    function refreshHomepageCarousel() {
        if (!homepageCarouselTracks.length) return;
        setHomepageCarouselSpeed();
        homepageCarouselTracks.forEach(track => {
            const style = window.getComputedStyle(track);
            if (style.display === 'none') return;
            track.style.animation = 'none';
            void track.offsetHeight;
            track.style.animation = '';
        });
    }

    if (homepageCarouselTracks.length) {
        refreshHomepageCarousel();
        window.addEventListener('pageshow', refreshHomepageCarousel);
        window.addEventListener('orientationchange', refreshHomepageCarousel);
        window.addEventListener('resize', refreshHomepageCarousel);
        let carouselRefreshTimeout;
        window.addEventListener('scroll', function () {
            if (window.innerWidth > 768) return;
            clearTimeout(carouselRefreshTimeout);
            carouselRefreshTimeout = setTimeout(refreshHomepageCarousel, 120);
        }, { passive: true });
    }

    document.querySelectorAll('[data-track]').forEach(element => {
        element.addEventListener('click', function () {
            const eventName = this.dataset.track;
            pushTrackingEvent(eventName, {
                service: this.dataset.service || undefined,
                messenger: this.dataset.messenger || undefined,
                text: this.textContent.trim()
            });
        });
    });

    document.querySelectorAll('.quick-order-form').forEach(formElement => {
        formElement.addEventListener('submit', function (event) {
            event.preventDefault();
        });

        const triggerButton = formElement.querySelector('.open-modal-button');
        if (triggerButton) {
            triggerButton.addEventListener('click', function () {
                const quickName = formElement.querySelector('input[name="name"]');
                const quickPhone = formElement.querySelector('input[name="phone"]');
                if (form && quickName && quickPhone) {
                    const modalName = form.querySelector('#name');
                    const modalPhone = form.querySelector('#phone');
                    if (modalName) modalName.value = quickName.value;
                    if (modalPhone) modalPhone.value = quickPhone.value;
                }
                pushTrackingEvent('form_submit', {
                    form_id: formElement.dataset.trackForm || 'quick_order',
                    lead_type: 'pre_fill'
                });
            });
        }
    });

    document.querySelectorAll('a[href^="tel:"]').forEach(link => {
        link.addEventListener('click', function () {
            pushTrackingEvent('phone_click', { source: 'tel_link' });
        });
    });

    document.querySelectorAll('a[href*="t.me"], a[href*="viber://"]').forEach(link => {
        link.addEventListener('click', function () {
            pushTrackingEvent('messenger_click', { source: this.href });
        });
    });
});

const localStyleSheet = Array.from(document.styleSheets).find(sheet => {
    try {
        return Boolean(sheet.cssRules);
    } catch (error) {
        return false;
    }
});

if (localStyleSheet) {
    localStyleSheet.insertRule(`
        @keyframes ripple {
            to {
                transform: scale(2);
                opacity: 0;
            }
        }
    `, localStyleSheet.cssRules.length);
}
