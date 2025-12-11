package main

import (
	"errors"
	"fmt"
	"os"
	"proWeb/cmd/cli/cmd"
	"proWeb/internal/cliUtils"
)

func main() {
	if err := cmd.Execute(); err != nil {
		var appErr *cliUtils.AppError
		if errors.As(err, &appErr) {
			fmt.Println(appErr.ToString())
			os.Exit(appErr.ExitCode())
		}

		fmt.Println("непредвиденная ошибка:", err)
		os.Exit(3)
	}
}
