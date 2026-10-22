(() => {
    const root = document.documentElement;
    const overlay = document.querySelector("[data-site-preloader]");

    if (!overlay || !root.classList.contains("preloader-active")) {
        overlay?.remove();
        return;
    }

    const count = overlay.querySelector("[data-site-preloader-count]");
    const bar = overlay.querySelector("[data-site-preloader-bar]");
    const startedAt = performance.now();
    let displayedProgress = 0;
    let pageReady = document.readyState === "complete";
    let finished = false;

    document.body.classList.add("preloader-lock");

    const markReady = () => {
        pageReady = true;
    };

    if (!pageReady) {
        window.addEventListener("load", markReady, { once: true });
    }

    window.setTimeout(markReady, 850);

    const finish = () => {
        if (finished) return;
        finished = true;
        overlay.dataset.phase = "fade-ui";

        window.setTimeout(() => {
            overlay.dataset.phase = "reveal";
        }, 90);

        window.setTimeout(() => {
            root.classList.remove("preloader-active");
            document.body.classList.remove("preloader-lock");
            overlay.remove();
            try {
                window.sessionStorage.setItem("avayusstroi_intro_seen", "1");
            } catch (error) {}
        }, 800);
    };

    const renderProgress = currentTime => {
        const elapsed = currentTime - startedAt;
        const loadingTarget = Math.min(94, 12 + elapsed * 0.18);
        const targetProgress = pageReady ? 100 : loadingTarget;
        displayedProgress += (targetProgress - displayedProgress) * 0.28;

        const roundedProgress = Math.min(100, Math.round(displayedProgress));
        if (count) count.textContent = String(roundedProgress);
        if (bar) bar.style.transform = `scaleX(${roundedProgress / 100})`;

        if (pageReady && elapsed >= 260 && displayedProgress >= 98.6) {
            if (count) count.textContent = "100";
            if (bar) bar.style.transform = "scaleX(1)";
            finish();
            return;
        }

        window.requestAnimationFrame(renderProgress);
    };

    window.requestAnimationFrame(renderProgress);
})();

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
    const projectImageCaption = document.getElementById("projectImageCaption");
    const projectImagePosition = document.getElementById("projectImagePosition");
    const projectImagePreviousButton = document.querySelector(".project-image-nav--previous");
    const projectImageNextButton = document.querySelector(".project-image-nav--next");
    const projectImages = Array.from(document.querySelectorAll(".project-carousel-image:not([aria-hidden='true'])"));
    const copyNotification = document.getElementById("copy-notification");
    const prefersDarkScheme = window.matchMedia("(prefers-color-scheme: dark)");
    const prefersReducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)");
    const marketingStorageKey = "avayusstroi_marketing_context";
    const trackingParamKeys = ["utm_source", "utm_medium", "utm_campaign", "utm_content", "utm_term", "gclid", "fbclid", "yclid"];

    const syncTheme = () => {
        document.documentElement.setAttribute("data-theme", prefersDarkScheme.matches ? "dark" : "light");
    };

    syncTheme();
    if (typeof prefersDarkScheme.addEventListener === "function") {
        prefersDarkScheme.addEventListener("change", syncTheme);
    } else if (typeof prefersDarkScheme.addListener === "function") {
        prefersDarkScheme.addListener(syncTheme);
    }

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

    const revealElements = Array.from(document.querySelectorAll([
        "[data-aos]",
        ".feature-item",
        ".service-item",
        ".workflow-step",
        ".trusted-clients__item",
        ".box",
        ".contact-card",
        ".price-card",
        ".case-card",
        ".problem-grid article",
        ".reuse-grid article",
        ".faq-section details"
    ].join(", ")));

    if (revealElements.length && !prefersReducedMotion.matches && "IntersectionObserver" in window) {
        document.documentElement.classList.add("motion-ready");

        revealElements.forEach((element, index) => {
            element.classList.add("reveal-item");
            element.style.setProperty("--reveal-delay", `${(index % 4) * 70}ms`);
        });

        const revealObserver = new IntersectionObserver((entries, observer) => {
            entries.forEach(entry => {
                if (!entry.isIntersecting) return;
                entry.target.classList.add("is-revealed");
                observer.unobserve(entry.target);
            });
        }, { rootMargin: "0px 0px -8%", threshold: 0.08 });

        revealElements.forEach(element => revealObserver.observe(element));
    }

    let activeProjectImageIndex = -1;
    let projectImageTrigger = null;

    const showProjectImage = index => {
        if (!projectImageModal || !projectImageModalContent || !projectImages.length) return;

        activeProjectImageIndex = (index + projectImages.length) % projectImages.length;
        const sourceImage = projectImages[activeProjectImageIndex];
        const description = sourceImage.alt || "Фото выполненного объекта";

        projectImageModal.classList.add("is-loading");
        projectImageModalContent.src = sourceImage.currentSrc || sourceImage.src;
        projectImageModalContent.alt = description;
        if (projectImageCaption) projectImageCaption.textContent = description;
        if (projectImagePosition) projectImagePosition.textContent = `${activeProjectImageIndex + 1} / ${projectImages.length}`;

        if (projectImageModalContent.complete) {
            window.requestAnimationFrame(() => projectImageModal.classList.remove("is-loading"));
        }
    };

    const openProjectImageModal = sourceImage => {
        if (!projectImageModal || !projectImageModalContent || !sourceImage) return;

        projectImageTrigger = sourceImage;
        showProjectImage(projectImages.indexOf(sourceImage));
        projectImageModal.classList.add("is-open");
        projectImageModal.setAttribute("aria-hidden", "false");
        syncScrollLock();
        window.setTimeout(() => projectImageCloseButton?.focus(), 0);
    };

    const closeProjectImageModal = () => {
        if (!projectImageModal || !projectImageModalContent) return;

        projectImageModal.classList.remove("is-open");
        projectImageModal.classList.remove("is-loading");
        projectImageModal.setAttribute("aria-hidden", "true");
        projectImageModalContent.src = "";
        activeProjectImageIndex = -1;
        syncScrollLock();
        projectImageTrigger?.focus();
        projectImageTrigger = null;
    };

    projectImages.forEach(image => {
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

    projectImageModalContent?.addEventListener("load", () => {
        projectImageModal?.classList.remove("is-loading");
    });

    projectImagePreviousButton?.addEventListener("click", () => {
        showProjectImage(activeProjectImageIndex - 1);
    });

    projectImageNextButton?.addEventListener("click", () => {
        showProjectImage(activeProjectImageIndex + 1);
    });

    if (projectImageModal) {
        projectImageModal.addEventListener("click", event => {
            const target = event.target;
            if (target instanceof Element && target.closest("[data-close-project-modal='true']")) {
                closeProjectImageModal();
            }
        });
    }

    const marqueeTracks = document.querySelectorAll(".projects .carousel-track");
    let marqueeResizeFrame = 0;

    const setMarqueeSpeed = () => {
        marqueeTracks.forEach(track => {
            const group = track.querySelector(".carousel-group");
            if (!group) return;

            const pixelsPerSecond = window.innerWidth <= 768 ? 42 : 58;
            const duration = Math.max(24, Math.min(90, group.scrollWidth / pixelsPerSecond));
            track.style.setProperty("--carousel-duration", `${duration.toFixed(2)}s`);
        });
    };

    const scheduleMarqueeMeasurement = () => {
        window.cancelAnimationFrame(marqueeResizeFrame);
        marqueeResizeFrame = window.requestAnimationFrame(setMarqueeSpeed);
    };

    if (marqueeTracks.length) {
        setMarqueeSpeed();
        window.addEventListener("load", setMarqueeSpeed, { once: true });
        window.addEventListener("orientationchange", scheduleMarqueeMeasurement);
        window.addEventListener("resize", scheduleMarqueeMeasurement);

        if ("ResizeObserver" in window) {
            const marqueeResizeObserver = new ResizeObserver(scheduleMarqueeMeasurement);
            marqueeTracks.forEach(track => {
                const group = track.querySelector(".carousel-group");
                if (group) marqueeResizeObserver.observe(group);
            });
        }
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
    if (clientCounters.length) {
        const numberFormatter = new Intl.NumberFormat("ru-BY", { maximumFractionDigits: 0 });
        const animateCounter = counter => {
            const target = Number(counter.dataset.target) || 0;
            if (prefersReducedMotion.matches) {
                counter.textContent = numberFormatter.format(target);
                return;
            }

            let currentValue = 0;
            let velocity = 0;
            let previousTime = performance.now();
            const stiffness = 100;
            const damping = 20;

            const updateCounter = currentTime => {
                const delta = Math.min((currentTime - previousTime) / 1000, 0.034);
                previousTime = currentTime;

                const acceleration = (target - currentValue) * stiffness - velocity * damping;
                velocity += acceleration * delta;
                currentValue += velocity * delta;
                counter.textContent = numberFormatter.format(Math.max(0, Math.min(target, Math.round(currentValue))));

                if (Math.abs(target - currentValue) > 0.5 || Math.abs(velocity) > 0.5) {
                    window.requestAnimationFrame(updateCounter);
                    return;
                }

                counter.textContent = numberFormatter.format(target);
                counter.closest(".trusted-clients__item--summary")?.classList.add("is-counted");
            };

            counter.textContent = "0";
            window.requestAnimationFrame(updateCounter);
        };

        if (!("IntersectionObserver" in window)) {
            clientCounters.forEach(counter => {
                counter.textContent = numberFormatter.format(Number(counter.dataset.target) || 0);
            });
        } else {
            const counterObserver = new IntersectionObserver((entries, observer) => {
                entries.forEach(entry => {
                    if (!entry.isIntersecting) return;
                    animateCounter(entry.target);
                    observer.unobserve(entry.target);
                });
            }, { threshold: 0.45 });

            clientCounters.forEach(counter => counterObserver.observe(counter));
        }
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
        if (projectImageModal && projectImageModal.classList.contains("is-open")) {
            if (event.key === "ArrowLeft") {
                showProjectImage(activeProjectImageIndex - 1);
                return;
            }
            if (event.key === "ArrowRight") {
                showProjectImage(activeProjectImageIndex + 1);
                return;
            }
            if (event.key === "Escape") {
                closeProjectImageModal();
                return;
            }
        }

        if (event.key === "Escape" && modal && modal.classList.contains("is-open")) {
            closeModal();
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
