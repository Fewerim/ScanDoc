package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"proWeb/internal/cliUtils"
	"proWeb/internal/exitCodes"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// logs - команда, открывающая папку с логами
func (a *App) logs(cmd *cobra.Command, args []string) error {
	const operation = "cli.openLog"

	a.Log.Info(operation, "открытие файла с логами")
	if err := openLog(a.Cfg.LogPath); err != nil {
		info := fmt.Sprintf("ошибка при открытии файла с логами: %v", err)
		a.Log.Error(operation, info, exitCodes.InternalError)
		return cliUtils.InternalError(info)
	}
	color.Blue("Файл с логами открыт")
	return nil
}

// openLog - открывает папку с логами
func openLog(pathToLog string) error {
	fullPath, err := filepath.Abs(pathToLog)
	if err != nil {
		return fmt.Errorf("неверный путь к файлу с логами")
	}

	cmd := exec.Command("notepad.exe", fullPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("не удалось открыть: %v", err)
	}

	return nil
}

func newOpenLogCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "open_log",
		Short:         "открывает файл с логами",
		Example:       "scandoc.exe openLog",
		RunE:          a.logs,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	return cmd
}
