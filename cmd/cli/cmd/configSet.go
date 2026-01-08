package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"proWeb/internal/config"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "C:\\Users\\Administrator\\Desktop\\githubTest\\PP2025-ProWeb\\config\\config.yaml"

// ConfigSet - команда, позволяющая менять значения конфига внутри файла
func ConfigSet(configPath string, port int, pythonExecutable, pythonScript, storagePath, logPath string) error {
	if port == 0 {
		return fmt.Errorf("порт обязателен для указания")
	}
	if pythonExecutable == "" {
		return fmt.Errorf("python-executable обязателен для указания")
	}
	if pythonScript == "" {
		return fmt.Errorf("python-script обязателен для указания")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("конфигурационный файл не найден: %s", configPath)
	}

	cfg, err := loadConfig(configPath)
	if err != nil {
		return fmt.Errorf("не удалось загрузить конфиг: %w", err)
	}

	cfg.Port = port
	cfg.PythonExecutable = pythonExecutable
	cfg.PythonScript = pythonScript

	if storagePath != "" {
		cfg.StoragePath = storagePath
	}

	if logPath != "" {
		cfg.LogPath = logPath
	}

	if err := saveConfig(cfg, configPath); err != nil {
		return fmt.Errorf("не удалось сохранить конфиг: %w", err)
	}

	fmt.Println("Конфигурация успешно обновлена")
	return nil
}

// loadConfig - чтение конфига из файла
func loadConfig(path string) (*config.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка разбора YAML: %w", err)
	}

	return &cfg, nil
}

// setupDefaultConfig - устанавливает базовые значения для конфига
func setupDefaultConfig() (cfg *config.Config, err error) {
	panic("implement me")
}

// saveConfig - сохраняет конфиг в файл
func saveConfig(cfg *config.Config, path string) error {
	var yamlLines []string

	yamlLines = append(yamlLines, fmt.Sprintf("port: %d", cfg.Port))
	yamlLines = append(yamlLines, fmt.Sprintf("python_executable: \"%s\"", escapeYAMLString(cfg.PythonExecutable)))
	yamlLines = append(yamlLines, fmt.Sprintf("python_script: \"%s\"", escapeYAMLString(cfg.PythonScript)))
	if cfg.StoragePath != "" {
		yamlLines = append(yamlLines, fmt.Sprintf("storage_path: \"%s\"", escapeYAMLString(cfg.StoragePath)))
	}

	if cfg.LogPath != "" {
		yamlLines = append(yamlLines, fmt.Sprintf("log_path: \"%s\"", escapeYAMLString(cfg.LogPath)))
	}

	yamlContent := strings.Join(yamlLines, "\n")

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return fmt.Errorf("ошибка создания директории: %w", err)
	}

	if err := os.WriteFile(path, []byte(yamlContent), 0777); err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}

	return nil
}

// escapeYAMLString - экранирует специальные символы в строке
func escapeYAMLString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

func NewConfigSetCmd() *cobra.Command {
	var (
		port             int
		pythonExecutable string
		pythonScript     string
		storagePath      string
		logPath          string
		configPath       string
	)

	cmd := &cobra.Command{
		Use:     "config_set",
		Short:   "Изменить настройки конфигурационного файла",
		Long:    "Изменить настройки конфига. Не все флаги обязательны - меняйте только то, что нужно.",
		Example: "scanner.exe config_set --port 8080\nscanner.exe config_set --config my-config.yaml --python-executable \"python3\"",
		RunE: func(cmd *cobra.Command, args []string) error {
			finalConfigPath := configPath
			if finalConfigPath == "" {
				finalConfigPath = defaultConfigPath
			}

			return ConfigSet(finalConfigPath, port, pythonExecutable, pythonScript, storagePath, logPath)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "Путь к конфигурационному файлу (необязательно, по умолчанию: config/config.yaml)")
	cmd.Flags().IntVarP(&port, "port", "p", 0, "Порт приложения(обязательно)")
	cmd.Flags().StringVar(&pythonExecutable, "python-executable", "", "Путь к Python интерпретатору(обязательно)")
	cmd.Flags().StringVar(&pythonScript, "python-script", "", "Путь к Python скрипту(обязательно)")
	cmd.Flags().StringVarP(&storagePath, "storage-path", "s", "", "Путь к хранилищу данных(необязательно)")
	cmd.Flags().StringVarP(&logPath, "log-path", "l", "", "Путь к логам(необязательно)")

	return cmd
}
