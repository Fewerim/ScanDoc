package appCmds

import (
	appUtils2 "proWeb/internal/appUtils"
	"proWeb/internal/exitCodes"
	"proWeb/internal/tesseract"

	"github.com/fatih/color"
)

// InitApp - подтягивает все необходимые зависимости для корректной работы python скрипта
// также устанавливает хранилище для обработанных файлов
// operation - имя операции (GUI | CLI)
func (a *App) InitApp(operation string) error {
	color.Blue("Команда init начала свое выполнение, проверка зависимостей и инициализация окружения")
	a.Log.Info(operation, "проверка и создание локального хранилища")
	if err := a.CheckStorageJSON(); err != nil {
		a.Log.Error(operation, err.Error(), exitCodes.InternalError)
		return appUtils2.InternalError(err.Error())
	}
	a.Log.Info(operation, "локальное хранилище успешно установлено")

	a.Log.Info(operation, "проверка наличия tesseract в PATH")
	if err := tesseract.CheckTesseract(); err != nil {
		a.Log.Error(operation, "tesseract не добавлен в PATH", exitCodes.UserError)
		return appUtils2.UserError("tesseract не добавлен в PATH")
	}

	a.Log.Info(operation, "проверка наличия кодировки UTF-8")
	if err := appUtils2.CheckIsAutorunCorrect(); err != nil {
		a.Log.Error(operation, err.Error(), appUtils2.GetExitCode(err, exitCodes.UserError))
		return err
	}

	a.Log.Info(operation, "начало установки зависимостей")
	color.Blue("Начало установки зависимостей")
	if err := appUtils2.InstallRequirements(a.Cfg.PythonScript); err != nil {
		a.Log.Error(operation, err.Error(), appUtils2.GetExitCode(err, exitCodes.InternalError))
		return err
	}
	result := appUtils2.CreateInitResult("зависимости успешно установлены")

	appUtils2.NewSuccess(&result).PrintSuccess()
	a.Log.Info(operation, "зависимости успешно установлены")
	return nil
}
