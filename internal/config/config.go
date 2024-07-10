package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"environment" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"` //если не будет установлен путь, приложение не запустится
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// MustLoad - функция прочитает файл с конфига, создаст и заполнит объект Config
func MustLoad() *Config {
	// читаем значение переменной окружения CONFIG_PATH из local.yaml
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env variables: %s", err)
	}

	configPath := os.Getenv("CONFIG_PATH") // где находится файл с конфигом, берем из переменной окружения. Можно сделать по флагу
	if configPath == "" {                  // если не находим путь к конфигу, завершаем процесс
		log.Fatal("CONFIG_PATH is not set")
	}

	// проверяем существует ли такой файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Stat возвращает FileInfo, описывающий указанный файл. Если возникнет ошибка, она будет иметь тип *PathError.
		// Is Not Exist возвращает логическое значение, указывающее, известна ли ошибка, сообщающая о том, что файл или каталог не существует.

		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		// ReadConfig считывает конфигурационный файл и анализирует его в зависимости от тегов в
		//предоставленной структуре. Затем он считывает и анализирует

		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
