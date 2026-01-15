package appCmds

import (
	"fmt"
	"os"
	"path/filepath"
	"proWeb/internal/appUtils"
	"proWeb/internal/logger"
	"proWeb/internal/storage"
	"proWeb/lib/config"
)

//type AppCmds interface {
//	GetConfig() *config.Config
//	GetCfgPath() string
//	GetLog() logger.Logger
//	SetCfgPath(path string)
//	StartApp(nameApp string) error
//	CheckStorageJSON() error
//}

type App struct {
	CfgPath string
	Cfg     *config.Config
	Log     logger.Logger
}

func NewApp() *App {
	return &App{
		CfgPath: "",
		Cfg:     nil,
		Log:     nil,
	}
}

//func (a *App) StartApp(nameApp string) error {
//	cfg := a.GetConfig()
//
//	a.LoadConfig()
//	a.SetupLogger(cfg.LogPath)
//
//	if err := a.InitPythonVenv(); err != nil {
//		return err
//	}
//	if err := a.CheckPythonScripts(); err != nil {
//		return err
//	}
//	if err := a.CheckStorageJSON(); err != nil {
//		return err
//	}
//
//	a.Log.Info(nameApp, "Успешный старт")
//	return nil
//}
//
//func (a *App) SetCfgPath(path string) {
//	a.CfgPath = path
//}
//
//func (a *App) GetConfig() *config.Config {
//	return a.Cfg
//}
//
//func (a *App) GetLog() logger.Logger {
//	return a.Log
//}
//
//func (a *App) GetCfgPath() string {
//	return a.CfgPath
//}

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
	venvPath := filepath.Join(a.Cfg.PythonVenvPath, ".venv")

	if _, err := os.Stat(venvPath); os.IsNotExist(err) {
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
	storage.InitStorage(a.Cfg.StoragePath)

	if !storage.CheckStorageExists() {
		if err := storage.CreateStorageJSON(); err != nil {
			return fmt.Errorf("ошибка создания локального хранилища: %v", err)
		}
	}
	return nil
}
