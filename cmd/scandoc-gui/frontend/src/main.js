import { OpenLog, StartInit, CheckInitStatus, GetFilesFromStorage, ReadFileFromStorage, SaveFileToStorage } from "../wailsjs/go/main/App"
import { EventsOn } from "../wailsjs/runtime/runtime"

let selectedFile = null;
let navigationStack = ['menuPage'];
let isProcessing = false;
let isInitialized = false;
let originalEditContent = "";
let currentEditFilename = null;

const menuButtons = [
    {id: "initBtn", page: "initPage", event: "initPage-clicked"},
    {id: "run_onceBtn", page: "run_oncePage", event: "run_oncePage-clicked"},
    {id: "run_multiBtn", page: "run_multiPage", event: "run_multiPage-clicked"},
    {id: "resultsBtn", page: "resultsPage", event: "resultsPage-clicked"}
]

document.addEventListener('DOMContentLoaded', bindAllButtons)

document.addEventListener('DOMContentLoaded', () => {
    EventsOn("init_status", (...args) => {
        const status = args[0];
        const errorMsg = args[1] || null;
        const statusEl = document.getElementById("status");

        if (statusEl && status) {
            statusEl.textContent = getStatusText(status);
            statusEl.className = `status-${status}`;
        }

        if (status === "process") {
            setProcessing(true);
        } else if (status === "success" || status === "already-init" || status === "done") {
            setProcessing(false);
            isInitialized = true;
            updateAllButtons();
        } else {
            setProcessing(false);
        }

        if (status === "error") {
            setProcessing(false);
        }
    });

    EventsOn("processing_start", () => {
        setProcessing(true);
    });

    EventsOn("processing_end", () => {
        setProcessing(false);
    });

    setupMenuButton("menuBtn", 'menuPage', "menu-clicked");
    menuButtons.forEach(({id, page, event}) => {
        setupMenuButton(id, page, event);
    });

    setupButton("editBtn", async () => {
        if (!selectedFile || isProcessing) {
            if (!selectedFile) alert('Сначала выберите файл!');
            return;
        }

        currentEditFilename = selectedFile;

        pushPage('editPage');
        await openEditorForFile(currentEditFilename);
    });


    document.querySelectorAll('[id*="exitBtn"]').forEach(btn => {
        btn.onclick = () => window.runtime.Quit();
    });

    document.querySelectorAll('h1.nameApp').forEach(scandoc => {
        scandoc.style.cursor = 'pointer';
        scandoc.style.userSelect = 'none';
        scandoc.onclick = (e) => pushPage('mainPage');
        scandoc.onmouseenter = () => scandoc.style.opacity = '0.85';
        scandoc.onmouseleave = () => scandoc.style.opacity = '1';
    });

    document.querySelectorAll('#cancelEditBtn').forEach(btn => {
        btn.onclick = (e) => {
            e.stopPropagation();
            if (isProcessing) return;

            const editorEl = document.getElementById("editor");
            if (editorEl) editorEl.value = originalEditContent;

            showNotice("Изменения отменены", "warn", 1600);
        };
    });

    document.querySelectorAll('#delBtn').forEach(btn => {
        btn.onclick = (e) => {

            showNotice("Файл удален", "warn", 1600);
        };
    });

    document.querySelectorAll('#saveEditBtn').forEach(btn => {
        btn.onclick = async (e) => {
            e.stopPropagation();
            if (isProcessing) return; // блокировка только когда init/run

            if (!currentEditFilename) {
                showNotice("Файл не выбран", "error", 2200);
                return;
            }

            const editorEl = document.getElementById("editor");
            const newContent = editorEl ? editorEl.value : "";

            try {
                await SaveFileToStorage(currentEditFilename, newContent);
                originalEditContent = newContent;
                showNotice("Файл сохранён", "success", 1800);
            } catch (err) {
                console.error("SaveFileToStorage error:", err);
                showNotice("Ошибка сохранения: " + (err?.message || err), "error", 2600);
            }
        };
    });


    setupPageObserver();
    checkInitOnPageLoad();
});

function setupMenuButton(id, page, event) {
    const btn = document.getElementById(id);
    if (btn) {
        btn.onclick = () => {
            pushPage(page);
            window.runtime.EventsEmit(event);
        };
        updateRunButtonState(btn);
    } else {
        console.warn(`Menu button #${id} not found`);
    }
}

