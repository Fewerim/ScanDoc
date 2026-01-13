import { OpenLog, StartInit, CheckInitStatus, GetFilesFromStorage, ReadFileFromStorage } from "../wailsjs/go/main/App"
import { EventsOn } from "../wailsjs/runtime/runtime"

const menuButtons = [
    {id: "initBtn", page: "initPage", event: "initPage-clicked"},
    {id: "run_onceBtn", page: "run_oncePage", event: "run_oncePage-clicked"},
    {id: "run_multiBtn", page: "run_multiPage", event: "run_multiPage-clicked"},
    {id: "resultsBtn", page: "resultsPage", event: "resultsPage-clicked"},
    {id: "backBtn", page: "menuPage", event: "back-clicked"}
]

document.addEventListener('DOMContentLoaded', bindAllButtons)

// Инициализация после загрузки DOM
document.addEventListener('DOMContentLoaded', () => {
    const menuBtn = document.getElementById('menuBtn');
    const backBtn = document.getElementById('backBtn');

    EventsOn("init_status", (...args) => {

        const status = args[0]
        const errorMsg = args[1] || null

        const statusEl = document.getElementById("status")
        if (statusEl && status) {
            statusEl.textContent = getStatusText(status)
            statusEl.className = `status-${status}`
        }
    })

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

function updateInitStatus(status, errorMsg) {
    const statusEl = document.getElementById("status")
    const initBtn = document.getElementById("startInitBtn")

    if (!statusEl || !initBtn) return

    statusEl.textContent = getStatusText(status)
    statusEl.className = `status-${status}`

    if (status === "success" || status === "already-init" || status === "done") {
        initBtn.disabled = true
        initBtn.textContent = "Инициализация завершена"
        initBtn.classList.add("completed")
    } else {
        initBtn.disabled = false
        initBtn.textContent = "Начать инициализацию"
        initBtn.classList.remove("completed")
    }

    if (status === "error" && errorMsg) {
        statusEl.textContent += `: ${errorMsg}`
        console.error("Init error:", errorMsg)
    }
}

function getStatusText(status) {
    const texts = {
        "ready": "Готов к инициализации",
        "process": "Выполняется...",
        "success": "Успешно завершено",
        "already-init": "Инициализация уже выполнена",
        "done": "Готово"
    }
    return texts[status] || status
}

function bindAllButtons() {
    // openLogBtn на ВСЕХ страницах
    document.querySelectorAll('#openLogBtn').forEach(btn => {
        btn.onclick = async (e) => {
            e.stopPropagation()
            await OpenLog()
        }
    })

    // backBtn на ВСЕХ страницах
    document.querySelectorAll('#backBtn').forEach(btn => {
        btn.onclick = (e) => {
            e.stopPropagation()
            showPage('menuPage')
        }
    })

    // startInitBtn
    document.querySelectorAll('#startInitBtn').forEach(btn => {
        btn.onclick = async (e) => {
            e.stopPropagation()
            await StartInit()
        }
    })
}

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
    bindAllButtons()

    document.querySelectorAll('.page').forEach(page => {
        page.classList.add('hidden');
    });
    document.getElementById(pageId).classList.remove('hidden');
    window.history.pushState({page: pageId}, '', `#${pageId}`);

    if (pageId === "initPage") {
        checkInitOnPageLoad()
    }
    if (pageId === 'resultsPage') {
        loadFiles()
    }
}

async function checkInitOnPageLoad() {
    try {
        const status = await CheckInitStatus()
        updateInitStatus(status)
    } catch (error) {
        updateInitStatus("error", error.message)
    }
}

async function loadFiles() {
    try {
        const files = await GetFilesFromStorage();
        const list = document.getElementById('filesList');
        list.innerHTML = '';

        if (!files || files.length === 0) {
            list.innerHTML = '<div style="padding:20px;text-align:center;color:#64748b;">Файлы не найдены</div>';
            return;
        }

        files.forEach(file => {
            const filename = typeof file === 'object' ? file.name || file.filename : file;
            const item = document.createElement('div');
            item.className = 'file-item';
            item.dataset.filename = filename;
            item.innerHTML = `<strong>${filename}</strong>`;
            item.addEventListener('click', () => loadFileContent(filename));
            list.appendChild(item);
        });
    } catch (error) {
        console.error('Ошибка:', error);
    }
}

function resetPreview() {
    document.getElementById('previewTitle').textContent = 'Выберите файл для просмотра';
    document.getElementById('fileContent').textContent = '';
}

function setupPageObserver() {
    const resultsPage = document.getElementById('resultsPage');
    if (!resultsPage) return;

    const observer = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
            if (mutation.type === 'attributes' && mutation.attributeName === 'class') {
                if (resultsPage.classList.contains('hidden')) {
                    // Страница скрыта — сбрасываем preview
                    resetPreview();
                } else {
                    // Страница показана — загружаем файлы
                    loadFiles();
                }
            }
        });
    });

    // Наблюдаем за классом hidden
    observer.observe(resultsPage, { attributes: true });
}

async function loadFileContent(filename) {
    try {
        document.getElementById('previewTitle').textContent = `Файл: ${filename}`;
        const content = await ReadFileFromStorage(filename);  // Ваша функция
        document.getElementById('fileContent').textContent = content || 'Файл пустой';
    } catch (error) {
        document.getElementById('fileContent').textContent = 'Ошибка: ' + error;
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const observer = new MutationObserver(() => {
        if (!document.getElementById('resultsPage').classList.contains('hidden')) {
            loadFiles();
        }
    });
    observer.observe(document.getElementById('resultsPage'), { attributes: true });
});

document.addEventListener('DOMContentLoaded', () => {
    setupPageObserver();
});

