package appCmds

import (
	"proWeb/lib/appUtils"
	"proWeb/lib/exitCodes"
	"proWeb/lib/tesseract"

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
		return appUtils.InternalError(err.Error())
	}
	a.Log.Info(operation, "локальное хранилище успешно установлено")

	a.Log.Info(operation, "проверка наличия tesseract в PATH")
	if err := tesseract.CheckTesseract(); err != nil {
		a.Log.Error(operation, "tesseract не добавлен в PATH", exitCodes.UserError)
		return appUtils.UserError("tesseract не добавлен в PATH")
	}

	a.Log.Info(operation, "проверка наличия кодировки UTF-8")
	if err := appUtils.CheckIsAutorunCorrect(); err != nil {
		a.Log.Error(operation, err.Error(), appUtils.GetExitCode(err, exitCodes.UserError))
		return err
	}

	a.Log.Info(operation, "начало установки зависимостей")
	color.Blue("Начало установки зависимостей")
	if err := appUtils.InstallRequirements(a.Cfg.PythonScript); err != nil {
		a.Log.Error(operation, err.Error(), appUtils.GetExitCode(err, exitCodes.InternalError))
		return err
	}
	result := appUtils.CreateInitResult("зависимости успешно установлены")

	appUtils.NewSuccess(&result).PrintSuccess()
	a.Log.Info(operation, "зависимости успешно установлены")
	return nil
}
