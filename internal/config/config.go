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
	Mail  Mail        `yaml:"mail"`
	Token Token       `yaml:"token"`
}

type Application struct {
	Env         string        `yaml:"env" env-required:"true"`
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"5m"`
}

type Database struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	Name     string `yaml:"name" env-required:"true"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWD" env-required:"true"`
}

type Mail struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `env:"MAIL_USER" env-required:"true"`
	Password string `env:"MAIL_PASSWD" env-required:"true"`
}

type Token struct {
	AccessTTL  time.Duration `yaml:"access_ttl" env-default:"1m"`
	RefreshTTL time.Duration `yaml:"refresh_ttl" env-default:"5m"`
}

func MustLoad() Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("empty config-path-env")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("the configuration file does not exist at the specified path: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed with reading config: %s", err)
	}
	return cfg
}
