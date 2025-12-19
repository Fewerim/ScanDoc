package cmd

import (
	"fmt"
	"proWeb/internal/cliUtils"
	"proWeb/internal/files"

	"github.com/spf13/cobra"
)

// initApp - подтягивает все необходимые зависимости для корректной работы python скрипта
// также устанавливает хранилище для обработанных файлов
func (a *App) initApp(cmd *cobra.Command, args []string) error {
	const operation = "cli.initApp"

	a.Log.Info(operation, "создание локального хранилища для обработанных файлов")
	files.InitStorage(a.Cfg.StoragePath)

	if !files.StorageExists() {
		//a.Log.Info(operation, "локальное хранилище не существует, создание нового")
		if err := files.CreateStorageJSON(); err != nil {
			a.Log.Error(operation, fmt.Sprintf("ошибка создания локального хранилища: %v", err), 3)
			return cliUtils.InternalError("ошибка создания локального хранилища")
		}
		a.Log.Info(operation, "локальное хранилище успешно создано")
	}

	a.Log.Info(operation, "начало установки зависимостей")
	if err := cliUtils.InstallRequirements(a.Cfg.PythonScript); err != nil {
		a.Log.Error(operation, "зависимости не были установлены", 3)
		return err
	}

	a.Log.Info(operation, "зависимости установлены")
	return nil
}

func newInitAppCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Устанавливает необходимые зависимости для корректной работы приложения",
		Args:  cobra.NoArgs,
		RunE:  a.initApp,
	}
}
