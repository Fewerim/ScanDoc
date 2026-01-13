package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

const (
	DefaultPathToConfig = "./config/config.yaml"

	DefaultPort          = 3210
	DefaultPyVenv        = ""
	DefaultPyExecutable  = ".venv/Scripts/python.exe"
	DefaultPyScript      = "/lib/service/scanPy"
	DefaultPathToStorage = "storageJSONs"
	DefaultPathToLog     = "log/info.log"
)

type Config struct {
	Port             int    `yaml:"port" required:"true"`
	PythonVenvPath   string `yaml:"python_venv" required:"true"`
	PythonExecutable string `yaml:"python_executable"`
	PythonScript     string `yaml:"python_script"`
	StoragePath      string `yaml:"storage_path,omitempty"`
	LogPath          string `yaml:"log_path,omitempty"`
}

// NewDefaultConfig - конструктор дефолтного конфига
func NewDefaultConfig() *Config {
	return &Config{
		Port:             DefaultPort,
		PythonVenvPath:   DefaultPyVenv,
		PythonExecutable: DefaultPyExecutable,
		PythonScript:     DefaultPyScript,
		StoragePath:      DefaultPathToStorage,
		LogPath:          DefaultPathToLog,
	}
}

// SetupDefaultConfig - устанавливает дефолтные значения для конфига
func (cfg *Config) SetupDefaultConfig() {
	projectRoot, err := FindProjectRoot(".")
	if err != nil {
		return
	}

	cfg.Port = DefaultPort
	cfg.PythonVenvPath = filepath.Join(projectRoot, DefaultPyVenv)
	cfg.PythonExecutable = filepath.Join(projectRoot, DefaultPyExecutable)
	cfg.PythonScript = filepath.Join(projectRoot, DefaultPyScript)
	cfg.StoragePath = filepath.Join(projectRoot, DefaultPathToStorage)
	cfg.LogPath = filepath.Join(projectRoot, DefaultPathToLog)
}

// SaveConfig - сохраняет конфиг в файл
func (cfg *Config) SaveConfig(path string) error {
	var yamlLines []string

	yamlLines = append(yamlLines, fmt.Sprintf("port: %d", cfg.Port))
	yamlLines = append(yamlLines, fmt.Sprintf("python_venv: \"%s\"", escapeYAMLString(cfg.PythonVenvPath)))
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

// MustLoad - читает конфиг и возвращает структуру конфига для работы приложения
func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("путь к конфигу пустой")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("пути к конфигу не найдено: " + path)
	}

	cfg := NewDefaultConfig()

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("ошибка при чтении конфига: " + err.Error())
	}

	return cfg
}

// LoadConfig - чтение конфига из файла
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка разбора YAML: %w", err)
	}

	return &cfg, nil
}

// MustLoadWithPath - читает конфиг по входящему пути и возвращает структуру конфига для работы приложения
func MustLoadWithPath(pathToConfig string) *Config {
	if _, err := os.Stat(pathToConfig); os.IsNotExist(err) {
		panic("пути к конфигу не найдено: " + pathToConfig)
	}

	cfg := NewDefaultConfig()

	if err := cleanenv.ReadConfig(pathToConfig, cfg); err != nil {
		panic("ошибка при чтении конфига: " + err.Error())
	}

	return cfg
}

// DefaultConfigPath - возвращает путь к конфиг файлу, чтобы приложение смогло инициализировать конфиг
func DefaultConfigPath() string {
	var rootDir string

	exePath, err := os.Executable()
	if err != nil {
		return DefaultPathToConfig
	}

	if strings.Contains(exePath, "\\scandoc-gui\\build") {
		exeDir := filepath.Dir(exePath)
		buildDir := filepath.Dir(exeDir)
		guiDir := filepath.Dir(buildDir)
		cmdDir := filepath.Dir(guiDir)
		rootDir = filepath.Dir(cmdDir)
	} else {
		exeDir := filepath.Dir(exePath)
		rootDir = filepath.Dir(exeDir)
	}

	return filepath.Join(rootDir, "config", "config.yaml")
}

// CheckConfigPathExists - проверяет наличие конфига по пути
func CheckConfigPathExists(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("конфигурационный файл не найден: %s", configPath)
	}
	return nil
}

// fetchConfigPath - достает путь к конфигу через флаг в командной строке
// priority: flag > env > default
// default: DefaultPathToConfig
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "config file path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	if res == "" {
		res = DefaultPathToConfig
	}

	return res
}

// escapeYAMLString - экранирует специальные символы в строке
func escapeYAMLString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

func FindProjectRoot(start string) (string, error) {
	dir := start
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			projectRoot, err := filepath.Abs(dir)
			if err != nil {
				return "", err
			}
			return projectRoot, nil
		}
		dir = filepath.Dir(dir)
	}
	return "", fmt.Errorf("корень проекта не найден")
}
