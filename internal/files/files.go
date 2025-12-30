package files

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
var nameStorage = "storageJSONs"

func InitStorage(storage string) {
	if storage != "" {
		nameStorage = storage
	}
}

// CreateStorageJSON - создает локальное хранилище для хранения JSON
func CreateStorageJSON() error {
	if err := os.Mkdir(nameStorage, 0777); err != nil {
		return fmt.Errorf("ошибка при создании хранилища для JSON файлов: %v", err)
	}
	log.Printf("Хранилище для обработанных файлов (%s) успешно создано", nameStorage)
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

// SaveFileToStorage - сохраняет файл без перезаписи в локальное хранилище
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
		return fmt.Errorf("ошибка при удалении файла %s: %v", filePath, err)
	}
	log.Printf("файл по пути %s был успешно удален", filePath)
	return nil
}

// ReadJSONFile - читает файл по пути
func ReadJSONFile(filePath string, result interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("ошибка при чтении файла: %v", err)
	}

	if err = json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	return nil
}

// GetListFilesInDirectory - выдает весь список файлов лежащих в хранилище
func GetListFilesInDirectory(directory string) ([]string, error) {
	fullPath := filepath.Join(nameStorage, directory)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении директории %s: %v", directory, err)
	}

	fileNames := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}
	return fileNames, nil
}

// DownloadFile - скачивает файл с url и создает в target пути
func DownloadFile(url, target string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// ProjectRoot - возвращает путь к корню проекта, если не удалось получить путь возвращает ошибку
func ProjectRoot() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	binDir := filepath.Dir(exePath)
	return filepath.Dir(binDir), nil
}
