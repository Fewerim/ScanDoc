package files

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// addExtensionJSON - добавляет расширение .json к названию файла, если его нет
func addExtensionJSON(filename string) string {
	if !strings.HasSuffix(filename, ".json") {
		return filename + ".json"
	}
	return filename
}

// getUniqFileName - создает и возвращает уникальное имя файла и путь
func getUniqFileName(fileName, fullDirectory string, overwrite bool) (string, string) {
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	ext := filepath.Ext(fileName)
	filePath := filepath.Join(fullDirectory, fileName)

	if overwrite {
		return fileName, filePath
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fileName, filePath
	}

	for i := 1; ; i++ {
		newFileName := fmt.Sprintf("%s_(%d)%s", baseName, i, ext)
		newFilePath := filepath.Join(fullDirectory, newFileName)
		if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
			return newFileName, newFilePath
		}
	}
}

// saveFileToDirectory - создает файл и сохраняет/перезаписывает в хранилище в локальной папке
func saveFileToDirectory(fileName, directory string, data interface{}, overWrite bool) error {
	fileNameWithExtension := addExtensionJSON(fileName)
	fullDirectory := filepath.Join(nameStorage, directory)

	fileName, filePath := getUniqFileName(fileNameWithExtension, fullDirectory, overWrite)

	if err := os.MkdirAll(fullDirectory, 0777); err != nil {
		return fmt.Errorf("ошибка создания %s директории: %v", directory, err)
	}

	jsonData, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return fmt.Errorf("ошибка маршаллинга json: %v", err)
	}

	if err = os.WriteFile(filePath, jsonData, 0666); err != nil {
		return fmt.Errorf("ошибка записи файла: %v", err)
	}

	action := "created"
	if overWrite {
		action = "overwritten"
	}
	log.Printf("%s файл %s находится в директории %s", action, fileName, fullDirectory)
	return nil
}
