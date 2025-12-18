package cmd

import (
	"proWeb/internal/cliUtils"

	"github.com/spf13/cobra"
)

// initApp - подтягивает все необходимые зависимости для корректной работы python скрипта
func (a *App) initApp(cmd *cobra.Command, args []string) error {
	const operation = "cli.initApp"

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
