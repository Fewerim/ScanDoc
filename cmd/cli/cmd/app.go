package cmd

import (
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
