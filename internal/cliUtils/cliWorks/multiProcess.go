package cliWorks

import (
	"fmt"
	"path/filepath"
	"proWeb/internal/cliUtils"
	"proWeb/internal/config"
	"proWeb/internal/files"
	"strings"
	"sync"
)

const maxParallelOperations = 5

type MultiProcessResult struct {
	Results []Result
}

// MultiProcessFiles - обрабатывает сразу все файлы из директории, возвращает результаты выполнения CLI командб ошибку при подключении сервера или проверки дтректории, слайс ошибок, возникших при обработке конкретных файлов
func MultiProcessFiles(directoryPath string, cfg *config.Config) (MultiProcessResult, error, []error) {
	filePaths, errorInfo := cliUtils.GetFilesFromDirectory(directoryPath)
	if errorInfo != "" {
		return MultiProcessResult{}, cliUtils.UserError(errorInfo), nil
	}

	cmd, err := cliUtils.StartPythonServer(cfg.Port, cfg.PythonExecutable, cfg.PythonScript)
	if err != nil {
		if err.Error() == cliUtils.ErrorNoPython {
			info := fmt.Sprintf("python не установлен или его нет в PATH, обратитесь к администратору")
			return MultiProcessResult{}, cliUtils.InternalError(info), nil
		}

		info := fmt.Sprintf("ошибка при старте сервера: %v", err)
		return MultiProcessResult{}, cliUtils.ServerError(info), nil
	}

	defer cliUtils.KillServer(cmd)

	maxWorkers := maxParallelOperations
	semaphore := make(chan struct{}, maxWorkers)

	results := make(chan Result, len(filePaths))
	errorsFileProcessing := make(chan string, len(filePaths))

	var wg sync.WaitGroup

	for _, filePath := range filePaths {
		wg.Add(1)

		semaphore <- struct{}{}

		go func(filePath string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			fileName := filepath.Base(filePath)
			extension := filepath.Ext(fileName)
			fileNameWithoutExt := strings.TrimSuffix(fileName, extension)

			err := cliUtils.ValidateExtensionFile(filePath)
			if err != nil {
				info := fmt.Sprintf("Расширение файла %s не поддерживается: %v", fileNameWithoutExt, err)
				errorsFileProcessing <- info
				return
			}

			data, err := cliUtils.SendFileToServer(filePath, cfg.Port)
			if err != nil {
				info := fmt.Sprintf("ошибка при отправке файла %s: %v", fileNameWithoutExt, err)
				errorsFileProcessing <- info
				return
			}

			errNew := files.SaveFileToStorage(fileNameWithoutExt, data)
			if errNew != nil {
				info := fmt.Sprintf("ошибка при попытке сохранить файл %s: %v", fileNameWithoutExt, errNew)
				errorsFileProcessing <- info
				return
			}

			result := createResult(fileNameWithoutExt)
			results <- result
		}(filePath)
	}

	wg.Wait()
	close(results)
	close(errorsFileProcessing)

	var multiProcessResult MultiProcessResult
	var allErrorsFileProcessing []error

	for result := range results {
		multiProcessResult.Results = append(multiProcessResult.Results, result)
	}

	for errs := range errorsFileProcessing {
		allErrorsFileProcessing = append(allErrorsFileProcessing, cliUtils.ServerError(errs))
	}

	return multiProcessResult, nil, allErrorsFileProcessing
}
