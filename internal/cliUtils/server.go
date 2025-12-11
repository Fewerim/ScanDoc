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
	"proWeb/internal/files"
	"time"
)

const (
	timeoutStart  = 30 * time.Second
	ErrorNoPython = "no python"
	pathPython    = "python3"
)

// Result - результат выполнения CLI команды (имя созданного файла, время создания)
type Result struct {
	FileName  string
	CreatedAt time.Time
}

// Data - заглушка имитирующая некие данные, т.е файл json
type Data struct {
	Something []byte
}

// createResult - конструктор для создания результата выполнения CLI команды
func createResult(fileName string) Result {
	return Result{
		FileName:  fileName,
		CreatedAt: time.Now(),
	}
}

// ProcessOnceFile - подключение к серверу, отправка файла, обработка результата, сохранение в локальное хранилище
func ProcessOnceFile(filePath, createdNameFile string, port int) (Result, error) {
	err := startPythonServer(port)
	if err != nil {
		if err.Error() == ErrorNoPython {
			info := fmt.Sprintf("python не установлен или его нет в PATH, обратитесь к администратору")
			return Result{}, InternalError(info)
		}

		info := fmt.Sprintf("ошибка при старте сервера: %v", err)
		return Result{}, ServerError(info)
	}

	data, err := sendFileToServer(filePath, port)
	if err != nil {
		info := fmt.Sprintf("ошибка при отправке файла: %v", err)
		return Result{}, ServerError(info)
	}

	err = files.SaveFileToStorage(createdNameFile, data)
	if err != nil {
		info := fmt.Sprintf("ошибка при попытке сохранить файл: %v", err)
		return Result{}, ServerError(info)
	}

	result := createResult(createdNameFile)
	return result, nil
}

// startPythonServer - создает соединение с сервером
func startPythonServer(port int) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutStart)
	defer cancel()

	_, err := exec.LookPath(pathPython)
	if err != nil {
		return errors.New(ErrorNoPython)
	}

	cmd := exec.CommandContext(ctx, pathPython, "server.py", "--port", fmt.Sprintf("%d", port))

	if err = cmd.Start(); err != nil {
		return fmt.Errorf("не удалось запустить сервер: %v", err)
	}

	addr := fmt.Sprintf("http://localhost:%v/health", port)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			_ = cmd.Process.Kill()
			return fmt.Errorf("сервер не успел запуститься вовремя")
		case <-ticker.C:
			resp, err := http.Get(addr)
			if err == nil && resp.StatusCode == http.StatusOK {
				return nil
			}
		}
	}
}

// sendFileToServer - отправляет файл на сервер и возвращает результат выполнения сервиса (файл)
func sendFileToServer(filePath string, port int) (Data, error) {
	body, contentType, err := buildMultipartBody(filePath)
	if err != nil {
		return Data{}, err
	}

	resp, err := recognizeRequest(port, body, contentType)
	if err != nil {
		return Data{}, err
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

	part, err := w.CreateFormFile("file", filepath.Base(filePath))
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

// recognizeRequest - отправка HTTP запроса на энд-поинт recognize
func recognizeRequest(port int, body io.Reader, contentType string) (*http.Response, error) {
	url := fmt.Sprintf("http://localhost:%v/recognize", port)
	req, err := http.NewRequest("POST", url, body)
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
func decodeDataResponse(r io.Reader) (Data, error) {
	var data Data
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return Data{}, fmt.Errorf("ошибка парсинга JSON-ответа: %v", err)
	}
	return data, nil
}
