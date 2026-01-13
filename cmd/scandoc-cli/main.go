package main

import (
	"errors"
	"fmt"
	"os"
	"proWeb/cmd/scandoc-cli/cmd"
	"proWeb/lib/appUtils"
	"proWeb/lib/exitCodes"
	"proWeb/lib/logger"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "config_set" {
		runConfigSetOnly()
		return
	}

	runStandardApp()
}

// handleError - ловит ошибки и выводит статус выхода
func handleError(err error, log logger.Logger) {
	const operation = "scandoc-CLI.main.handleError"

	var appErr *appUtils.AppError
	if errors.As(err, &appErr) {
		fmt.Println(appErr.ToString())
		os.Exit(appErr.ExitCode())
	}

	if log != nil {
		log.Error(operation, fmt.Sprintf("непредвиденная ошибка: %v", err), appUtils.GetExitCode(err, exitCodes.InternalError))
	}

	msg := fmt.Sprintf("непредвиденная ошибка: %v", err)
	color.Red(msg)
	os.Exit(exitCodes.InternalError)
}

// catchPanic - ловит панику, если такая случилась
func catchPanic() {
	if r := recover(); r != nil {
		var msg string
		if err, ok := r.(error); ok {
			msg = err.Error()
		} else {
			msg = fmt.Sprintf("%v", r)
		}
		msg = strings.TrimSuffix(msg, ": Access is denied.")

		if strings.Contains(msg, "nil pointer") {
			color.Red("программа завершилась с паникой: возможно, команда получила не все аргументы. проверьте ввод")
		} else {
			color.Red("программа завершилась с паникой: " + msg)
		}
		os.Exit(3)
	}
}

// runStandardApp - запускает приложение в обычном режиме
func runStandardApp() {
	const op = ".main"
	defer catchPanic()

	a := cmd.NewApp()

	if err := a.Execute(); err != nil {
		handleError(err, a.App.Log)
		return
	}

	if a.App.Log != nil {
		a.App.Log.Info(a.Name+op, "приложение успешно завершено")
	}
}

// runConfigSetOnly - запускает приложение только для установки значений конфига
func runConfigSetOnly() {
	rootCmd := &cobra.Command{
		Use:   "scandoc.exe",
		Short: "ScanDoc - CLI для распознавания документов",
	}

	rootCmd.AddCommand(cmd.NewConfigSetCmd())

	if err := rootCmd.Execute(); err != nil {
		color.Red("Ошибка: %v", err)
		os.Exit(1)
	}
}
