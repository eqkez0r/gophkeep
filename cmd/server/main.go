package main

import (
	"context"
	"github.com/eqkez0r/gophkeep/internal/app"
	"github.com/eqkez0r/gophkeep/internal/storage"
	"github.com/eqkez0r/gophkeep/pkg/logger"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	log, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}
	store, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	a, err := app.New(log, store)
	if err != nil {
		log.Fatal(err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go a.Start(ctx, wg)
	<-ctx.Done()
	a.Stop(wg)
	wg.Wait()
}
