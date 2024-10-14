package storage

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type config struct {
	DatabaseType string
	DatabaseURL  string
}

func initCfg() (*config, error) {
	const defaultPath = "./config.yaml"

	cfg := &config{}

	err := cleanenv.ReadConfig(defaultPath, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
