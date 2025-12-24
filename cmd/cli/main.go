package main

import (
	"errors"
	"fmt"
	"os"
	"proWeb/cmd/cli/cmd"
	"proWeb/internal/cliUtils"
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

	app.Log.Info(op, "приложение успешно завершено")
}

// handleError - ловит ошибки и выводит статус выхода
func handleError(err error, log logger.Logger) {
	const operation = "cli.main.handleError"

	var appErr *cliUtils.AppError
	if errors.As(err, &appErr) {
		log.Error(operation, err.Error(), appErr.ExitCode())
		fmt.Println(appErr.ToString())
		os.Exit(appErr.ExitCode())
	}

	log.Error(operation, fmt.Sprintf("непредвиденная ошибка: %v", err), 3)
	fmt.Println("непредвиденная ошибка:", err)
	os.Exit(3)
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
