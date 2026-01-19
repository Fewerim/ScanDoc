package appCmds

import (
	"errors"
	"fmt"
	"proWeb/internal/appUtils"
	"proWeb/internal/appUtils/appWorks"
	"proWeb/internal/exitCodes"
	"proWeb/internal/tesseract"
	"time"
)

// MultiFiles - проверяет входные данные, создает подключение к серверу,
// обрабатывает и сохраняет файлы локально
func (a *App) MultiFiles(operation, directory, createdFolderName string) (err error) {
	appUtils.InfoMessage("Команда run_multi начала свое выполнение, валидация входных данных и проверка предварительных условий")

	if err := tesseract.CheckTesseract(); err != nil {
		a.log.Error(operation, "tesseract не добавлен в PATH", exitCodes.UserError)
		return appUtils.UserError("tesseract не добавлен в PATH")
	}

	initUsed, err := appUtils.CheckInitWasUsed(a.cfg.PythonVenvPath)
	if err != nil {
		a.log.Error(operation, "ошибка чтения папки с зависимостями", exitCodes.InternalError)
		return appUtils.InternalError("ошибка чтения папки с зависимостями")
	}
	if !initUsed {
		a.log.Error(operation, "приложение не было инициализировано. Необходимо запустить команду init", exitCodes.UserError)
		return appUtils.UserError("приложение не было инициализировано. Необходимо запустить команду init")
	}

	if err := appUtils.CheckIsAutorunCorrect(); err != nil {
		a.log.Error(operation, err.Error(), appUtils.GetExitCode(err, exitCodes.UserError))
		return err
	}

	start := time.Now()

	a.log.Info(operation, "Команда начала свое выполнение")

	defer func() {
		if r := recover(); r != nil {
			err = appUtils.InternalError(fmt.Sprintf("внутренняя ошибка: %v", r))
			a.log.Error(operation, "операция завершена с паникой", exitCodes.InternalError)
		}
	}()

	if err := a.CheckStorageJSON(); err != nil {
		a.log.Error(operation, err.Error(), appUtils.GetExitCode(err, exitCodes.InternalError))
		return appUtils.InternalError(err.Error())
	}

	if err = appUtils.CheckExistsPath(directory); err != nil {
		a.log.Error(operation, err.Error(), appUtils.GetExitCode(err, exitCodes.UserError))
		return err
	}

	a.log.Info(operation, "начало обработки директории файлов")
	appUtils.InfoMessage("Начало обработки директории файлов")

	result, err, errs := appWorks.MultiProcessFiles(directory, createdFolderName, a.cfg, a.storage)
	if err != nil {
		a.log.Error(operation, err.Error(), appUtils.GetExitCode(err, exitCodes.ServerError))
		return err
	}

	if errs != nil {
		for _, fileErr := range errs {
			a.log.Error(operation, fileErr.Err.Error(), appUtils.GetExitCode(err, exitCodes.ServerError), fileErr.FileName)
		}
	}
	elapsed := time.Since(start)
	result.SetElapsedTime(elapsed)

	if len(result.Results) == 0 {
		err = errors.New("в директории нет подходящих файлов для обработки")
		a.log.Error(operation, err.Error(), appUtils.GetExitCode(err, exitCodes.InternalError))
		return err
	}

	appUtils.NewSuccess(&result).PrintSuccess()
	if errs != nil {
		appUtils.FilesNotProcessed(errs)
	}

	a.log.Info(operation, fmt.Sprintf("операция завершена, время выполнения: %.3fs", result.GetElapsedTime()))
	return nil
}
