package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	DefaultPathToConfig = "./config/config.yaml"

	defaultPort          = 3210
	defaultPyExecutable  = ".venv/Scripts/python.exe"
	defaultPyScript      = "./internal/service/scanPy"
	defaultPathToStorage = "storageJSONs"
	defaultPathToLog     = "log/config.log"
)

type Config struct {
	Port             int    `yaml:"port" required:"true"`
	PythonExecutable string `yaml:"python_executable"`
	PythonScript     string `yaml:"python_script"`
	StoragePath      string `yaml:"storage_path"`
	LogPath          string `yaml:"log_path"`
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

	cfg := &Config{
		Port:             defaultPort,
		PythonExecutable: defaultPyExecutable,
		PythonScript:     defaultPyScript,
		StoragePath:      defaultPathToStorage,
		LogPath:          defaultPathToLog,
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("ошибка при чтении конфига: " + err.Error())
	}

	return cfg
}

// MustLoadWithPath - читает конфиг по входящему пути и возвращает структуру конфига для работы приложения
func MustLoadWithPath(pathToConfig string) *Config {
	if _, err := os.Stat(pathToConfig); os.IsNotExist(err) {
		panic("пути к конфигу не найдено: " + pathToConfig)
	}

	cfg := &Config{
		Port:             defaultPort,
		PythonExecutable: defaultPyExecutable,
		PythonScript:     defaultPyScript,
		StoragePath:      defaultPathToStorage,
		LogPath:          defaultPathToLog,
	}

	if err := cleanenv.ReadConfig(pathToConfig, cfg); err != nil {
		panic("ошибка при чтении конфига: " + err.Error())
	}

	return cfg
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
