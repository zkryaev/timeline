package config

import (
	"log"
	"os"
	"time"
	"timeline/internal/utils/envars"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App        Application `yaml:"app"`
	DB         Database
	Mail       Mail
	Token      Token `yaml:"token"`
	S3         S3
	Prometheus Prometheus
	Analytics  AnalyticsService
}

type Settings struct {
	UseLocalBackData    bool `yaml:"use_local_back_data"`
	EnableAuthorization bool `yaml:"enable_authorization"`
	EnableMedia         bool `yaml:"enable_media"`
	EnableMail          bool `yaml:"enable_mail"`
	EnableMetrics       bool `yaml:"enable_metrics"`
	EnableAnalytics     bool `yaml:"enable_analytics"`
}

type Application struct {
	Server   HTTPServer
	Settings Settings `yaml:"settings" env-required:"true"`
	Stage    string   `yaml:"stage"`
}

type HTTPServer struct {
	Host        string        `env:"APP_HTTP_HOST" env-default:"localhost"`
	Port        string        `env:"APP_HTTP_PORT" env-required:"true"`
	Timeout     time.Duration `env:"APP_HTTP_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"APP_HTTP_IDLE_TIMEOUT" env-default:"5m"`
}

type Database struct {
	Protocol string `env:"DB" env-required:"true"`
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     string `env:"PGB_PORT" env-required:"true"`
	Name     string `env:"DB_NAME" env-required:"true"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWD" env-required:"true"`
	SSLmode  string `env:"DB_SSLMODE" env-required:"true"`
}

type Mail struct {
	Service  string `env:"MAIL_SERVICE" env-required:"true"`
	Host     string `env:"MAIL_HOST" env-required:"true"`
	Port     int    `env:"MAIL_PORT" env-required:"true"`
	User     string `env:"MAIL_USER" env-required:"true"`
	Password string `env:"MAIL_PASSWD" env-required:"true"`
}

type Token struct {
	AccessTTL  time.Duration `yaml:"access_ttl" env-default:"1m"`
	RefreshTTL time.Duration `yaml:"refresh_ttl" env-default:"5m"`
}

type S3 struct {
	Name          string `env:"S3" env-required:"true"`
	Host          string `env:"S3_HOST" env-required:"true"`
	User          string `env:"S3_ROOT_USER" env-required:"true"`
	Password      string `env:"S3_ROOT_PASSWORD" env-required:"true"`
	DefaultBucket string `env:"S3_DEFAULT_BUCKET" env-required:"true"`
	DataPort      string `env:"S3_DATA_PORT" env-default:"9000"`
	ConsolePort   string `env:"S3_CONSOLE_PORT" env-default:"9001"`
	SSLmode       bool   `env:"S3_SSLMODE" env-default:"false"`
}

type AnalyticsService struct {
	Host string `env:"ANALYTICS_HOST" env-default:"localhost"`
	Port string `env:"ANALYTICS_PORT" env-required:"true"`
}

type Prometheus struct {
	Host string `env:"PROMETHEUS_PRODUCER_HOST" env-required:"true"`
	Port string `env:"PROMETHEUS_PRODUCER_PORT" env-required:"true"`
}

func MustLoad() *Config {
	configPath := envars.GetPathByEnv("CONFIG_PATH")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("the cfg file doesn't exist at the path: ", configPath)
	}
	cfg := &Config{}
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatal("failed read config: ", err.Error())
	}
	return cfg
}
