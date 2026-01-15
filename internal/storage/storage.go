package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"proWeb/internal/appUtils"
	"proWeb/internal/files"
)

var defaultNameStorage = "storageJSONs"

// InitStorage - инициализирует название локального хранилище
func InitStorage(storage string) {
	if storage != "" {
		defaultNameStorage = storage
	}
}

// CreateStorageJSON - создает локальное хранилище для хранения обработанных файлов
func CreateStorageJSON() error {
	if err := files.CreateFolder(defaultNameStorage); err != nil {
		return err
	}
	appUtils.InfoMessage(fmt.Sprintf("Хранилище для обработанных файлов (%s) успешно создано", defaultNameStorage))
	return nil
}

// CheckStorageExists - проверяет наличие локального хранилища
func CheckStorageExists() bool {
	_, err := os.Stat(defaultNameStorage)
	return !os.IsNotExist(err)
}

// GetStorageFiles - достает все файлы из хранилища и подпапок
func GetStorageFiles(pathToStorage string) ([]files.File, error) {
	results, err := files.GetFilesFromDirectory(pathToStorage)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения файлов из хранилища")
	}
	return results, nil
}

// ReadFileFromStorage - читает содержимое файла из локального хранилища
func ReadFileFromStorage(pathToStorage, fileName string) (string, error) {
	filePath := filepath.Join(pathToStorage, fileName)

	content, err := files.ReadFileFromDirectory(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// SaveFileToStorage - сохраняет файл в локальное хранилище
func SaveFileToStorage(pathToStorage, folder, filename string, content interface{}, overwrite bool) error {
	fullPath := filepath.Join(pathToStorage, folder)

	if err := files.SaveFileToDirectory(filename, fullPath, content, overwrite); err != nil {
		return fmt.Errorf("ошибка сохранения файла в локальное хранилище: %v", err)
	}
	return nil
}

// DeleteFileFromStorage - удаляет файл из локального хранилища
func DeleteFileFromStorage(pathToStorage, fileName string) error {
	fullPath := filepath.Join(pathToStorage, fileName)

	if err := files.DeleteFileFromDirectory(fullPath); err != nil {
		return fmt.Errorf("ошибка удаления файла из локального хранилища")
	}
	return nil
}
