package files

import (
	"fmt"
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

func CreateFolder(nameFolder string) error {
	if err := os.Mkdir(nameFolder, 0777); err != nil {
		return fmt.Errorf("ошибка при создании папки %v", err)
	}

	return nil
}
