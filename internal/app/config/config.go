package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"os"
)

type AppConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
	FileStorage   string `env:"FILE_STORAGE_PATH" envDefault:"./data/storage.txt"`
}

func (config *AppConfig) Init() error {
	fmt.Println(os.Args)
	if err := env.Parse(config); err != nil {
		fmt.Println("unable to load server settings", err)
		return err
	}

	if config.BaseURL == "" || config.FileStorage == "" || config.ServerAddress == "" {
		if err := env.Parse(&config); err != nil {
			fmt.Println("unable to load server settings", err)
			return err
		}
	}

	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:8080/"
	}

	if config.BaseURL[len(config.BaseURL)-1:] != "/" {
		config.BaseURL += "/"
	}
	return nil
}

func NewConfig() *AppConfig {
	return &AppConfig{}
}
