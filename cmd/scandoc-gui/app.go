package main

import (
	"context"
	"fmt"
	"os"
	"proWeb/lib/appUtils"
	"proWeb/lib/config"
	"proWeb/lib/files"
	"proWeb/lib/logger"
)

// App struct
type App struct {
	ctx     context.Context
	CfgPath string
	Log     logger.Logger
	Cfg     *config.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		ctx:     context.Background(),
		CfgPath: "",
		Log:     nil,
		Cfg:     nil,
	}
}

// LoadConfig - загружает конфиг по пути из поля App, если путь пустой, устанавливает его
func (a *App) LoadConfig() {
	path := a.CfgPath
	if path == "" {
		path = config.DefaultConfigPath()
	}

	a.CfgPath = path
	a.Cfg = config.MustLoadWithPath(path)
}

// SetupLogger - устанавливает логгер, который пишет в файл по переданному пути
func (a *App) SetupLogger(pathToFile string) {
	a.Log = logger.MustSetup(pathToFile)
}

// InitPythonVenv - инициализирует venv для python, если файла нет, создает
func (a *App) InitPythonVenv() error {
	if _, err := os.Stat(".venv"); os.IsNotExist(err) {
		if err = appUtils.CreateVenv(a.Cfg.PythonVenvPath); err != nil {
			return err
		}
	}
	return nil
}

// CheckPythonScripts - проверяет наличие пути к python скрипту и наличие venv файла для успешного запуска скрипта
func (a *App) CheckPythonScripts() error {
	pyVenv := a.Cfg.PythonExecutable
	pyScript := a.Cfg.PythonScript

	if _, err := os.Stat(pyVenv); os.IsNotExist(err) {
		return fmt.Errorf("python из venv не найден: %s", pyVenv)
	}
	if _, err := os.Stat(pyScript); os.IsNotExist(err) {
		return fmt.Errorf("серверный скрипт не найден: %s", pyScript)
	}
	return nil
}

// CheckStorageJSON - проверяет наличие локального хранилища, если его нет, создает новое
func (a *App) CheckStorageJSON() error {
	files.InitStorage(a.Cfg.StoragePath)

	if !files.StorageExists() {
		if err := files.CreateStorageJSON(); err != nil {
			return fmt.Errorf("ошибка создания локального хранилища: %v", err)
		}
	}
	return nil
}

func (a *App) SetupApp(ctx context.Context) error {
	const op = "scandoc-gui.startup"
	a.LoadConfig()
	a.SetupLogger(a.Cfg.LogPath)

	if err := a.InitPythonVenv(); err != nil {
		return err
	}

	if err := a.CheckPythonScripts(); err != nil {
		return err
	}

	a.ctx = ctx
	a.Log.Info(op, "успешный старт")
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
