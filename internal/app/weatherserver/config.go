package weatherserver

import (
	"encoding/json"
	"github.com/valen0k/wb_test_content/internal/app/store"
	"os"
)

type Config struct {
	BindAddr string `json:"bind_addr"`
	LogLevel string `json:"log_level"`
	Store    *store.Config
	APIKey   string `json:"api_key"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: "localhost:8080",
		LogLevel: "debug",
		Store:    store.NewConfig(),
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
