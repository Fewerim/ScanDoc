package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"proWeb/internal/cliUtils"
	"proWeb/internal/exitCodes"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func (a *App) clearTerminal(cmd *cobra.Command, args []string) error {
	const operation = "cli.clearTerminal"

	a.Log.Info(operation, "очистка терминала")

	if err := clearConsole(); err != nil {
		a.Log.Error(operation, fmt.Sprintf("не удалось очистить терминал: %v", err), exitCodes.InternalError)
		return cliUtils.InternalError(fmt.Sprintf("не удалось очистить терминал: %v", err))
	}

	a.Log.Info(operation, "терминал успешно очищен")
	color.Blue("терминал очищен")
	return nil
}

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
