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
)

func main() {
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
