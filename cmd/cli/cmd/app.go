package cmd

import "proWeb/internal/logger"

type App struct {
	l logger.Logger
}

// NewApp - конструктор для приложения
func NewApp(l logger.Logger) *App {
	return &App{l: l}
}

// Execute - делегирует запуск CLI приложения, вызываясь на экземпляре App
func (a *App) Execute() error {
	return rootCmd.Execute()
}
