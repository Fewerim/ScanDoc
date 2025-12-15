package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "proweb-cli",
		Short: "CLI для распознавания бухгалтерских документов",
	}

	configPath string
)

const defaultConfigPath = "./config/config.yaml"

// initCommands - инициализирует CLI команды, предварительно обработав флаг для получения пути к конфигу
// если флаг не был введен, используется дефолтный путь.
func (a *App) initCommands() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file path")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		path, err := resolveConfigPath()
		if err != nil {
			return err
		}

		a.LoadConfig(path)

		a.SetupLogger(a.Cfg.LogPath)

		return nil
	}

	rootCmd.AddCommand(newHelperCmd(a))
	rootCmd.AddCommand(newRunOnceCmd(a))
}

// Execute - делегирует запуск CLI приложения, вызываясь на экземпляре App
func (a *App) Execute() error {
	a.initCommands()
	return rootCmd.Execute()
}

// resolveConfigPath - получает путь к конфигу из переменной окружения, флага или устанавливает дефолтное значение пути
func resolveConfigPath() (string, error) {
	path := configPath
	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	if path == "" {
		path = defaultConfigPath
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("конфиг файл по пути %s не найден", path)
	}

	return path, nil
}
