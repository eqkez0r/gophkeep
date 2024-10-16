package app

import (
	"context"
	"github.com/eqkez0r/gophkeep/internal/services"
	"github.com/eqkez0r/gophkeep/internal/services/auth"
	"github.com/eqkez0r/gophkeep/internal/services/gophkeep"
	"github.com/eqkez0r/gophkeep/internal/storage"
	"go.uber.org/zap"
	"sync"
)

type App struct {
	logger   *zap.SugaredLogger
	services []services.Service
}

func New(
	logger *zap.SugaredLogger,
	store storage.Storage,
) (*App, error) {
	cfg, err := initConfig()
	if err != nil {
		return nil, err
	}
	s := []services.Service{
		auth.New(logger, store, cfg.AuthServiceAddr),
		gophkeep.New(logger, store, cfg.KeeperServiceAddr),
	}
	return &App{
		logger:   logger,
		services: s,
	}, nil
}

func (a *App) Start(ctx context.Context, wg *sync.WaitGroup) {
	a.logger.Info("start services")
	for _, s := range a.services {
		wg.Add(1)
		go s.Run(ctx, wg)
	}
}

func (a *App) Stop(wg *sync.WaitGroup) {
	a.logger.Info("stop services")
	for _, s := range a.services {
		s.GracefulShutdown()
		wg.Done()
	}
	a.logger.Info("all services was finished")
	wg.Done()
}
