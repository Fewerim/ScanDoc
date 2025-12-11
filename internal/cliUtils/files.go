package cliUtils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	docxFormat = ".docx"
	pdfFormat  = ".pdf"
	jpgFormat  = ".jpg"
	pngFormat  = ".png"
	xlsxFormat = ".xlsx"
)

// Для проверки поддерживаемых форматов
var allowedFormats = map[string]struct{}{
	docxFormat: {},
	pdfFormat:  {},
	jpgFormat:  {},
	pngFormat:  {},
	xlsxFormat: {},
}

// ValidateExtensionFile - проверяет есть ли поддержка расширения прикрепленного файла
func ValidateExtensionFile(filePath string) error {
	if filePath == "" {
		return errors.New("не указано имя файла или путь")
	}

	ext := filepath.Ext(strings.ToLower(filePath))
	if _, ok := allowedFormats[ext]; !ok {
		return errors.New("расширение файла не поддерживается")
	}

	return nil
}

func CheckExistsFile(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return errors.New("файл с таким именем не найден (проверьте путь к файлу)")
	}
	if info.IsDir() {
		return errors.New("во флаге был передан только путь, без файла")
	}

	return nil
}
