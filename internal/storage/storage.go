package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"proWeb/internal/appUtils"
	"proWeb/internal/files"
)

const defaultStorage = "storageJSONs"

type Storage struct {
	BasePath string
}

func New(path string) *Storage {
	if path == "" {
		path = defaultStorage
	}

	return &Storage{BasePath: path}
}

// Init - инициализирует хранилище
func (s *Storage) Init(path string) {
	if path != "" {
		s.BasePath = path
	}
}

// Create - создает хранилище
func (s *Storage) Create() error {
	if err := files.CreateFolder(s.BasePath); err != nil {
		return err
	}
	appUtils.InfoMessage(fmt.Sprintf("Хранилище для обработанных файлов (%s) успешно создано", s.BasePath))
	return nil
}

// CheckExists - проверяет наличие хранилища
func (s *Storage) CheckExists() bool {
	_, err := os.Stat(s.BasePath)
	return !os.IsNotExist(err)
}

// ReadFile - читает файл из хранилища
func (s *Storage) ReadFile(fileName string) (string, error) {
	filePath := filepath.Join(s.BasePath, fileName)
	content, err := files.ReadFileFromDirectory(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// SaveFile - сохраняет файл в хранилище
func (s *Storage) SaveFile(folder, filename string, content any, overwrite bool) error {
	fullPath := filepath.Join(s.BasePath, folder)

	if err := files.SaveFileToDirectory(filename, fullPath, content, overwrite); err != nil {
		return fmt.Errorf("ошибка сохранения файла в локальное хранилище: %v", err)
	}
	return nil
}

// GetFiles - достает все файлы из хранилища и подпапок
func (s *Storage) GetFiles() ([]files.File, error) {
	results, err := files.GetFilesFromDirectory(s.BasePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения файлов из хранилища")
	}
	return results, nil
}

// DeleteFile - удаляет файл из локального хранилища
func (s *Storage) DeleteFile(fileName string) error {
	fullPath := filepath.Join(s.BasePath, fileName)

	if err := files.DeleteFileFromDirectory(fullPath); err != nil {
		return fmt.Errorf("ошибка удаления файла из локального хранилища")
	}
	return nil
}
