package appCmds

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"proWeb/internal/appUtils"
	"proWeb/internal/exitCodes"
	"proWeb/internal/storage"

	"github.com/fatih/color"
)

// OpenStorage - команда приложения. Открывает хранилище
func (a *App) OpenStorage(operation string, clearFlag bool) error {
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

// GetFilesFromStorage - команда приложения. Получает все файлы из хранилища
func (a *App) GetFilesFromStorage(operation string) ([]storage.File, error) {
	a.Log.Info(operation, "Получение файлов из хранилища")
	files, err := storage.GetStorageFiles(a.Cfg.StoragePath)
	if err != nil {
		a.Log.Error(operation, err.Error(), exitCodes.InternalError)
		return nil, appUtils.InternalError(err.Error())
	}
	a.Log.Info(operation, "Файлы из хранилища успешно получены")
	return files, nil
}

// ReadFileFromStorage - команда приложения. Читает содержимое файла из хранилища
func (a *App) ReadFileFromStorage(operation, fileName string) (string, error) {
	a.Log.Info(operation, fmt.Sprintf("Чтение содержимого файла: %s", fileName))
	content, err := storage.ReadFileFromStorage(a.Cfg.StoragePath, fileName)
	if err != nil {
		a.Log.Error(operation, err.Error(), exitCodes.InternalError)
		return "", appUtils.InternalError(err.Error())
	}
	a.Log.Info(operation, fmt.Sprintf("Файл %s успешно прочитан", fileName))
	return content, nil
}

// SaveFileToStorage - команда приложения. Сохраняет файл в хранилище
func (a *App) SaveFileToStorage(operation, fileName string, content string) error {
	a.Log.Info(operation, fmt.Sprintf("Сохранение файла %v в хранилище", fileName))
	err := storage.SaveFileToStorage(a.Cfg.StoragePath, fileName, content)
	if err != nil {
		a.Log.Error(operation, err.Error(), exitCodes.InternalError)
		return appUtils.InternalError(err.Error())
	}
	a.Log.Info(operation, fmt.Sprintf("Успешное сохранение файла %v в хранилище", fileName))
	return nil
}

// DeleteFileFromStorage - команда приложения. Удаляет файл из хранилища
func (a *App) DeleteFileFromStorage(operation, fileName string) error {
	a.Log.Info(operation, fmt.Sprintf("Удаление файла %v из хранилища", fileName))
	err := storage.DeleteFileFromStorage(a.Cfg.StoragePath, fileName)
	if err != nil {
		a.Log.Error(operation, err.Error(), exitCodes.InternalError)
		return appUtils.InternalError(err.Error())
	}
	a.Log.Info(operation, fmt.Sprintf("Файл %v успешно удален из хранилища", fileName))
	return nil
}

// openStorage - открывает локально хранилище в explorer
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
