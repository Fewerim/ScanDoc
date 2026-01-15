package appCmds

import (
	"proWeb/internal/appUtils"
	"proWeb/internal/exitCodes"
	"proWeb/internal/tesseract"
)

// InitApp - подтягивает все необходимые зависимости для корректной работы python скрипта
// также устанавливает хранилище для обработанных файлов
// operation - имя операции (GUI | CLI)
func (a *App) InitApp(operation string) error {
	appUtils.InfoMessage("Команда init начала свое выполнение, проверка зависимостей и инициализация окружения")
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
	appUtils.InfoMessage("Начало установки зависимостей")
	if err := appUtils.InstallRequirements(a.Cfg.PythonVenvPath, a.Cfg.PythonScript); err != nil {
		a.Log.Error(operation, err.Error(), appUtils.GetExitCode(err, exitCodes.InternalError))
		return err
	}
	result := appUtils.CreateInitResult("зависимости успешно установлены")

	appUtils.NewSuccess(&result).PrintSuccess()
	a.Log.Info(operation, "зависимости успешно установлены")
	return nil
}

// CheckInit - проверяет, был ли init до вызова функции
func (a *App) CheckInit(operation string) (bool, error) {
	a.Log.Info(operation, "Проверка инициализации приложения")

	ok, err := appUtils.CheckInitWasUsed(a.Cfg.PythonVenvPath)
	if err != nil {
		a.Log.Error(operation, err.Error(), exitCodes.InternalError)
		return false, err
	}
	return ok, nil
}
