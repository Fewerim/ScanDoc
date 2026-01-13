import { OpenLog, StartInit } from "../wailsjs/go/main/App"

const menuButtons = [
    {id: "initBtn", page: "initPage", event: "initPage-clicked"},
    {id: "run_onceBtn", page: "run_oncePage", event: "run_oncePage-clicked"},
    {id: "run_multiBtn", page: "run_multiPage", event: "run_multiPage-clicked"},
    {id: "resultsBtn", page: "resultsPage", event: "resultsPage-clicked"},
    {id: "backBtn", page: "menuPage", event: "back-clicked"}
]

document.getElementById("openLogBtn").addEventListener("click", async () => {
    try {
        await OpenLog()
    } catch (error) {
        console.error("Ошибка:", error)
    }
})

document.getElementById("startInitBtn").addEventListener("click", async () => {
    try {
        await StartInit()
    } catch (error) {
        console.error("Ошибка:", error)
    }
})

// document.getElementById("openStorageBtn").addEventListener("click", async () => {
//     try {
//         await OpenStorage()
//     } catch (error) {
//         console.error("Ошибка:", error)
//     }
// })

// Инициализация после загрузки DOM
document.addEventListener('DOMContentLoaded', () => {
    const menuBtn = document.getElementById('menuBtn');
    const backBtn = document.getElementById('backBtn');

    setupButton("menuBtn", ()=>{
        showPage('menuPage');
        window.runtime.EventsEmit("menu-clicked");
    })

    menuButtons.forEach(({id, page, event}) => {
        setupButton(id, () => {
            showPage(page)
            window.runtime.EventsEmit(event)
        })
    })

    document.querySelectorAll('[id*="exitBtn"]').forEach(btn => {
        btn.onclick = () => window.runtime.Quit();
    });

    document.querySelectorAll('h1.nameApp').forEach(scandoc => {
        scandoc.style.cursor = 'pointer';
        scandoc.style.userSelect = 'none';
        scandoc.onclick = (e) => showPage('mainPage');
        scandoc.onmouseenter = () => scandoc.style.opacity = '0.85';
        scandoc.onmouseleave = () => scandoc.style.opacity = '1';
    });
});

// setupButton - устанавливает кнопку
function setupButton(id, handler) {
    const btn = document.getElementById(id)
    if (btn) {
        btn.onclick = handler;
    } else {
        console.warn(`Button #${id} not found`);
    }
}

// showPage - показывает страницу, удаляя у нее класс 'hidden'
function showPage(pageId) {
    document.querySelectorAll('.page').forEach(page => {
        page.classList.add('hidden');
    });
    document.getElementById(pageId).classList.remove('hidden');
    window.history.pushState({page: pageId}, '', `#${pageId}`);
}

