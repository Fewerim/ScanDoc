package storage

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type File struct {
	Name     string `json:"name"`
	Ext      string `json:"ext"`
	Size     string `json:"size"`
	Modified string `json:"modified"`
}

func newFile(name, ext, modified, size string) File {
	return File{
		Name:     name,
		Ext:      ext,
		Size:     size,
		Modified: modified,
	}
}

// GetStorageFiles - все файлы из storage и подпапок
func GetStorageFiles(pathToStorage string) ([]File, error) {
	var results []File

	err := filepath.WalkDir(pathToStorage, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		relPath, err := filepath.Rel(pathToStorage, path)
		if err != nil {
			relPath = d.Name()
		}

		results = append(results, newFile(
			relPath,
			filepath.Ext(relPath),
			info.ModTime().Format("2006-01-02 15:04"),
			formatSize(info.Size()),
		),
		)
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(results, func(i, j int) bool {
		ti, _ := time.Parse("2006-01-02 15:04", results[i].Modified)
		tj, _ := time.Parse("2006-01-02 15:04", results[j].Modified)
		return ti.After(tj)
	})

	return results, nil
}

// ReadFileFromStorage - читает содержимое файла из локального хранилища
func ReadFileFromStorage(pathToStorage, fileName string) (string, error) {
	filePath := filepath.Join(pathToStorage, fileName)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения файла из локального хранилища")
	}
	return string(content), nil
}

// SaveFileToStorage - сохраняет файл в локальное хранилище
func SaveFileToStorage(pathToStorage, filename, content string) error {
	fullPath := filepath.Join(pathToStorage, filename)

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("ошибка сохранения файла в локальное хранилище: %v", err)
	}
	return nil
}

// DeleteFileFromStorage - удаляет файл из локального хранилища
func DeleteFileFromStorage(pathToStorage, fileName string) error {
	fullPath := filepath.Join(pathToStorage, fileName)

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("ошибка удаления файла из локального хранилища")
	}
	return nil
}

// formatSize - форматирует размер в человеческий вид
func formatSize(size int64) string {
	var suffixes = []string{"B", "KB", "MB", "GB"}
	if size == 0 {
		return "0 B"
	}
	i := int(math.Log2(float64(size)) / 10)
	return fmt.Sprintf("%.1f %s", float64(size)/math.Pow(1024, float64(i)), suffixes[i])
}
