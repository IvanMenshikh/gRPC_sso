package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// негласное ключевое слово Must в Go используется для функций, которые должны вернуть успешный результат или завершить программу с ошибкой (panic).
func MustLoad() *Config {
	path := fetchConfigPath()

	// Проверяем на пустой путь
	if path == "" {
		panic("config path is empty (Пустой путь)")
	}

	// Проверяем, что файл существует
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist (Файла по такому пути не существует): " + path)
	}

	var cfg Config

	// Читаем конфиг из файла по указанному пути
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config (Ошибка чтения): " + err.Error())
	}

	return &cfg
}

// fetchConfigPath извлекает путь к конфигурационному файлу из аргументов командной строки или переменных окружения.
// Приоритет: flag > env > default (flag больше всего приоритетен, затем env, затем дефолтное значение).
// Defolt value - пустая строка
func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "Path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
