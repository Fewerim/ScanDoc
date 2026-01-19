package appWorks

import (
	"context"
	"fmt"
	"proWeb/internal/appUtils"
	"proWeb/internal/files"
	"proWeb/internal/storage"
	"proWeb/lib/config"
	"time"
)

const processTimeout = 180 * time.Second

// ProcessOnceFile - подключение к серверу, отправка файла, обработка результата, сохранение в локальное хранилище
func ProcessOnceFile(filePath, createdNameFile string, cfg *config.Config, storage *storage.Storage) (appUtils.OnceProcessResult, error) {
	_, cancel := context.WithTimeout(context.Background(), processTimeout)
	defer cancel()

	cmd, err := appUtils.StartPythonServer(cfg.Port, cfg.PythonExecutable, cfg.PythonScript)
	if err != nil {
		if err.Error() == appUtils.ErrorNoPython {
			info := fmt.Sprintf("python не установлен или его нет в PATH, обратитесь к администратору")
			return appUtils.OnceProcessResult{}, appUtils.InternalError(info)
		}

		info := fmt.Sprintf("ошибка при старте сервера: %v", err)
		return appUtils.OnceProcessResult{}, appUtils.ServerError(info)
	}

	defer appUtils.KillServer(cmd)

	data, docType, err := appUtils.SendFileToServer(filePath, cfg.Port)
	if err != nil {
		info := fmt.Sprintf("ошибка при отправке файла: %v", err)
		return appUtils.OnceProcessResult{}, appUtils.ServerError(info)
	}

	//err = files.SaveFileToDirectory(createdNameFile, "", data, files.NotOverwrite)
	err = storage.SaveFile("", createdNameFile, data, files.NotOverwrite)
	if err != nil {
		info := fmt.Sprintf("ошибка при попытке сохранить файл: %v", err)
		return appUtils.OnceProcessResult{}, appUtils.ServerError(info)
	}

	onceProcessResult := appUtils.CreateOnceProcessResult(createdNameFile, docType, cfg.StoragePath)
	return onceProcessResult, nil
}
