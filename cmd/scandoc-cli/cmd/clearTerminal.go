package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"proWeb/lib/cliUtils"
	"proWeb/lib/exitCodes"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// clearTerminal - команда, очищающая терминал
func (a *App) clearTerminal(cmd *cobra.Command, args []string) error {
	const operation = "scandoc-cli.clearTerminal"

	a.Log.Info(operation, "очистка терминала")

	if err := clearConsole(); err != nil {
		info := fmt.Sprintf("не удалось очистить терминал: %v", err)
		a.Log.Error(operation, info, exitCodes.InternalError)
		return cliUtils.InternalError(info)
	}

	a.Log.Info(operation, "терминал успешно очищен")
	color.Blue("Терминал очищен")
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

func newClearCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:           "clear_terminal",
		Short:         "Очистить терминал",
		Long:          "Очищает экран терминала (работает в CMD, PowerShell, bash, zsh)",
		Example:       "scandoc.exe clear\nочистит экран терминала",
		Args:          cobra.NoArgs,
		RunE:          a.clearTerminal,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
}
