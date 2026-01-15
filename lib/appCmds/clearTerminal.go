package appCmds

import (
	"fmt"
	"os"
	"os/exec"
	"proWeb/internal/appUtils"
	"proWeb/internal/exitCodes"
	"runtime"
)

// ClearTerminal - команда, очищающая терминал
func (a *App) ClearTerminal(operation string) error {
	a.log.Info(operation, "очистка терминала")

	if err := clearConsole(); err != nil {
		info := fmt.Sprintf("не удалось очистить терминал: %v", err)
		a.log.Error(operation, info, exitCodes.InternalError)
		return appUtils.InternalError(info)
	}

	a.log.Info(operation, "терминал успешно очищен")
	appUtils.InfoMessage("Терминал очищен")
	return nil
}

// clearConsole - очищает консоль от текста
func clearConsole() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
