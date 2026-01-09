package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"proWeb/internal/cliUtils"
	"proWeb/internal/exitCodes"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// logs - команда, открывающая папку с логами
func (a *App) logs(clearFlag bool) error {
	const operation = "cli.openLog"
	if clearFlag {
		if err := clearLog(a.Cfg.LogPath); err != nil {
			info := fmt.Sprintf("ошибка удаления: %v", err)
			a.Log.Error(operation, info, exitCodes.InternalError)
			return cliUtils.InternalError(info)
		}
		a.Log.Info(operation, "файла с логами был очищен")
	}

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

// clearLog - удаляет файл с логами (создается сразу новый)
func clearLog(pathToLog string) error {
	fullPath, err := filepath.Abs(pathToLog)
	if err != nil {
		return fmt.Errorf("неверный путь к файлу с логами")
	}

	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("не удалось обрезать лог: %v", err)
	}
	f.Close()

	return nil
}

func newOpenLogCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "open_log",
		Short:   "Открывает файл с логами",
		Example: "scandoc.exe openLog",
		RunE: func(cmd *cobra.Command, args []string) error {
			clearFlag, err := cmd.Flags().GetBool("clear")
			if err != nil {
				return fmt.Errorf("ошибка чтения флага: %v", err)
			}
			return a.logs(clearFlag)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().BoolP("clear", "c", false, "очистка файла с логами")
	return cmd
}
