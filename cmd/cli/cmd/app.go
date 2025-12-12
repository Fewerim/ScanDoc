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
func NewApp(l logger.Logger) *App {
	return &App{Log: l}
}

// LoadConfig - загружает конфиг
func (a *App) LoadConfig(path string) {
	a.Cfg = config.MustLoadWithPath(path)
}
