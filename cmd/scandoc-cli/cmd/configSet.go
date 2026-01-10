package cmd

import (
	"fmt"
	"proWeb/internal/cliUtils"
	"proWeb/internal/config"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// configSet - команда, позволяющая менять значения конфига внутри файла
func configSet(configPath string, port int, pythonExecutable, pythonScript, storagePath, logPath string) error {
	if err := checkRequiredFlags(port, pythonExecutable, pythonScript); err != nil {
		return cliUtils.UserError(err.Error())
	}

	if err := config.CheckConfigPathExists(configPath); err != nil {
		return cliUtils.InternalError(err.Error())
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return cliUtils.InternalError(err.Error())
	}

	cfg.Port = port
	cfg.PythonExecutable = pythonExecutable
	cfg.PythonScript = pythonScript

	if storagePath != "" {
		cfg.StoragePath = storagePath
	}
	if cfg.StoragePath == "" {
		cfg.StoragePath = config.DefaultPathToStorage
	}

	if logPath != "" {
		cfg.LogPath = logPath
	}
	if cfg.LogPath == "" {
		cfg.LogPath = config.DefaultPathToLog
	}

	if err := cfg.SaveConfig(configPath); err != nil {
		return err
	}

	color.Blue("Конфигурация успешно обновлена")
	return nil
}

// setupDefaultConfig - устанавливает базовые значения для конфига
func setupDefaultConfig(configPath string) error {
	if err := config.CheckConfigPathExists(configPath); err != nil {
		return cliUtils.InternalError(err.Error())
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return cliUtils.InternalError(err.Error())
	}

	cfg.SetupDefaultConfig()

	if err := cfg.SaveConfig(configPath); err != nil {
		return cliUtils.InternalError(err.Error())
	}

	color.Blue("Конфигурация успешно обновлена")
	return nil
}

// checkRequiredFlags - проверяет наличие обязательных флагов для установки конфига
func checkRequiredFlags(port int, pyexe, pyScript string) error {
	if port == 0 {
		return fmt.Errorf("порт обязателен для указания")
	}
	if pyexe == "" {
		return fmt.Errorf("python-executable обязателен для указания")
	}
	if pyScript == "" {
		return fmt.Errorf("python-script обязателен для указания")
	}

	return nil
}

func NewConfigSetCmd() *cobra.Command {
	var (
		port             int
		pythonExecutable string
		pythonScript     string
		storagePath      string
		logPath          string
		configPath       string
		useDefault       bool
	)

	cmd := &cobra.Command{
		Use:     "config_set",
		Short:   "Изменить настройки конфигурационного файла",
		Long:    "Изменить настройки конфига. Не все флаги обязательны - меняйте только то, что нужно.",
		Example: "scandoc.exe config_set --port 8080\nscandoc.exe config_set --config my-config.yaml --python-executable \"python3\"",
		RunE: func(cmd *cobra.Command, args []string) error {
			finalConfigPath := configPath
			if finalConfigPath == "" {
				finalConfigPath = config.DefaultConfigPath()
			}

			if useDefault {
				return setupDefaultConfig(finalConfigPath)
			}

			return configSet(finalConfigPath, port, pythonExecutable, pythonScript, storagePath, logPath)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "Путь к конфигурационному файлу (необязательно, по умолчанию: config/config.yaml)")
	cmd.Flags().BoolVar(&useDefault, "default", false, "Установить значения конфига по умолчанию(необязательно)")
	cmd.Flags().IntVarP(&port, "port", "p", 0, "Порт приложения(обязательно)")
	cmd.Flags().StringVar(&pythonExecutable, "python-executable", "", "Путь к Python интерпретатору(обязательно)")
	cmd.Flags().StringVar(&pythonScript, "python-script", "", "Путь к Python скрипту(обязательно)")
	cmd.Flags().StringVarP(&storagePath, "storage-path", "s", "", "Путь к хранилищу данных(необязательно)")
	cmd.Flags().StringVarP(&logPath, "log-path", "l", "", "Путь к логам(необязательно)")

	return cmd
}
