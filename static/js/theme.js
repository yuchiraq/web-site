document.addEventListener("DOMContentLoaded", function () {
    const themeToggle = document.getElementById("theme-toggle");
    const body = document.body;

    function applyTheme(theme) {
        if (theme === "dark") {
            body.classList.add("dark-mode");
            themeToggle.textContent = "☀️";
        } else {
            body.classList.remove("dark-mode");
            themeToggle.textContent = "🌙";
        }
    }

    // Проверяем системные настройки и localStorage
    const prefersDarkScheme = window.matchMedia("(prefers-color-scheme: dark)");
    let savedTheme = localStorage.getItem("theme");
    
    if (!savedTheme) {
        savedTheme = prefersDarkScheme.matches ? "dark" : "light";
        localStorage.setItem("theme", savedTheme);
    }
    applyTheme(savedTheme);

    themeToggle.addEventListener("click", function () {
        let newTheme = body.classList.contains("dark-mode") ? "light" : "dark";
        localStorage.setItem("theme", newTheme);
        applyTheme(newTheme);
    });

    prefersDarkScheme.addEventListener("change", (e) => {
        if (!localStorage.getItem("theme")) {
            applyTheme(e.matches ? "dark" : "light");
        }
    });
});
