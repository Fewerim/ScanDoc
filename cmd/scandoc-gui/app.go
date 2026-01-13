package main

import (
	"context"
	"proWeb/lib/appCmds"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const NameApp = "scandoc-GUI"

// App struct
type App struct {
	Name string
	ctx  context.Context
	app  *appCmds.App
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		Name: NameApp,
		ctx:  context.Background(),
		app:  appCmds.New(),
	}
}

func (a *App) SetupApp(ctx context.Context) error {
	const op = ".startup"
	a.app.LoadConfig()
	a.app.SetupLogger(a.app.Cfg.LogPath)

	if err := a.app.InitPythonVenv(); err != nil {
		return err
	}

	if err := a.app.CheckPythonScripts(); err != nil {
		return err
	}

	if err := a.app.CheckStorageJSON(); err != nil {
		return err
	}

	a.ctx = ctx
	a.app.Log.Info(a.Name+op, "успешный старт")
	return nil
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	if err := a.SetupApp(ctx); err != nil {
		panic(err)
	}
}

// OpenLog - открывает папку с логами
func (a *App) OpenLog() {
	const op = ".open_log"

	if err := a.app.OpenLogFolder(a.Name+op, false); err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
	}
}

// OpenStorage - открывает локальное хранилище
func (a *App) OpenStorage() {
	const op = ".open_storage"

	if err := a.app.OpenStorage(a.Name+op, false); err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
	}
}

func (a *App) StartInit() {
	const op = ".init"

	if ok, _ := a.app.CheckInit(a.Name + op); ok {
		runtime.EventsEmit(a.ctx, "init_status", "already-init")
		return
	}

	runtime.EventsEmit(a.ctx, "init_status", "process")

	if err := a.app.InitApp(a.Name + op); err != nil {
		runtime.EventsEmit(a.ctx, "init_status", "error")
		return
	}
	runtime.EventsEmit(a.ctx, "init_status", "success")
}

// CheckInitStatus - проверяет, проинициализировано приложение или нет
func (a *App) CheckInitStatus() string {
	const op = ".check_init_status"

	ok, err := a.app.CheckInit(a.Name + op)
	if err != nil {
		return "error"
	}
	if ok {
		return "already-init"
	}
	return "ready"
}

func (a *App) GetFilesFromStorage() interface{} {
	const op = ".get_files_from_storage"

	files, err := a.app.GetFilesFromStorage(a.Name + op)
	if err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
		return nil
	}

	return files
}

func (a *App) ReadFileFromStorage(fileName string) string {
	const op = ".read_file_from_storage"

	content, err := a.app.ReadFileFromStorage(a.Name+op, fileName)
	if err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
		return ""
	}
	return content
}
