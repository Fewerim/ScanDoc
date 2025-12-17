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
	timeoutStart  = 180 * time.Second
	ErrorNoPython = "no python"
	serverPath    = "internal/service/scanPy/src/run_api.py"
)

// Data - заглушка имитирующая некие данные, т.е файл json
type Data struct {
	Filename    string            `json:"filename"`
	DocType     string            `json:"doc_type"`
	ProcessedAt string            `json:"processed_at"`
	Fields      map[string]string `json:"fields"`
}

// StartPythonServer - создает соединение с сервером
func StartPythonServer(port int, pythonExecutable string) (*exec.Cmd, error) {
	if _, err := exec.LookPath(pythonExecutable); err != nil {
		return nil, errors.New(ErrorNoPython)
	}

	cmd := exec.Command(pythonExecutable, serverPath, "--port", fmt.Sprintf("%d", port))

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("не удалось запустить сервер: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutStart)
	defer cancel()

	addr := fmt.Sprintf("http://127.0.0.1:%d/health", port)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			_ = cmd.Process.Kill()
			return nil, fmt.Errorf("сервер не успел запуститься вовремя")
		case <-ticker.C:
			resp, err := http.Get(addr)
			if err == nil && resp.StatusCode == http.StatusOK {
				return cmd, nil
			}
		}
	}
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

	body, contentType, err := buildMultipartBody(filePath)
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
func buildMultipartBody(filePath string) (body *bytes.Buffer, contentType string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer f.Close()

	body = &bytes.Buffer{}
	w := multipart.NewWriter(body)

	//TODO: file | image
	part, err := w.CreateFormFile("image", filepath.Base(filePath))
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
	url := fmt.Sprintf("http://127.0.0.1:%v/scan", port)
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

// decodeDataResponse - парсинг JSON ответа
func decodeDataResponse(r io.Reader) (interface{}, error) {
	var data interface{}
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return data, fmt.Errorf("ошибка парсинга JSON-ответа: %v", err)
	}
	return data, nil
}
