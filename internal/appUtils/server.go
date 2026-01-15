package appUtils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	parser2 "proWeb/internal/parser"
	"proWeb/internal/typesJSON/typesUtils"
	"time"
)

const (
	serverStartTimeout       = 15 * time.Second
	healthCheckAttempts      = 15
	healthCheckInterval      = 1 * time.Second
	healthCheckClientTimeout = 5 * time.Second
	ErrorNoPython            = "python не найден"
)

// StartPythonServer - создает соединение с сервером
func StartPythonServer(port int, pythonExecutable, pathToScript string) (*exec.Cmd, error) {
	if _, err := exec.LookPath(pythonExecutable); err != nil {
		return nil, InternalError(ErrorNoPython)
	}
	portStr := fmt.Sprintf("%d", port)
	pathToScript = filepath.Join(pathToScript, "src/run_api.py")

	cmd := exec.Command(pythonExecutable, pathToScript, "--port", portStr)

	if err := cmd.Start(); err != nil {
		info := fmt.Sprintf("не удалось запустить сервер: %v", err)
		return nil, ServerError(info)
	}

	ctx, cancel := context.WithTimeout(context.Background(), serverStartTimeout)
	defer cancel()

	if err := healthCheck(ctx, port); err != nil {
		cmd.Process.Kill()
		return nil, ServerError(err.Error())
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
func SendFileToServer(filePath string, port int) (interface{}, string, error) {
	var data interface{}

	resp, err := scanRequest(port, filePath)
	if err != nil {
		return data, "", InternalError(err.Error())
	}
	defer resp.Body.Close()

	t, err := typesUtils.GetDoctype(resp)
	if err != nil {
		return data, "", InternalError(err.Error())
	}

	//TODO: получить из заголовков ответа тип документа и использовать подходящую структуру для сериализации

	decodeData, err := decodeDataResponse(resp.Body)
	if err != nil {
		return data, "", InternalError(err.Error())
	}

	decodeDataWithTable, err := parser2.UpdateTableInData(decodeData)
	if err != nil {
		if !errors.Is(err, parser2.ErrNotFoundTable) {
			return data, "", InternalError(err.Error())
		}
		decodeDataWithTable = decodeData
	}

	if t == "torg12" {
		decodeDataWithTable, err = parser2.AddTotalFields(decodeDataWithTable)
		if err != nil {
			return data, "", InternalError(err.Error())
		}
	}

	return decodeDataWithTable, t, nil
}

// scanRequest - отправка HTTP запроса на энд-поинт scan
func scanRequest(port int, filePath string) (*http.Response, error) {
	url := fmt.Sprintf("http://localhost:%v/scan?url=%s", port, url.QueryEscape(filePath))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка формирования запроса: %v", err)
	}
	//req.Header.Set("Content-Type", contentType)

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
	client := &http.Client{Timeout: healthCheckClientTimeout}

	maxAttempts := healthCheckAttempts

	for i := 0; i < maxAttempts; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		url := fmt.Sprintf("http://localhost:%v/health", port)

		resp, err := client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			if resp.Body != nil {
				resp.Body.Close()
			}

			return nil
		}
		time.Sleep(healthCheckInterval)
	}

	return fmt.Errorf("сервер не успел запуститься за %d попыток, попробуйте еще раз", maxAttempts)
}

// decodeDataResponse - парсинг JSON ответа
func decodeDataResponse(r io.Reader) (interface{}, error) {
	var data interface{}

	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return data, fmt.Errorf("ошибка парсинга JSON-ответа: %v", err)
	}
	return data, nil
}
