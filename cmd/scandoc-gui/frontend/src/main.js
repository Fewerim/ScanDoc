import {
    OpenLog,
    StartInit,
    CheckInitStatus,
    GetFilesFromStorage,
    ReadFileFromStorage,
    SaveFileToStorage,
    DeleteFileFromStorage,
    RunMultiFile,
    RunOnceFile
} from "../wailsjs/go/main/App"
import { EventsOn } from "../wailsjs/runtime/runtime"

let selectedFile = null
let navigationStack = ["menuPage"]
let isProcessing = false
let isInitialized = false
let originalEditContent = ""
let currentEditFilename = null
let isLoadingFiles = false
let inputsBound = false
let noticeTimer = null
let initRequestedByButton = false

const menuButtons = [
    { id: "initBtn", page: "initPage", event: "initPage-clicked" },
    { id: "run_onceBtn", page: "run_oncePage", event: "run_oncePage-clicked" },
    { id: "run_multiBtn", page: "run_multiPage", event: "run_multiPage-clicked" },
    { id: "resultsBtn", page: "resultsPage", event: "resultsPage-clicked" }
]

document.addEventListener("DOMContentLoaded", () => {
    registerEvents()
    setupMenu()
    setupStaticUiHandlers()
    setupPageObserver()
    bindInputsOnce()
    bindAllButtons()
    checkInitOnPageLoad()
})

function byId(id) {
    return document.getElementById(id)
}

function qsa(sel) {
    return document.querySelectorAll(sel)
}

function bindOnClick(sel, handler) {
    qsa(sel).forEach(el => (el.onclick = handler))
}

function setBtnState(btn, disabled, title = "") {
    if (!btn) return
    btn.disabled = !!disabled
    btn.style.opacity = btn.disabled ? "0.5" : "1"
    btn.style.cursor = btn.disabled ? "not-allowed" : "pointer"
    btn.title = title
}

function stop(e) {
    e?.stopPropagation?.()
}

function registerEvents() {
    EventsOn("init_status", status => {
        const statusEl = byId("status")
        if (statusEl && status) {
            statusEl.textContent = getStatusText(status)
            statusEl.className = `status-${status}`
        }

        if (status === "process") setProcessing(true)
        else setProcessing(false)

        if (status === "success" || status === "already-init" || status === "done") {
            isInitialized = true
            if (initRequestedByButton && (status === "success" || status === "already-init" || status === "done")) {
                initRequestedByButton = false
                showNotice(
                    status === "already-init" ? "Инициализация уже выполнена" : "Инициализация завершена",
                    "success",
                    2200
                )
            }
            updateAllButtons()
        }
    })

    EventsOn("run_once_status", status => {
        const statusEl = byId("status_run_once")
        if (statusEl) {
            statusEl.textContent = status
            statusEl.className = `status-${status}`
        }

        if (status === "success") showNotice("Файл успешно обработан", "success", 1800)
        if (status === "error") showNotice("Ошибка обработки (смотрите лог)", "error", 2400)

        updateRunOnceStartButton()
    })

    EventsOn("run_multi_status", (status, errMsg = null) => {
        const statusEl = byId("status_run_multi")
        if (statusEl) {
            statusEl.textContent = status
            statusEl.className = `status-${status}`
        }

        if (status === "success") showNotice("Папка успешно обработана", "success", 1800)
        if (status === "error") showNotice(errMsg ? ("Ошибка обработки (смотрите лог)" + errMsg) : "Ошибка обработки (смотрите лог)", "error", 2600)

        updateRunMultiStartButton()
    })

    EventsOn("processing_start", () => setProcessing(true))
    EventsOn("processing_end", () => setProcessing(false))
}

function setupMenu() {
    setupMenuButton("menuBtn", "menuPage", "menu-clicked")
    menuButtons.forEach(({ id, page, event }) => setupMenuButton(id, page, event))
}

