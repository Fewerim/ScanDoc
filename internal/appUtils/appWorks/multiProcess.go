package appWorks

import (
	"fmt"
	"path/filepath"
	"proWeb/internal/appUtils"
	"proWeb/internal/files"
	"proWeb/internal/storage"
	"proWeb/lib/config"
	"strings"
	"sync"
)

const maxParallelOperations = 5

// MultiProcessFiles - подключение к серверу, обработка файлов в директории, сохранение результатов в локальное хранилище.
// Обрабатывает сразу все файлы из директории, возвращает результаты выполнения CLI команды, ошибку при подключении сервера или проверки директории,
// слайс ошибок, возникших при обработке конкретных файлов
func MultiProcessFiles(directoryPath, folderName string, cfg *config.Config, storage *storage.Storage) (appUtils.MultiProcessResult, error, []appUtils.FileError) {
	filePaths, errorInfo := appUtils.GetFilesFromDirectory(directoryPath)
	if errorInfo != "" {
		return appUtils.MultiProcessResult{}, appUtils.UserError(errorInfo), nil
	}

	cmd, err := appUtils.StartPythonServer(cfg.Port, cfg.PythonExecutable, cfg.PythonScript)
	if err != nil {
		if err.Error() == appUtils.ErrorNoPython {
			info := fmt.Sprintf("python не установлен или его нет в PATH, обратитесь к администратору")
			return appUtils.MultiProcessResult{}, appUtils.InternalError(info), nil
		}

		info := fmt.Sprintf("ошибка при старте сервера: %v", err)
		return appUtils.MultiProcessResult{}, appUtils.ServerError(info), nil
	}

	defer appUtils.KillServer(cmd)

	if len(filePaths) > 5 {
		appUtils.InfoMessage("Система одновременно может обрабатывать не более 5 файлов, если файлов больше, то на это требуется больше времени")
	}

	maxWorkers := maxParallelOperations
	semaphore := make(chan struct{}, maxWorkers)

	results := make(chan appUtils.Result, len(filePaths))
	errorsFileProcessing := make(chan appUtils.FileError, len(filePaths))

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

			err := appUtils.ValidateExtensionFile(filePath)
			if err != nil {
				info := fmt.Sprintf("Расширение файла не поддерживается: %v", err)
				errorsFileProcessing <- appUtils.FileError{fileName, appUtils.UserError(info)}
				return
			}

			if ext := filepath.Ext(filePath); ext == ".pdf" {
				filePath, err = appUtils.GetJpgFromPdf(filePath)
				if err != nil {
					info := err.Error()
					errorsFileProcessing <- appUtils.FileError{fileName, appUtils.UserError(info)}
					return
				}
			}

			data, docType, err := appUtils.SendFileToServer(filePath, cfg.Port)
			if err != nil {
				info := fmt.Sprintf("ошибка при отправке файла: %v", err)
				errorsFileProcessing <- appUtils.FileError{fileName, appUtils.ServerError(info)}
				return
			}

			//errNew := files.SaveFileToDirectory(fileNameWithoutExt, folderName, data, files.NotOverwrite)
			errNew := storage.SaveFile(folderName, fileNameWithoutExt, data, files.NotOverwrite)
			if errNew != nil {
				info := fmt.Sprintf("ошибка при попытке сохранить файл: %v", errNew)
				errorsFileProcessing <- appUtils.FileError{fileName, appUtils.ServerError(info)}
				return
			}

			result := appUtils.CreateResultWithFolder(fileNameWithoutExt, docType, cfg.StoragePath, folderName)
			results <- result
		}(filePath)
	}

	wg.Wait()
	close(results)
	close(errorsFileProcessing)

	multiProcessResult := appUtils.CreateMultiProcessResult()
	var allErrorsFileProcessing []appUtils.FileError

	for result := range results {
		multiProcessResult.SetResult(result)
	}

	for errs := range errorsFileProcessing {
		allErrorsFileProcessing = append(allErrorsFileProcessing, errs)
	}

	return multiProcessResult, nil, allErrorsFileProcessing
}
