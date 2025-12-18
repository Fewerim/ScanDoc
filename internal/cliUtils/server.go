package cliUtils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	ProcessTimeout           = 180 * time.Second
	ServerStartTimeout       = 10 * time.Second
	HealthCheckAttempts      = 10
	HealthCheckInterval      = 1 * time.Second
	HealthCheckClientTimeout = 5 * time.Second
	ErrorNoPython            = "python не найден"
)

// Data - заглушка имитирующая некие данные, т.е файл json
type Data struct {
	Filename    string            `json:"filename"`
	DocType     string            `json:"doc_type"`
	ProcessedAt string            `json:"processed_at"`
	Fields      map[string]string `json:"fields"`
}

// StartPythonServer - создает соединение с сервером
func StartPythonServer(port int, pythonExecutable, pathToScript string) (*exec.Cmd, error) {
	if _, err := exec.LookPath(pythonExecutable); err != nil {
		return nil, errors.New(ErrorNoPython)
	}
	portStr := fmt.Sprintf("%d", port)

	cmd := exec.Command(pythonExecutable, pathToScript, "--port", portStr)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("не удалось запустить сервер: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), ServerStartTimeout)
	defer cancel()

	if err := healthCheck(ctx, port); err != nil {
		cmd.Process.Kill()
		return nil, err
	}

	return cmd, nil
}

// KillServer - принудительное закрытие соединения
func KillServer(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}
}

// SendFileToServer - отправляет файл на сервер и возвращает результат выполнения сервиса (файл)
func SendFileToServer(filePath string, port int) (interface{}, error) {
	var data interface{}
	fieldNameOfRequest := setupFieldName(filePath)

	body, contentType, err := buildMultipartBody(filePath, fieldNameOfRequest)
	if err != nil {
		return data, err
	}

	resp, err := scanRequest(port, body, contentType)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	//TODO: получить из заголовков ответа тип документа и использовать подходящую структуру для сериализации

	return decodeDataResponse(resp.Body)
}

// buildMultipartBody - создание тела для запроса
func buildMultipartBody(filePath, fieldName string) (body *bytes.Buffer, contentType string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer f.Close()

	body = &bytes.Buffer{}
	w := multipart.NewWriter(body)

	//TODO: file | image
	part, err := w.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return nil, "", fmt.Errorf("ошибка форматирования multipart: %v", err)
	}
	if _, err = io.Copy(part, f); err != nil {
		return nil, "", fmt.Errorf("ошибка чтения файла: %v", err)
	}
	if err = w.Close(); err != nil {
		return nil, "", fmt.Errorf("ошибка закрытия multipart: %v", err)
	}

	return body, w.FormDataContentType(), nil
}

// scanRequest - отправка HTTP запроса на энд-поинт scan
func scanRequest(port int, body io.Reader, contentType string) (*http.Response, error) {
	url := fmt.Sprintf("http://localhost:%v/scan", port)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("ошибка формирования запроса: %v", err)
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к серверу: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("сервер вернул статус: %v", resp.StatusCode)
	}

	return resp, nil
}

// healthCheck - проверяет запустился ли сервер, если нет, возвращает ошибку
func healthCheck(ctx context.Context, port int) error {
	client := &http.Client{Timeout: HealthCheckClientTimeout}

	maxAttempts := HealthCheckAttempts

	for i := 0; i < maxAttempts; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		url := fmt.Sprintf("http://localhost:%d/health", port)

		resp, err := client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			if resp.Body != nil {
				resp.Body.Close()
			}

			return nil
		}
		time.Sleep(HealthCheckInterval)
	}

	return fmt.Errorf("сервер не успел запуститься за %d попыток", maxAttempts)
}

// decodeDataResponse - парсинг JSON ответа
func decodeDataResponse(r io.Reader) (interface{}, error) {
	var data interface{}

	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return data, fmt.Errorf("ошибка парсинга JSON-ответа: %v", err)
	}
	return data, nil
}

// setupFieldName - возвращает тип файла для поля запроса (image или file)
func setupFieldName(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg", ".png":
		return "image"
	default:
		return "file"
	}
}
