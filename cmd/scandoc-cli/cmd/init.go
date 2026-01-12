package cmd

import (
	"proWeb/lib/cliUtils"
	"proWeb/lib/exitCodes"
	"proWeb/lib/tesseract"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// initApp - подтягивает все необходимые зависимости для корректной работы python скрипта
// также устанавливает хранилище для обработанных файлов
func (a *App) initApp(cmd *cobra.Command, args []string) error {
	const operation = "scandoc-cli.initApp"

	color.Blue("Команда init начала свое выполнение, проверка зависимостей и инициализация окружения")
	a.Log.Info(operation, "проверка и создание локального хранилища")
	if err := a.CheckStorageJSON(); err != nil {
		a.Log.Error(operation, err.Error(), exitCodes.InternalError)
		return cliUtils.InternalError(err.Error())
	}
	a.Log.Info(operation, "локальное хранилище успешно установлено")

	a.Log.Info(operation, "проверка наличия tesseract в PATH")
	if err := tesseract.CheckTesseract(); err != nil {
		a.Log.Error(operation, "tesseract не добавлен в PATH", exitCodes.UserError)
		return cliUtils.UserError("tesseract не добавлен в PATH")
	}

	a.Log.Info(operation, "проверка наличия кодировки UTF-8")
	if err := cliUtils.CheckIsAutorunCorrect(); err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.UserError))
		return err
	}

	a.Log.Info(operation, "начало установки зависимостей")
	color.Blue("Начало установки зависимостей")
	if err := cliUtils.InstallRequirements(a.Cfg.PythonScript); err != nil {
		a.Log.Error(operation, err.Error(), cliUtils.GetExitCode(err, exitCodes.InternalError))
		return err
	}
	result := cliUtils.CreateInitResult("зависимости успешно установлены")

	cliUtils.NewSuccess(&result).PrintSuccess()
	a.Log.Info(operation, "зависимости успешно установлены")
	return nil
}

func newInitAppCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:           "init",
		Short:         "Устанавливает необходимые зависимости для корректной работы приложения",
		Example:       "scandoc.exe init\nустановит необходимые зависимости и локальное хранилище",
		Args:          cobra.NoArgs,
		RunE:          a.initApp,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
}
