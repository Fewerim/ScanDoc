package cliWorks

import (
	"fmt"
	"proWeb/internal/cliUtils"
	"proWeb/internal/config"
	"proWeb/internal/files"
	"time"
)

// Result - результат выполнения CLI команды (имя созданного файла, время создания)
type Result struct {
	FileName  string
	CreatedAt time.Time
}

// ProcessOnceFile - подключение к серверу, отправка файла, обработка результата, сохранение в локальное хранилище
func ProcessOnceFile(filePath, createdNameFile string, cfg *config.Config) (Result, error) {
	cmd, err := cliUtils.StartPythonServer(cfg.Port, cfg.PythonExecutable)
	if err != nil {
		if err.Error() == cliUtils.ErrorNoPython {
			info := fmt.Sprintf("python не установлен или его нет в PATH, обратитесь к администратору")
			return Result{}, cliUtils.InternalError(info)
		}

		info := fmt.Sprintf("ошибка при старте сервера: %v", err)
		return Result{}, cliUtils.ServerError(info)
	}

	defer cliUtils.KillServer(cmd)

	data, err := cliUtils.SendFileToServer(filePath, cfg.Port)
	if err != nil {
		info := fmt.Sprintf("ошибка при отправке файла: %v", err)
		return Result{}, cliUtils.ServerError(info)
	}

	err = files.SaveFileToStorage(createdNameFile, data)
	if err != nil {
		info := fmt.Sprintf("ошибка при попытке сохранить файл: %v", err)
		return Result{}, cliUtils.ServerError(info)
	}

	result := createResult(createdNameFile)
	return result, nil
}

// createResult - конструктор для создания результата выполнения CLI команды
func createResult(fileName string) Result {
	return Result{
		FileName:  fileName,
		CreatedAt: time.Now(),
	}
}
