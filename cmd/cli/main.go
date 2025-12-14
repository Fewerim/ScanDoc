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
	log, err := logger.NewFileLog("app.log")
	if err != nil {
		log = logger.New(os.Stdout)
		handleError(err, log)
	}

	log.Info("main", "старт приложения")

	app := cmd.NewApp(log)

	if errs := app.Execute(); errs != nil {
		handleError(errs, log)
	}

	log.Info("main", "приложения успешно завершено")
}

// handleError - ловит ошибки и выводит статус выхода
func handleError(err error, log logger.Logger) {
	const operation = "main.handleError"

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
