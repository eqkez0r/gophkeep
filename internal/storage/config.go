package storage

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type config struct {
	DatabaseType string `json:"database_type" yaml:"database_type" env:"database_type"`
	DatabaseURL  string `json:"database_url" yaml:"database_url" env:"database_url"`
}

func initConfig() (*config, error) {
	const defaultPath = "./config.yaml"

	cfg := &config{}

	err := cleanenv.ReadConfig(defaultPath, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
