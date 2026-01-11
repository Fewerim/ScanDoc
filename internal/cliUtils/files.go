package cliUtils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	//docxFormat = ".docx" в данный момент не поддерживается
	//pdfFormat  = ".pdf"  в данный момент не поддерживается
	jpgFormat  = ".jpg"
	pngFormat  = ".png"
	jpegFormat = ".jpeg"
	//xlsxFormat = ".xlsx" в данный момент не поддерживается
)

// Для проверки поддерживаемых форматов
var allowedFormats = map[string]struct{}{
	//docxFormat: {},
	//pdfFormat:  {},
	jpgFormat:  {},
	pngFormat:  {},
	jpegFormat: {},
	//xlsxFormat: {},
}

// ValidateExtensionFile - проверяет есть ли поддержка расширения прикрепленного файла
func ValidateExtensionFile(filePath string) error {
	if filePath == "" {
		return UserError("не указано имя файла или путь")
	}

	ext := filepath.Ext(strings.ToLower(filePath))
	if _, ok := allowedFormats[ext]; !ok {
		return UserError("расширение файла не поддерживается")
	}

	return nil
}

// CheckExistsFile - проверяет наличие файла в системе
func CheckExistsFile(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return UserError("файл с таким именем не найден (проверьте путь к файлу)")
	}
	if info.IsDir() {
		return UserError("во флаге был передан только путь, без файла")
	}

	return nil
}

// CheckExistsPath - проверяет наличие пути в системе
func CheckExistsPath(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return UserError("такой путь не найден (проверьте путь)")
	}
	if !info.IsDir() {
		return UserError("во флаге был передан не путь, а что-то другое")
	}
	return nil
}

// GetFilesFromDirectory - проверяет наличие директории в системе, возвращает слайс путей к каждому файлу из этой директории и текст ошибки
func GetFilesFromDirectory(directoryPath string) ([]string, string) {
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		info := fmt.Sprintf("директория не существует: %s", directoryPath)
		return []string{}, info
	}

	fileInfo, err := os.Stat(directoryPath)
	if err != nil {
		info := fmt.Sprintf("ошибка при проверке директории %s: %v", directoryPath, err)
		return []string{}, info
	}

	if !fileInfo.IsDir() {
		info := fmt.Sprintf("указанный путь не является директорией: %s", directoryPath)
		return []string{}, info
	}

	files, err := os.ReadDir(directoryPath)
	if err != nil {
		info := fmt.Sprintf("ошибка чтения директории: %v", err)
		return []string{}, info
	}

	var filePaths []string
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(directoryPath, file.Name())
			filePaths = append(filePaths, filePath)
		}
	}

	if len(filePaths) == 0 {
		info := fmt.Sprintf("в директории нет файлов для обработки: %s", directoryPath)
		return []string{}, info
	}

	return filePaths, ""
}