function setupStaticUiHandlers() {
    qsa('[id*="exitBtn"]').forEach(btn => (btn.onclick = () => window.runtime.Quit()))

    qsa("h1.nameApp").forEach(scandoc => {
        scandoc.style.cursor = "pointer"
        scandoc.style.userSelect = "none"
        scandoc.onclick = () => pushPage("mainPage")
        scandoc.onmouseenter = () => (scandoc.style.opacity = "0.85")
        scandoc.onmouseleave = () => (scandoc.style.opacity = "1")
    })

    bindOnClick("#cancelEditBtn", e => {
        stop(e)
        if (isProcessing) return
        const editorEl = byId("editor")
        if (editorEl) editorEl.value = originalEditContent
        showNotice("Изменения отменены", "warn", 1600)
    })

    bindOnClick("#saveEditBtn", async e => {
        stop(e)
        if (isProcessing) return
        if (!currentEditFilename) {
            showNotice("Файл не выбран", "error", 2200)
            return
        }

        const editorEl = byId("editor")
        const newContent = editorEl ? editorEl.value : ""

        try {
            await SaveFileToStorage(currentEditFilename, newContent)
            originalEditContent = newContent
            showNotice("Файл сохранён", "success", 1800)
        } catch (err) {
            console.error("SaveFileToStorage error:", err)
            showNotice("Ошибка сохранения: " + (err?.message || err), "error", 2600)
        }
    })

    bindOnClick("#delBtn", async e => {
        stop(e)
        if (isProcessing) return
        if (!selectedFile) {
            updateDeleteButton(false)
            return
        }

        const filenameToDelete = selectedFile
        try {
            await DeleteFileFromStorage(filenameToDelete)
            showNotice("Файл удален", "warn", 1600)
            selectedFile = null
            resetPreview()
            updateEditButton(false)
            updateDeleteButton(false)
            await loadFiles()
        } catch (err) {
            console.error("DeleteFileFromStorage error:", err)
            showNotice("Ошибка удаления: " + (err?.message || err), "error", 2400)
        }
    })
}

function setupMenuButton(id, page, event) {
    const btn = byId(id)
    if (!btn) {
        console.warn(`Menu button #${id} not found`)
        return
    }

    btn.onclick = () => {
        pushPage(page)
        window.runtime.EventsEmit(event)
    }

    updateRunButtonState(btn)
}

function bindInputsOnce() {
    if (inputsBound) return
    inputsBound = true

    const fileInput = byId("fileInput")
    const filenameInput = byId("filenameInput")
    const dirPathInput = byId("dirPathInput")
    const folderNameCreate = byId("folderNameCreate")

    fileInput?.addEventListener("input", updateRunOnceStartButton)
    filenameInput?.addEventListener("input", updateRunOnceStartButton)
    dirPathInput?.addEventListener("input", updateRunMultiStartButton)
    folderNameCreate?.addEventListener("input", updateRunMultiStartButton)
}

function bindAllButtons() {
    bindOnClick("#openLogBtn", async e => {
        stop(e)
        if (isProcessing) return
        await OpenLog()
    })

    bindOnClick("#backBtn", e => {
        stop(e)
        goBack()
    })

    bindOnClick("#open_resultsBtn", e => {
        stop(e)
        if (isProcessing) return
        pushPage("resultsPage")
    })

    bindOnClick("#startInitBtn", async e => {
        stop(e)
        if (isProcessing) return
        initRequestedByButton = true
        await StartInit()
    })

    const editBtn = byId("editBtn")
    if (editBtn) {
        editBtn.onclick = async e => {
            stop(e)
            if (!selectedFile || isProcessing) {
                if (!selectedFile) alert("Сначала выберите файл!")
                return
            }
            currentEditFilename = selectedFile
            pushPage("editPage")
            await openEditorForFile(currentEditFilename)
        }
    }

    bindOnClick("#start_run_onceBtn", async e => {
        stop(e)
        if (isProcessing) return

        const filePath = (byId("fileInput")?.value || "").trim()
        const outName = (byId("filenameInput")?.value || "").trim()

        if (!filePath) {
            showNotice("Укажите путь к файлу", "error", 2200)
            updateRunOnceStartButton()
            return
        }
        if (!outName) {
            showNotice("Укажите название результата", "error", 2200)
            updateRunOnceStartButton()
            return
        }

        try {
            await RunOnceFile(filePath, outName)
        } catch (err) {
            console.error("RunOnceFile error:", err)
            showNotice("Ошибка запуска обработки", "error", 2400)
        }
    })

    bindOnClick("#start_run_multiBtn", async e => {
        stop(e)
        if (isProcessing) return

        const dir = (byId("dirPathInput")?.value || "").trim()
        const folder = (byId("folderNameCreate")?.value || "").trim()

        if (!dir) {
            showNotice("Укажите путь к папке", "error", 2200)
            updateRunMultiStartButton()
            return
        }
        if (!folder) {
            showNotice("Укажите название папки результата", "error", 2200)
            updateRunMultiStartButton()
            return
        }

        try {
            await RunMultiFile(dir, folder)
        } catch (err) {
            console.error("RunMultiFile error:", err)
            showNotice("Ошибка запуска обработки", "error", 2400)
        }
    })

    updateRunOnceStartButton()
    updateRunMultiStartButton()
    updateAllButtons()
}

function updateRunOnceStartButton() {
    const btn = byId("start_run_onceBtn")
    if (!btn) return

    const path = (byId("fileInput")?.value || "").trim()
    const outName = (byId("filenameInput")?.value || "").trim()

    setBtnState(btn, isProcessing || !path || !outName)
}

