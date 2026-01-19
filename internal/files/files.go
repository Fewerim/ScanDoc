package files

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	Overwrite    = true
	NotOverwrite = false
)

type File struct {
	Name     string `json:"name"`
	Ext      string `json:"ext"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
}

func newFile(name, ext, modified string, size int64) File {
	return File{
		Name:     name,
		Ext:      ext,
		Size:     size,
		Modified: modified,
	}
}

// SaveFileToDirectory - создает файл и сохраняет/перезаписывает в хранилище в локальной папке
func SaveFileToDirectory(fileName, directory string, data interface{}, overwrite bool) error {
	fileNameWithExtension := addExtensionJSON(fileName)

	fileName, filePath := getUniqFileName(fileNameWithExtension, directory, overwrite)

	if err := os.MkdirAll(directory, 0777); err != nil {
		return fmt.Errorf("ошибка создания %s директории: %v", directory, err)
	}

	jsonData, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return fmt.Errorf("ошибка маршаллинга json: %v", err)
	}

	if err = os.WriteFile(filePath, jsonData, 0666); err != nil {
		return fmt.Errorf("ошибка записи файла: %v", err)
	}
	return nil
}

// DeleteFileFromDirectory - удаляет файл по пути
func DeleteFileFromDirectory(filePath string) error {
	if err := os.RemoveAll(filePath); err != nil {
		return fmt.Errorf("ошибка при удалении файла %s: %v", filePath, err)
	}
	return nil
}

// ReadFileFromDirectory - читает файл из необходимой директории
func ReadFileFromDirectory(filePath string) (content []byte, err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении файла: %v", err)
	}
	return data, nil
}

// GetFilesFromDirectory - выдает весь список файлов лежащих в хранилище
func GetFilesFromDirectory(directory string) ([]File, error) {
	var results []File

	err := filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
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

		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			relPath = d.Name()
		}

		results = append(results, newFile(
			relPath,
			filepath.Ext(relPath),
			info.ModTime().Format("2006-01-02 15:04"),
			info.Size(),
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

// FindProjectRoot - возвращает путь к корню проекта, если не удалось получить путь возвращает ошибку
func FindProjectRoot(start string) (string, error) {
	dir := start
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			projectRoot, err := filepath.Abs(dir)
			if err != nil {
				return "", err
			}
			return projectRoot, nil
		}
		dir = filepath.Dir(dir)
	}
	return "", fmt.Errorf("корень проекта не найден")
}
