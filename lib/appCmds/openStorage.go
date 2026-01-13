package appCmds

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"proWeb/lib/appUtils"
	"proWeb/lib/exitCodes"

	"github.com/fatih/color"
)

// Storage - команда, открывающая папку storage
func (a *App) Storage(operation string, clearFlag bool) error {
	storagePath := a.Cfg.StoragePath

	if clearFlag {
		if err := clearStorage(storagePath); err != nil {
			info := fmt.Sprintf("ошибка очистки локального хранилища: %v", err)
			a.Log.Error(operation, info, exitCodes.InternalError)
			return appUtils.InternalError(info)
		}
		a.Log.Info(operation, "локальное хранилище было очищено")
		color.Blue("Локальное хранилище было очищено")
	}

	a.Log.Info(operation, "открытие локального хранилища")
	if err := openStorage(storagePath); err != nil {
		info := fmt.Sprintf("ошибка при открытии локального хранилища: %v", err)
		a.Log.Error(operation, info, exitCodes.InternalError)
		return appUtils.InternalError(info)
	}

	color.Blue("Локальное хранилище открыто")
	return nil
}

func openStorage(storagePath string) error {
	fullPath, err := filepath.Abs(storagePath)
	if err != nil {
		return fmt.Errorf("неверный путь к папке storage")
	}

	if err := appUtils.CheckExistsPath(fullPath); err != nil {
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("не удалось создать папку storage: %v", err)
		}
	}

	cmd := exec.Command("explorer.exe", fullPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("не удалось открыть папку storage: %v", err)
	}
	return nil
}

// clearStorage - очищает папку storage
func clearStorage(storagePath string) error {
	fullPath, err := filepath.Abs(storagePath)
	if err != nil {
		return fmt.Errorf("неверный путь к папке storage")
	}

	if err := appUtils.CheckExistsPath(fullPath); err != nil {
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("не удалось создать папку storage: %v", err)
		}
		return nil
	}

	dir, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("не удалось открыть папку storage: %v", err)
	}
	defer dir.Close()

	files, err := dir.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("не удалось прочитать содержимое папки storage: %v", err)
	}

	for _, file := range files {
		filePath := filepath.Join(fullPath, file)

		info, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		if info.IsDir() {
			if err := os.RemoveAll(filePath); err != nil {
				return fmt.Errorf("не удалось удалить подпапку %s: %v", file, err)
			}
		} else {
			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("не удалось удалить файл %s: %v", file, err)
			}
		}
	}

	return nil
}