function updateRunMultiStartButton() {
    const btn = byId("start_run_multiBtn")
    if (!btn) return

    const dir = (byId("dirPathInput")?.value || "").trim()
    const folder = (byId("folderNameCreate")?.value || "").trim()

    setBtnState(btn, isProcessing || !dir || !folder)
}

function updateRunButtonState(btn) {
    if (!btn) return

    if (!isInitialized && (btn.id === "run_onceBtn" || btn.id === "run_multiBtn")) {
        btn.disabled = true
        btn.style.opacity = "0.7"
        btn.style.cursor = "not-allowed"
        btn.title = "Сначала выполните инициализацию!"
        btn.dataset.disabledReason = "requires-init"

        if (!btn.querySelector(".init-badge")) {
            const badge = document.createElement("span")
            badge.className = "init-badge"
            badge.textContent = "требуется инициализация"
            badge.style.cssText =
                "display: block; font-size: 0.8em; font-weight: bold; color: #D5C9F1FF; margin-top: 2px; opacity: 1;"
            btn.appendChild(badge)
        }
        return
    }

    btn.disabled = false
    btn.style.opacity = "1"
    btn.style.cursor = "pointer"
    btn.title = ""
    delete btn.dataset.disabledReason

    const badge = btn.querySelector(".init-badge")
    if (badge) badge.remove()
}

function updateDeleteButton(enabled) {
    const delBtn = byId("delBtn")
    if (!delBtn) return

    if (isProcessing) {
        setBtnState(delBtn, true)
        return
    }

    const resultsPage = byId("resultsPage")
    if (resultsPage?.classList.contains("hidden")) {
        setBtnState(delBtn, true)
        return
    }

    setBtnState(delBtn, !(enabled && selectedFile))
}

function setProcessing(processing) {
    isProcessing = processing
    document.body.classList.toggle("processing", isProcessing)
    updateAllButtons()
    updateRunOnceStartButton()
    updateRunMultiStartButton()
}

function updateAllButtons() {
    qsa("#run_onceBtn, #run_multiBtn").forEach(updateRunButtonState)

    const buttons = document.querySelectorAll("button")
    const inputs = document.querySelectorAll('input[type="button"], input[type="submit"]')
    const links = document.querySelectorAll("a[href]")
    const selects = document.querySelectorAll("select")
    const textareas = document.querySelectorAll("textarea")
    const onclicks = document.querySelectorAll("[onclick]")

    ;[...buttons, ...inputs, ...links, ...selects, ...textareas, ...onclicks].forEach(el => {
        const menuIds = ["menuBtn", "initBtn", "run_onceBtn", "run_multiBtn", "resultsBtn"]
        if (menuIds.includes(el.id) || el.matches("h1.nameApp") || el.id?.includes("exitBtn")) return

        if (isProcessing) {
            el.disabled = true
            el.style.opacity = "0.5"
            el.style.cursor = "not-allowed"
            el.style.pointerEvents = "none"
            el.title = "Выполнение в процессе..."
            return
        }

        el.disabled = false
        el.style.opacity = "1"
        el.style.cursor = ""
        el.style.pointerEvents = ""
        el.title = ""
    })
}

function pushPage(pageId) {
    navigationStack.push(pageId)
    showPage(pageId)
}

function goBack() {
    if (navigationStack.length > 1) {
        navigationStack.pop()
        const previousPage = navigationStack[navigationStack.length - 1]
        showPage(previousPage)
    } else {
        showPage("menuPage")
    }
}

function showPage(pageId) {
    const resultsPage = byId("resultsPage")
    if (resultsPage && !resultsPage.classList.contains("hidden")) {
        if (pageId !== "resultsPage") resetResultsState()
    }

    bindAllButtons()

    qsa(".page").forEach(page => page.classList.add("hidden"))
    byId(pageId)?.classList.remove("hidden")

    if (pageId === "initPage") checkInitOnPageLoad()
    else if (pageId === "run_oncePage") resetRunOncePage()
    else if (pageId === "run_multiPage") resetRunMultiPage()

    window.history.pushState({ page: pageId }, "", `#${pageId}`)
}

function resetRunOncePage() {
    const pathEl = byId("fileInput")
    const nameEl = byId("filenameInput")
    const statusEl = byId("status_run_once")

    if (pathEl) pathEl.value = ""
    if (nameEl) nameEl.value = ""

    if (statusEl) {
        statusEl.textContent = "ready"
        statusEl.className = "status-ready"
    }

    updateRunOnceStartButton()
}

function resetRunMultiPage() {
    const pathEl = byId("dirPathInput")
    const nameEl = byId("folderNameCreate")
    const statusEl = byId("status_run_multi")

    if (pathEl) pathEl.value = ""
    if (nameEl) nameEl.value = ""

    if (statusEl) {
        statusEl.textContent = "ready"
        statusEl.className = "status-ready"
    }

    updateRunMultiStartButton()
}

