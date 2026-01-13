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

	if err := a.app.InitApp(a.Name + op); err != nil {
		runtime.LogErrorf(a.ctx, err.Error())
	}
}