function updateRunButtonState(btn) {
    if (!isInitialized && (btn.id === 'run_onceBtn' || btn.id === 'run_multiBtn')) {
        btn.disabled = true;
        btn.style.opacity = '0.7';
        btn.style.cursor = 'not-allowed';
        btn.title = 'Сначала выполните инициализацию!';
        btn.dataset.disabledReason = 'requires-init';

        if (!btn.querySelector('.init-badge')) {
            const badge = document.createElement('span');
            badge.className = 'init-badge';
            badge.textContent = 'требуется инициализация';
            badge.style.cssText =
                'display: block; ' +
                'font-size: 0.8em; font-weight: bold; ' +
                'color: #D5C9F1FF; ' +
                'margin-top: 2px; ' +
                'opacity: 1;';
            btn.appendChild(badge);
        }

    } else {
        btn.disabled = false;
        btn.style.opacity = '1';
        btn.style.cursor = 'pointer';
        btn.title = '';
        delete btn.dataset.disabledReason;

        const badge = btn.querySelector('.init-badge');
        if (badge) badge.remove();
    }
}


function setProcessing(processing) {
    isProcessing = processing;
    document.body.classList.toggle('processing', isProcessing);
    updateAllButtons();
}

// ✅ Безопасная функция без сложных селекторов
function updateAllButtons() {
    // Обновляем RUN кнопки с подсказкой
    document.querySelectorAll('#run_onceBtn, #run_multiBtn').forEach(updateRunButtonState);

    // ПРОСТЫЕ селекторы
    const buttons = document.querySelectorAll('button');
    const inputs = document.querySelectorAll('input[type="button"], input[type="submit"]');
    const links = document.querySelectorAll('a[href]');
    const selects = document.querySelectorAll('select');
    const textareas = document.querySelectorAll('textarea');
    const onclicks = document.querySelectorAll('[onclick]');

    // Объединяем все элементы
    [...buttons, ...inputs, ...links, ...selects, ...textareas, ...onclicks].forEach(el => {
        // Исключаем меню кнопки
        const menuIds = ['menuBtn', 'initBtn', 'run_onceBtn', 'run_multiBtn', 'resultsBtn'];
        if (menuIds.includes(el.id) || el.matches('h1.nameApp') || el.id?.includes('exitBtn')) {
            return;
        }

        if (isProcessing) {
            el.disabled = true;
            el.style.opacity = '0.5';
            el.style.cursor = 'not-allowed';
            el.style.pointerEvents = 'none';
            el.title = 'Выполнение в процессе...';
        } else {
            el.disabled = false;
            el.style.opacity = '1';
            el.style.cursor = '';
            el.style.pointerEvents = '';
            el.title = '';
        }
    });
}

function pushPage(pageId) {
    navigationStack.push(pageId);
    showPage(pageId);
}

function goBack() {
    if (navigationStack.length > 1) {
        navigationStack.pop();
        const previousPage = navigationStack[navigationStack.length - 1];
        showPage(previousPage);
    } else {
        showPage('menuPage');
    }
}

function showPage(pageId) {
    if (document.getElementById('resultsPage') && !document.getElementById('resultsPage').classList.contains('hidden')) {
        if (pageId !== 'resultsPage') {
            resetResultsState();
        }
    }

    bindAllButtons();
    document.querySelectorAll('.page').forEach(page => {
        page.classList.add('hidden');
    });
    document.getElementById(pageId)?.classList.remove('hidden');

    switch (pageId) {
        case 'initPage':
            checkInitOnPageLoad();
            break;
        case 'resultsPage':
            loadFiles();
            break;
    }

    window.history.pushState({page: pageId}, '', `#${pageId}`);
}

function bindAllButtons() {
    document.querySelectorAll('#openLogBtn').forEach(btn => {
        btn.onclick = async (e) => {
            e.stopPropagation();
            if (isProcessing) return;
            await OpenLog();
        };
    });

    document.querySelectorAll('#backBtn').forEach(btn => {
        btn.onclick = (e) => {
            e.stopPropagation();
            goBack();
        };
    });

    document.querySelectorAll('#open_resultsBtn').forEach(btn => {
        btn.onclick = (e) => {
            e.stopPropagation();
            if (isProcessing) return;
            pushPage('resultsPage');
        };
    });

    document.querySelectorAll('#startInitBtn').forEach(btn => {
        btn.onclick = async (e) => {
            e.stopPropagation();
            if (isProcessing) return;
            await StartInit();
        };
    });


}

function setupButton(id, handler) {
    const btn = document.getElementById(id);
    if (btn) {
        btn.onclick = (e) => {
            if (isProcessing && id !== 'editBtn') return;
            handler(e);
        };
    } else {
        console.warn(`Button #${id} not found`);
    }
}

async function loadFiles() {
    selectedFile = null;
    updateEditButton(false);

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
            item.style.cursor = isProcessing ? 'not-allowed' : 'pointer';
            item.style.opacity = isProcessing ? '0.5' : '1';
            item.addEventListener('click', () => {
                if (isProcessing) return;
                selectedFile = filename;
                loadFileContent(filename);
                updateEditButton(true);
            });
            list.appendChild(item);
        });
    } catch (error) {
        console.error('Ошибка:', error);
    }
}

