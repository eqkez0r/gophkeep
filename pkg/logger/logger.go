package logger

import "go.uber.org/zap"

func New() (*zap.SugaredLogger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableStacktrace = true
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
