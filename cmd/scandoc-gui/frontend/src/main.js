// Функция показа страницы (остаётся)
function showPage(pageId) {
    document.querySelectorAll('.page').forEach(page => {
        page.classList.add('hidden');
    });
    document.getElementById(pageId).classList.remove('hidden');
    window.history.pushState({page: pageId}, '', `#${pageId}`);
}

// Инициализация после загрузки DOM
document.addEventListener('DOMContentLoaded', () => {
    const menuBtn = document.getElementById('menuBtn');
    const exitBtn = document.getElementById('exitBtn');

    menuBtn.onclick = async () => {
        console.log("Menu clicked");
        showPage('menuPage');
        window.runtime.EventsEmit("menu-clicked");
    };

    document.querySelectorAll('[id*="exitBtn"]').forEach(btn => {
        btn.onclick = () => window.runtime.Quit();
    });

    document.querySelectorAll('h1.nameApp').forEach(scandoc => {
        scandoc.style.cursor = 'pointer';
        scandoc.style.userSelect = 'none';
        scandoc.onclick = (e) => {
            console.log('SCANDOC clicked → mainPage');
            showPage('mainPage');
            scandoc.style.transform = 'scale(0.97)';
            setTimeout(() => scandoc.style.transform = '', 200);
        };

        scandoc.onmouseenter = () => scandoc.style.opacity = '0.85';
        scandoc.onmouseleave = () => scandoc.style.opacity = '1';
    });
});

