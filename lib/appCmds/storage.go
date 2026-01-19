package appCmds

import (
	"fmt"
	"os"
	"path/filepath"
	"proWeb/internal/appUtils"
	"proWeb/internal/appUtils/command"
	"proWeb/internal/exitCodes"
	"proWeb/internal/files"
)

// OpenStorage - команда приложения. Открывает хранилище
func (a *App) OpenStorage(operation string, clearFlag bool) error {
	storagePath := a.cfg.StoragePath

	if clearFlag {
		if err := clearStorage(storagePath); err != nil {
			info := fmt.Sprintf("ошибка очистки локального хранилища: %v", err)
			a.log.Error(operation, info, exitCodes.InternalError)
			return appUtils.InternalError(info)
		}
		a.log.Info(operation, "локальное хранилище было очищено")
		appUtils.InfoMessage("Локальное хранилище было очищено")
	}

	a.log.Info(operation, "открытие локального хранилища")
	if err := openStorage(storagePath); err != nil {
		info := fmt.Sprintf("ошибка при открытии локального хранилища: %v", err)
		a.log.Error(operation, info, exitCodes.InternalError)
		return appUtils.InternalError(info)
	}

	appUtils.InfoMessage("Локальное хранилище открыто")
	return nil
}

// CheckStorageJSON - проверяет наличие локального хранилища, если его нет, создает новое
func (a *App) CheckStorageJSON() error {
	a.storage.Init(a.cfg.StoragePath)

	if !a.storage.CheckExists() {
		if err := a.storage.Create(); err != nil {
			return fmt.Errorf("ошибка создания локального хранилища: %v", err)
		}
	}
	return nil
}

// GetFilesFromStorage - команда приложения. Получает все файлы из хранилища
func (a *App) GetFilesFromStorage(operation string) ([]files.File, error) {
	a.log.Info(operation, "Получение файлов из хранилища")
	results, err := a.storage.GetFiles()
	if err != nil {
		a.log.Error(operation, err.Error(), exitCodes.InternalError)
		return nil, appUtils.InternalError(err.Error())
	}
	a.log.Info(operation, "Файлы из хранилища успешно получены")
	return results, nil
}

// ReadFileFromStorage - команда приложения. Читает содержимое файла из хранилища
func (a *App) ReadFileFromStorage(operation, fileName string) (string, error) {
	a.log.Info(operation, fmt.Sprintf("Чтение содержимого файла: %s", fileName))
	content, err := a.storage.ReadFile(fileName)
	if err != nil {
		a.log.Error(operation, err.Error(), exitCodes.InternalError)
		return "", appUtils.InternalError(err.Error())
	}
	a.log.Info(operation, fmt.Sprintf("Файл %s успешно прочитан", fileName))
	return content, nil
}

// SaveFileToStorage - команда приложения. Сохраняет файл в хранилище
func (a *App) SaveFileToStorage(operation, fileName, folder string, content interface{}) error {
	a.log.Info(operation, fmt.Sprintf("Сохранение файла %v в хранилище", fileName))
	err := a.storage.SaveFile(folder, fileName, content, files.Overwrite)
	if err != nil {
		a.log.Error(operation, err.Error(), exitCodes.InternalError)
		return appUtils.InternalError(err.Error())
	}
	a.log.Info(operation, fmt.Sprintf("Успешное сохранение файла %v в хранилище", fileName))
	return nil
}

// DeleteFileFromStorage - команда приложения. Удаляет файл из хранилища
func (a *App) DeleteFileFromStorage(operation, fileName string) error {
	a.log.Info(operation, fmt.Sprintf("Удаление файла %v из хранилища", fileName))
	err := a.storage.DeleteFile(fileName)
	if err != nil {
		a.log.Error(operation, err.Error(), exitCodes.InternalError)
		return appUtils.InternalError(err.Error())
	}
	a.log.Info(operation, fmt.Sprintf("Файл %v успешно удален из хранилища", fileName))
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

	cmd := command.Command("explorer.exe", fullPath)
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
