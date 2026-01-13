package appCmds

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"proWeb/internal/appUtils"
	"proWeb/internal/exitCodes"

	"github.com/fatih/color"
)

// Logs - команда, открывающая папку с логами
func (a *App) Logs(operation string, clearFlag bool) error {
	if clearFlag {
		if err := clearLog(a.Cfg.LogPath); err != nil {
			info := fmt.Sprintf("ошибка удаления: %v", err)
			a.Log.Error(operation, info, exitCodes.InternalError)
			return appUtils.InternalError(info)
		}
		a.Log.Info(operation, "файла с логами был очищен")
		color.Blue("Файл с логами был очищен")
	}

	a.Log.Info(operation, "открытие файла с логами")
	if err := openLog(a.Cfg.LogPath); err != nil {
		info := fmt.Sprintf("ошибка при открытии файла с логами: %v", err)
		a.Log.Error(operation, info, exitCodes.InternalError)
		return appUtils.InternalError(info)
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
