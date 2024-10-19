package services

import (
	"context"
	"sync"
)

type Service interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	GracefulShutdown()
}
