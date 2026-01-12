package cliWorks

import (
	"context"
	"fmt"
	"proWeb/lib/cliUtils"
	"proWeb/lib/config"
	"proWeb/lib/files"
	"time"
)

const processTimeout = 180 * time.Second

// ProcessOnceFile - подключение к серверу, отправка файла, обработка результата, сохранение в локальное хранилище
func ProcessOnceFile(filePath, createdNameFile string, cfg *config.Config) (cliUtils.OnceProcessResult, error) {
	_, cancel := context.WithTimeout(context.Background(), processTimeout)
	defer cancel()

	cmd, err := cliUtils.StartPythonServer(cfg.Port, cfg.PythonExecutable, cfg.PythonScript)
	if err != nil {
		if err.Error() == cliUtils.ErrorNoPython {
			info := fmt.Sprintf("python не установлен или его нет в PATH, обратитесь к администратору")
			return cliUtils.OnceProcessResult{}, cliUtils.InternalError(info)
		}

		info := fmt.Sprintf("ошибка при старте сервера: %v", err)
		return cliUtils.OnceProcessResult{}, cliUtils.ServerError(info)
	}

	defer cliUtils.KillServer(cmd)

	data, docType, err := cliUtils.SendFileToServer(filePath, cfg.Port)
	if err != nil {
		info := fmt.Sprintf("ошибка при отправке файла: %v", err)
		return cliUtils.OnceProcessResult{}, cliUtils.ServerError(info)
	}

	err = files.SaveFileToStorage(createdNameFile, data)
	if err != nil {
		info := fmt.Sprintf("ошибка при попытке сохранить файл: %v", err)
		return cliUtils.OnceProcessResult{}, cliUtils.ServerError(info)
	}

	onceProcessResult := cliUtils.CreateOnceProcessResult(createdNameFile, docType, cfg.StoragePath)
	return onceProcessResult, nil
}
