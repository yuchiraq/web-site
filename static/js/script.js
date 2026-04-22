document.addEventListener("DOMContentLoaded", () => {
    const body = document.body;
    const header = document.querySelector(".header");
    const hero = document.querySelector(".hero");
    const navLinks = Array.from(document.querySelectorAll(".nav a"));

    const phoneIcon = document.getElementById("openPhoneForm");
    const modal = document.getElementById("phoneFormModal");
    const modalBackdrop = modal ? modal.querySelector(".modal-backdrop") : null;
    const closeModalButton = modal ? modal.querySelector(".close") : null;
    const form = document.getElementById("contactForm");
    const responseMessage = document.getElementById("responseMessage");
    const submitButton = form ? form.querySelector('button[type="submit"]') : null;
    const submitText = submitButton ? submitButton.querySelector(".submit-text") : null;
    const loadingSpinner = submitButton ? submitButton.querySelector(".loading-spinner") : null;

    const projectImageModal = document.getElementById("projectImageModal");
    const projectImageModalContent = document.getElementById("projectImageModalContent");
    const projectImageCloseButton = document.querySelector(".project-image-close");
    const copyNotification = document.getElementById("copy-notification");

    const prefersDarkScheme = window.matchMedia("(prefers-color-scheme: dark)");
    const prefersReducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)");
    const marketingStorageKey = "avayusstroi_marketing_context";
    const trackingParamKeys = ["utm_source", "utm_medium", "utm_campaign", "utm_content", "utm_term", "gclid", "fbclid", "yclid"];
    document.documentElement.setAttribute("data-theme", prefersDarkScheme.matches ? "dark" : "light");

    let activeLeadContext = {
        source: window.location.pathname,
        ctaText: ""
    };

    const getStoredMarketingContext = () => {
        try {
            const rawValue = window.localStorage.getItem(marketingStorageKey);
            if (!rawValue) return {};

            const parsedValue = JSON.parse(rawValue);
            return parsedValue && typeof parsedValue === "object" ? parsedValue : {};
        } catch (error) {
            return {};
        }
    };

    const resolveMarketingContext = () => {
        const params = new URLSearchParams(window.location.search);
        const storedContext = getStoredMarketingContext();
        const nextContext = { ...storedContext };
        let hasFreshTrackingParams = false;

        trackingParamKeys.forEach(paramKey => {
            const value = params.get(paramKey);
            if (!value) return;

            nextContext[paramKey] = value;
            hasFreshTrackingParams = true;
        });

        if (!nextContext.initialReferrer && document.referrer) {
            nextContext.initialReferrer = document.referrer;
        }

        if (!nextContext.landingPage || hasFreshTrackingParams) {
            nextContext.landingPage = window.location.href;
        }

        try {
            window.localStorage.setItem(marketingStorageKey, JSON.stringify(nextContext));
        } catch (error) {
            // Ignore storage errors and continue with in-memory context.
        }

        return nextContext;
    };

    const marketingContext = resolveMarketingContext();

    const getMarketingContext = () => ({
        landingPage: marketingContext.landingPage || window.location.href,
        referrer: marketingContext.initialReferrer || document.referrer || "",
        utmSource: marketingContext.utm_source || "",
        utmMedium: marketingContext.utm_medium || "",
        utmCampaign: marketingContext.utm_campaign || "",
        utmContent: marketingContext.utm_content || "",
        utmTerm: marketingContext.utm_term || "",
        gclid: marketingContext.gclid || "",
        fbclid: marketingContext.fbclid || "",
        yclid: marketingContext.yclid || ""
    });

    const setActiveLeadContext = (source, ctaText = "") => {
        activeLeadContext = {
            source: source || window.location.pathname,
            ctaText
        };
    };

    const syncScrollLock = () => {
        const shouldLock =
            body.classList.contains("modal-open") ||
            Boolean(projectImageModal && projectImageModal.classList.contains("is-open"));
        body.classList.toggle("scroll-locked", shouldLock);
    };

    const pushTrackingEvent = (eventName, params = {}) => {
        if (!eventName) return;

        window.dataLayer = window.dataLayer || [];
        const payload = {
            event: eventName,
            event_category: "lead",
            page_path: window.location.pathname,
            ...getMarketingContext(),
            ...params
        };

        window.dataLayer.push(payload);

        if (typeof window.gtag === "function") {
            window.gtag("event", eventName, params);
        }

        if (typeof window.ym === "function") {
            window.ym(100205864, "reachGoal", eventName, params);
        }
    };

    if (navLinks.length) {
        const currentPath = window.location.pathname.replace(/\/+$/, "") || "/";
        navLinks.forEach(link => {
            const linkPath = new URL(link.href, window.location.origin).pathname.replace(/\/+$/, "") || "/";
            const isCurrent = linkPath === "/"
                ? currentPath === "/"
                : currentPath === linkPath || currentPath.startsWith(`${linkPath}/`);

            link.classList.toggle("is-current", isCurrent);
            if (isCurrent) {
                link.setAttribute("aria-current", "page");
            } else {
                link.removeAttribute("aria-current");
            }
        });
    }

    if (window.location.pathname === "/") {
        const heroTitle = document.querySelector(".hero h1");
        const heroLead = document.querySelector(".hero p");

        if (heroTitle) {
            heroTitle.innerHTML = 'Автономная канализация и <span class="hero-accent">водопонижение</span> под ключ';
        }

        if (heroLead) {
            heroLead.textContent = "Подбираем решение, выполняем монтаж и сдаем объект по Бресту и области с понятной сметой и гарантией на работы.";
        }
    }

    let scrollTicking = false;
    let lastScroll = 0;
    const updateOnScroll = () => {
        const currentScroll = window.scrollY;

        if (header) {
            if (currentScroll > lastScroll && currentScroll > 100) {
                header.classList.add("scrolled");
            } else if (currentScroll < lastScroll || currentScroll <= 24) {
                header.classList.remove("scrolled");
            }
        }

        if (hero && !prefersReducedMotion.matches && window.innerWidth > 960) {
            hero.style.backgroundPosition = `center ${Math.round(currentScroll * 0.18)}px`;
        } else if (hero) {
            hero.style.backgroundPosition = "center";
        }

        lastScroll = currentScroll;
        scrollTicking = false;
    };

    window.addEventListener("scroll", () => {
        if (scrollTicking) return;
        scrollTicking = true;
        window.requestAnimationFrame(updateOnScroll);
    }, { passive: true });

    updateOnScroll();

    const openModal = () => {
        if (!modal) return;

        modal.classList.add("is-open");
        modal.setAttribute("aria-hidden", "false");
        body.classList.add("modal-open");
        syncScrollLock();
    };

    const closeModal = () => {
        if (!modal) return;

        modal.classList.remove("is-open");
        modal.setAttribute("aria-hidden", "true");
        body.classList.remove("modal-open");
        syncScrollLock();
    };

    if (phoneIcon) {
        phoneIcon.addEventListener("click", event => {
            event.preventDefault();
            setActiveLeadContext("phone_floating_button", "floating_phone");
            openModal();
            pushTrackingEvent("phone_click", { source: "phone_floating_button" });
        });
    }

    document.querySelectorAll(".open-modal-button").forEach(button => {
        button.addEventListener("click", event => {
            event.preventDefault();
            setActiveLeadContext(button.dataset.service || button.dataset.track || "modal", button.textContent.trim());
            openModal();
        });
    });

    if (closeModalButton) {
        closeModalButton.addEventListener("click", closeModal);
    }

    if (modalBackdrop) {
        modalBackdrop.addEventListener("click", closeModal);
    }

    if (form) {
        form.addEventListener("submit", event => {
            event.preventDefault();

            const formData = new FormData(form);
            const name = String(formData.get("name") || "").trim();
            const phone = String(formData.get("phone") || "").trim();

            if (!name || !phone) {
                if (responseMessage) {
                    responseMessage.textContent = "Заполните имя и телефон, пожалуйста.";
                    responseMessage.style.color = "#c13a2f";
                }
                return;
            }

            if (submitText) submitText.style.display = "none";
            if (loadingSpinner) loadingSpinner.style.display = "inline-block";
            if (submitButton) submitButton.disabled = true;

            const marketingPayload = getMarketingContext();

            fetch("/submit", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    name,
                    phone,
                    page: window.location.href,
                    source: activeLeadContext.source,
                    ctaText: activeLeadContext.ctaText,
                    ...marketingPayload
                })
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error("Ошибка сети");
                    }
                    return response.json();
                })
                .then(data => {
                    if (data.success) {
                        if (responseMessage) {
                            responseMessage.textContent = "Спасибо! Мы свяжемся с вами в ближайшее время.";
                            responseMessage.style.color = "#2f7d4a";
                        }

                        pushTrackingEvent("form_submit", {
                            form_id: "contactForm",
                            lead_type: "callback"
                        });

                        form.reset();
                        window.setTimeout(closeModal, 1800);
                        return;
                    }

                    throw new Error("Ошибка отправки");
                })
                .catch(error => {
                    console.error("Ошибка:", error);
                    if (responseMessage) {
                        responseMessage.textContent = "Не удалось отправить заявку. Попробуйте еще раз.";
                        responseMessage.style.color = "#c13a2f";
                    }
                })
                .finally(() => {
                    if (submitText) submitText.style.display = "inline-block";
                    if (loadingSpinner) loadingSpinner.style.display = "none";
                    if (submitButton) submitButton.disabled = false;
                });
        });
    }

    document.querySelectorAll('a[href^="#"]:not([href="#"])').forEach(anchor => {
        anchor.addEventListener("click", event => {
            const targetSelector = anchor.getAttribute("href");
            if (!targetSelector) return;

            const target = document.querySelector(targetSelector);
            if (!target) return;

            event.preventDefault();
            target.scrollIntoView({ behavior: prefersReducedMotion.matches ? "auto" : "smooth", block: "start" });
        });
    });

    document.querySelectorAll("button, .cta-button, .service-button, .order-button, .contact-link, .messenger-badge").forEach(button => {
        button.addEventListener("click", event => {
            const target = event.currentTarget;
            if (!(target instanceof HTMLElement)) return;

            const ripple = document.createElement("span");
            const rect = target.getBoundingClientRect();
            const size = Math.max(rect.width, rect.height);
            const x = event.clientX - rect.left - size / 2;
            const y = event.clientY - rect.top - size / 2;

            ripple.style.cssText = `
                position: absolute;
                background: rgba(255, 255, 255, 0.26);
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

            target.appendChild(ripple);
            ripple.addEventListener("animationend", () => ripple.remove());
        });
    });

    const lazyImages = document.querySelectorAll("img[data-src]");
    if (lazyImages.length && "IntersectionObserver" in window) {
        const imageObserver = new IntersectionObserver((entries, observer) => {
            entries.forEach(entry => {
                if (!entry.isIntersecting) return;
                const img = entry.target;
                img.setAttribute("src", img.getAttribute("data-src"));
                img.onload = () => img.removeAttribute("data-src");
                observer.unobserve(img);
            });
        }, { threshold: 0.1 });

        lazyImages.forEach(img => imageObserver.observe(img));
    }

    const seasonalImages = document.querySelectorAll("[data-summer-src][data-winter-src]");
    if (seasonalImages.length) {
        const month = new Date().getMonth() + 1;
        const isWinter = month === 12 || month <= 2;
        seasonalImages.forEach(image => {
            image.src = isWinter ? image.dataset.winterSrc : image.dataset.summerSrc;
        });
    }

    const animatedElements = document.querySelectorAll(".animate");
    if (animatedElements.length && "IntersectionObserver" in window) {
        const animationObserver = new IntersectionObserver((entries, observer) => {
            entries.forEach(entry => {
                if (!entry.isIntersecting) return;
                entry.target.classList.add("visible");
                observer.unobserve(entry.target);
            });
        }, { threshold: 0.1 });

        animatedElements.forEach(element => animationObserver.observe(element));
    }

    const openProjectImageModal = sourceImage => {
        if (!projectImageModal || !projectImageModalContent || !sourceImage) return;

        projectImageModalContent.src = sourceImage.currentSrc || sourceImage.src;
        projectImageModalContent.alt = sourceImage.alt || "Увеличенное изображение проекта";
        projectImageModal.classList.add("is-open");
        projectImageModal.setAttribute("aria-hidden", "false");
        syncScrollLock();
    };

    const closeProjectImageModal = () => {
        if (!projectImageModal || !projectImageModalContent) return;

        projectImageModal.classList.remove("is-open");
        projectImageModal.setAttribute("aria-hidden", "true");
        projectImageModalContent.src = "";
        syncScrollLock();
    };

    document.querySelectorAll(".project-carousel-image").forEach(image => {
        if (image.getAttribute("aria-hidden") === "true") return;

        image.addEventListener("click", () => openProjectImageModal(image));
        image.addEventListener("keydown", event => {
            if (event.key === "Enter" || event.key === " ") {
                event.preventDefault();
                openProjectImageModal(image);
            }
        });
    });

    if (projectImageCloseButton) {
        projectImageCloseButton.addEventListener("click", closeProjectImageModal);
    }

    if (projectImageModal) {
        projectImageModal.addEventListener("click", event => {
            const target = event.target;
            if (target instanceof Element && target.closest("[data-close-project-modal='true']")) {
                closeProjectImageModal();
            }
        });
    }

    const homepageCarouselTracks = document.querySelectorAll(".projects .carousel-track");

    const setHomepageCarouselSpeed = () => {
        homepageCarouselTracks.forEach(track => {
            const visibleItemsCount = track.querySelectorAll(".project-carousel-image:not([aria-hidden='true'])").length;
            if (!visibleItemsCount) return;

            const isMobile = window.innerWidth <= 768;
            const perItemDuration = isMobile ? 0.55 : 1.05;
            const minDuration = isMobile ? 7 : 16;
            const maxDuration = isMobile ? 14 : 38;
            let duration = Math.round(visibleItemsCount * perItemDuration);
            duration = Math.max(minDuration, Math.min(maxDuration, duration));

            if (track.closest(".carousel-row-secondary") && !isMobile) {
                duration += 1;
            }

            if (isMobile) {
                duration = Math.max(7, Math.min(14, Math.round(visibleItemsCount * 0.55)));
            }

            track.style.setProperty("--carousel-duration", `${duration}s`);
        });
    };

    const refreshHomepageCarousel = () => {
        if (!homepageCarouselTracks.length || prefersReducedMotion.matches) return;

        setHomepageCarouselSpeed();
        homepageCarouselTracks.forEach(track => {
            const style = window.getComputedStyle(track);
            if (style.display === "none") return;

            track.style.animation = "none";
            void track.offsetHeight;
            track.style.animation = "";
        });
    };

    if (homepageCarouselTracks.length) {
        refreshHomepageCarousel();
        window.addEventListener("pageshow", refreshHomepageCarousel);
        window.addEventListener("orientationchange", refreshHomepageCarousel);
        window.addEventListener("resize", refreshHomepageCarousel);
    }

    document.querySelectorAll("[data-track]").forEach(element => {
        element.addEventListener("click", function () {
            pushTrackingEvent(this.dataset.track, {
                service: this.dataset.service || undefined,
                messenger: this.dataset.messenger || undefined,
                text: this.textContent.trim()
            });
        });
    });

    document.querySelectorAll(".quick-order-form").forEach(formElement => {
        formElement.addEventListener("submit", event => {
            event.preventDefault();
        });

        const triggerButton = formElement.querySelector(".open-modal-button");
        if (!triggerButton) return;

        triggerButton.addEventListener("click", () => {
            const quickName = formElement.querySelector('input[name="name"]');
            const quickPhone = formElement.querySelector('input[name="phone"]');

            if (form && quickName && quickPhone) {
                const modalName = form.querySelector("#name");
                const modalPhone = form.querySelector("#phone");
                if (modalName) modalName.value = quickName.value;
                if (modalPhone) modalPhone.value = quickPhone.value;
            }

            pushTrackingEvent("form_submit", {
                form_id: formElement.dataset.trackForm || "quick_order",
                lead_type: "pre_fill"
            });
        });
    });

    document.querySelectorAll('a[href^="tel:"]').forEach(link => {
        link.addEventListener("click", function () {
            if (this.dataset.track) return;
            pushTrackingEvent("phone_click", { source: "tel_link" });
        });
    });

    document.querySelectorAll('a[href*="t.me"], a[href*="viber://"]').forEach(link => {
        link.addEventListener("click", function () {
            if (this.dataset.track) return;
            pushTrackingEvent("messenger_click", { source: this.href });
        });
    });

    const clientCounters = document.querySelectorAll(".trusted-clients .counter[data-target]");
    if (clientCounters.length && "IntersectionObserver" in window) {
        const animateCounter = counter => {
            const target = Number(counter.dataset.target) || 0;
            const stepSize = 1;
            const intervalMs = 100;
            let currentValue = 0;

            if (prefersReducedMotion.matches) {
                counter.textContent = String(target);
                return;
            }

            counter.textContent = "0";

            const timerId = window.setInterval(() => {
                currentValue = Math.min(currentValue + stepSize, target);
                counter.textContent = String(currentValue);

                if (currentValue >= target) {
                    window.clearInterval(timerId);
                }
            }, intervalMs);
        };

        const counterObserver = new IntersectionObserver((entries, observer) => {
            entries.forEach(entry => {
                if (!entry.isIntersecting) return;
                animateCounter(entry.target);
                observer.unobserve(entry.target);
            });
        }, { threshold: 0.45 });

        clientCounters.forEach(counter => counterObserver.observe(counter));
    }

    const showCopyNotification = () => {
        if (!copyNotification) return;

        copyNotification.classList.add("is-visible");
        window.clearTimeout(showCopyNotification.timeoutId);
        showCopyNotification.timeoutId = window.setTimeout(() => {
            copyNotification.classList.remove("is-visible");
        }, 1800);
    };

    const getCopyCardText = card => {
        const clone = card.cloneNode(true);
        clone.querySelectorAll(".contact-actions, .contact-copy-hint, .messenger-list").forEach(element => element.remove());
        return clone.innerText.replace(/\n{3,}/g, "\n\n").trim();
    };

    const copyText = async text => {
        if (!text) return false;

        if (navigator.clipboard && window.isSecureContext) {
            try {
                await navigator.clipboard.writeText(text);
                return true;
            } catch (error) {
                // Fall back to a manual copy strategy below.
            }
        }

        const textArea = document.createElement("textarea");
        textArea.value = text;
        textArea.setAttribute("readonly", "");
        textArea.style.position = "fixed";
        textArea.style.top = "0";
        textArea.style.left = "0";
        textArea.style.opacity = "0";
        textArea.style.pointerEvents = "none";
        document.body.appendChild(textArea);
        textArea.focus({ preventScroll: true });
        textArea.select();
        textArea.setSelectionRange(0, textArea.value.length);

        let didCopy = false;
        try {
            didCopy = document.execCommand("copy");
        } catch (error) {
            didCopy = false;
        }

        textArea.remove();
        return didCopy;
    };

    document.querySelectorAll("[data-copy-card]").forEach(card => {
        const copyCard = () => {
            copyText(getCopyCardText(card))
                .then(didCopy => {
                    if (didCopy) {
                        showCopyNotification();
                        return;
                    }

                    console.error("Не удалось скопировать текст.");
                })
                .catch(error => console.error("Не удалось скопировать текст:", error));
        };

        card.addEventListener("click", event => {
            const target = event.target;
            if (target instanceof Element && target.closest("a, button")) return;
            copyCard();
        });

        card.addEventListener("keydown", event => {
            if (event.key !== "Enter" && event.key !== " ") return;
            const target = event.target;
            if (target instanceof Element && target.closest("a, button")) return;
            event.preventDefault();
            copyCard();
        });
    });

    document.addEventListener("keydown", event => {
        if (event.key !== "Escape") return;

        if (modal && modal.classList.contains("is-open")) {
            closeModal();
        }

        if (projectImageModal && projectImageModal.classList.contains("is-open")) {
            closeProjectImageModal();
        }
    });
});

const rippleStyle = document.createElement("style");
rippleStyle.textContent = `
    @keyframes ripple {
        to {
            transform: scale(2);
            opacity: 0;
        }
    }
`;
document.head.appendChild(rippleStyle);
