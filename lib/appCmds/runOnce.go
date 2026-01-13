package appCmds

import (
	"fmt"
	"path/filepath"
	cliUtils "proWeb/lib/appUtils"
	"proWeb/lib/appUtils/appWorks"
	"proWeb/lib/exitCodes"
	"proWeb/lib/tesseract"
	"time"

	"github.com/fatih/color"
)

// OnceFile - проверяет входные данные, создает подключение к серверу,
// обрабатывает и сохраняет файл локально
func (a *App) OnceFile(operation, filePath, createdFileName string) (err error) {
	color.Blue("Команда run_once начала свое выполнение, валидация входных данных и проверка предварительных условий")

	if err := tesseract.CheckTesseract(); err != nil {
		a.Log.Error(operation, "tesseract не добавлен в PATH", exitCodes.UserError)
		return cliUtils.UserError("tesseract не добавлен в PATH")
	}

	initUsed, err := cliUtils.CheckInitWasUsed()
	if err != nil {
		a.Log.Error(operation, "ошибка чтения папки с зависимостями", exitCodes.InternalError)
		return cliUtils.InternalError("ошибка чтения папки с зависимостями")
	}
	if !initUsed {
		a.Log.Error(operation, "приложение не было инициализировано. Необходимо запустить команду init", exitCodes.UserError)
		return cliUtils.UserError("приложение не было инициализировано. Необходимо запустить команду init")
	}

	if err := cliUtils.CheckIsAutorunCorrect(); err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.UserError))
		return err
	}

	a.Log.Info(operation, "Команда начала свое выполнение")

	start := time.Now()

	defer func() {
		if r := recover(); r != nil {
			err = cliUtils.InternalError(fmt.Sprintf("внутренняя ошибка: %v", r))
			a.Log.Error(operation, "операция завершена с паникой", exitCodes.InternalError)
		}
	}()

	if filePath == "" {
		a.Log.Error(operation, "не указан путь к файлу", exitCodes.UserError)
		return cliUtils.UserError("укажите путь к файлу")
	}

	if createdFileName == "" {
		a.Log.Error(operation, "не указано имя нового файла", exitCodes.UserError)
		return cliUtils.UserError("не указано имя нового файла")
	}

	if err := a.CheckStorageJSON(); err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.InternalError))
		return cliUtils.InternalError(err.Error())
	}

	if err = cliUtils.ValidateExtensionFile(filePath); err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.UserError), filepath.Base(filePath))
		return err
	}

	if err = cliUtils.CheckExistsFile(filePath); err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.UserError))
		return err
	}

	a.Log.Info(operation, "начало обработки файла")
	color.Blue("Начало обработки файла")

	result, err := appWorks.ProcessOnceFile(filePath, createdFileName, a.Cfg)
	if err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.ServerError), filepath.Base(filePath))
		return err
	}
	elapsed := time.Since(start)
	result.SetElapsedTime(elapsed)

	cliUtils.NewSuccess(&result).PrintSuccess()
	a.Log.Info(operation, fmt.Sprintf("операция завершена, время выполнения: %.3fs", result.GetElapsedTime()))
	return nil
}
