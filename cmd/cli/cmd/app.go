package cmd

import (
	"fmt"
	"os"
	"proWeb/internal/cliUtils"
	"proWeb/internal/config"
	"proWeb/internal/logger"
)

type App struct {
	Log logger.Logger
	Cfg *config.Config
}

// NewApp - конструктор для приложения
func NewApp() *App {
	return &App{Log: nil, Cfg: nil}
}

// LoadConfig - загружает конфиг из переданного пути
func (a *App) LoadConfig(path string) {
	a.Cfg = config.MustLoadWithPath(path)
}

// SetupLogger - устанавливает логгер, который пишет в файл по переданному пути
func (a *App) SetupLogger(path string) {
	a.Log = logger.MustSetup(path)
}

// InitPythonVenv - инициализирует venv для python, если файла нет, создает
func (a *App) InitPythonVenv() error {
	if _, err := os.Stat(".venv"); os.IsNotExist(err) {
		if err = cliUtils.CreateVenv(); err != nil {
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
