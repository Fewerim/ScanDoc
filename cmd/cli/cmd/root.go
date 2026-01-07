package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "scanner.exe",
		Short: "ScanDoc - CLI для распознавания бухгалтерских документов",
	}

	configFlag string
)

// initCommands - инициализирует CLI команды, предварительно обработав флаг для получения пути к конфигу
// если флаг не был введен, используется дефолтный путь.
func (a *App) initCommands() {
	rootCmd.PersistentFlags().StringVar(&configFlag, "config", "", "путь к конфигурации приложения")
	rootCmd.LocalFlags().BoolP("help", "h", false, "показать справку по команде")

	rootCmd.DisableAutoGenTag = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Lookup("help") != nil && cmd.Flags().Lookup("help").Changed ||
			cmd.CommandPath() == "proweb-cli" && len(args) == 0 {
			return nil
		}

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
	rootCmd.SetHelpCommand(newHelperCmd(a))
	rootCmd.AddCommand(newInstallTesseract(a))
	rootCmd.AddCommand(newInitAppCmd(a))
	rootCmd.AddCommand(newRunOnceCmd(a))
	rootCmd.AddCommand(newMultiRunCmd(a))
}

// Execute - делегирует запуск CLI приложения, вызываясь на экземпляре App
func (a *App) Execute() error {
	a.initCommands()
	return rootCmd.Execute()
}
