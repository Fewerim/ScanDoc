package main

import (
	"errors"
	"fmt"
	"os"
	"proWeb/cmd/cli/cmd"
	"proWeb/internal/cliUtils"
	"proWeb/internal/exitCodes"
	"proWeb/internal/logger"
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
	const operation = "cli.main.handleError"

	var appErr *cliUtils.AppError
	if errors.As(err, &appErr) {
		fmt.Println(appErr.ToString())
		os.Exit(appErr.ExitCode())
	}

	if log != nil {
		log.Error(operation, fmt.Sprintf("непредвиденная ошибка: %v", err), cliUtils.GetExitCode(err, exitCodes.InternalError))
	}

	msg := fmt.Sprintf("непредвиденная ошибка: %v", err)
	color.Red(msg)
	os.Exit(exitCodes.InternalError)
}

// catchPanic - ловит панику, если такая случилась
func catchPanic() {
	if r := recover(); r != nil {
		msg := fmt.Sprintf("%v", r)
		if strings.Contains(msg, "nil pointer") {
			pan := fmt.Sprintf("программа завершилась с паникой: возможно, команда получила не все аргументы. проверьте ввод")
			color.Red(pan)
		} else {
			color.Red("программа завершилась с паникой: ", r)
		}
		os.Exit(3)
	}
}

// runStandardApp - запускает приложение в обычном режиме
func runStandardApp() {
	const op = "cli.main"
	defer catchPanic()

	app := cmd.NewApp()

	if err := app.Execute(); err != nil {
		handleError(err, app.Log)
		return
	}

	if app.Log != nil {
		app.Log.Info(op, "приложение успешно завершено")
	}
}

// runConfigSetOnly - запускает приложение только для установки значений конфига
func runConfigSetOnly() {
	rootCmd := &cobra.Command{
		Use:   "scanner.exe",
		Short: "ScanDoc - CLI для распознавания документов",
	}

	rootCmd.AddCommand(cmd.NewConfigSetCmd())

	if err := rootCmd.Execute(); err != nil {
		color.Red("Ошибка: %v", err)
		os.Exit(1)
	}
}
