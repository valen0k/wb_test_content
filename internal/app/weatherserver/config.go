package weatherserver

import (
	"encoding/json"
	"os"
)

type Config struct {
	BindAddr    string `json:"bind_addr"`
	LogLevel    string `json:"log_level"`
	DatabaseURL string `json:"database_url"`
	APIKey      string `json:"api_key"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: "localhost:8080",
		LogLevel: "debug",
	}
}

func (c *Config) DecodeConfigFile(configFile string) error {
	file, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(file, c); err != nil {
		return err
	}

	return nil
}
