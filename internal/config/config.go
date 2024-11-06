package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App   Application `yaml:"app"`
	DB    Database    `yaml:"db"`
	Token Token       `yaml:"token"`
}

type Application struct {
	Env         string        `yaml:"env" env-required:"true"`
	Host        string        `yaml:"host" env:"HOST" env-default:"localhost"`
	Port        string        `yaml:"port" env:"PORT" env-default:"5432"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"5m"`
}

type Database struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"DB_NAME" env-required:"true"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `env:"DB_PASSWD" env-required:"true"`
}

type Token struct {
	AccessTTL  time.Duration `yaml:"access_ttl" env-default:"1m"`
	RefreshTTL time.Duration `yaml:"refresh_ttl" env-default:"5m"`
}

func MustLoad() Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return cfg
}
