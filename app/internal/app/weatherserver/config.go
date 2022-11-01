package weatherserver

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	BindIP     string `env:"BIND_IP" env-default:"0.0.0.0"`
	Port       string `env:"PORT" env-default:"8081"`
	LogLevel   string `env:"LOG_LEVEL" env-default:"info"`
	APIKey     string `env:"API_KEY" env-required:"true"`
	PostgreSQL Database
}

type Database struct {
	Username string `env:"PSQL_USERNAME" env-required:"true"`
	Password string `env:"PSQL_PASSWORD" env-required:"true"`
	Host     string `env:"PSQL_HOST" env-required:"true"`
	Port     string `env:"PSQL_PORT" env-required:"true"`
	Database string `env:"PSQL_DATABASE" env-required:"true"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		log.Println("gather config")

		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			log.Fatalln(err)
		}
	})

	return instance
}
