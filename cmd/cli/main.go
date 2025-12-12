package main

import (
	"errors"
	"fmt"
	"os"
	"proWeb/internal/cliUtils"
)

func main() {

}

// handleError - ловит ошибки и выводит статус выхода
func handleError(err error) {
	var appErr *cliUtils.AppError
	if errors.As(err, &appErr) {
		fmt.Println(appErr.ToString())
		os.Exit(appErr.ExitCode())
	}

	fmt.Println("непредвиденная ошибка:", err)
	os.Exit(3)
}
