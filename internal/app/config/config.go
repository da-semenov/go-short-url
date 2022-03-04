package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"
)

type AppConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
	FileStorage   string `env:"FILE_STORAGE_PATH" envDefault:"./data/storage.txt"`
}

func (config *AppConfig) Init() error {
	pflag.StringVarP(&config.ServerAddress, "a", "a", config.ServerAddress, "Http-server address")
	pflag.StringVarP(&config.BaseURL, "b", "b", config.BaseURL, "Base URL")
	pflag.StringVarP(&config.FileStorage, "f", "f", config.FileStorage, "File storage path")
	pflag.Parse()
	if err := env.Parse(config); err != nil {
		fmt.Println("unable to load server settings", err)
		return err
	}
	return nil
}

func NewConfig() *AppConfig {
	return &AppConfig{}
}
