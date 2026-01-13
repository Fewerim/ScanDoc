package appCmds

import (
	"errors"
	"fmt"
	appUtils2 "proWeb/internal/appUtils"
	"proWeb/internal/appUtils/appWorks"
	"proWeb/internal/exitCodes"
	"proWeb/internal/tesseract"
	"time"

	"github.com/fatih/color"
)

// MultiFiles - проверяет входные данные, создает подключение к серверу,
// обрабатывает и сохраняет файлы локально
func (a *App) MultiFiles(operation, directory, createdFolderName string) (err error) {
	color.Blue("Команда run_multi начала свое выполнение, валидация входных данных и проверка предварительных условий")

	if err := tesseract.CheckTesseract(); err != nil {
		a.Log.Error(operation, "tesseract не добавлен в PATH", exitCodes.UserError)
		return appUtils2.UserError("tesseract не добавлен в PATH")
	}

	initUsed, err := appUtils2.CheckInitWasUsed()
	if err != nil {
		a.Log.Error(operation, "ошибка чтения папки с зависимостями", exitCodes.InternalError)
		return appUtils2.InternalError("ошибка чтения папки с зависимостями")
	}
	if !initUsed {
		a.Log.Error(operation, "приложение не было инициализировано. Необходимо запустить команду init", exitCodes.UserError)
		return appUtils2.UserError("приложение не было инициализировано. Необходимо запустить команду init")
	}

	if err := appUtils2.CheckIsAutorunCorrect(); err != nil {
		a.Log.Error(operation, err.Error(), appUtils2.GetExitCode(err, exitCodes.UserError))
		return err
	}

	start := time.Now()

	a.Log.Info(operation, "Команда начала свое выполнение")

	defer func() {
		if r := recover(); r != nil {
			err = appUtils2.InternalError(fmt.Sprintf("внутренняя ошибка: %v", r))
			a.Log.Error(operation, "операция завершена с паникой", exitCodes.InternalError)
		}
	}()

	if err := a.CheckStorageJSON(); err != nil {
		a.Log.Error(operation, err.Error(), appUtils2.GetExitCode(err, exitCodes.InternalError))
		return appUtils2.InternalError(err.Error())
	}

	if err = appUtils2.CheckExistsPath(directory); err != nil {
		a.Log.Error(operation, err.Error(), appUtils2.GetExitCode(err, exitCodes.UserError))
		return err
	}

	a.Log.Info(operation, "начало обработки директории файлов")
	color.Blue("Начало обработки директории файлов")

	result, err, errs := appWorks.MultiProcessFiles(directory, a.Cfg, createdFolderName)
	if err != nil {
		a.Log.Error(operation, err.Error(), appUtils2.GetExitCode(err, exitCodes.ServerError))
		return err
	}

	if errs != nil {
		for _, fileErr := range errs {
			a.Log.Error(operation, fileErr.Err.Error(), appUtils2.GetExitCode(err, exitCodes.ServerError), fileErr.FileName)
		}
	}
	elapsed := time.Since(start)
	result.SetElapsedTime(elapsed)

	if len(result.Results) == 0 {
		err = errors.New("в директории нет подходящих файлов для обработки")
		a.Log.Error(operation, err.Error(), appUtils2.GetExitCode(err, exitCodes.InternalError))
		return err
	}

	appUtils2.NewSuccess(&result).PrintSuccess()
	if errs != nil {
		appUtils2.FilesNotProcessed(errs)
	}

	a.Log.Info(operation, fmt.Sprintf("операция завершена, время выполнения: %.3fs", result.GetElapsedTime()))
	return nil
}
