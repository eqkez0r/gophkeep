package app

import "github.com/ilyakaznacheev/cleanenv"

type config struct {
	AuthServiceAddr   string `json:"auth_service_addr" yaml:"auth_service_addr"`
	KeeperServiceAddr string `json:"keeper_service_addr" yaml:"keeper_service_addr"`
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
