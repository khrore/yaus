package config

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Env        string     `yaml:"env" env:"YAUS_ENV" env-required:"true"`
	DBPath     string     `yaml:"db_path" env:"YAUS_DB_PATH" env-required:"true"`
	HTTPServer HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env:"YAUS_ADDR" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env:"YAUS_TIMEOUT" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_duration" env:"YAUS_IDLE_TIMEOUT" env-required:"true"`
	User        string        `yaml:"user" env:"YAUS_USER" env-required:"true"`
	Password    string        `yaml:"password" env:"YAUS_PASSWORD" env-required:"true"`
}

const (
	DefaultEnv         = "local"
	DefaultDBPath      = "~/.local/share/yaus/yaus.db"
	DefaultAddress     = "localhost:8082"
	DefaultTimeout     = 4 * time.Second
	DefaultIdleTimeout = time.Minute
	DefaultUser        = "user"
	DefaultPassword    = "1234"
)

func MustLoad() Config {
	configPath := os.Getenv("YAUS_CONFIG_PATH")
	if configPath == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			log.Fatal(err)
		}

		configPath = filepath.Join(userConfigDir, "yaus")
		os.Mkdir(configPath, 0755)

		os.Setenv("YAUS_CONFIG_PATH", configPath)
	}

	configFile := filepath.Join(configPath, "config.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		config := Config{
			DefaultEnv,
			DefaultDBPath,
			HTTPServer{
				DefaultAddress,
				DefaultTimeout,
				DefaultIdleTimeout,
				DefaultUser,
				DefaultPassword,
			},
		}
		yaml, err := yaml.Marshal(config)
		if err != nil {
			log.Fatalf("cannot create data for config: %s", err)
		}
		if err := os.WriteFile(configFile, yaml, 0755); err != nil {
			log.Fatalf("cannot create config file: %s", err)
		}
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configFile, &cfg); err != nil && err != io.EOF {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}
