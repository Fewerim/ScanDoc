package files

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

/*
Пакет используется для работы с файлами JSON и локальным хранилищем
*/
type TestJson struct {
	Name       string   `json:"name"`
	Age        int      `json:"age"`
	Profession string   `json:"profession"`
	Skills     []string `json:"skills"`
}

// Название локального хранилища
// В будущем добавить чтобы .gitignore сам обновлялся если меняется nameStorage
const nameStorage = "storageJSONs"

// CreateStorageJSON - создает локальное хранилище для хранения JSON
func CreateStorageJSON() error {
	if err := os.Mkdir(nameStorage, 0777); err != nil {
		return fmt.Errorf("error creating storage JSON directory: %v", err)
	}
	log.Printf("Successfully created %s directory", nameStorage)
	return nil
}

// StorageExists - проверяет наличие локального хранилища
func StorageExists() bool {
	_, err := os.Stat(nameStorage)
	return !os.IsNotExist(err)
}

// SaveFileToDirectory - сохраняет файл без перезаписи, если такой есть (создается уникальный)
func SaveFileToDirectory(fileName, directory string, data interface{}) error {
	return saveFileToDirectory(fileName, directory, data, false)
}

func SaveFileToStorage(fileName string, data interface{}) error {
	return saveFileToDirectory(fileName, "", data, false)
}

// OverwriteFileToDirectory - сохраняет файл c перезаписью, если такой есть (перезаписывает существующий)
func OverwriteFileToDirectory(fileName, directory string, data interface{}) error {
	return saveFileToDirectory(fileName, directory, data, true)
}

// DeleteFileFromDirectory - удаляет файл по пути
func DeleteFileFromDirectory(filePath string) error {
	if err := os.RemoveAll(filePath); err != nil {
		return fmt.Errorf("error deleting file %s: %v", filePath, err)
	}
	log.Printf("Successfully deleted file %s", filePath)
	return nil
}

// ReadJSONFile - читает файл по пути
func ReadJSONFile(filePath string, result interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	if err = json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	return nil
}

// GetListFilesInDirectory - выдает весь список файлов лежащих в хранилище
func GetListFilesInDirectory(directory string) ([]string, error) {
	fullPath := filepath.Join(nameStorage, directory)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %v", directory, err)
	}

	fileNames := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}
	return fileNames, nil
}
