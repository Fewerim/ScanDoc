package appCmds

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"proWeb/internal/appUtils"
	"proWeb/internal/exitCodes"
)

// OpenLogFolder - команда, открывающая папку с логами
func (a *App) OpenLogFolder(operation string, clearFlag bool) error {
	if clearFlag {
		if err := clearLog(a.cfg.LogPath); err != nil {
			info := fmt.Sprintf("ошибка удаления: %v", err)
			a.log.Error(operation, info, exitCodes.InternalError)
			return appUtils.InternalError(info)
		}
		a.log.Info(operation, "файл с логами был очищен")
		appUtils.InfoMessage("Файл с логами был очищен")
	}

	a.log.Info(operation, "открытие файла с логами")
	if err := openLog(a.cfg.LogPath); err != nil {
		info := fmt.Sprintf("ошибка при открытии файла с логами: %v", err)
		a.log.Error(operation, info, exitCodes.InternalError)
		return appUtils.InternalError(info)
	}

	appUtils.InfoMessage("Файл с логами открыт")
	return nil
}

// CleanLogFolder - очистка папки с логами
func (a *App) CleanLogFolder(operation string) error {
	if err := clearLog(a.cfg.LogPath); err != nil {
		info := fmt.Sprintf("ошибка удаления: %v", err)
		a.log.Error(operation, info, exitCodes.InternalError)
		return appUtils.InternalError(info)
	}
	a.log.Info(operation, "файл с логами был очищен")
	appUtils.InfoMessage("Файл с логами был очищен")
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
