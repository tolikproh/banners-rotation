package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/tolikproh/banners-rotation/internal/cnst"
)

type Config struct {
	HTTPServer `yaml:"httpserver"`
	Logger     `yaml:"logger"`
	Storage    `yaml:"storage"`
	RabbitMQ   `yaml:"rabbitmq"`
}

type HTTPServer struct {
	Host string `yaml:"host" env:"BANNERS_HTTP_SERVER_HOST"`
	Port string `yaml:"port" env:"BANNERS_HTTP_SERVER_PORT"`
}

type Logger struct {
	Level string `yaml:"level" env:"BANNERS_LOG_LEVEL"`
}

type Storage struct {
	Conn string `yaml:"conn" env:"BANNERS_STORAGE_CONN"`
}

type RabbitMQ struct {
	Address  string `yaml:"address" env:"BANNERS_RABBITMQ_ADDRESS"`
	Exchange string `yaml:"exchange" env:"BANNERS_RABBITMQ_EXCHANGE"`
	Queue    string `yaml:"queue" env:"BANNERS_RABBITMQ_QUEUE"`
}

func defaultConfig() *Config {
	cfg := new(Config)

	// Http Server
	cfg.HTTPServer.Host = "localhost"
	cfg.HTTPServer.Port = "8080"

	// Logger
	cfg.Logger.Level = cnst.LoggerLevelInfo

	// Storage
	cfg.Storage.Conn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	// Bracker RabbitMQ
	cfg.RabbitMQ.Address = "amqp://guest:guest@localhost:5672/"
	cfg.RabbitMQ.Exchange = "banners-rotation"
	cfg.RabbitMQ.Queue = "banners-rotation"

	return cfg
}

// NewConfig Set Default.
func NewConfig(path, file string) *Config {
	cfg := defaultConfig()
	cleanenv.ReadConfig(path+"/"+file, &cfg)
	cleanenv.ReadEnv(cfg)
	return cfg
}
