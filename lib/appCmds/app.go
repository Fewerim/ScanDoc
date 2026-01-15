package appCmds

import (
	"fmt"
	"os"
	"path/filepath"
	"proWeb/internal/appUtils"
	"proWeb/internal/files"
	"proWeb/internal/logger"
	"proWeb/internal/storage"
	"proWeb/lib/config"
)

type AppCmds interface {
	GetConfig() *config.Config
	GetCfgPath() string
	GetLog() logger.Logger
	SetCfgPath(path string)
	StartApp(nameApp string) error
	CheckStorageJSON() error
	ClearTerminal(operation string) error
	InitApp(operation string) error
	CheckInit(operation string) (bool, error)
	InstallTesseract() error
	OpenLogFolder(operation string, clearFlag bool) error
	CleanLogFolder(operation string) error
	MultiFiles(operation, directory, createdFolderName string) error
	OnceFile(operation, filePath, createdFileName string) error
	OpenStorage(operation string, clearFlag bool) error
	GetFilesFromStorage(operation string) ([]files.File, error)
	ReadFileFromStorage(operation, fileName string) (string, error)
	SaveFileToStorage(operation, fileName, folder string, content interface{}) error
	DeleteFileFromStorage(operation, fileName string) error
}

type App struct {
	cfgPath string
	storage *storage.Storage
	cfg     *config.Config
	log     logger.Logger
}

func NewApp() *App {
	return &App{
		cfgPath: "",
		storage: nil,
		cfg:     nil,
		log:     nil,
	}
}

// StartApp - загружает и проверяет необходимые модули приложения
func (a *App) StartApp(operation string) error {
	a.loadConfig()
	a.setupLogger(a.cfg.LogPath)
	a.setupStorage(a.cfg.StoragePath)

	if err := a.initPythonVenv(); err != nil {
		return err
	}
	if err := a.checkPythonScripts(); err != nil {
		return err
	}
	if err := a.CheckStorageJSON(); err != nil {
		return err
	}

	a.log.Info(operation, "успешный старт")
	return nil
}

// GetConfig - возвращает конфиг приложения
func (a *App) GetConfig() *config.Config {
	return a.cfg
}

// GetLog - возвращает логгер приложения
func (a *App) GetLog() logger.Logger {
	return a.log
}

// GetCfgPath - возвращает путь к конфигу
func (a *App) GetCfgPath() string {
	return a.cfgPath
}

// SetCfgPath - устанавливает путь к конфигу
func (a *App) SetCfgPath(path string) {
	a.cfgPath = path
}

// loadConfig - загружает конфиг по пути из поля App, если путь пустой, устанавливает его
func (a *App) loadConfig() {
	path := a.cfgPath
	if path == "" {
		path = config.DefaultConfigPath()
	}

	a.cfgPath = path
	a.cfg = config.MustLoadWithPath(path)
}

// setupLogger - устанавливает логгер, который пишет в файл по переданному пути
func (a *App) setupLogger(pathToFile string) {
	a.log = logger.MustSetup(pathToFile)
}

// setupStorage - устанавливает локальное хранилище
func (a *App) setupStorage(pathToStorage string) {
	a.storage = storage.New(pathToStorage)
}

// initPythonVenv - инициализирует venv для python, если файла нет, создает
func (a *App) initPythonVenv() error {
	venvPath := filepath.Join(a.cfg.PythonVenvPath, ".venv")

	if _, err := os.Stat(venvPath); os.IsNotExist(err) {
		if err = appUtils.CreateVenv(a.cfg.PythonVenvPath); err != nil {
			return err
		}
	}
	return nil
}

// checkPythonScripts - проверяет наличие пути к python скрипту и наличие venv файла для успешного запуска скрипта
func (a *App) checkPythonScripts() error {
	pyVenv := a.cfg.PythonExecutable
	pyScript := a.cfg.PythonScript

	if _, err := os.Stat(pyVenv); os.IsNotExist(err) {
		return fmt.Errorf("python из venv не найден: %s", pyVenv)
	}
	if _, err := os.Stat(pyScript); os.IsNotExist(err) {
		return fmt.Errorf("серверный скрипт не найден: %s", pyScript)
	}
	return nil
}
