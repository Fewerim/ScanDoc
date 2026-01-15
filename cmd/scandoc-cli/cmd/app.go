package cmd

import (
	"proWeb/lib/appCmds"
)

const NameApp = "scandoc-CLI"

type AppCLI struct {
	Name    string
	AppCmds appCmds.AppCmds
}

// NewApp - конструктор для приложения
func NewApp() *AppCLI {
	return &AppCLI{Name: NameApp, AppCmds: appCmds.NewApp()}
}