async function loadFiles() {
    if (isLoadingFiles) return
    isLoadingFiles = true

    selectedFile = null
    updateEditButton(false)
    updateDeleteButton(false)

    try {
        const files = await GetFilesFromStorage()
        const list = byId("filesList")
        if (!list) return
        list.innerHTML = ""

        if (!files || files.length === 0) {
            list.innerHTML =
                '<div style="padding:20px;text-align:center;color:rgba(248, 250, 252, 0.55);font-weight: bold;">Файлы не найдены</div>'
            return
        }

        files.forEach(file => {
            const filename = typeof file === "object" ? file.name || file.filename : file
            const item = document.createElement("div")
            item.className = "file-item"
            item.dataset.filename = filename
            item.innerHTML = `<strong>${filename}</strong>`
            item.style.cursor = isProcessing ? "not-allowed" : "pointer"
            item.style.opacity = isProcessing ? "0.5" : "1"

            item.addEventListener("click", () => {
                if (isProcessing) return
                qsa("#filesList .file-item").forEach(el => el.classList.remove("is-selected"))
                item.classList.add("is-selected")

                selectedFile = filename
                loadFileContent(filename)
                updateEditButton(true)
                updateDeleteButton(true)
            })

            list.appendChild(item)
        })
    } finally {
        isLoadingFiles = false
    }
}

async function loadFileContent(filename) {
    try {
        byId("previewTitle").textContent = `Файл: ${filename}`
        const content = await ReadFileFromStorage(filename)
        byId("fileContent").textContent = content || "Файл пустой"
    } catch (error) {
        byId("fileContent").textContent = "Ошибка: " + error
    }
}

function updateEditButton(enabled) {
    const editBtn = byId("editBtn")
    if (!editBtn) return

    if (isProcessing) {
        setBtnState(editBtn, true)
        return
    }

    const resultsPage = byId("resultsPage")
    if (resultsPage?.classList.contains("hidden")) {
        setBtnState(editBtn, true)
        selectedFile = null
        return
    }

    if (enabled && selectedFile) {
        setBtnState(editBtn, false)
        return
    }

    setBtnState(editBtn, true)
    selectedFile = null
    resetPreview()
}

function resetPreview() {
    const previewTitle = byId("previewTitle")
    const fileContent = byId("fileContent")
    if (previewTitle) previewTitle.textContent = "Выберите файл для просмотра"
    if (fileContent) fileContent.textContent = ""
}

function resetResultsState() {
    selectedFile = null
    updateEditButton(false)
    updateDeleteButton(false)
}

function updateInitStatus(status, errorMsg) {
    const statusEl = byId("status")
    const initBtn = byId("startInitBtn")
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
        ready: "Готов к инициализации",
        process: "Выполняется...",
        success: "Успешно завершено",
        "already-init": "Инициализация уже выполнена",
        done: "Готово"
    }
    return texts[status] || status
}

function setupPageObserver() {
    const resultsPage = byId("resultsPage")
    if (!resultsPage) return

    const observer = new MutationObserver(mutations => {
        mutations.forEach(mutation => {
            if (mutation.type === "attributes" && mutation.attributeName === "class") {
                if (resultsPage.classList.contains("hidden")) resetPreview()
                else loadFiles()
            }
        })
    })

    observer.observe(resultsPage, { attributes: true })
}

async function checkInitOnPageLoad() {
    try {
        const status = await CheckInitStatus()
        if (status === "success" || status === "already-init" || status === "done") isInitialized = true
        updateInitStatus(status)
        updateAllButtons()
    } catch (error) {
        updateInitStatus("error", error.message)
    }
}

async function openEditorForFile(filename) {
    try {
        currentEditFilename = filename

        const titleEl = byId("editFilename")
        if (titleEl) titleEl.textContent = `Файл: ${filename}`

        const content = await ReadFileFromStorage(filename)
        originalEditContent = content || ""

        const editorEl = byId("editor")
        if (editorEl) editorEl.value = originalEditContent
    } catch (e) {
        showNotice("Ошибка чтения файла: " + (e?.message || e), "error", 2400)
        console.error("ReadFileFromStorage error:", e)
    }
}

function showNotice(text, type = "success", timeoutMs = 1800) {
    const old = byId("appNotice")
    if (old) old.remove()
    if (noticeTimer) clearTimeout(noticeTimer)

    const el = document.createElement("div")
    el.id = "appNotice"
    el.className = `notice ${type}`
    el.textContent = text

    document.body.appendChild(el)
    requestAnimationFrame(() => el.classList.add("show"))

    noticeTimer = setTimeout(() => {
        el.classList.remove("show")
        setTimeout(() => el.remove(), 250)
    }, timeoutMs)
}
