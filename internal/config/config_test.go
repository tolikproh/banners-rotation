package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/tolikproh/banners-rotation/internal/cnst"
)

type configData struct {
	name string
	path string
	file string
	env  map[string]string
	exp  *Config
}

func initData() *[]configData {
	return &[]configData{
		{
			name: "default config",
			path: "",
			file: "",
			env:  map[string]string{},
			exp: &Config{
				HTTPServer: HTTPServer{
					Host: "localhost",
					Port: "8080",
				},
				Logger: Logger{
					Level: cnst.LoggerLevelInfo,
				},
				Storage: Storage{
					Conn: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
				},
				RabbitMQ: RabbitMQ{
					Address:  "amqp://guest:guest@localhost:5672/",
					Exchange: "banners-rotation",
					Queue:    "banners-rotation",
				},
			},
		},
		{
			name: "config file yaml",
			path: "./testdata",
			file: "config_test.yaml",
			env:  map[string]string{},
			exp: &Config{
				HTTPServer: HTTPServer{
					Host: "banners.ru",
					Port: "8000",
				},
				Logger: Logger{
					Level: cnst.LoggerLevelError,
				},
				Storage: Storage{
					Conn: "postgres://user:password@localhost:5432/banners_db?sslmode=disable",
				},
				RabbitMQ: RabbitMQ{
					Address:  "amqp://test:pswd@127.0.0.1:5561/",
					Exchange: "banners_exchange_config",
					Queue:    "banners_queue_config",
				},
			},
		},
		{
			name: "config environment",
			path: "./testdata",
			file: "config_test.yaml",
			env: map[string]string{
				"BANNERS_HTTP_SERVER_HOST":  "env.net",
				"BANNERS_HTTP_SERVER_PORT":  "1234",
				"BANNERS_LOG_LEVEL":         "warn",
				"BANNERS_STORAGE_CONN":      "http://memory.com",
				"BANNERS_RABBITMQ_ADDRESS":  "amqp://env:envpswd@localhost:4321/",
				"BANNERS_RABBITMQ_EXCHANGE": "banners_env_exchange",
				"BANNERS_RABBITMQ_QUEUE":    "banners_env_queue",
			},
			exp: &Config{
				HTTPServer: HTTPServer{
					Host: "env.net",
					Port: "1234",
				},
				Logger: Logger{
					Level: cnst.LoggerLevelWarn,
				},
				Storage: Storage{
					Conn: "http://memory.com",
				},
				RabbitMQ: RabbitMQ{
					Address:  "amqp://env:envpswd@localhost:4321/",
					Exchange: "banners_env_exchange",
					Queue:    "banners_env_queue",
				},
			},
		},
		{
			name: "config file and secret environment",
			path: "./testdata",
			file: "config_test.yaml",
			env: map[string]string{
				"BANNERS_STORAGE_CONN":     "postgres://user:secret@localhost:5555",
				"BANNERS_RABBITMQ_ADDRESS": "amqp://env:envpswd@localhost:5678/",
			},
			exp: &Config{
				HTTPServer: HTTPServer{
					Host: "banners.ru",
					Port: "8000",
				},
				Logger: Logger{
					Level: cnst.LoggerLevelError,
				},
				Storage: Storage{
					Conn: "postgres://user:secret@localhost:5555",
				},
				RabbitMQ: RabbitMQ{
					Address:  "amqp://env:envpswd@localhost:5678/",
					Exchange: "banners_exchange_config",
					Queue:    "banners_queue_config",
				},
			},
		},
	}
}

func TestConfig(t *testing.T) {
	testCases := initData()
	os.Clearenv()
	for _, tc := range *testCases {
		t.Run(tc.name, func(t *testing.T) {
			for env, val := range tc.env {
				os.Setenv(env, val)
			}
			defer os.Clearenv()

			cfg := NewConfig(tc.path, tc.file)

			if !reflect.DeepEqual(cfg, tc.exp) {
				t.Errorf("wrong data %v, exp %v", cfg, tc.exp)
			}
		})
	}
}
