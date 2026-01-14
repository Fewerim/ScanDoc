package appWorks

import (
	"fmt"
	"path/filepath"
	appUtils2 "proWeb/internal/appUtils"
	"proWeb/internal/files"
	"proWeb/lib/config"
	"strings"
	"sync"

	"github.com/fatih/color"
)

const maxParallelOperations = 5

// MultiProcessFiles - подключение к серверу, обработка файлов в директории, сохранение результатов в локальное хранилище.
// Обрабатывает сразу все файлы из директории, возвращает результаты выполнения CLI команды, ошибку при подключении сервера или проверки директории,
// слайс ошибок, возникших при обработке конкретных файлов
func MultiProcessFiles(directoryPath string, cfg *config.Config, folderName string) (appUtils2.MultiProcessResult, error, []appUtils2.FileError) {
	filePaths, errorInfo := appUtils2.GetFilesFromDirectory(directoryPath)
	if errorInfo != "" {
		return appUtils2.MultiProcessResult{}, appUtils2.UserError(errorInfo), nil
	}

	cmd, err := appUtils2.StartPythonServer(cfg.Port, cfg.PythonExecutable, cfg.PythonScript)
	if err != nil {
		if err.Error() == appUtils2.ErrorNoPython {
			info := fmt.Sprintf("python не установлен или его нет в PATH, обратитесь к администратору")
			return appUtils2.MultiProcessResult{}, appUtils2.InternalError(info), nil
		}

		info := fmt.Sprintf("ошибка при старте сервера: %v", err)
		return appUtils2.MultiProcessResult{}, appUtils2.ServerError(info), nil
	}

	defer appUtils2.KillServer(cmd)

	if len(filePaths) > 5 {
		color.Blue("Система одновременно может обрабатывать не более 5 файлов, если файлов больше, то на это требуется больше времени")
	}

	maxWorkers := maxParallelOperations
	semaphore := make(chan struct{}, maxWorkers)

	results := make(chan appUtils2.Result, len(filePaths))
	errorsFileProcessing := make(chan appUtils2.FileError, len(filePaths))

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

			err := appUtils2.ValidateExtensionFile(filePath)
			if err != nil {
				info := fmt.Sprintf("Расширение файла не поддерживается: %v", err)
				errorsFileProcessing <- appUtils2.FileError{fileName, appUtils2.UserError(info)}
				return
			}

			if ext := filepath.Ext(filePath); ext == ".pdf" {
				filePath, err = appUtils2.GetJpgFromPdf(filePath)
				if err != nil {
					info := err.Error()
					errorsFileProcessing <- appUtils2.FileError{fileName, appUtils2.UserError(info)}
					return
				}
			}

			data, docType, err := appUtils2.SendFileToServer(filePath, cfg.Port)
			if err != nil {
				info := fmt.Sprintf("ошибка при отправке файла: %v", err)
				errorsFileProcessing <- appUtils2.FileError{fileName, appUtils2.ServerError(info)}
				return
			}

			errNew := files.SaveFileToDirectory(fileNameWithoutExt, folderName, data)
			if errNew != nil {
				info := fmt.Sprintf("ошибка при попытке сохранить файл: %v", errNew)
				errorsFileProcessing <- appUtils2.FileError{fileName, appUtils2.ServerError(info)}
				return
			}

			result := appUtils2.CreateResultWithFolder(fileNameWithoutExt, docType, cfg.StoragePath, folderName)
			results <- result
		}(filePath)
	}

	wg.Wait()
	close(results)
	close(errorsFileProcessing)

	multiProcessResult := appUtils2.CreateMultiProcessResult()
	var allErrorsFileProcessing []appUtils2.FileError

	for result := range results {
		multiProcessResult.SetResult(result)
	}

	for errs := range errorsFileProcessing {
		allErrorsFileProcessing = append(allErrorsFileProcessing, errs)
	}

	return multiProcessResult, nil, allErrorsFileProcessing
}
