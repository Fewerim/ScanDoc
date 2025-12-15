package main

import (
	"errors"
	"fmt"
	"os"
	"proWeb/cmd/cli/cmd"
	"proWeb/internal/cliUtils"
	"proWeb/internal/logger"
)

func main() {
	const op = "cli.main"
	defer catchPanic()

	app := cmd.NewApp()
	if err := app.Execute(); err != nil {
		handleError(err, app.Log)
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

func catchPanic() {
	if r := recover(); r != nil {
		fmt.Println("программа завершилась с паникой: ", r)
		os.Exit(3)
	}
}