async function loadFileContent(filename) {
    try {
        document.getElementById('previewTitle').textContent = `Файл: ${filename}`;
        const content = await ReadFileFromStorage(filename);
        document.getElementById('fileContent').textContent = content || 'Файл пустой';
    } catch (error) {
        document.getElementById('fileContent').textContent = 'Ошибка: ' + error;
    }
}

function updateEditButton(enabled) {
    const editBtn = document.getElementById('editBtn');
    if (!editBtn) return;

    if (isProcessing) {
        editBtn.disabled = true;
        editBtn.style.opacity = '0.5';
        editBtn.style.cursor = 'not-allowed';
        return;
    }

    const resultsPage = document.getElementById('resultsPage');
    if (resultsPage?.classList.contains('hidden')) {
        editBtn.disabled = true;
        editBtn.style.opacity = '0.5';
        editBtn.style.cursor = 'not-allowed';
        selectedFile = null;
        return;
    }

    if (enabled && selectedFile) {
        editBtn.disabled = false;
        editBtn.style.opacity = '1';
        editBtn.style.cursor = 'pointer';
    } else {
        editBtn.disabled = true;
        editBtn.style.opacity = '0.5';
        editBtn.style.cursor = 'not-allowed';
        selectedFile = null;
        resetPreview();
    }
}

function resetPreview() {
    const previewTitle = document.getElementById('previewTitle');
    const fileContent = document.getElementById('fileContent');
    if (previewTitle) previewTitle.textContent = 'Выберите файл для просмотра';
    if (fileContent) fileContent.textContent = '';
}

function resetResultsState() {
    selectedFile = null;
    updateEditButton(false);
}

function updateInitStatus(status, errorMsg) {
    const statusEl = document.getElementById("status");
    const initBtn = document.getElementById("startInitBtn");

    if (!statusEl || !initBtn) return;

    statusEl.textContent = getStatusText(status);
    statusEl.className = `status-${status}`;

    if (status === "success" || status === "already-init" || status === "done") {
        initBtn.disabled = true;
        initBtn.textContent = "Инициализация завершена";
        initBtn.classList.add("completed");
    } else {
        initBtn.disabled = false;
        initBtn.textContent = "Начать инициализацию";
        initBtn.classList.remove("completed");
    }

    if (status === "error" && errorMsg) {
        statusEl.textContent += `: ${errorMsg}`;
        console.error("Init error:", errorMsg);
    }
}

function getStatusText(status) {
    const texts = {
        "ready": "Готов к инициализации",
        "process": "Выполняется...",
        "success": "Успешно завершено",
        "already-init": "Инициализация уже выполнена",
        "done": "Готово"
    };
    return texts[status] || status;
}

function setupPageObserver() {
    const resultsPage = document.getElementById('resultsPage');
    if (!resultsPage) return;

    const observer = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
            if (mutation.type === 'attributes' && mutation.attributeName === 'class') {
                if (resultsPage.classList.contains('hidden')) {
                    resetPreview();
                } else {
                    loadFiles();
                }
            }
        });
    });
    observer.observe(resultsPage, { attributes: true });
}

async function checkInitOnPageLoad() {
    try {
        const status = await CheckInitStatus();
        if (status === "success" || status === "already-init" || status === "done") {
            isInitialized = true;
        }
        updateInitStatus(status);
        updateAllButtons();
    } catch (error) {
        updateInitStatus("error", error.message);
    }
}


async function openEditorForFile(filename) {
    try {
        currentEditFilename = filename;

        const titleEl = document.getElementById("editFilename");
        if (titleEl) titleEl.textContent = `Файл: ${filename}`;

        const content = await ReadFileFromStorage(filename);
        originalEditContent = content || "";

        const editorEl = document.getElementById("editor");
        if (editorEl) editorEl.value = originalEditContent;

    } catch (e) {
        console.error("ReadFileFromStorage error:", e);
        showNotice("Ошибка чтения файла: " + (e?.message || e), "error", 2400);
    }
}


let noticeTimer = null;

function showNotice(text, type = "success", timeoutMs = 1800) {
    const old = document.getElementById("appNotice");
    if (old) old.remove();
    if (noticeTimer) clearTimeout(noticeTimer);

    const el = document.createElement("div");
    el.id = "appNotice";
    el.className = `notice ${type}`;
    el.textContent = text;

    document.body.appendChild(el);
    requestAnimationFrame(() => el.classList.add("show"));

    noticeTimer = setTimeout(() => {
        el.classList.remove("show");
        setTimeout(() => el.remove(), 250);
    }, timeoutMs);
}
