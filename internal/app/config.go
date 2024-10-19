package app

import "github.com/ilyakaznacheev/cleanenv"

type config struct {
	KeeperServiceAddr string `json:"keeper_service_addr" yaml:"keeper_service_addr"`
}

func initConfig(path string) (*config, error) {
	const defaultPath = "./config.yaml"
	if path == "" {
		path = defaultPath
	}
	cfg := &config{}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
