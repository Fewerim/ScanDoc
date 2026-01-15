package cmd

import (
	"proWeb/lib/appCmds"
)

const NameApp = "scandoc-CLI"

type AppCLI struct {
	Name string
	App  *appCmds.App
}

// NewApp - конструктор для приложения
func NewApp() *AppCLI {
	return &AppCLI{Name: NameApp, App: appCmds.NewApp()}
}
