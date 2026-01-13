package main

import (
	"context"
	"fmt"
	"proWeb/lib/appCmds"
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
