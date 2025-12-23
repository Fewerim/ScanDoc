package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "proweb-cli",
		Short: "CLI для распознавания бухгалтерских документов",
	}

	configFlag string
)

// initCommands - инициализирует CLI команды, предварительно обработав флаг для получения пути к конфигу
// если флаг не был введен, используется дефолтный путь.
func (a *App) initCommands() {
	rootCmd.PersistentFlags().StringVar(&configFlag, "config", "", "config file path")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if configFlag == "" {
			a.CfgPath = configFlag
		}

		a.LoadConfig()
		a.SetupLogger(a.Cfg.LogPath)

		if err := a.InitPythonVenv(); err != nil {
			return err
		}

		if err := a.CheckPythonScripts(); err != nil {
			return err
		}

		return nil
	}

	rootCmd.AddCommand(newHelperCmd(a))
	rootCmd.AddCommand(newInitAppCmd(a))
	rootCmd.AddCommand(newRunOnceCmd(a))
	rootCmd.AddCommand(newMultyRunCmd(a))
}

// Execute - делегирует запуск CLI приложения, вызываясь на экземпляре App
func (a *App) Execute() error {
	a.initCommands()
	return rootCmd.Execute()
}
