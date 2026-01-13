package appWorks

import (
	"context"
	"fmt"
	appUtils2 "proWeb/internal/appUtils"
	"proWeb/internal/files"
	"proWeb/lib/config"
	"time"
)

const processTimeout = 180 * time.Second

// ProcessOnceFile - подключение к серверу, отправка файла, обработка результата, сохранение в локальное хранилище
func ProcessOnceFile(filePath, createdNameFile string, cfg *config.Config) (appUtils2.OnceProcessResult, error) {
	_, cancel := context.WithTimeout(context.Background(), processTimeout)
	defer cancel()

	cmd, err := appUtils2.StartPythonServer(cfg.Port, cfg.PythonExecutable, cfg.PythonScript)
	if err != nil {
		if err.Error() == appUtils2.ErrorNoPython {
			info := fmt.Sprintf("python не установлен или его нет в PATH, обратитесь к администратору")
			return appUtils2.OnceProcessResult{}, appUtils2.InternalError(info)
		}

		info := fmt.Sprintf("ошибка при старте сервера: %v", err)
		return appUtils2.OnceProcessResult{}, appUtils2.ServerError(info)
	}

	defer appUtils2.KillServer(cmd)

	data, docType, err := appUtils2.SendFileToServer(filePath, cfg.Port)
	if err != nil {
		info := fmt.Sprintf("ошибка при отправке файла: %v", err)
		return appUtils2.OnceProcessResult{}, appUtils2.ServerError(info)
	}

	err = files.SaveFileToStorage(createdNameFile, data)
	if err != nil {
		info := fmt.Sprintf("ошибка при попытке сохранить файл: %v", err)
		return appUtils2.OnceProcessResult{}, appUtils2.ServerError(info)
	}

	onceProcessResult := appUtils2.CreateOnceProcessResult(createdNameFile, docType, cfg.StoragePath)
	return onceProcessResult, nil
}
