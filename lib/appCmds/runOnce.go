package appCmds

import (
	"fmt"
	"path/filepath"
	"proWeb/internal/appUtils"
	"proWeb/internal/appUtils/appWorks"
	"proWeb/internal/exitCodes"
	"proWeb/internal/tesseract"
	"time"
)

// OnceFile - проверяет входные данные, создает подключение к серверу,
// обрабатывает и сохраняет файл локально
func (a *App) OnceFile(operation, filePath, createdFileName string) (err error) {
	appUtils.InfoMessage("Команда run_once начала свое выполнение, валидация входных данных и проверка предварительных условий")

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

	a.log.Info(operation, "Команда начала свое выполнение")

	start := time.Now()

	defer func() {
		if r := recover(); r != nil {
			err = appUtils.InternalError(fmt.Sprintf("внутренняя ошибка: %v", r))
			a.log.Error(operation, "операция завершена с паникой", exitCodes.InternalError)
		}
	}()

	if filePath == "" {
		a.log.Error(operation, "не указан путь к файлу", exitCodes.UserError)
		return appUtils.UserError("укажите путь к файлу")
	}

	if createdFileName == "" {
		a.log.Error(operation, "не указано имя нового файла", exitCodes.UserError)
		return appUtils.UserError("не указано имя нового файла")
	}

	if err := a.CheckStorageJSON(); err != nil {
		a.log.Error(operation, err.Error(), appUtils.GetExitCode(err, exitCodes.InternalError))
		return appUtils.InternalError(err.Error())
	}

	if err = appUtils.ValidateExtensionFile(filePath); err != nil {
		a.log.Error(operation, err.Error(), appUtils.GetExitCode(err, exitCodes.UserError), filepath.Base(filePath))
		return appUtils.UserError(err.Error())
	}

	if err = appUtils.CheckExistsFile(filePath); err != nil {
		a.log.Error(operation, err.Error(), appUtils.GetExitCode(err, exitCodes.UserError))
		return appUtils.UserError(err.Error())
	}

	if ext := filepath.Ext(filePath); ext == ".pdf" {
		filePath, err = appUtils.GetJpgFromPdf(filePath)
		if err != nil {
			a.log.Error(operation, err.Error(), appUtils.GetExitCode(err, exitCodes.UserError))
			return appUtils.UserError(err.Error())
		}
	}

	a.log.Info(operation, "начало обработки файла")
	appUtils.InfoMessage("Начало обработки файла")

	result, err := appWorks.ProcessOnceFile(filePath, createdFileName, a.cfg, a.storage)
	if err != nil {
		a.log.Error(operation, err.Error(), appUtils.GetExitCode(err, exitCodes.ServerError), filepath.Base(filePath))
		return err
	}
	elapsed := time.Since(start)
	result.SetElapsedTime(elapsed)

	appUtils.NewSuccess(&result).PrintSuccess()
	a.log.Info(operation, fmt.Sprintf("операция завершена, время выполнения: %.3fs", result.GetElapsedTime()))
	return nil
}
